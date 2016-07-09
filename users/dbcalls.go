package users

import (
	"coderangers/dao"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func Login(username, password string) (userID int, ok bool) {
	db, err := dao.Open()
	if err != nil {
		return 0, false
	}

	var hashedPassword string
	err = db.QueryRow("SELECT id, hashed_password FROM user_account WHERE username=?", username).Scan(&userID, &hashedPassword)
	if err != nil {
		return 0, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return userID, err == nil
}

func Register(username, password, email, lastName, firstName string, admin bool) (int, error) {
	db, err := dao.Open()
	if err != nil {
		return 0, err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 0)

	result, err := db.Exec("INSERT INTO user_account (username, hashed_password, email, last_name, first_name, admin, date_joined) VALUES (?, ?, ?, ?, ?, ?, ?)",
		username, hashedPassword, email, lastName, firstName, admin, time.Now())
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
