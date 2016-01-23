package users

import (
	".././cookies"
	".././dao"
	//".././data"
	".././templating"
	"fmt"
	"net/http"
	"strconv"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if cookies.IsLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		templating.RenderPage(w, "register", nil)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")
		accessLevel := r.FormValue("access_level")
		if accessLevel != "" && !dao.IsAdmin(r) {
			fmt.Fprintf(w, `
				<body style="background: black; text-align: center;">
					<video src="/images/gandalf.mp4" autoplay loop>You Shall Not Pass!</video>
				</body>
			`)
			return
		}
		admin := accessLevel == "admin"

		userID, err := Register(username, password, email, admin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cookies.Login(r, w, userID, username)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		templating.ErrorPage(w, http.StatusMethodNotAllowed)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if cookies.IsLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		templating.RenderPage(w, "login", nil)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")

		userID, ok := Login(username, password)
		if !ok {
			http.Error(w, "Invalid username or password.", http.StatusBadRequest)
			return
		}

		cookies.Login(r, w, userID, username)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		templating.ErrorPage(w, http.StatusMethodNotAllowed)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		cookies.Logout(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		templating.ErrorPage(w, http.StatusMethodNotAllowed)
	}
}

func ViewProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if !cookies.IsLoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		userID, _ := cookies.GetUserID(r)
		userData, _ := GetUserData(userID)
		templating.RenderPage(w, "viewprofile", userData)
	default:
		templating.ErrorPage(w, http.StatusMethodNotAllowed)
	}
}

func ViewUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		userID, err := strconv.Atoi(r.URL.Path[len("/profile/"):])
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		userData, err := GetUserData(userID)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		templating.RenderPage(w, "viewprofile", userData)
	default:
		templating.ErrorPage(w, http.StatusMethodNotAllowed)
	}
}
