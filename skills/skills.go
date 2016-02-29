package skills

import (
	".././dao"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

const (
	Received            = "received"
	Compiling           = "compiling"
	Running             = "running"
	Judging             = "judging"
	Inqueue             = "inqueue"
	Accepted            = "accepted"
	PresentationError   = "presentation error"
	WrongAnswer         = "wrong answer"
	CompileError        = "compile error"
	RuntimeError        = "runtime error"
	TimeLimitExceeded   = "time limit exceeded"
	MemoryLimitExceeded = "memory limit exceeded"
	OutputLimitExceeded = "output limit exceeded"
	SubmissionError     = "submission error"
	RestrictedFunction  = "restricted function"
	CantBeJudged        = "can't be judged"
)

type Skill struct {
	ID                       string
	Title                    string
	Description              string
	NumberOfProblemsToUnlock int
	Prerequisites            []string
}

type Problem struct {
	Index        int
	Title        string
	Description  string
	Difficulty   int
	SkillID      string
	SampleInput  string
	SampleOutput string
	UvaID        string
	Input        string
	Output       string
	TimeLimit    int
	MemoryLimit  int
	Solved       bool
}

func AddSamples() {
	s := Skill{
		ID:                       "1",
		Title:                    "Math",
		Description:              "Mathematic Problems",
		NumberOfProblemsToUnlock: 3,
	}
	addSkill(s)
	s = Skill{
		ID:                       "2",
		Title:                    "Ad Hoc",
		Description:              "Implementation Problems",
		NumberOfProblemsToUnlock: 3,
	}
	addSkill(s)
}

func addSkill(skill Skill) error {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	var where string
	if skill.Prerequisites != nil {
		where = "WHERE id = \"" + skill.Prerequisites[0] + "\""
		if len(skill.Prerequisites) >= 1 {
			for i, s := range skill.Prerequisites {
				if i != 0 {
					where += " or id = \"" + s + "\""
				}
			}
		}
		var count int
		db.QueryRow("SELECT COUNT(*) FROM SKILLS " + where).Scan(&count)

		if len(skill.Prerequisites) > 0 {
			if count != len(skill.Prerequisites) {
				return errors.New("Skill Prerequisites not in database.")
			}
		}
	}
	_, err = db.Exec("INSERT INTO skills (id, title, description, number_of_problems_to_unlock) VALUES (?, ?, ?, ?)",
		skill.ID, skill.Title, skill.Description, skill.NumberOfProblemsToUnlock)

	if err != nil {
		return err
	}

	if skill.Prerequisites != nil {
		var values string
		for i := 0; i < len(skill.Prerequisites); i++ {
			if i == len(skill.Prerequisites)-1 {
				values += " (" + skill.ID + ", " + skill.Prerequisites[i] + ");"
			} else {
				values += " (" + skill.ID + ", " + skill.Prerequisites[i] + "),"
			}
		}

		_, err = db.Exec("INSERT INTO prerequisites (skill_id, prerequisite_id) VALUES " + values)

		if err != nil {
			return err
		}
	}
	return nil
}

func editSkill(skill Skill, originalID string) error {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	var where string
	if skill.Prerequisites != nil {
		where = "WHERE id = \"" + skill.Prerequisites[0] + "\""
		if len(skill.Prerequisites) >= 1 {
			for i, s := range skill.Prerequisites {
				if i != 0 {
					where += " or id = \"" + s + "\""
				}
			}
		}
		var count int
		db.QueryRow("SELECT COUNT(*) FROM SKILLS " + where).Scan(&count)

		if len(skill.Prerequisites) > 0 {
			if count != len(skill.Prerequisites) {
				return errors.New("Skill Prerequisites not in database.")
			}
		}
	}
	_, err = db.Exec("UPDATE skills SET id = ?, title = ?, description = ?, number_of_problems_to_unlock = ? WHERE id = ?",
		skill.ID, skill.Title, skill.Description, skill.NumberOfProblemsToUnlock, originalID)

	if err != nil {
		return err
	}

	if skill.Prerequisites != nil {
		var values string
		for i := 0; i < len(skill.Prerequisites); i++ {
			if i == len(skill.Prerequisites)-1 {
				values += " (" + skill.ID + ", " + skill.Prerequisites[i] + ");"
			} else {
				values += " (" + skill.ID + ", " + skill.Prerequisites[i] + "),"
			}
		}
		_, err = db.Exec("DELETE FROM prerequisites WHERE skill_id = ?", skill.ID)
		if err != nil {
			return err
		}
		_, err = db.Exec("INSERT INTO prerequisites (skill_id, prerequisite_id) VALUES " + values)

		if err != nil {
			return err
		}
	}
	return nil
}

func GetAllSkills() (skills []Skill, err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return skills, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, description, number_of_problems_to_unlock FROM skills")
	if err != nil {
		return skills, err
	}

	var skill Skill
	for rows.Next() {
		rows.Scan(&skill.ID, &skill.Title, &skill.Description, &skill.NumberOfProblemsToUnlock)
		skills = append(skills, skill)
	}
	return skills, nil
}

func getProblemsInSkill(skillID string) (problems []Problem, err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT problems.id, problems.title, problems.description, difficulty, skill_id, time_limit, memory_limit, sample_input,
                        sample_output, input, output, uva_id 
                        FROM problems, skills, inputoutput WHERE problems.id = inputoutput.problem_id AND skill_id = skills.id AND skill_id = ?`, skillID)
	if err != nil {
		return
	}
	for rows.Next() {
		var problem Problem
		err = rows.Scan(&problem.Index, &problem.Title, &problem.Description, &problem.Difficulty, &problem.SkillID, &problem.TimeLimit,
			&problem.MemoryLimit, &problem.SampleInput, &problem.SampleOutput, &problem.Input, &problem.Output, &problem.UvaID)
		problems = append(problems, problem)
	}
	return
}

func getProblemsInSkillForUser(skillID string, userID int) (problems []Problem, err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT problems.id, problems.title, problems.description, difficulty, skill_id, time_limit, memory_limit, sample_input,
                        sample_output, input, output, uva_id, verdict is not null  
                        FROM problems, skills, inputoutput 
                        LEFT JOIN (SELECT * FROM submissions WHERE verdict = ? AND user_id = ?) AS submissions ON (problems.id = submissions.problem_id) 
                        WHERE problems.id = inputoutput.problem_id AND skill_id = skills.id AND skill_id = ?       
                        `, Accepted, userID, skillID)
	if err != nil {
		return
	}
	for rows.Next() {
		var problem Problem
		err = rows.Scan(&problem.Index, &problem.Title, &problem.Description, &problem.Difficulty, &problem.SkillID, &problem.TimeLimit,
			&problem.MemoryLimit, &problem.SampleInput, &problem.SampleOutput, &problem.Input, &problem.Output, &problem.UvaID, &problem.Solved)

		problems = append(problems, problem)
	}
	return
}

