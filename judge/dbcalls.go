package judge

import (
	".././dao"
	".././problems"
	".././skills"
	".././users"
	"errors"
	"fmt"

	"time"
)

func AddProblem(problem problems.Problem) (err error) {

	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return
	}

	result, err := tx.Exec("INSERT INTO problems (title, description, difficulty, skill_id, uva_id, time_limit, memory_limit, sample_input, sample_output) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		problem.Title, problem.Description, problem.Difficulty, problem.SkillID, problem.UvaID, problem.TimeLimit, problem.MemoryLimit, problem.SampleInput, problem.SampleOutput)
	if err != nil {
		tx.Rollback()
		return
	}

	problemID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return
	}
	if problem.UvaID != "" {
		_, err = tx.Exec("INSERT INTO inputoutput (problem_id, input_number, input, output) VALUES (?, ?, ?, ?)",
			problemID, 1, problem.Input, problem.Output)

		if err != nil {
			tx.Rollback()
			return
		}
	}

	tags := "INSERT INTO tags (problem_id, tag) VALUES "
	for i := 0; i < len(problem.Tags); i++ {
		if i == len(problem.Tags)-1 {
			tags += " (" + fmt.Sprint(problemID) + `, "` + problem.Tags[i] + `"); `
		} else {
			tags += " (" + fmt.Sprint(problemID) + `, "` + problem.Tags[i] + `"), `
		}
	}
	if problem.Tags != nil {
		_, err = tx.Exec(tags)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

func editProblem(problem problems.Problem) error {

	db, err := dao.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE problems SET title = ?, description = ?, difficulty = ?, skill_id = ?, uva_id = ?, time_limit = ?, memory_limit = ?, sample_input = ?, sample_output = ? WHERE id = ?",
		problem.Title, problem.Description, problem.Difficulty, problem.SkillID, problem.UvaID, problem.TimeLimit, problem.MemoryLimit, problem.SampleInput, problem.SampleOutput, problem.Index)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE inputoutput SET input = ?, output = ? WHERE problem_id = ? AND input_number = ?",
		problem.Input, problem.Output, problem.Index, 1)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM tags WHERE problem_id = ?", problem.Index)

	if err != nil {
		tx.Rollback()
		return err
	}

	tags := "INSERT INTO tags (problem_id, tag) VALUES "
	for i := 0; i < len(problem.Tags); i++ {
		if i == len(problem.Tags)-1 {
			tags += " (" + fmt.Sprint(problem.Index) + ", " + problem.Tags[i] + "); "
		} else {
			tags += " (" + fmt.Sprint(problem.Index) + ", " + problem.Tags[i] + "), "
		}
	}
	if problem.Tags != nil {
		_, err = tx.Exec(tags)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func deleteProblem(problemID int) {

	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return
	}

	_, err = tx.Exec("DELETE FROM problems WHERE id = ?", problemID)
	if err != nil {
		tx.Rollback()
		return
	}

	_, err = tx.Exec("DELETE FROM inputoutput where problem_id = ? AND input_number = ?",
		problemID, 1)

	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
}

func addSubmission(submission Submission, userID int) (int, error) {
	db, err := dao.Open()
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

	return int(submissionID), nil
}

func getSubmissions() (submissions []Submission, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT submissions.id, problem_id, title, username, verdict, user_account.id, IFNULL(runtime, 0), IFNULL(uva_submission_id, 0) FROM problems, submissions, user_account " +
		"WHERE submissions.problem_id = problems.id AND user_account.id = submissions.user_id " +
		"ORDER BY timestamp DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		var submission Submission
		err = rows.Scan(&submission.ID, &submission.ProblemIndex, &submission.ProblemTitle, &submission.Username, &submission.Verdict, &submission.UserID, &submission.Runtime, &submission.UvaSubmissionID)
		if err != nil {
			return
		}
		submissions = append(submissions, submission)
	}

	return
}

func GetSubmission(id int) (submission Submission) {
	db, err := dao.Open()
	if err != nil {
		return submission
	}
	defer db.Close()
	db.QueryRow("SELECT submissions.id, problem_id, username, verdict, user_account.id, uva_submission_id, runtime FROM submissions, user_account "+
		"WHERE user_account.id = submissions.user_id and submissions.id = ?", id).Scan(&submission.ID, &submission.ProblemIndex,
		&submission.Username, &submission.Verdict, &submission.UserID, &submission.UvaSubmissionID, &submission.Runtime)

	return submission
}

func usedSubmissionID(id int) bool {
	db, err := dao.Open()
	if err != nil {
		return true
	}
	defer db.Close()
	var count int
	db.QueryRow("SELECT COUNT(*) FROM submissions WHERE uva_submission_id = ?", id).Scan(&count)
	if count == 0 {
		return false
	} else {
		return true
	}
}

func acceptedAlready(userID, problemID int) bool {
	db, err := dao.Open()
	if err != nil {
		return false
	}
	defer db.Close()

	var count int
	db.QueryRow("SELECT COUNT(*) FROM submissions, user_account "+
		"WHERE user_account.id = submissions.user_id AND verdict = ?"+
		"AND submissions.problem_id = ? AND user_id = ?", problems.Accepted, problemID, userID).Scan(&count)

	if count == 0 {
		return false
	}

	return true
}

