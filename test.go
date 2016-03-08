package main

import (
	"./dao"
	"./judge"
	"./leaderboards"
	"./skills"
	"./templating"
	"./users"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	err := dao.CreateDB()
	fmt.Println(err)

	if err == nil {
		users.Register("admin", "admin", "frzsk@yahoo.com", true)
		// users.RegisterAndFakeData("FCsean", "FCsean", false, 500, 50)
		// users.RegisterAndFakeData("gopherzapper_", "gopherzapper_", false, 12300, 1230)
		// users.RegisterAndFakeData("DarkMega12", "DarkMega12", false, 100000, 10000)
		//users.RegisterAndFakeData("gmg", "gmg", false, 3230, 323)
		judge.AddSamples()
		skills.AddSamples()
	}
	templating.InitTemplates()
	wd, _ := os.Getwd()
	judge.DIR = filepath.Join(wd, "submissions")
	os.Mkdir(judge.DIR, 0777)
	judge.InitQueues()
	http.HandleFunc("/", judge.HomeHandler)
	http.HandleFunc("/problems", judge.ProblemsHandler)
	http.HandleFunc("/add-problem", judge.AddHandler)
	http.HandleFunc("/edit/", judge.EditHandler)
	http.HandleFunc("/delete", judge.DeleteHandler)
	http.HandleFunc("/view/", judge.ViewHandler)
	http.HandleFunc("/submit/", judge.SubmitHandler)
	http.HandleFunc("/submissions/", judge.SubmissionsHandler)
	http.HandleFunc("/register", users.RegisterHandler)
	http.HandleFunc("/login", users.LoginHandler)
	http.HandleFunc("/logout", users.LogoutHandler)
	http.HandleFunc("/add-skill", skills.AddSkillHandler)
	http.HandleFunc("/edit-skill/", skills.EditSkillHandler)

	http.HandleFunc("/leaderboards", leaderboards.LeaderboardsHandler)
	http.HandleFunc("/profile", users.ViewProfileHandler)
	http.HandleFunc("/profile/", users.ViewUserProfileHandler)
	http.HandleFunc("/skill/", skills.SkillHandler)
	http.HandleFunc("/skill-tree/", skills.SkillTreeHandler)
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	// http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts"))))

	http.ListenAndServe(":80", nil)
}
