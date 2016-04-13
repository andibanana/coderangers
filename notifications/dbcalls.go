package notifications

import (
	"coderangers/dao"
)

func AddNotification(submissionID, userID int) (err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}

	_, err = db.Exec(`INSERT INTO notifications (submission_id, user_id, viewed) 
                    VALUES (?, ?, ?);`, submissionID, userID, false)

	return
}

func SetViewedNotification(submissionID int) (err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}

	_, err = db.Exec(`UPDATE notifications
                    SET viewed = ?
                    WHERE submission_id = ?;`, true, submissionID)

	return
}

func GetUnviewedNotifications(userID int) (submissionID int, err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}

	err = db.QueryRow(`SELECT submission_id 
                    FROM notifications
                    WHERE NOT viewed AND user_id = ?;`, userID).Scan(&submissionID)

	return
}
