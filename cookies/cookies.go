package cookies

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var cookies = sessions.NewCookieStore([]byte("813629774771309960518707211349999998"))

func GetUserID(req *http.Request) (userID int, ok bool) {
	session, _ := cookies.Get(req, "session")
	val := session.Values["user_id"]
	userID, ok = val.(int)
	return
}

func IsLoggedIn(req *http.Request) bool {
	_, ok := GetUserID(req)
	return ok
}

func Login(r *http.Request, w http.ResponseWriter, userID int, username string) {
	session, _ := cookies.Get(r, "session")
	session.Values["user_id"] = userID
	session.Values["username"] = username
	session.Save(r, w)
}

func Logout(r *http.Request, w http.ResponseWriter) {
	session, _ := cookies.Get(r, "session")
	delete(session.Values, "user_id")
	delete(session.Values, "username")
	session.Save(r, w)
}
