package skills

import (
	".././dao"
	".././problems"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

type Skill struct {
	ID                       string
	Title                    string
	Description              string
	NumberOfProblemsToUnlock int
	Prerequisites            []string
	Mastered                 bool
	Learned                  bool
	Solved                   int
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

func getProblemsInSkill(skillID string) (problemsInSkill []problems.Problem, err error) {
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
		var problem problems.Problem
		err = rows.Scan(&problem.Index, &problem.Title, &problem.Description, &problem.Difficulty, &problem.SkillID, &problem.TimeLimit,
			&problem.MemoryLimit, &problem.SampleInput, &problem.SampleOutput, &problem.Input, &problem.Output, &problem.UvaID)
		problemsInSkill = append(problemsInSkill, problem)
	}
	return
}

func getProblemsInSkillForUser(skillID string, userID int) (problemsInSkill []problems.Problem, err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT problems.id, problems.title, problems.description, difficulty, skill_id, time_limit, memory_limit, sample_input,
                        sample_output, input, output, uva_id, verdict is not null  
                        FROM problems, skills, inputoutput 
                        LEFT JOIN (SELECT DISTINCT problem_id, verdict FROM submissions WHERE verdict = ? AND user_id = ?) AS submissions ON (problems.id = submissions.problem_id) 
                        WHERE problems.id = inputoutput.problem_id AND skill_id = skills.id AND skill_id = ?       
                        `, problems.Accepted, userID, skillID)
	if err != nil {
		return
	}
	for rows.Next() {
		var problem problems.Problem
		err = rows.Scan(&problem.Index, &problem.Title, &problem.Description, &problem.Difficulty, &problem.SkillID, &problem.TimeLimit,
			&problem.MemoryLimit, &problem.SampleInput, &problem.SampleOutput, &problem.Input, &problem.Output, &problem.UvaID, &problem.Solved)

		problemsInSkill = append(problemsInSkill, problem)
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

func getUserDataOnSkills(userID int) (skills map[string]Skill, err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()
	skills = make(map[string]Skill)
	rows, err := db.Query(`SELECT id, title, description, number_of_problems_to_unlock, IFNULL(solved, 0) as solved, IFNULL(solved >= number_of_problems, 0) AS mastered, IFNULL(solved >= number_of_problems_to_unlock, 0) AS unlocked FROM 
                          (SELECT COUNT(DISTINCT problems.id) as number_of_problems, skills.title, skills.id, number_of_problems_to_unlock, skills.description 
                          FROM skills, problems 
                          WHERE skills.id = problems.skill_id
                          GROUP BY skills.id) AS skills
                        LEFT JOIN
                          (SELECT COUNT(DISTINCT problem_id) as solved, skill_id 
                          FROM problems, submissions 
                          WHERE problems.id = submissions.problem_id AND user_id = ?
                          GROUP BY skill_id) AS solved
                          ON (skills.id = solved.skill_id);`, userID)
	if err != nil {
		return
	}
	for rows.Next() {
		var skill Skill
		err = rows.Scan(&skill.ID, &skill.Title, &skill.Description, &skill.NumberOfProblemsToUnlock, &skill.Solved, &skill.Mastered, &skill.Learned)
		if err != nil {
			return
		}
		skills[skill.ID] = skill
	}
	return
}

func getUserDataOnSkill(userID int, skillID string) (skill Skill, err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(`SELECT id, title, description, number_of_problems_to_unlock, IFNULL(solved, 0) as solved, IFNULL(solved >= number_of_problems, 0) AS mastered, IFNULL(solved >= number_of_problems_to_unlock, 0) AS unlocked FROM 
                          (SELECT COUNT(DISTINCT problems.id) as number_of_problems, skills.title, skills.id, number_of_problems_to_unlock, skills.description
                          FROM skills, problems 
                          WHERE skills.id = problems.skill_id AND skills.id = ?
                          GROUP BY skills.id) AS skills
                        LEFT JOIN
                          (SELECT COUNT(DISTINCT problem_id) as solved, skill_id 
                          FROM problems, submissions 
                          WHERE problems.id = submissions.problem_id AND user_id = ? AND skill_id = ?
                          GROUP BY skill_id) AS solved
                          ON (skills.id = solved.skill_id);`, skillID, userID, skillID).Scan(&skill.ID,
		&skill.Title, &skill.Description, &skill.NumberOfProblemsToUnlock, &skill.Solved, &skill.Mastered, &skill.Learned)
	if err != nil {
		return
	}

	return
}
