package skills

import (
	".././dao"
	".././problems"
	"errors"
)

type Skill struct {
	ID                       string
	Title                    string
	Description              string
	NumberOfProblemsToUnlock int
	NumberOfProblems         int
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

	s = Skill{
		ID:                       "A",
		Title:                    "Introduction to Competitive Programming",
		Description:              "Trivial problems focused on familiarizing yourself with the software",
		NumberOfProblemsToUnlock: 2,
		//Prerequisites:			  []string{"a", "b", "c"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "B",
		Title:                    "Ad Hoc 101",
		Description:              "Problems that can be solved with basic programming skills... I hope...",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"A"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "C1",
		Title:                    "Simple Math",
		Description:              "Problems involving basic math problems such as multiplication and fractions",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"B"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "C2",
		Title:                    "Garbage in, Garbage out",
		Description:              "Memory is cheap but not infinite, plus we need to cut down on defense spending where we can if we wanna keep the free coffee at the mess hall",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"B"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "D1",
		Title:                    "More Math",
		Description:              "A lot of people don't like math. I intend to change that",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"C1"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "D2",
		Title:                    "Text Twist",
		Description:              "'RACE CAR' read backwards is actually 'RACE CAR'... who knew?",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"C2"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "D3",
		Title:                    "Try Try Again",
		Description:              "If you keep hitting the compile button, it's bound to work eventually right?",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"C2"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "E",
		Title:                    "Back to Basics",
		Description:              "I hope you still know how to pitch a tent",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"C1", "D2", "D3"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "F1",
		Title:                    "Even More Math",
		Description:              "As if there wasn't enough numbers already, they added letters",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"D1", "E"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "F2",
		Title:                    "Know your Data Structures I",
		Description:              "There is one rule in this organization... actually a lot more but this one is important; Keep your data organized or die",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"E"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "F3",
		Title:                    "Greed is Good",
		Description:              "Follow the money, and hopefully it leads to more money",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"D3"},
	}
	addSkill(s)
}

