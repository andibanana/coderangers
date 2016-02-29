package skills

import (
	".././cookies"
	".././dao"
	".././templating"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func SkillHandler(w http.ResponseWriter, r *http.Request) {
	skill := r.URL.Path[len("/skill/"):]
	skills, err := getSkill(skill)
	if err != nil {
		templating.ErrorPage(w, 404)
		return
	}
	loggedIn := cookies.IsLoggedIn(r)
	var problemsInSkill []Problem
	if loggedIn {
		userID, _ := cookies.GetUserID(r)
		problemsInSkill, err = getProblemsInSkillForUser(skill, userID)
	} else {
		problemsInSkill, err = getProblemsInSkill(skill)
	}
	if err != nil {
		fmt.Println(err)
		templating.ErrorPage(w, 404)
		return
	}

	data := struct {
		ProblemList []Problem
		Skill       Skill
		IsAdmin     bool
		IsLoggedIn  bool
	}{
		problemsInSkill,
		skills,
		dao.IsAdmin(r),
		loggedIn,
	}
	templating.RenderPageWithBase(w, "skill", data)
}

func SkillTreeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		IsLoggedIn := cookies.IsLoggedIn(r)
		if !IsLoggedIn {
			templating.ErrorPage(w, 404)
			return
		}
		userID, _ := cookies.GetUserID(r)
		unlockedSkills, err := GetUnlockedSkills(userID)
		if err != nil {
			templating.ErrorPage(w, 404)
			return
		}
		data := struct {
			UnlockedSkills map[string]bool
			IsLoggedIn     bool
			IsAdmin        bool
		}{
			unlockedSkills,
			IsLoggedIn,
			dao.IsAdmin(r),
		}
		templating.RenderPageWithBase(w, "skill-tree", data)
	}
}

func AddSkillHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		IsLoggedIn := cookies.IsLoggedIn(r)
		if !IsLoggedIn {
			templating.ErrorPage(w, 404)
			break
		}
		IsAdmin := dao.IsAdmin(r)
		if !IsAdmin {
			templating.ErrorPage(w, 404)
			break
		}
		skills, err := GetAllSkills()
		if err != nil {
			templating.ErrorPage(w, 404)
			break
		}
		templating.RenderPage(w, "addskill", skills)
	case "POST":
		IsLoggedIn := cookies.IsLoggedIn(r)
		if !IsLoggedIn {
			templating.ErrorPage(w, 404)
			break
		}
		IsAdmin := dao.IsAdmin(r)
		if !IsAdmin {
			templating.ErrorPage(w, 404)
			break
		}
		NumberOfProblemsToUnlock, err := strconv.Atoi(r.FormValue("number_of_problems_to_unlock"))
		if err != nil {
			templating.ErrorPage(w, 400)
			break
		}

		prereq := strings.Replace(r.FormValue("prerequisites"), " ", "", -1)
		arrprereq := strings.Split(prereq, ",")
		if arrprereq[0] == "" {
			arrprereq = nil
		}
		skill := &Skill{
			ID:                       r.FormValue("id"),
			Title:                    r.FormValue("title"),
			Description:              r.FormValue("description"),
			NumberOfProblemsToUnlock: NumberOfProblemsToUnlock,
			Prerequisites:            arrprereq,
		}
		addSkill(*skill)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func EditSkillHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		IsLoggedIn := cookies.IsLoggedIn(r)
		if !IsLoggedIn {
			templating.ErrorPage(w, 404)
			return
		}
		IsAdmin := dao.IsAdmin(r)
		if !IsAdmin {
			templating.ErrorPage(w, 404)
			return
		}
		skills, err := GetAllSkills()
		if err != nil {
			templating.ErrorPage(w, 404)
			return
		}
		id := r.URL.Path[len("/edit-skill/"):]
		skill, err := getSkill(id)
		if err != nil {
			templating.ErrorPage(w, 404)
			return
		}
		data := struct {
			Skills []Skill
			Skill  Skill
		}{
			skills,
			skill,
		}
		templating.RenderPage(w, "editskill", data)
	case "POST":
		IsLoggedIn := cookies.IsLoggedIn(r)
		if !IsLoggedIn {
			templating.ErrorPage(w, 404)
			return
		}
		IsAdmin := dao.IsAdmin(r)
		if !IsAdmin {
			templating.ErrorPage(w, 404)
			return
		}
		NumberOfProblemsToUnlock, err := strconv.Atoi(r.FormValue("number_of_problems_to_unlock"))
		if err != nil {
			templating.ErrorPage(w, 400)
			return
		}

		prereq := strings.Replace(r.FormValue("prerequisites"), " ", "", -1)
		arrprereq := strings.Split(prereq, ",")
		if arrprereq[0] == "" {
			arrprereq = nil
		}
		skill := &Skill{
			ID:                       r.FormValue("id"),
			Title:                    r.FormValue("title"),
			Description:              r.FormValue("description"),
			NumberOfProblemsToUnlock: NumberOfProblemsToUnlock,
			Prerequisites:            arrprereq,
		}
		editSkill(*skill, r.URL.Path[len("/edit-skill/"):])

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
