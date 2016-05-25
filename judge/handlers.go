package judge

import (
	"coderangers/cookies"
	"coderangers/dao"
	"coderangers/notifications"
	"coderangers/problems"
	"coderangers/skills"
	"coderangers/templating"
	"coderangers/users"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const Limit = 25

func stringToArray(input string) []string {
	cleaned := strings.Replace(input, " ", "", -1)
	arrInput := strings.Split(cleaned, ",")
	if arrInput[0] == "" {
		arrInput = nil
	}
	return arrInput
}

func ProblemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		problemList, err := GetProblems()
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
		}
		data := struct {
			ProblemList []problems.Problem
			IsAdmin     bool
			IsLoggedIn  bool
		}{
			problemList,
			dao.IsAdmin(r),
			cookies.IsLoggedIn(r),
		}
		templating.RenderPageWithBase(w, "viewproblems", data)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		index, err := strconv.Atoi(r.URL.Path[len("/edit-problem/"):])
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		problem, err := GetProblem(index)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		templating.RenderPage(w, "editproblem", problem)
	case "POST":
		time_limit, err := strconv.Atoi(r.FormValue("time_limit"))
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		memory_limit, err := strconv.Atoi(r.FormValue("memory_limit"))
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		difficulty, err := strconv.Atoi(r.FormValue("difficulty"))
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		index, err := strconv.Atoi(r.FormValue("problem_id"))
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		p := &problems.Problem{
			Index:        index,
			Title:        r.FormValue("title"),
			Description:  r.FormValue("description"),
			SkillID:      r.FormValue("skill"),
			Difficulty:   difficulty,
			UvaID:        r.FormValue("uva_id"),
			Input:        r.FormValue("input"),
			Output:       r.FormValue("output"),
			SampleInput:  r.FormValue("sample_input"),
			SampleOutput: r.FormValue("sample_output"),
			TimeLimit:    time_limit,
			MemoryLimit:  memory_limit,
			Tags:         stringToArray(r.FormValue("tags")),
		}
		err = editProblem(*p)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/view/"+r.FormValue("problem_id"), http.StatusFound)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		skills, err := skills.GetAllSkills()
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		templating.RenderPage(w, "addproblem", skills)
	case "POST":
		time_limit, err := strconv.Atoi(r.FormValue("time_limit"))
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		memory_limit, err := strconv.Atoi(r.FormValue("memory_limit"))
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		difficulty, err := strconv.Atoi(r.FormValue("difficulty"))
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		tags := stringToArray(r.FormValue("tags"))
		p := &problems.Problem{
			Index:        -1,
			Title:        r.FormValue("title"),
			Description:  r.FormValue("description"),
			SkillID:      r.FormValue("skill"),
			Difficulty:   difficulty,
			UvaID:        r.FormValue("uva_id"),
			Input:        r.FormValue("input"),
			Output:       r.FormValue("output"),
			SampleInput:  r.FormValue("sample_input"),
			SampleOutput: r.FormValue("sample_output"),
			TimeLimit:    time_limit,
			MemoryLimit:  memory_limit,
			Tags:         tags,
		}
		err = AddProblem(*p)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/problems/", http.StatusFound)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		problemID, err := strconv.Atoi(r.FormValue("problem_id"))
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = deleteProblem(problemID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/problems/", http.StatusFound)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		index, err := strconv.Atoi(r.URL.Path[len("/view/"):])
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		var problem problems.Problem
		var userID int
		if cookies.IsLoggedIn(r) {
			userID, _ = cookies.GetUserID(r)
			err = users.AddViewedProblem(userID, index)
			if err != nil {
				// log.Println(err)
			}
			problem, err = GetUserProblem(index, userID)
			if err != nil {
				templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			problem, err = GetProblem(index)
			if err != nil {
				templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		skill, err := skills.GetSkill(problem.SkillID)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		unlockedSkills, err := skills.GetUnlockedSkills(userID)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		code, language, err := getLastCodeInSubmission(userID, index)
		if err != nil {
			// log.Println(err)
		}

		otherUser, err := GetUserWhoRecentlySolvedProblem(userID, problem.Index)
		hasOtherUser := true
		if err != nil {
			hasOtherUser = false
		}

		solveCount, err := getNumberOtherUsersSolved(userID, problem.Index)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := struct {
			Problem      problems.Problem
			Skill        skills.Skill
			Locked       bool
			IsAdmin      bool
			IsLoggedIn   bool
			Code         string
			Language     string
			OtherUser    users.UserData
			HasOtherUser bool
			SolveCount   int
		}{
			problem,
			skill,
			!unlockedSkills[skill.ID],
			dao.IsAdmin(r),
			cookies.IsLoggedIn(r),
			code,
			language,
			otherUser,
			hasOtherUser,
			solveCount,
		}

		templating.RenderPageWithBase(w, "viewproblem", data)
		// perhaps have a JS WARNING..
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if !cookies.IsLoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		index, _ := strconv.Atoi(r.URL.Path[len("/submit/"):])
		problem, err := GetProblem(index)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		userID, _ := cookies.GetUserID(r)
		unlockedSkills, err := skills.GetUnlockedSkills(userID)
		if !unlockedSkills[problem.SkillID] {
			templating.ErrorPage(w, "Skill not unlocked.", http.StatusUnauthorized)
			return
		}
		d, _ := ioutil.TempDir(DIR, "")
		lang := r.FormValue("language")
		if len(r.FormValue("code")) == 0 {
			templating.ErrorPage(w, "Empty code", http.StatusUnauthorized)
			return
		}
		if lang == Java {
			ioutil.WriteFile(filepath.Join(d, "Main.java"), []byte(r.FormValue("code")), 0600)
		} else if lang == C {
			ioutil.WriteFile(filepath.Join(d, "Main.c"), []byte(r.FormValue("code")), 0600)
		}
		s := &Submission{
			UserID:       userID,
			ProblemIndex: index,
			Directory:    d,
			Verdict:      problems.Received,
			Language:     lang,
			Timestamp:    time.Now(),
		}
		submissionID, err := addSubmission(*s, userID)
		if err != nil {
			log.Println(err)
		}
		s.ID = submissionID
		log.Println("added to db ", s.ID)
		s.ProblemTitle = problem.Title
		user, err := users.GetUserData(userID)
		if err != nil {
			log.Println(err)
		}
		s.Username = user.Username
		message, err := json.Marshal(s)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("send message ", s.ID)
			notifications.SendMessageTo(s.UserID, string(message), notifications.Submissions)
			log.Println("sent message ", s.ID)
		}
		go addToSubmissionQueue(s)
		http.Redirect(w, r, "/submissions/", http.StatusFound)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}

func SubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		index, err := strconv.Atoi(r.URL.Path[len("/submissions/"):])
		if err != nil {
			index = 0
		}
		submissions, count, err := getSubmissions(Limit, index*Limit)
		if err != nil {
			log.Print(err)
		}
		data := struct {
			Submissions []Submission
			IsLoggedIn  bool
			IsAdmin     bool
			Max         int
			Index       int
		}{
			submissions,
			cookies.IsLoggedIn(r),
			dao.IsAdmin(r),
			(int(math.Ceil(float64(count) / Limit))),
			index,
		}
		templating.RenderPageWithBase(w, "submissions", data)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}

func MySubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		index, err := strconv.Atoi(r.URL.Path[len("/my-submissions/"):])
		if err != nil {
			index = 0
		}
		if !cookies.IsLoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		userID, _ := cookies.GetUserID(r)
		submissions, count, err := getUserSubmissions(userID, Limit, index*Limit)
		if err != nil {
			log.Print(err)
		}
		data := struct {
			Submissions []Submission
			IsLoggedIn  bool
			IsAdmin     bool
			Max         int
			Index       int
		}{
			submissions,
			cookies.IsLoggedIn(r),
			dao.IsAdmin(r),
			(int(math.Ceil(float64(count) / Limit))),
			index,
		}
		templating.RenderPageWithBase(w, "my-submissions", data)
	default:
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if !cookies.IsLoggedIn(r) {
			//render skill-tree without data
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		userID, _ := cookies.GetUserID(r)
		allSkills, err := skills.GetUserDataOnSkills(userID)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		unlockedSkills, err := skills.GetUnlockedSkills(userID)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}
		var skill *skills.Skill
		var suggestSkill bool
		for _, element := range allSkills {
			if element.Mastered {
				continue
			}
			if element.Learned || unlockedSkills[element.ID] {
				skill = element
				suggestSkill = true
				problems, err := skills.GetProblemsInSkill(skill.ID)
				skill.NumberOfProblems = len(problems)
				if err != nil {
					templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
					return
				}
				break
			}
		}
		userData, err := users.GetUserData(userID)
		if err != nil {
			templating.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}

		suggestProblem := false
		var problem problems.Problem

		unsolvedUnlockedProblems, err := GetUnsolvedUnlockedProblem(userID)
		if len(unsolvedUnlockedProblems) != 0 {
			problem = unsolvedUnlockedProblems[rand.Intn(len(unsolvedUnlockedProblems))]
			suggestProblem = true
		}

		homeMessage := [2]string{
			"Just like the old saying goes, coding a day keeps the doctor away. So keep Practicing!",
			"Statistics says that coding everyday increases one's coding skills, 99.9 percent of the time",
		}
		message := homeMessage[rand.Intn(len(homeMessage))]

		user, err := GetUserWhoRecentlySolvedProblem(userID, problem.Index)
		hasOtherUser := true
		if err != nil {
			hasOtherUser = false
		}

		data := struct {
			IsLoggedIn     bool
			IsAdmin        bool
			Skill          skills.Skill
			SuggestSkill   bool
			UserData       users.UserData
			SuggestProblem bool
			Problem        problems.Problem
			HasOtherUser   bool
			OtherUser      users.UserData
			Message        string
		}{
			cookies.IsLoggedIn(r),
			dao.IsAdmin(r),
			*skill,
			suggestSkill,
			userData,
			suggestProblem,
			problem,
			hasOtherUser,
			user,
			message,
		}
		templating.RenderPageWithBase(w, "home", data)
	default:
		log.Println(r.Method)
		templating.ErrorPage(w, "", http.StatusMethodNotAllowed)
	}
}

func RandomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if !cookies.IsLoggedIn(r) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "")
			return
		}
		userID, _ := cookies.GetUserID(r)
		unsolvedUnlockedProblems, err := GetUnsolvedUnlockedProblem(userID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "")
		}
		if len(unsolvedUnlockedProblems) != 0 {
			problem := unsolvedUnlockedProblems[rand.Intn(len(unsolvedUnlockedProblems))]
			message, err := json.Marshal(problem)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "")
				return
			}
			fmt.Fprint(w, string(message))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "")
		}
	}
}
