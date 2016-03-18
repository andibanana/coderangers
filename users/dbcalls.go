package users

import (
	".././dao"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func Login(username, password string) (userID int, ok bool) {
	db, err := dao.Open()
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

func Register(username, password, email string, admin bool) (int, error) {
	db, err := dao.Open()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 0)

	result, err := db.Exec("INSERT INTO user_account (username, hashed_password, email, admin, date_joined) VALUES (?, ?, ?, ?, ?)",
		username, hashedPassword, email, admin, time.Now())
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(userID), nil
}

func changePassword(userID int, password string) (err error) {
	db, err := dao.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 0)

	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE user_account SET hashed_password = ? WHERE id = ?;",
		hashedPassword, userID)

	if err != nil {
		return err
	}

	return err
}

func RegisterAndFakeData(username, password string, admin bool, xp, coins int) (int, error) {
	db, err := dao.Open()
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

	_, err = tx.Exec("INSERT INTO user_data (user_id, submitted_count, accepted_count, attempted_count, viewed_problems_count, experience, daily_challenge, coins) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", userID, 0, 0, 0, 0, xp, 0, coins)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return int(userID), nil
}
