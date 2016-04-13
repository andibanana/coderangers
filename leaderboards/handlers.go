package leaderboards

import (
	"coderangers/cookies"
	"coderangers/dao"
	"coderangers/templating"
	"net/http"
)

func LeaderboardsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users, err := GetTopUsers(100, 0)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		ID, _ := cookies.GetUserID(r)
		data := struct {
			Leaderboards []User
			IsAdmin      bool
			IsLoggedIn   bool
			UserID       int
		}{
			users,
			dao.IsAdmin(r),
			cookies.IsLoggedIn(r),
			ID,
		}
		templating.RenderPageWithBase(w, "leaderboards", data)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}
