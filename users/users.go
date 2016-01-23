package users

import (
	".././dao"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	ViewedProblems = "viewed_problems_count"
	Accepted       = "accepted_count"
	Submitted      = "submitted_count"
	Experience     = "experience"
)

type UserData struct {
	Username   string
	Experience int
	Coins      int
	Submitted  int
	Accepted   int
	Attempted  int
}

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

	err = db.QueryRow("SELECT username, experience, submitted_count, accepted_count, attempted_count "+
		"FROM user_data, user_account WHERE user_id = id AND user_id = ?", userID).Scan(&data.Username,
		&data.Experience, &data.Submitted, &data.Accepted, &data.Attempted)

	return
}

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

func UpdateAttemptedCount(userID int) error {
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
	err = tx.QueryRow("SELECT COUNT(DISTINCT problem_id) FROM submissions WHERE user_id = ?", userID).Scan(&count)

	_, err = tx.Exec("UPDATE user_data SET attempted_count = ? where user_id = ?", count, userID)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

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
	err = tx.QueryRow("SELECT "+toUpdate+" FROM user_data WHERE user_id = ?", userID).Scan(&count)

	count += 1
	_, err = tx.Exec("UPDATE user_data SET "+toUpdate+" = ? where user_id = ?", count, userID)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
