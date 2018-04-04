package main

import (
	"coderangers/connections"
	"coderangers/dao"
	"coderangers/emails"
	"coderangers/judge"
	"coderangers/leaderboards"
	"coderangers/notifications"
	"coderangers/skills"
	"coderangers/templating"
	"coderangers/users"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

/*
TODO:
- reorganize / structure the files better (ex a package for submissions?)
- deal with errors! (log if fatal etc)
- use files for the problems and submissions. or use a db.. json. preload. edits?
- seperate files for templates.. a little design?
- implement the memory stuff and all? hm.. checks.. errors.. validation..
*/

var config = getConfig()

type Configuration struct {
	AdminUsername string
	AdminPassword string
	AdminEmail    string
	Create        bool
	Domain        []string
	Https         bool
}

func getConfig() Configuration {
	file, _ := os.Open("admin.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("ADMIN CONFIG ERROR", err)
	}
	return configuration
}

func page(content string) string {
	return "<html><body>\n" + content + "\n</body></html>"
}

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error log file", err)
	}
	defer f.Close()
	log.SetOutput(f)
	err = dao.CreateDB()
	if err == nil {
		dao.AddTables()
		skills.AddSamples()
		err = judge.AddSamples()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("New database created")
	}
	if config.Create {
		_, err = users.Register(config.AdminUsername, config.AdminPassword, config.AdminEmail, "", "", true)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println("New admin account created")
		}
	}
	templating.InitTemplates()
	wd, _ := os.Getwd()
	judge.DIR = filepath.Join(wd, "submissions")
	os.Mkdir(judge.DIR, 0777)
	judge.InitQueues()

	mux := http.NewServeMux()

	mux.HandleFunc("/", judge.HomeHandler)
	mux.HandleFunc("/problems", judge.ProblemsHandler)
	mux.HandleFunc("/random-problem", judge.RandomHandler)
	mux.HandleFunc("/add-problem", judge.AddHandler)
	mux.HandleFunc("/edit-problem/", judge.EditHandler)
	mux.HandleFunc("/delete", judge.DeleteHandler)
	mux.HandleFunc("/view/", judge.ViewHandler)
	mux.HandleFunc("/submit/", judge.SubmitHandler)
	mux.HandleFunc("/submissions/", judge.SubmissionsHandler)
	mux.HandleFunc("/my-submissions/", judge.MySubmissionsHandler)
	mux.HandleFunc("/skill-summary/", judge.SkillSummaryHandler)
	mux.HandleFunc("/runtime-error/", judge.RuntimeErrorHandler)
	mux.HandleFunc("/register", users.RegisterHandler)
	mux.HandleFunc("/login", users.LoginHandler)
	mux.HandleFunc("/logout", users.LogoutHandler)
	mux.HandleFunc("/add-name", users.AddNameHandler)
	mux.HandleFunc("/change-password", users.ChangePasswordHandler)
	mux.HandleFunc("/add-skill", skills.AddSkillHandler)
	mux.HandleFunc("/edit-skill/", skills.EditSkillHandler)

	mux.HandleFunc("/leaderboards", leaderboards.LeaderboardsHandler)
	mux.HandleFunc("/"+notifications.Notifications, notifications.InitHandler().ServeHTTP)
	mux.HandleFunc("/"+notifications.Submissions, notifications.InitHandler().ServeHTTP)
	mux.HandleFunc("/viewed-notification", notifications.ViewedHandler)

	mux.HandleFunc("/connect", connections.CheckHandler)
	mux.HandleFunc("/profile", users.ViewProfileHandler)
	mux.HandleFunc("/profile/", users.ViewUserProfileHandler)
	mux.HandleFunc("/skill/", skills.SkillHandler)
	mux.HandleFunc("/skill-tree/", skills.SkillTreeHandler)
	mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	// http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts"))))

	mux.HandleFunc("/email-logo", emails.EmailLogoHandler)
	emails.SendEmailsEvery(24 * time.Hour)

	fmt.Println("RESEND: ", judge.ResendReceivedAndCheckInqueue())

	log.Println("Start")

	if config.Https {

		fmt.Println(config.Domain)
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.Domain...), //your domain here
			Cache:      autocert.DirCache("certs"),               //folder for storing certificates
		}

		server := &http.Server{
			Addr: ":https",
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
			Handler: mux,
		}

		fmt.Println("serving https and redirecting http")
		go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		server.ListenAndServeTLS("", "")
	} else {
		fmt.Println("serving http")
		http.ListenAndServe(":80", mux)
	}

	db, _ := dao.Open()
	db.Close()
}
