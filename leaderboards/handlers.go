package leaderboards

import (
	".././templating"
	"net/http"
)

func LeaderboardsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users, _ := GetTopUsers(10, 0)
		templating.RenderPage(w, "leaderboards", users)
	default:
		templating.ErrorPage(w, http.StatusMethodNotAllowed)
	}
}
