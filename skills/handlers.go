package skills

import (
	"coderangers/cookies"
	"coderangers/dao"
	"coderangers/problems"
	"coderangers/templating"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func SkillHandler(w http.ResponseWriter, r *http.Request) {
	skill := r.URL.Path[len("/skill/"):]
	var skills Skill
	loggedIn := cookies.IsLoggedIn(r)
	var err error
	if loggedIn {
		userID, _ := cookies.GetUserID(r)
		skills, err = GetUserDataOnSkill(userID, skill)
	} else {
		skills, err = GetSkill(skill)
	}
	if err != nil {
		templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}
	var problemsInSkill []problems.Problem
	var userID int
	if loggedIn {
		userID, _ := cookies.GetUserID(r)
		problemsInSkill, err = getProblemsInSkillForUser(skill, userID)
	} else {
		problemsInSkill, err = GetProblemsInSkill(skill)
	}
	if err != nil {
		templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}
	unlockedSkills, err := GetUnlockedSkills(userID)
	if err != nil {
		templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
		return
	}
	data := struct {
		ProblemList []problems.Problem
		Skill       Skill
		IsAdmin     bool
		IsLoggedIn  bool
		Locked      bool
	}{
		problemsInSkill,
		skills,
		dao.IsAdmin(r),
		loggedIn,
		!unlockedSkills[skill],
	}
	templating.RenderPageWithBase(w, "skill", data)
}

func SkillTreeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		IsLoggedIn := cookies.IsLoggedIn(r)
		if !IsLoggedIn {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		userID, _ := cookies.GetUserID(r)
		unlockedSkills, err := GetUnlockedSkills(userID)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		skillsData, err := GetUserDataOnSkills(userID)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		for index, _ := range skillsData {
			probs, err := GetProblemsInSkill(skillsData[index].ID)
			if err != nil {
				templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
				return
			}
			skillsData[index].NumberOfProblems = len(probs)
		}
		data := struct {
			UnlockedSkills map[string]bool
			IsLoggedIn     bool
			IsAdmin        bool
			SkillsData     map[string]*Skill
		}{
			unlockedSkills,
			IsLoggedIn,
			dao.IsAdmin(r),
			skillsData,
		}
		templating.RenderPageWithBase(w, "skill-tree", data)
	}
}

func AddSkillHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		IsLoggedIn := cookies.IsLoggedIn(r)
		if !IsLoggedIn {
			templating.ErrorPage(w, "Not logged in.", http.StatusBadRequest)
			break
		}
		IsAdmin := dao.IsAdmin(r)
		if !IsAdmin {
			templating.ErrorPage(w, "Not an admin.", http.StatusBadRequest)
			break
		}
		skills, err := GetAllSkills()
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			break
		}
		templating.RenderPage(w, "addskill", skills)
	case "POST":
		IsLoggedIn := cookies.IsLoggedIn(r)
		if !IsLoggedIn {
			templating.ErrorPage(w, "Not Logged in.", http.StatusBadRequest)
			break
		}
		IsAdmin := dao.IsAdmin(r)
		if !IsAdmin {
			templating.ErrorPage(w, "Not an admin.", http.StatusBadRequest)
			break
		}
		NumberOfProblemsToUnlock, err := strconv.Atoi(r.FormValue("number_of_problems_to_unlock"))
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
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
			templating.ErrorPage(w, "Not logged in.", http.StatusBadRequest)
			return
		}
		IsAdmin := dao.IsAdmin(r)
		if !IsAdmin {
			templating.ErrorPage(w, "Not an admin.", http.StatusBadRequest)
			return
		}
		skills, err := GetAllSkills()
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		id := r.URL.Path[len("/edit-skill/"):]
		skill, err := GetSkill(id)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
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
			templating.ErrorPage(w, "Not Logged In!", http.StatusBadRequest)
			return
		}
		IsAdmin := dao.IsAdmin(r)
		if !IsAdmin {
			templating.ErrorPage(w, "Not an admin.", http.StatusBadRequest)
			return
		}
		NumberOfProblemsToUnlock, err := strconv.Atoi(r.FormValue("number_of_problems_to_unlock"))
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
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
		err = editSkill(*skill, r.URL.Path[len("/edit-skill/"):])
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
