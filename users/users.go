package users

import (
	".././dao"
	".././problems"
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
}

func GetUserData(userID int) (data UserData, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, username, email FROM user_account WHERE id = ?", userID).Scan(&data.ID, &data.Username, &data.Email)

	err = db.QueryRow(`SELECT SUM(difficulty) FROM
                    (SELECT DISTINCT problem_id, difficulty FROM submissions, problems 
                    WHERE problems.id = submissions.problem_id AND user_id = ? AND verdict = ?);`, userID, problems.Accepted).Scan(&data.Experience)

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
	defer db.Close()

	_, err = db.Exec("INSERT INTO viewed_problems (user_id, problem_id) VALUES (?, ?)",
		userID, problemID)

	if err != nil {
		return err
	}
	return nil
}
