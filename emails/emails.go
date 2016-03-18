package emails

import (
	".././dao"
	".././judge"
	".././problems"
	".././skills"
	".././users"
	"log"
	"net/smtp"
	"strconv"
	"time"
)

func SendEmailsToInactive() (err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT user_account.id, username, email, max(timestamp) FROM
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
		submittime, err = time.Parse("2006-01-02 15:04:05.999999999Z07:00", timestamp)
		if err != nil {
			return
		}
		duration := time.Since(submittime)
		if duration.Hours()/24 >= 3 {
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
			message := "<h1>Hi " + username + "</h1>"
			message += "You've been inactive for a few days!<br>We want you back, here are a few things you can do!<br>"
			if len(unsolvedProblems) != 0 {
				var user users.UserData
				problem, err = judge.GetProblem(unsolvedProblems[0])
				message += `<h2>` + problem.Title + `</h2>`
				message += `You can try to solve this problem!<br>`
				user, err = judge.GetUserWhoRecentlySolvedProblem(unsolvedProblems[0], userID)
				if err == nil && len(user.Username) != 0 {
					message += user.Username + ` recently solved this.<br>`
				}
			}
			if suggestSkill {
				message += "<h2>" + skill.Title + "</h2>"
				message += skill.Description + "<br>"
				if skill.Learned {
					message += "You should try to master this skill. Solve " + strconv.Itoa(skill.NumberOfProblems-skill.Solved) + " more problems to master the skill.<br>"
				} else {
					message += "You should learn this skill. Solve " + strconv.Itoa(skill.NumberOfProblemsToUnlock-skill.Solved) + " more problems to learn the skill.<br>"
				}
			}
			go SendEmail(email, "Hi "+username+"! You've been inactive for a few days!", message)
		}
	}
	return
}

func SendEmailsEvery(interval time.Duration) {
	SendEmailsToInactive()
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err := SendEmailsToInactive()
				log.Println(err)
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
