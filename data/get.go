package data

import (
	".././dao"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func GetSpecificUserData(userID int, toGet string) (int, error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var value int
	db.QueryRow("SELECT ? FROM user_data WHERE user_id = ?", toGet, userID).Scan(&value)

	return value, err
}

func GetUserData(userID int) (data UserData, err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow("SELECT username, experience, coins, submitted_count, accepted_count, attempted_count "+
		"FROM user_data, user_account WHERE user_id = id AND user_id = ?", userID).Scan(&data.Username,
		&data.Experience, &data.Coins, &data.Submitted, &data.Accepted, &data.Attempted)

	return
}
