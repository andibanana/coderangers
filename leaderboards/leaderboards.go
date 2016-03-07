package leaderboards

import (
	".././dao"
	".././problems"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID         int
	Username   string
	Experience int
}

func GetTopUsers(limit, offset int) (users []User, err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return users, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT user_id, username, SUM(difficulty) AS experience FROM
                          (SELECT DISTINCT problem_id, difficulty, user_id FROM submissions, problems 
                          WHERE problems.id = submissions.problem_id AND verdict = ?) AS submitted, user_account
                        WHERE submitted.user_id = user_account.id
                        GROUP BY user_id
                        ORDER BY experience DESC
                        LIMIT ? OFFSET ?;`, problems.Accepted, limit, offset)

	if err != nil {
		return users, err
	}

	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Username, &user.Experience)
		users = append(users, user)
	}

	return users, nil
}
