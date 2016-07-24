package leaderboards

import (
	"coderangers/dao"
	"coderangers/problems"
)

type User struct {
	Rank       int
	ID         int
	Username   string
	Experience int
}

func GetTopUsers(limit, offset int) (users []User, err error) {
	db, err := dao.Open()
	if err != nil {
		return users, err
	}

	rows, err := db.Query(`SELECT user_account.id, username, IFNULL(SUM(difficulty), 0) AS experience FROM
                          user_account
                        LEFT JOIN
                          (SELECT DISTINCT problem_id, difficulty, user_id, verdict FROM problems LEFT JOIN submissions ON problems.id = submissions.problem_id 
                          WHERE verdict = ?) AS submitted 
                        ON submitted.user_id = user_account.id 
                        GROUP BY user_account.id
                        ORDER BY experience DESC
                        LIMIT ? OFFSET ?;`, problems.Accepted, limit, offset)

	if err != nil {
		return users, err
	}

	rank := 1
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Username, &user.Experience)
		user.Rank = rank
		rank += 1
		users = append(users, user)
	}

	return users, nil
}
