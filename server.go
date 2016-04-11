package main

import (
	"coderangers/dao"
	"coderangers/emails"
	"coderangers/judge"
	"coderangers/leaderboards"
	"coderangers/notifications"
	"coderangers/skills"
	"coderangers/templating"
	"coderangers/users"
	"fmt"
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

func page(content string) string {
	return "<html><body>\n" + content + "\n</body></html>"
}

func main() {
	log.SetFlags(log.Llongfile)
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error log file", err)
	}
	defer f.Close()
	log.SetOutput(f)
	err = dao.CreateDB()
	fmt.Println(err)

	if err == nil {
		_, err = users.Register("admin", "admin", "frzsk@yahoo.com", true)
		if err != nil {
			log.Println(err)
		}
		skills.AddSamples()
		err = judge.AddSamples()
		if err != nil {
			log.Println(err)
		}
	}
	templating.InitTemplates()
	wd, _ := os.Getwd()
	judge.DIR = filepath.Join(wd, "submissions")
	os.Mkdir(judge.DIR, 0777)
	judge.InitQueues()
	fmt.Println("RESEND: ", judge.ResendReceivedAndCheckInqueue())

	mux := http.NewServeMux()

	mux.HandleFunc("/", judge.HomeHandler)
	mux.HandleFunc("/problems", judge.ProblemsHandler)
	mux.HandleFunc("/add-problem", judge.AddHandler)
	mux.HandleFunc("/edit-problem/", judge.EditHandler)
	mux.HandleFunc("/delete", judge.DeleteHandler)
	mux.HandleFunc("/view/", judge.ViewHandler)
	mux.HandleFunc("/submit/", judge.SubmitHandler)
	mux.HandleFunc("/submissions/", judge.SubmissionsHandler)
	mux.HandleFunc("/register", users.RegisterHandler)
	mux.HandleFunc("/login", users.LoginHandler)
	mux.HandleFunc("/logout", users.LogoutHandler)
	mux.HandleFunc("/change-password", users.ChangePasswordHandler)
	mux.HandleFunc("/add-skill", skills.AddSkillHandler)
	mux.HandleFunc("/edit-skill/", skills.EditSkillHandler)

	mux.HandleFunc("/leaderboards", leaderboards.LeaderboardsHandler)
	mux.HandleFunc("/"+notifications.Notifications, notifications.InitHandler().ServeHTTP)
	mux.HandleFunc("/"+notifications.Submissions, notifications.InitHandler().ServeHTTP)
	mux.HandleFunc("/profile", users.ViewProfileHandler)
	mux.HandleFunc("/profile/", users.ViewUserProfileHandler)
	mux.HandleFunc("/skill/", skills.SkillHandler)
	mux.HandleFunc("/skill-tree/", skills.SkillTreeHandler)
	mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	// http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts"))))

	emails.SendEmailsEvery(24 * time.Hour)

	fmt.Println("serving")
	log.Println("Start")
	http.ListenAndServe(":80", mux)
	db, _ := dao.Open()
	db.Close()
}
