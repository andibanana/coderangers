package achievements

import (
	".././dao"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Achievement struct {
	Title       string
	Description string
	Image       string
}

func GetAchievements(userID int) (achievements []Achievement, err error) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, title, IFNULL(solved >= number_of_problems, 0) AS mastered, IFNULL(solved >= number_of_problems_to_unlock, 0) AS unlocked FROM 
                          (SELECT COUNT(DISTINCT problems.id) as number_of_problems, skills.title, skills.id, number_of_problems_to_unlock 
                          FROM skills, problems 
                          WHERE skills.id = problems.skill_id
                          GROUP BY skills.id) AS skills
                        LEFT JOIN
                          (SELECT COUNT(DISTINCT problem_id) as solved, skill_id 
                          FROM problems, submissions 
                          WHERE problems.id = submissions.problem_id AND user_id = ?
                          GROUP BY skill_id) AS solved
                          ON (skills.id = solved.skill_id);`, userID)

	for rows.Next() {
		var id string
		var title string
		var mastered bool
		var unlocked bool
		var achievement Achievement
		err = rows.Scan(&id, &title, &mastered, &unlocked)
		if err != nil {
			return
		}
		achievement.Image = id + ".png"
		if unlocked {
			achievement.Title = "Learned skill " + title + " (" + id + ")"
			achievement.Description = "Learned skill " + title + " (" + id + ")"
			achievements = append(achievements, achievement)
		}
		if mastered {
			achievement.Title = "Mastered skill " + title + " (" + id + ")"
			achievement.Description = "Mastered skill " + title + " (" + id + ")"
			achievements = append(achievements, achievement)
		}
	}
	return

}
