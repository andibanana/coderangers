package connections

import (
	"coderangers/cookies"
	"coderangers/judge"
	"coderangers/notifications"
	"log"
	"net/http"
)

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.IsLoggedIn(r) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, _ := cookies.GetUserID(r)

	subID, err := notifications.GetUnviewedNotifications(userID)
	if err != nil {
		log.Println(err)
		return
	}

	judge.ResendNotification(subID)
}
