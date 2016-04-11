package achievements

import (
	"coderangers/dao"
	"coderangers/problems"
	"coderangers/skills"
)

type Achievement struct {
	Title       string
	Description string
	Image       string
}

func GetAchievements(userID int) (achievements []Achievement, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}

	rows, err := db.Query(`SELECT id, title, IFNULL(solved >= number_of_problems, 0) AS mastered, IFNULL(solved >= number_of_problems_to_unlock, 0) AS unlocked FROM 
                          (SELECT COUNT(DISTINCT problems.id) as number_of_problems, skills.title, skills.id, number_of_problems_to_unlock 
                          FROM skills, problems 
                          WHERE skills.id = problems.skill_id
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
		var id string
		var title string
		var mastered bool
		var unlocked bool
		var achievement Achievement
		err = rows.Scan(&id, &title, &mastered, &unlocked)
		if err != nil {
			return
		}
		if unlocked {
			achievement.Image = "learned/" + id + ".png"
			achievement.Title = "Learned skill " + title + " (" + id + ")"
			achievement.Description = "Learned skill " + title + " (" + id + ")"
			achievements = append(achievements, achievement)
		}
		if mastered {
			achievement.Image = "mastered/" + id + ".png"
			achievement.Title = "Mastered skill " + title + " (" + id + ")"
			achievement.Description = "Mastered skill " + title + " (" + id + ")"
			achievements = append(achievements, achievement)
		}
	}
	return

}

func CheckNewAchievementsInSkill(userID, submissionID int, skillID string) (achievements []Achievement, err error) {
	solvedWithout, err := skills.GetSolvedInSkillWithoutSubmission(userID, submissionID, skillID)
	if err != nil {
		return
	}
	solved, err := skills.GetSolvedInSkill(userID, skillID)
	if err != nil {
		return
	}
	if solved == solvedWithout {
		return
	}
	skill, err := skills.GetSkill(skillID)
	if err != nil {
		return
	}
	var problems []problems.Problem
	if solved >= skill.NumberOfProblemsToUnlock {
		var achievement Achievement
		achievement.Image = "learned/" + skill.ID + ".png"
		achievement.Title = "Learned skill " + skill.Title + " (" + skill.ID + ")"
		achievement.Description = "Learned skill " + skill.Title + " (" + skill.ID + ")"
		achievements = append(achievements, achievement)
		problems, err = skills.GetProblemsInSkill(skillID)
		if err != nil {
			return
		}
		if solved == len(problems) {
			achievement.Image = "mastered/" + skill.ID + ".png"
			achievement.Title = "Mastered skill " + skill.Title + " (" + skill.ID + ")"
			achievement.Description = "Mastered skill " + skill.Title + " (" + skill.ID + ")"
			achievements = append(achievements, achievement)
		}
	}
	return
}
