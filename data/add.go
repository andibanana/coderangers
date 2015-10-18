package data

import (
	".././dao"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func AddViewedProblem(userID, problemID int) error {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO viewed_problems (user_id, problem_id) VALUES (?, ?)",
		userID, problemID)

	if err != nil {
		return err
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM viewed_problems WHERE user_id = ?", userID).Scan(&count)
	updateViewedProblemCount(userID, count)
	return nil
}
