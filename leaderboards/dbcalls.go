package leaderboards

import (
	".././dao"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func GetTopUsers(limit, offset int) (users []User, err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return users, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT user_account.id, username, experience "+
		"FROM user_account, user_data WHERE user_account.id = user_data.user_id "+
		"ORDER BY experience DESC "+
		"LIMIT ? OFFSET ?", limit, offset)

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
