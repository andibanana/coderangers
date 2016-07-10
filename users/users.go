package users

import (
	"coderangers/dao"
	"coderangers/problems"
)

const (
	ViewedProblems = "viewed_problems_count"
	Accepted       = "accepted_count"
	Submitted      = "submitted_count"
	Experience     = "experience"
)

type UserData struct {
	ID         int
	Username   string
	Email      string
	Experience int
	Coins      int
	Submitted  int
	Accepted   int
	Attempted  int
	LastName   string
	FirstName  string
}

func GetUserData(userID int) (data UserData, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}

	err = db.QueryRow("SELECT id, username, email, last_name, first_name FROM user_account WHERE id = ?", userID).Scan(&data.ID, &data.Username, &data.Email, &data.LastName, &data.FirstName)

	err = db.QueryRow(`SELECT SUM(difficulty) FROM
                    (SELECT DISTINCT problem_id, difficulty FROM submissions, problems 
                    WHERE problems.id = submissions.problem_id AND user_id = ? AND verdict = ?) AS solved;`, userID, problems.Accepted).Scan(&data.Experience)

	err = db.QueryRow(`SELECT COUNT(*) FROM submissions
                    WHERE user_id = ?;`, userID).Scan(&data.Submitted)

	err = db.QueryRow(`SELECT COUNT(DISTINCT problem_id) FROM submissions
                    WHERE user_id = ?;`, userID).Scan(&data.Attempted)

	err = db.QueryRow(`SELECT COUNT(DISTINCT problem_id) FROM submissions
                    WHERE verdict = ? AND user_id = ?;`, problems.Accepted, userID).Scan(&data.Accepted)

	return
}

func AddViewedProblem(userID, problemID int) error {
	db, err := dao.Open()
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO viewed_problems (user_id, problem_id) VALUES (?, ?)",
		userID, problemID)

	if err != nil {
		return err
	}
	return nil
}

func addName(userID int, lastName, firstName string) error {
	db, err := dao.Open()
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE user_account SET last_name = ?, first_name = ? WHERE id = ?;", lastName, firstName, userID)

	if err != nil {
		return err
	}

	return nil
}

func GetNames() (users map[string]UserData, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}

	users = make(map[string]UserData)

	rows, err := db.Query(`SELECT IFNULL(first_name, ""), IFNULL(last_name, "") , username
                        FROM user_account;`)

	if err != nil {
		return
	}

	for rows.Next() {
		var user UserData
		err = rows.Scan(&user.FirstName, &user.LastName, &user.Username)
		if err != nil {
			return
		}
		users[user.Username] = user
	}

	return
}
