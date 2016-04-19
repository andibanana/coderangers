package emails

import (
	"coderangers/dao"
	"coderangers/helper"
	"coderangers/judge"
	"coderangers/problems"
	"coderangers/skills"
	"coderangers/users"
	"fmt"
	"log"
	"math"
	"net/smtp"
	"strconv"
	"time"
)

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
		var submittime time.Time
		err = rows.Scan(&userID, &username, &email, &timestamp)
		if err != nil {
			return
		}
		submittime, err = helper.ParseTime(timestamp)
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
			message := `<img src="http://coderangers.pro/images/logoBlack.png" style="max-width:300px;">`
			message += "<h1>Hi " + username + "</h1>"
			message += "You've been inactive for a few days!<br>We want you back, here are a few things you can do!<br>"
			if len(unsolvedProblems) != 0 {
				var user users.UserData
				problem, err = judge.GetProblem(unsolvedProblems[0])
				message += `<div style="background-color:#DBDBDB;"><a href="http://coderangers.pro/view/` + fmt.Sprintf("%d", problem.Index) + `?mail=true"><h2>` + problem.Title + `</h2></a>`
				message += `You can try to solve this problem!<br>`
				user, err = judge.GetUserWhoRecentlySolvedProblem(userID, unsolvedProblems[0])
				if err == nil && len(user.Username) != 0 {
					message += `<a href="http://coderangers.pro/profile/` + fmt.Sprintf("%d", user.ID) + `?mail=true">` + user.Username + `</a> recently solved this.<br>`
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
					message += `<div style="background-color:#DBDBDB;"><a href="http://coderangers.pro/view/` + fmt.Sprintf("%d", problem.Index) + `?mail=true"><h2>` + problem.Title + `</h2></a>`
					message += `You can try to solve this problem!<br>`
					user, err = judge.GetUserWhoRecentlySolvedProblem(userID, problem.Index)
					if err == nil && len(user.Username) != 0 {
						message += `<a href="http://coderangers.pro/profile/` + fmt.Sprintf("%d", user.ID) + `?mail=true">` + user.Username + `</a> recently solved this.<br>`
					}
				}
			}
			if suggestSkill {
				message += `<div style="background-color:#DBDBDB;"><a href="http://coderangers.pro/skill/` + skill.ID + `?mail=true">` + `<div style="display:inline-block;"><img src="http://coderangers.pro/images/skill-icons/` + skill.ID + `.png" style="vertical-align:middle;max-width:100px;"></div><div style="display:inline-block;vertical-align:middle;"><h2 style="display:inline;">` + skill.Title + "</h2><br></a>"
				message += skill.Description + "<br></div><br>"
				if skill.Learned {
					message += "You should try to master this skill. Solve " + strconv.Itoa(skill.NumberOfProblems-skill.Solved) + " more problems to master the skill.<br>"
				} else {
					message += "You should learn this skill. Solve " + strconv.Itoa(skill.NumberOfProblemsToUnlock-skill.Solved) + " more problems to learn the skill.<br></div>"
				}
			}
			go SendEmail(email, "Hi "+username+"! You've been inactive for a few days!", message)
		}
	}
	return
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
		"CodeRangers1@gmail.com",
		"coderanger123",
		"smtp.gmail.com",
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	err = smtp.SendMail(
		"smtp.gmail.com:25",
		auth,
		"CodeRangers",
		[]string{to},
		[]byte("Subject: "+subject+"\r\n"+mime+body+"\r\n"),
	)
	return
}
