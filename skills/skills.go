package skills

import (
	".././dao"
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
	skill.Prerequisites = prerequisites
	return skill, nil
}