func getSkill(id string) (skill Skill, err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return skill, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, title, description, number_of_problems_to_unlock FROM skills WHERE id = ?", id).Scan(&skill.ID, &skill.Title, &skill.Description, &skill.NumberOfProblemsToUnlock)
	if err != nil {
		return skill, err
	}

	rows, err := db.Query("SELECT prerequisite_id FROM prerequisites WHERE skill_id = ?", id)
	if err != nil {
		return skill, err
	}
	var prerequisites []string
	var prereq string
	for rows.Next() {
		rows.Scan(&prereq)
		prerequisites = append(prerequisites, prereq)
	}
	if len(prerequisites) > 0 {
		skill.Prerequisites = prerequisites
	} else {
		skill.Prerequisites = nil
	}
	return skill, nil
}

func GetUnlockedSkills(userID int) (unlockedSkills map[string]bool, err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()
	unlockedSkills = make(map[string]bool)
	rows, err := db.Query(`SELECT id, prerequisite_id, achieved_id FROM
                    (SELECT id, prerequisite_id FROM skills LEFT JOIN prerequisites ON id = skill_id) AS prerequisite_table
                    LEFT JOIN
                    (SELECT id AS achieved_id FROM skills, (SELECT skill_id, COUNT(DISTINCT problem_id) as solved FROM submissions, problems 
                    WHERE user_id = ? AND problem_id = problems.id AND verdict = "accepted" GROUP BY skill_id) AS unique_solved 
                    WHERE skills.id = unique_solved.skill_id AND unique_solved.solved >= skills.number_of_problems_to_unlock) AS achieved_table
                    ON prerequisite_id = achieved_id;`, userID)
	if err != nil {
		return
	}

	var skillID string
	var prerequisiteID string
	var achievedID string
	for rows.Next() {
		rows.Scan(&skillID, &prerequisiteID, &achievedID)
		if _, ok := unlockedSkills[skillID]; !ok {
			unlockedSkills[skillID] = true
		}
		if prerequisiteID != achievedID {
			unlockedSkills[skillID] = false
		}
	}
	return
}
