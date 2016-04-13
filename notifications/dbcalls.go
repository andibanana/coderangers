package notifications

import (
	"coderangers/dao"
)

func AddNotification(submissionID int) (err error) {
	db, err := dao.Open()
	if err != nil {
		return
	}

	_, err = db.Exec(`INSERT INTO notifications 
                    VALUES (?, ?);`, submissionID, false)

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
