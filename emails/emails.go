package emails

import (
	".././dao"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func SendEmailsToInactive() (err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT user_account.id, email, max(timestamp) FROM
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
		var userID int
		var submittime time.Time
		err = rows.Scan(&userID, &email, &timestamp)
		if err != nil {
			return
		}
		submittime, err = time.Parse("2006-01-02 15:04:05.999999999Z07:00", timestamp)
		if err != nil {
			return
		}
		duration := time.Since(submittime)
		if duration.Hours()/24 >= 3 {
			//get unsolved problems
			//get skills
			//send email
		}
	}
	return
}
