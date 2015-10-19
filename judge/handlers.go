package judge

import (
	".././cookies"
	".././dao"
	".././data"
	".././templating"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
)

func ProblemsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := cookies.GetUserID(r)
	var dailyChallenge Problem
	if ok {
		dailyChallenge = getDailyChallenge(userID)
	}
	data := struct {
		ProblemList    []Problem
		IsAdmin        bool
		IsLoggedIn     bool
		DailyChallenge Problem
	}{
		GetProblems(),
		dao.IsAdmin(r),
		cookies.IsLoggedIn(r),
		dailyChallenge,
	}
	templating.RenderPage(w, "viewproblems", data)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		index, err := strconv.Atoi(r.URL.Path[len("/edit/"):])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		problem, err := GetProblem(index)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		templating.RenderPage(w, "editproblem", problem)
	case "POST":
		time_limit, err := strconv.Atoi(r.FormValue("time_limit"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		memory_limit, err := strconv.Atoi(r.FormValue("memory_limit"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		difficulty, err := strconv.Atoi(r.FormValue("difficulty"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		index, err := strconv.Atoi(r.FormValue("problem_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p := &Problem{
			Index:        index,
			Title:        r.FormValue("title"),
			Description:  r.FormValue("description"),
			Category:     r.FormValue("category"),
			Difficulty:   difficulty,
			Hint:         r.FormValue("hint"),
			Input:        r.FormValue("input"),
			Output:       r.FormValue("output"),
			SampleInput:  r.FormValue("sample_input"),
			SampleOutput: r.FormValue("sample_output"),
			TimeLimit:    time_limit,
			MemoryLimit:  memory_limit,
		}
		problemQueue <- p
		editProblem(*p)
		http.Redirect(w, r, "/view/"+r.FormValue("problem_id"), http.StatusFound)
	}
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templating.RenderPage(w, "addproblem", nil)
	case "POST":
		time_limit, err := strconv.Atoi(r.FormValue("time_limit"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		memory_limit, err := strconv.Atoi(r.FormValue("memory_limit"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		difficulty, err := strconv.Atoi(r.FormValue("difficulty"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p := &Problem{
			Index:        -1,
			Title:        r.FormValue("title"),
			Description:  r.FormValue("description"),
			Category:     r.FormValue("category"),
			Difficulty:   difficulty,
			Hint:         r.FormValue("hint"),
			Input:        r.FormValue("input"),
			Output:       r.FormValue("output"),
			SampleInput:  r.FormValue("sample_input"),
			SampleOutput: r.FormValue("sample_output"),
			TimeLimit:    time_limit,
			MemoryLimit:  memory_limit,
		}
		problemQueue <- p
		AddProblem(*p)
		http.Redirect(w, r, "/problems/", http.StatusFound)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		problemID, err := strconv.Atoi(r.FormValue("problem_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		deleteProblem(problemID)
		http.Redirect(w, r, "/problems/", http.StatusFound)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	index, err := strconv.Atoi(r.URL.Path[len("/view/"):])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	problem, err := GetProblem(index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var userID int
	if cookies.IsLoggedIn(r) {
		userID, _ = cookies.GetUserID(r)
		data.AddViewedProblem(userID, index)
	}
	submitted, verdictData := getProblemStatistics(index)
	rate := float64(verdictData.Accepted) / float64(submitted) * 100
	if verdictData.Accepted == 0 {
		rate = 0
	}
	data := struct {
		Problem     Problem
		HintBought  bool
		Submitted   int
		Rate        float64
		VerdictData VerdictData
	}{
		problem,
		boughtHintAlready(userID, index),
		submitted,
		rate,
		verdictData,
	}
	templating.RenderPage(w, "viewproblem", data)
	// perhaps have a JS WARNING..
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.IsLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	index, _ := strconv.Atoi(r.URL.Path[len("/submit/"):])
	_, err := GetProblem(index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	d, _ := ioutil.TempDir(DIR, "")
	ioutil.WriteFile(filepath.Join(d, "Main.java"), []byte(r.FormValue("code")), 0600)
	userID, _ := cookies.GetUserID(r)
	dailyChallenge := getDailyChallenge(userID)
	s := &Submission{
		UserID:         userID,
		ProblemIndex:   index,
		Directory:      d,
		Verdict:        Received,
		DailyChallenge: dailyChallenge.Index == index,
	}
	submissionID, err := addSubmission(*s, userID)
	data.UpdateAttemptedCount(userID)
	s.ID = submissionID
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	submissionQueue <- s
	http.Redirect(w, r, "/submissions/", http.StatusFound)
}

func SubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	templating.RenderPage(w, "submissions", getSubmissions())
}

func BuyHintHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if !cookies.IsLoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		problemID, err := strconv.Atoi(r.FormValue("p"))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		userID, _ := cookies.GetUserID(r)
		if !boughtHintAlready(userID, problemID) {
			if !buyHint(userID, problemID) {
				http.Error(w, "Not enough coins", http.StatusBadRequest)
				return
			}
		}
		http.Redirect(w, r, "/view/"+strconv.Itoa(problemID), http.StatusSeeOther)
		return

	default:
		templating.ErrorPage(w, http.StatusMethodNotAllowed)
	}
}