func addSkill(skill Skill) error {
	db, err := dao.Open()
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
				values += " (\"" + skill.ID + "\", \"" + skill.Prerequisites[i] + "\");"
			} else {
				values += " (\"" + skill.ID + "\", \"" + skill.Prerequisites[i] + "\"),"
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
	db, err := dao.Open()
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
	db, err := dao.Open()
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

func GetProblemsInSkill(skillID string) (problemsInSkill []problems.Problem, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT problems.id, problems.title, problems.description, difficulty, skill_id, time_limit, memory_limit, sample_input,
                        sample_output, IFNULL(input, ""), IFNULL(output, ""), uva_id  
                        FROM problems, skills
                          LEFT JOIN
                        inputoutput 
                        ON (problems.id = inputoutput.problem_id)
                        WHERE skill_id = skills.id AND skill_id = ?;`, skillID)
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
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT problems.id, problems.title, problems.description, difficulty, skill_id, time_limit, memory_limit, sample_input,
                        sample_output, IFNULL(input,"") , IFNULL(output,"") , uva_id, verdict is not null  
                        FROM problems, skills
                        LEFT JOIN inputoutput ON (problems.id = inputoutput.problem_id)
                        LEFT JOIN (SELECT DISTINCT problem_id, verdict FROM submissions WHERE verdict = ? AND user_id = ?) AS submissions ON (problems.id = submissions.problem_id) 
                        WHERE skill_id = skills.id AND skill_id = ?;`, problems.Accepted, userID, skillID)
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

func GetSkill(id string) (skill Skill, err error) {
	db, err := dao.Open()
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
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()
	unlockedSkills = make(map[string]bool)
	rows, err := db.Query(`SELECT id, prerequisite_id, achieved_id FROM
                    (SELECT id, prerequisite_id FROM skills LEFT JOIN prerequisites ON id = skill_id) AS prerequisite_table
                    LEFT JOIN
                    (SELECT id AS achieved_id FROM skills, (SELECT skill_id, COUNT(DISTINCT problem_id) as solved FROM submissions, problems 
                    WHERE user_id = ? AND problem_id = problems.id AND verdict = ? GROUP BY skill_id) AS unique_solved 
                    WHERE skills.id = unique_solved.skill_id AND unique_solved.solved >= skills.number_of_problems_to_unlock) AS achieved_table
                    ON prerequisite_id = achieved_id;`, userID, problems.Accepted)
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

func GetUserDataOnSkills(userID int) (skills map[string]*Skill, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()
	skills = make(map[string]*Skill)
	rows, err := db.Query(`SELECT id, title, description, number_of_problems_to_unlock, IFNULL(solved, 0) as solved, IFNULL(solved >= number_of_problems, 0) AS mastered, IFNULL(solved >= number_of_problems_to_unlock, 0) AS unlocked FROM 
                          (SELECT COUNT(DISTINCT problems.id) as number_of_problems, skills.title, skills.id, number_of_problems_to_unlock, skills.description 
                          FROM skills
                          LEFT JOIN problems 
                          ON (skills.id = problems.skill_id)
                          GROUP BY skills.id) AS skills
                        LEFT JOIN
                          (SELECT COUNT(DISTINCT problem_id) as solved, skill_id 
                          FROM problems, submissions 
                          WHERE problems.id = submissions.problem_id AND user_id = ? AND verdict = ?
                          GROUP BY skill_id) AS solved
                          ON (skills.id = solved.skill_id);`, userID, problems.Accepted)
	if err != nil {
		return
	}
	for rows.Next() {
		var skill Skill
		err = rows.Scan(&skill.ID, &skill.Title, &skill.Description, &skill.NumberOfProblemsToUnlock, &skill.Solved, &skill.Mastered, &skill.Learned)
		if err != nil {
			return
		}
		skills[skill.ID] = &skill
	}
	return
}

func GetUserDataOnSkill(userID int, skillID string) (skill Skill, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(`SELECT id, title, description, number_of_problems_to_unlock, IFNULL(solved, 0) as solved, IFNULL(solved >= number_of_problems, 0) AS mastered, IFNULL(solved >= number_of_problems_to_unlock, 0) AS unlocked FROM 
                      (SELECT COUNT(DISTINCT problems.id) as number_of_problems, skills.title, skills.id, number_of_problems_to_unlock, skills.description 
                      FROM skills
                      LEFT JOIN problems 
                      ON (skills.id = problems.skill_id)
                      WHERE skills.id = ?
                      GROUP BY skills.id) AS skills
                    LEFT JOIN
                      (SELECT COUNT(DISTINCT problem_id) as solved, skill_id 
                      FROM problems, submissions 
                      WHERE problems.id = submissions.problem_id AND user_id = ? AND skill_id = ? AND verdict = ?
                      GROUP BY skill_id) AS solved
                      ON (skills.id = solved.skill_id);`, skillID, userID, skillID, problems.Accepted).Scan(&skill.ID,
		&skill.Title, &skill.Description, &skill.NumberOfProblemsToUnlock, &skill.Solved, &skill.Mastered, &skill.Learned)
	if err != nil {
		return
	}

	return
}

func GetSolvedInSkillWithoutSubmission(userID, submissionID int, skillID string) (solvedCount int, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(`SELECT COUNT (DISTINCT problems.ID)
                    FROM problems, submissions 
                    WHERE problems.ID = submissions.problem_id AND submissions.user_id = ? AND skill_id = ? AND submissions.ID != ? AND verdict = ?;`, userID, skillID,
		submissionID, problems.Accepted).Scan(&solvedCount)
	if err != nil {
		return
	}

	return
}

func GetSolvedInSkill(userID int, skillID string) (solvedCount int, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(`SELECT COUNT (DISTINCT problems.ID)
                    FROM problems, submissions 
                    WHERE problems.ID = submissions.problem_id AND submissions.user_id = ? AND skill_id = ? AND verdict = ?;`, userID, skillID,
		problems.Accepted).Scan(&solvedCount)
	if err != nil {
		return
	}

	return
}
