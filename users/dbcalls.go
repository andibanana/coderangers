package users

import (
	".././dao"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func Login(username, password string) (userID int, ok bool) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return 0, false
	}
	defer db.Close()

	var hashedPassword string
	err = db.QueryRow("SELECT id, hashed_password FROM user_account WHERE username=?", username).Scan(&userID, &hashedPassword)
	if err != nil {
		return 0, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return userID, err == nil
}

func Register(username, password string, admin bool) (int, error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 0)

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec("INSERT INTO user_account (username, hashed_password, admin, date_joined) VALUES (?, ?, ?, ?)",
		username, hashedPassword, admin, time.Now())
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	userID, err := result.LastInsertId()

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec("INSERT INTO user_data (user_id, submitted_count, accepted_count, viewed_problems_count, experience, daily_challenge) VALUES (?, ?, ?, ?, ?, ?)", userID, 0, 0, 0, 0, 0)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return int(userID), nil
}
