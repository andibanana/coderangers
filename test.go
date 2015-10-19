package main

import (
	"./dao"
	"./judge"
	"./leaderboards"
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
		users.Register("admin", "admin", true)
		judge.AddSamples()
	}
	templating.InitTemplates()
	wd, _ := os.Getwd()
	judge.DIR = filepath.Join(wd, "submissions")
	os.Mkdir(judge.DIR, 0777)
	judge.InitQueues()
	http.HandleFunc("/", judge.ProblemsHandler)
	http.HandleFunc("/problems/", judge.ProblemsHandler)
	http.HandleFunc("/add/", judge.AddHandler)
	http.HandleFunc("/view/", judge.ViewHandler)
	http.HandleFunc("/submit/", judge.SubmitHandler)
	http.HandleFunc("/submissions/", judge.SubmissionsHandler)
	http.HandleFunc("/register", users.RegisterHandler)
	http.HandleFunc("/login", users.LoginHandler)
	http.HandleFunc("/logout", users.LogoutHandler)
	http.HandleFunc("/leaderboards", leaderboards.LeaderboardsHandler)
	http.HandleFunc("/buy_hint", judge.BuyHintHandler)

	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles"))))

	http.ListenAndServe(":80", nil)
}
