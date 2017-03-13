package emails

import (
	"coderangers/dao"
	"coderangers/helper"
	"coderangers/judge"
	"coderangers/problems"
	"coderangers/skills"
	"coderangers/users"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

var config = getConfig()

type Configuration struct {
	Email       string
	Password    string
	SmtpAuth    string
	SmtpAddress string
	Domain      string
}

type ExtraConfig struct {
	Domain string
}

func getConfig() Configuration {
	file, _ := os.Open("email.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("DOMAIN ERROR", err)
	}
	file, _ = os.Open("admin.json")
	decoder = json.NewDecoder(file)
	extraConfig := ExtraConfig{}
	err = decoder.Decode(&extraConfig)
	if err != nil {
		log.Fatal("DOMAIN ERROR", err)
	}
	configuration.Domain = extraConfig.Domain
	return configuration
}

func SendEmailsToInactive() (err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}

	rows, err := db.Query(`SELECT user_account.id, username, email, IFNULL(max(timestamp), "2006-01-02 15:04:05.999999999+08:00") FROM
                        user_account
                          LEFT JOIN
                        submissions
                          ON(user_account.id = submissions.user_id)
                        GROUP BY user_account.id;`)

	if err != nil {
		return
	}
	for rows.Next() {
		var timestamp string
		var email string
		var username string
		var userID int
		err = rows.Scan(&userID, &username, &email, &timestamp)
		if err != nil {
			return
		}
		err = SendInactiveEmail(userID, username, email, timestamp)
		if err != nil {
			return
		}
	}
	return
}

func SendInactiveEmail(userID int, username, email, timestamp string) (err error) {
	submittime, err := helper.ParseTime(timestamp)
	if err != nil {
		log.Println(err)
		return
	}
	duration := time.Since(submittime)
	var days = int(math.Floor(duration.Hours() / 24))
	if days%7 == 3 || (days != 0 && days%7 == 0) {
		var unsolvedProblems []int
		var userDataOnSkill map[string]*skills.Skill
		var problem problems.Problem
		var suggestSkill bool
		var skill *skills.Skill
		var unlockedSkills map[string]bool
		unsolvedProblems, err = judge.GetUnsolvedTriedProblems(userID)
		if err != nil {
			return
		}
		userDataOnSkill, err = skills.GetUserDataOnSkills(userID)
		if err != nil {
			return
		}
		unlockedSkills, err = skills.GetUnlockedSkills(userID)
		if err != nil {
			return
		}
		for _, element := range userDataOnSkill {
			if element.Mastered {
				continue
			}
			if element.Learned || unlockedSkills[element.ID] {
				skill = element
				suggestSkill = true
				var problems []problems.Problem
				problems, err = skills.GetProblemsInSkill(skill.ID)
				skill.NumberOfProblems = len(problems)
				if err != nil {
					return
				}
				break
			}
		}
		message := `<img src="http://` + config.Domain + `/email-logo?u=` + fmt.Sprintf("%v", userID) + `&t=` + fmt.Sprintf("%v", time.Now().Unix()) + `" style="max-width:300px;">`
		message += "<h1>Hi " + username + "</h1>"
		message += "You've been inactive for a few days!<br>We want you back, here are a few things you can do!<br>"
		if len(unsolvedProblems) != 0 {
			var user users.UserData
			problem, err = judge.GetProblem(unsolvedProblems[0])
			message += `<div style="background-color:#DBDBDB;"><a href="http://` + config.Domain + `/view/` + fmt.Sprintf("%d", problem.Index) + `?mail=true"><h2>` + problem.Title + `</h2></a>`
			message += `You can try to solve this problem!<br>`
			user, err = judge.GetUserWhoRecentlySolvedProblem(userID, unsolvedProblems[0])
			if err == nil && len(user.Username) != 0 {
				message += `<a href="http://` + config.Domain + `/profile/` + fmt.Sprintf("%d", user.ID) + `?mail=true">` + user.Username + `</a> recently solved this.<br>`
			}
			message += "</div>"
		} else {
			unlockedProblems, err := judge.GetUnsolvedUnlockedProblem(userID)
			if err != nil {
				return err
			}
			var user users.UserData
			if len(unlockedProblems) != 0 {
				problem = unlockedProblems[0]
				message += `<div style="background-color:#DBDBDB;"><a href="http://` + config.Domain + `/view/` + fmt.Sprintf("%d", problem.Index) + `?mail=true"><h2>` + problem.Title + `</h2></a>`
				message += `You can try to solve this problem!<br>`
				user, err = judge.GetUserWhoRecentlySolvedProblem(userID, problem.Index)
				if err == nil && len(user.Username) != 0 {
					message += `<a href="http://` + config.Domain + `/profile/` + fmt.Sprintf("%d", user.ID) + `?mail=true">` + user.Username + `</a> recently solved this.<br>`
				}
			}
		}
		if suggestSkill {
			message += `<div style="background-color:#DBDBDB;"><a href="http://` + config.Domain + `/skill/` + skill.ID + `?mail=true">` + `<div style="display:inline-block;"><img src="http://` + config.Domain + `/images/skill-icons/` + skill.ID + `.png" style="vertical-align:middle;max-width:100px;"></div><div style="display:inline-block;vertical-align:middle;"><h2 style="display:inline;">` + skill.Title + "</h2><br></a>"
			message += skill.Description + "<br></div><br>"
			if skill.Learned {
				message += "You should try to master this skill. Solve " + strconv.Itoa(skill.NumberOfProblems-skill.Solved) + " more problems to master the skill.<br>"
			} else {
				message += "You should learn this skill. Solve " + strconv.Itoa(skill.NumberOfProblemsToUnlock-skill.Solved) + " more problems to learn the skill.<br></div>"
			}
		}
		go SendEmail(email, "Hi "+username+"! You've been inactive for a few days!", message)
	}
	return nil
}
func SendEmailsEvery(interval time.Duration) {
	// SendEmailsToInactive()
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err := SendEmailsToInactive()
				if err != nil {
					log.Println(err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func SendEmail(to, subject, body string) (err error) {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		config.Email,
		config.Password,
		config.SmtpAuth,
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	err = smtp.SendMail(
		config.SmtpAddress,
		auth,
		"CodeRangers",
		[]string{to},
		[]byte("Subject: "+subject+"\r\n"+mime+body+"\r\n"),
	)
	return
}

func EmailLogoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u := r.FormValue("u")
		t := r.FormValue("t")
		i, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			log.Println(err)
		}
		dateSent := time.Unix(i, 0)
		db, err := dao.Open()
		if err != nil {
			log.Println(err)
		}
		_, err = db.Exec("INSERT INTO email_tracking VALUES(?, ?, ?)", u, dateSent, time.Now())
		if err != nil {
			// log.Println(err)
		}
		http.ServeFile(w, r, "images/logoBlack.png")
	}
}