func UpdateVerdict(id int, verdict string) error {
	db, err := dao.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE submissions SET verdict = ? WHERE id = ?", verdict, id)

	if err != nil {
		fmt.Println("UPDATE VERDICT: ", err)
		fmt.Println(verdict)
		return err
	}

	return nil
}

func UpdateRuntime(id int, runtime float64) error {
	db, err := dao.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE submissions SET runtime = ? WHERE id = ?", runtime, id)

	if err != nil {
		return err
	}

	return nil
}

func updateUvaSubmissionID(id, submissionID int) error {
	db, err := dao.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE submissions SET uva_submission_id = ? WHERE id = ?", submissionID, id)

	if err != nil {
		return err
	}

	return nil
}

func GetProblems() (problemList []problems.Problem) {
	db, err := dao.Open()
	if err != nil {
		return nil
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, title, description, difficulty, skill_id, time_limit, memory_limit, sample_input, sample_output, uva_id FROM problems")
	//, inputoutput " +
	//"WHERE problems.id = inputoutput.problem_id ")
	// fmt.Println(err)
	if err != nil {
		return nil
	}

	for rows.Next() {
		var problem problems.Problem
		rows.Scan(&problem.Index, &problem.Title, &problem.Description, &problem.Difficulty, &problem.SkillID, &problem.TimeLimit,
			&problem.MemoryLimit, &problem.SampleInput, &problem.SampleOutput, &problem.UvaID)
		problemList = append(problemList, problem)
	}

	return problemList
}

func GetRelatedProblems(userID, problemID int) (relatedProblems []problems.Problem, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()
	rows, err := db.Query(`
    SELECT id, title, description, difficulty, skill_id, time_limit, memory_limit, sample_input, sample_output, uva_id FROM problems WHERE id IN (

      SELECT DISTINCT problem_id FROM tags WHERE problem_id != ? AND tag IN 
      (SELECT DISTINCT tag FROM problems, tags WHERE problem_id = id AND id = ?)

      EXCEPT

      SELECT DISTINCT problem_id 
      FROM submissions 
      WHERE verdict = ? AND user_id = ?
    );`, problemID, problemID, problems.Accepted, userID)
	if err != nil {
		return
	}

	unlocked, err := skills.GetUnlockedSkills(userID)
	if err != nil {
		return
	}
	for rows.Next() {
		var problem problems.Problem
		err = rows.Scan(&problem.Index, &problem.Title, &problem.Description, &problem.Difficulty, &problem.SkillID, &problem.TimeLimit,
			&problem.MemoryLimit, &problem.SampleInput, &problem.SampleOutput, &problem.UvaID)
		if unlocked[problem.SkillID] {
			relatedProblems = append(relatedProblems, problem)
		}
	}

	return
}

func GetProblem(index int) (problems.Problem, error) {
	db, err := dao.Open()
	var problem problems.Problem
	if err != nil {
		return problem, err
	}
	defer db.Close()
	err = db.QueryRow("SELECT id, title, description, difficulty, skill_id, time_limit, memory_limit, sample_input, sample_output, input, output, uva_id FROM problems, inputoutput "+
		"WHERE problems.id = inputoutput.problem_id and problems.id = ?", index).Scan(&problem.Index, &problem.Title, &problem.Description,
		&problem.Difficulty, &problem.SkillID, &problem.TimeLimit, &problem.MemoryLimit, &problem.SampleInput,
		&problem.SampleOutput, &problem.Input, &problem.Output, &problem.UvaID)

	if err != nil {
		return problem, errors.New("No such problem")
	}

	rows, err := db.Query("SELECT tag FROM tags WHERE problem_id = ?", index)

	if err != nil {
		return problem, err
	}
	var tags []string
	var tag string
	for rows.Next() {
		rows.Scan(&tag)
		tags = append(tags, tag)
	}
	if len(tags) > 0 {
		problem.Tags = tags
	} else {
		problem.Tags = nil
	}
	return problem, nil
}

func GetUnsolvedTriedProblems(userID int) (unsolvedProblems []int, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT DISTINCT problem_id FROM submissions WHERE user_id = ?
                  EXCEPT
                  SELECT DISTINCT problem_id FROM submissions WHERE user_id = ? AND verdict = ?;`, userID, userID, problems.Accepted)
	for rows.Next() {
		var problem int
		err = rows.Scan(&problem)
		if err != nil {
			return
		}
		unsolvedProblems = append(unsolvedProblems, problem)
	}
	return
}

func GetUnsolvedProblems(userID int) (unsolvedProblems []int, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT DISTINCT id FROM problems 
                          EXCEPT
                        SELECT DISTINCT problem_id as id FROM submissions WHERE user_id = ? AND verdict = ?;`,
		userID, userID, problems.Accepted)
	for rows.Next() {
		var problem int
		err = rows.Scan(&problem)
		if err != nil {
			return
		}
		unsolvedProblems = append(unsolvedProblems, problem)
	}
	return
}

func GetUserWhoRecentlySolvedProblem(userID, problemID int) (user users.UserData, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(`SELECT user_account.id, username, email FROM user_account, submissions 
                      WHERE user_account.id = submissions.user_id AND submissions.problem_id = ? AND verdict = ? AND submissions.user_id != ?
                      ORDER BY timestamp DESC;`, problemID, problems.Accepted, userID).Scan(&user.ID, &user.Username, &user.Email)

	return
}
