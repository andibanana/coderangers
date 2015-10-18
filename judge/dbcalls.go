package judge

import (
	".././dao"
	".././data"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func AddProblem(problem Problem) {

	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return
	}

	result, err := tx.Exec("INSERT INTO problems (title, description, difficulty, category, hint, time_limit, memory_limit, sample_input, sample_output) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		problem.Title, problem.Description, problem.Difficulty, problem.Category, problem.Hint, problem.TimeLimit, problem.MemoryLimit, problem.SampleInput, problem.SampleOutput)
	if err != nil {
		tx.Rollback()
		return
	}

	problemID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return
	}

	_, err = tx.Exec("INSERT INTO inputoutput (problem_id, input_number, input, output) VALUES (?, ?, ?, ?)",
		problemID, 1, problem.Input, problem.Output)

	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
}

func addSubmission(submission Submission, userID int) (int, error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return -1, err
	}
	defer db.Close()

	if _, err := GetProblem(submission.ProblemIndex); err != nil {
		return -1, errors.New("No such problem")
	}
	result, err := db.Exec("INSERT INTO submissions (problem_id, user_id, directory, verdict, timestamp) VALUES (?, ?, ?, ?, ?)",
		submission.ProblemIndex, userID, submission.Directory, submission.Verdict, time.Now())

	if err != nil {
		return -1, err
	}

	submissionID, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	data.IncrementCount(userID, data.Submitted)

	return int(submissionID), nil
}

func getSubmissions() []Submission {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return nil
	}
	defer db.Close()

	rows, err := db.Query("SELECT submissions.id, problem_id, username, verdict FROM problems, submissions, user_account " +
		"WHERE problems.id = submissions.problem_id and user_account.id = submissions.user_id " +
		"ORDER BY timestamp DESC")

	var submissions []Submission
	for rows.Next() {
		var submission Submission
		rows.Scan(&submission.ID, &submission.ProblemIndex, &submission.Username, &submission.Verdict)
		submissions = append(submissions, submission)
	}

	return submissions
}

func acceptedAlready(userID, problemID int) bool {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return false
	}
	defer db.Close()

	var count int
	db.QueryRow("SELECT COUNT(*) FROM submissions, user_account "+
		"WHERE user_account.id = submissions.user_id AND verdict = ?"+
		"AND submissions.problem_id = ? AND user_id = ?", Accepted, problemID, userID).Scan(&count)

	if count == 0 {
		return false
	}

	return true
}

func UpdateVerdict(id int, verdict string) error {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE submissions SET verdict = ? WHERE id = ?", verdict, id)

	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func GetProblems() []Problem {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return nil
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, title, description, difficulty, category, time_limit, memory_limit, sample_input, sample_output, input, output FROM problems, inputoutput " +
		"WHERE problems.id = inputoutput.problem_id ")

	var problems []Problem
	for rows.Next() {
		var problem Problem
		rows.Scan(&problem.Index, &problem.Title, &problem.Description, &problem.Difficulty, &problem.Category, &problem.TimeLimit,
			&problem.MemoryLimit, &problem.SampleInput, &problem.SampleOutput, &problem.Input, &problem.Output)
		problems = append(problems, problem)
	}

	return problems
}

func GetProblem(index int) (Problem, error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	var problem Problem
	if err != nil {
		return problem, err
	}
	defer db.Close()
	err = db.QueryRow("SELECT id, title, description, difficulty, category, time_limit, memory_limit, sample_input, sample_output, input, output FROM problems, inputoutput "+
		"WHERE problems.id = inputoutput.problem_id and problems.id = ?", index).Scan(&problem.Index, &problem.Title, &problem.Description,
		&problem.Difficulty, &problem.Category, &problem.TimeLimit, &problem.MemoryLimit, &problem.SampleInput,
		&problem.SampleOutput, &problem.Input, &problem.Output)

	if err != nil {
		return problem, errors.New("No such problem")
	}
	return problem, nil
}

func getDailyChallenge(userID int) (problem Problem) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return problem
	}
	defer db.Close()

	experience, _ := data.GetUserData(userID, data.Experience)
	var difficulty string
	switch {
	case experience <= 100:
		difficulty = Easy
	case experience <= 150:
		difficulty = Medium
	case experience <= 200:
		difficulty = Hard
	}
	var problemID int
	currentTime := time.Now()
	day := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.Local)
	db.QueryRow("SELECT problem_id FROM daily_challenges WHERE day = ? and difficulty = ?", day, difficulty).Scan(&problemID)
	problem, _ = GetProblem(problemID)
	return problem
}

func AddDailyChallenge(time time.Time, difficulty string, problemID int) error {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO daily_challenges (day, difficulty, problem_id) VALUES (?, ?, ?)",
		time, difficulty, problemID)
	return err
}
