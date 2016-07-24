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
		view := r.FormValue("view")
		if len(view) == 0 {
			var less []User
			var index int
			for i, e := range users {
				if e.ID == ID {
					index = i
				}
			}

			for i, e := range users {
				if i < 3 || (i > index-5 && i < index+5) {
					less = append(less, e)
				}
			}
			users = less
		}

		data := struct {
			Leaderboards []User
			IsAdmin      bool
			IsLoggedIn   bool
			UserID       int
			All          bool
		}{
			users,
			dao.IsAdmin(r),
			cookies.IsLoggedIn(r),
			ID,
			len(view) != 0,
		}
		templating.RenderPageWithBase(w, "leaderboards", data)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}
