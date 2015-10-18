package data

import (
	".././dao"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	ViewedProblems = "viewed_problems_count"
	Accepted       = "accepted_count"
	Submitted      = "submitted_count"
)

func updateViewedProblemCount(userID, count int) error {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE user_data SET viewed_problems_count = ? WHERE user_id = ?",
		count, userID)

	if err != nil {
		return err
	}

	return nil
}

func IncrementCount(userID int, toUpdate string) error {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var count int
	tx.QueryRow("SELECT ? FROM user_data WHERE user_id = ?", toUpdate, userID).Scan(&count)

	count += 1
	_, err = tx.Exec("UPDATE user_data SET "+toUpdate+" = ? where user_id = ?", count, userID)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
