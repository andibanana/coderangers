package users

import (
	"coderangers/achievements"
	"coderangers/cookies"
	"coderangers/dao"
	"coderangers/notifications"
	"coderangers/templating"
	"encoding/json"
	"log"
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
			templating.ErrorPage(w, "You shall not pass", http.StatusUnauthorized)
			return
		}
		admin := accessLevel == "admin"

		userID, err := Register(username, password, email, admin)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}

		cookies.Login(r, w, userID, username)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
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
			templating.RenderPage(w, "login", "Invalid username or password.")
			return
		}

		cookies.Login(r, w, userID, username)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if cookies.IsLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		templating.RenderPage(w, "changepassword", nil)
	case "POST":
		username := r.FormValue("username")
		oldPassword := r.FormValue("old_password")
		newPassword := r.FormValue("new_password")
		cnewPassword := r.FormValue("cnew_password")

		if newPassword != cnewPassword {
			templating.RenderPage(w, "changepassword", "New password not the same as confirm new password.")
			return
		}

		userID, ok := Login(username, oldPassword)
		if !ok {
			templating.RenderPage(w, "changepassword", "Invalid username or password.")
			return
		}

		err := changePassword(userID, newPassword)

		if err != nil {
			templating.RenderPage(w, "changepassword", err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		userID, _ := cookies.GetUserID(r)
		cookies.Logout(r, w)
		data := struct {
			LoggedOut bool
		}{
			true,
		}
		message, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
		} else {
			notifications.SendMessageTo(userID, string(message), notifications.Notifications)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
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
		userData, err := GetUserData(userID)
		if err != nil {
			templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
			return
		}
		badges, unearned, err := achievements.GetAchievements(userID)
		if err != nil {
			templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
			return
		}
		data := struct {
			UserData     UserData
			IsAdmin      bool
			IsLoggedIn   bool
			Achievements []achievements.Achievement
			Unearned     []achievements.Achievement
		}{
			userData,
			dao.IsAdmin(r),
			cookies.IsLoggedIn(r),
			badges,
			unearned,
		}
		templating.RenderPageWithBase(w, "viewprofile", data)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
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
			templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
			return
		}
		badges, unearned, err := achievements.GetAchievements(userID)
		if err != nil {
			templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
			return
		}
		data := struct {
			UserData     UserData
			IsAdmin      bool
			IsLoggedIn   bool
			Achievements []achievements.Achievement
			Unearned     []achievements.Achievement
		}{
			userData,
			dao.IsAdmin(r),
			cookies.IsLoggedIn(r),
			badges,
			unearned,
		}
		templating.RenderPageWithBase(w, "viewprofile", data)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}
