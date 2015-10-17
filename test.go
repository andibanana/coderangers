package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
TODO:
- reorganize / structure the files better (ex a package for submissions?)
- deal with errors! (log if fatal etc)
- use files for the problems and submissions. or use a db.. json. preload. edits?
- seperate files for templates.. a little design?
- implement the memory stuff and all? hm.. checks.. errors.. validation..
*/

var DIR string

type Problem struct {
	Index        int
	Title        string
	Description  string
	Difficulty   int
	Category     string
	SampleInput  string
	SampleOutput string
	Hint         string
	Input        string
	Output       string
	TimeLimit    int
	MemoryLimit  int
}

const (
	Received  = "received"
	Compiling = "compiling"
	Running   = "running"
	Judging   = "judging"

	Accepted = "accepted"
	// PresentationError    = "presentation error"
	WrongAnswer       = "wrong answer"
	CompileError      = "compile error"
	RuntimeError      = "runtime error"
	TimeLimitExceeded = "time limit exceeded"
	// MemoryLimitExceeded  = "memory limit exceeded"
	// OutputLimitExceeded  = "output limit exceeded"
	// SubmissionError      = "submission error"
	// RestrictedFunction   = "restricted function"
	// CantBeJudged         = "can't be judged"
)

type Submission struct {
	Username     string
	ID           int
	ProblemIndex int
	Directory    string
	Verdict      string
}

type Error struct {
	Verdict string
	Details string
}

func (e Error) Error() string {
	return e.Verdict // + ":\n" + e.Details
}

var (
	problemList     []*Problem
	problemQueue    chan *Problem
	submissionList  []*Submission
	submissionQueue chan *Submission
)

func initQueues() {
	problemQueue = make(chan *Problem)
	go func() {
		for p := range problemQueue {
			p.Index = len(problemList)
			problemList = append(problemList, p)
		}
	}()

	submissionQueue = make(chan *Submission)
	go func() {
		for s := range submissionQueue {
			submissionList = append(submissionList, s)
			go s.judge()
		}
	}()
}

func page(content string) string {
	return "<html><body>\n" + content + "\n</body></html>"
}

func problemsHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		ProblemList []Problem
		IsAdmin     bool
		IsLoggedIn  bool
	}{
		getProblems(),
		isAdmin(r),
		isLoggedIn(r),
	}
	renderPage(w, "viewproblems", data)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
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
		addProblem(*p)
		http.Redirect(w, r, "/problems/", http.StatusFound)
	} else {
		renderPage(w, "addproblem", nil)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(r.URL.Path[len("/view/"):])
	problem, err := getProblem(index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderPage(w, "viewproblem", problem)
	// perhaps have a JS WARNING..
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	index, _ := strconv.Atoi(r.URL.Path[len("/submit/"):])
	d, _ := ioutil.TempDir(DIR, "")
	ioutil.WriteFile(filepath.Join(d, "Main.java"), []byte(r.FormValue("code")), 0600)
	s := &Submission{
		ProblemIndex: index,
		Directory:    d,
		Verdict:      Received,
	}
	userID, _ := getUserID(r)
	submissionID, err := addSubmission(*s, userID)
	s.ID = submissionID
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	submissionQueue <- s
	http.Redirect(w, r, "/submissions/", http.StatusFound)
}

func submissionsHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "submissions", getSubmissions())
}

func main() {
	createDB()
	initTemplates()
	wd, _ := os.Getwd()
	DIR = filepath.Join(wd, "submissions")
	os.Mkdir(DIR, 0777)
	initQueues()
	http.HandleFunc("/", problemsHandler)
	http.HandleFunc("/problems/", problemsHandler)
	http.HandleFunc("/add/", addHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/submit/", submitHandler)
	http.HandleFunc("/submissions/", submissionsHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
  
  http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles"))))
  
	http.ListenAndServe(":80", nil)
}

func (s *Submission) judge() {
	var err *Error

	p, _ := getProblem(s.ProblemIndex)

	s.Verdict = Compiling
	updateVerdict(s.ID, Compiling)

	err = s.compile()
	if err != nil {
		s.Verdict = err.Verdict
		updateVerdict(s.ID, err.Verdict)
		return
	}

	s.Verdict = Running
	updateVerdict(s.ID, Running)
	t := time.Now()
	output, err := s.run(p)
	d := time.Now().Sub(t)
	fmt.Println(d)
	if err != nil {
		s.Verdict = err.Verdict
		updateVerdict(s.ID, err.Verdict)
		return
	}

	s.Verdict = Judging
	updateVerdict(s.ID, Judging)

	if strings.Replace(output, "\r\n", "\n", -1) != strings.Replace(p.Output, "\r\n", "\n", -1) {
		// whitespace checks..? floats? etc.
		fmt.Println(output)
		s.Verdict = WrongAnswer
		updateVerdict(s.ID, WrongAnswer)
		return
	}

	s.Verdict = Accepted
	updateVerdict(s.ID, Accepted)
}

func (s Submission) compile() *Error {
	var stderr bytes.Buffer

	cmd := exec.Command("javac", "Main.java")
	cmd.Dir = s.Directory
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		return &Error{CompileError, stderr.String()}
	}

	return nil
}

func (s Submission) run(p Problem) (string, *Error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("java", "-Djava.security.manager", "Main") // "-Xmx20m"
	cmd.Dir = s.Directory
	cmd.Stdin = strings.NewReader(p.Input)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Start()
	timeout := time.After(time.Duration(p.TimeLimit) * time.Second)
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case <-timeout:
		cmd.Process.Kill()
		return "", &Error{TimeLimitExceeded, ""}
	case err := <-done:
		if err != nil {
			fmt.Println(stderr.String())
			return "", &Error{RuntimeError, stderr.String()}
		}
	}

	return stdout.String(), nil
}

const DatabaseURL = "file:database.sqlite?cache=shared&mode=rwc"

var cookies = sessions.NewCookieStore([]byte("813629774771309960518707211349999998"))

var templates *template.Template

func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if isLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		renderPage(w, "register", nil)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")

		accessLevel := r.FormValue("access_level")
		if accessLevel != "" && !isAdmin(r) {
			fmt.Fprintf(w, `
				<body style="background: black; text-align: center;">
					<video src="/images/gandalf.mp4" autoplay loop>You Shall Not Pass!</video>
				</body>
			`)
			return
		}
		admin := accessLevel == "admin"

		userID, err := register(username, password, admin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if accessLevel == "" {
			session, _ := cookies.Get(r, "session")
			session.Values["user_id"] = userID
			session.Values["username"] = username
			session.Save(r, w)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		errorPage(w, http.StatusMethodNotAllowed)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if isLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		renderPage(w, "login", nil)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")

		userID, ok := login(username, password)
		if !ok {
			http.Error(w, "Invalid username or password.", http.StatusBadRequest)
			return
		}

		session, _ := cookies.Get(r, "session")
		session.Values["user_id"] = userID
		session.Values["username"] = username
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		errorPage(w, http.StatusMethodNotAllowed)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		session, _ := cookies.Get(r, "session")
		delete(session.Values, "user_id")
		delete(session.Values, "username")
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		errorPage(w, http.StatusMethodNotAllowed)
	}
}

func login(username, password string) (userID int, ok bool) {
	db, err := sql.Open("sqlite3", DatabaseURL)
	if err != nil {
		return 0, false
	}
	defer db.Close()

	var hashedPassword string
	err = db.QueryRow("SELECT id, hashed_password FROM user_account WHERE username=?", username).Scan(&userID, &hashedPassword)
	if err != nil {
		return 0, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return userID, err == nil
}

func register(username, password string, admin bool) (int, error) {
	db, err := sql.Open("sqlite3", DatabaseURL)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 0)

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec("INSERT INTO user_account (username, hashed_password, admin, date_joined, experience) VALUES (?, ?, ?, ?, 0)",
		username, hashedPassword, admin, time.Now())
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return int(userID), nil
}

func addProblem(problem Problem) {

	db, err := sql.Open("sqlite3", DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return
	}

	result, err := tx.Exec("INSERT INTO problems (title, description, difficulty, category, hint, time_limit, memory_limit, sample_input, sample_output) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		problem.Title, problem.Description, problem.Difficulty, problem.Category, problem.Hint, problem.TimeLimit, problem.MemoryLimit, problem.SampleInput, problem.SampleOutput)
	if err != nil {
		tx.Rollback()
		return
	}

	problemID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return
	}

	_, err = tx.Exec("INSERT INTO inputoutput (problem_id, input_number, input, output) VALUES (?, ?, ?, ?)",
		problemID, 1, problem.Input, problem.Output)

	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
}

func addSubmission(submission Submission, userID int) (int, error) {
	db, err := sql.Open("sqlite3", DatabaseURL)
	if err != nil {
		return -1, err
	}
	defer db.Close()

	if _, err := getProblem(submission.ProblemIndex); err != nil {
		return -1, errors.New("No such problem")
	}
	result, err := db.Exec("INSERT INTO submissions (problem_id, user_id, directory, verdict, timestamp) VALUES (?, ?, ?, ?, ?)",
		submission.ProblemIndex, userID, submission.Directory, submission.Verdict, time.Now())

	if err != nil {
		return -1, err
	}

	submissionID, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	return int(submissionID), nil
}

func getSubmissions() []Submission {
	db, err := sql.Open("sqlite3", DatabaseURL)
	if err != nil {
		return nil
	}
	defer db.Close()

	rows, err := db.Query("SELECT submissions.id, problem_id, username, verdict FROM problems, submissions, user_account " +
		"WHERE problems.id = submissions.problem_id and user_account.id = submissions.user_id " +
		"ORDER BY timestamp ")

	var submissions []Submission
	for rows.Next() {
		var submission Submission
		rows.Scan(&submission.ID, &submission.ProblemIndex, &submission.Username, &submission.Verdict)
		submissions = append(submissions, submission)
	}

	return submissions
}

func updateVerdict(id int, verdict string) error {
	db, err := sql.Open("sqlite3", DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE submissions SET verdict = ? WHERE id = ?", verdict, id)

	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func getProblems() []Problem {
	db, err := sql.Open("sqlite3", DatabaseURL)
	if err != nil {
		return nil
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, title, description, difficulty, category, time_limit, memory_limit, sample_input, sample_output, input, output FROM problems, inputoutput " +
		"WHERE problems.id = inputoutput.problem_id ")

	var problems []Problem
	for rows.Next() {
		var problem Problem
		rows.Scan(&problem.Index, &problem.Title, &problem.Description, &problem.Difficulty, &problem.Category, &problem.TimeLimit,
			&problem.MemoryLimit, &problem.SampleInput, &problem.SampleOutput, &problem.Input, &problem.Output)
		problems = append(problems, problem)
	}

	return problems
}

func getProblem(index int) (Problem, error) {
	db, err := sql.Open("sqlite3", DatabaseURL)
	var problem Problem
	if err != nil {
		return problem, errors.New("DB Problem")
	}
	defer db.Close()
	err = db.QueryRow("SELECT id, title, description, difficulty, category, time_limit, memory_limit, sample_input, sample_output, input, output FROM problems, inputoutput "+
		"WHERE problems.id = inputoutput.problem_id and problems.id = ?", index).Scan(&problem.Index, &problem.Title, &problem.Description,
		&problem.Difficulty, &problem.Category, &problem.TimeLimit, &problem.MemoryLimit, &problem.SampleInput,
		&problem.SampleOutput, &problem.Input, &problem.Output)

	if err != nil {
		return problem, errors.New("No such problem")
	}
	return problem, nil
}

func createDB() error {
	db, err := sql.Open("sqlite3", DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE user_account (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
					
			username VARCHAR(50) UNIQUE NOT NULL,
			hashed_password CHARACTER(60) NOT NULL,
			admin BOOLEAN NOT NULL DEFAULT FALSE,
      date_joined DATE NOT NULL ,
      experience INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE problems (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
			
			title VARCHAR(100) NOT NULL,
      description VARCHAR(200) NOT NULL,
      category VARCHAR(200) NOT NULL,
      difficulty INTEGER,
      hint TEXT,
      time_limit INTEGER,
      memory_limit INTEGER,
      sample_input TEXT,
      sample_output TEXT
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE inputoutput (
      problem_id INTEGER,
      input_number INTEGER,
    
			input TEXT NOT NULL,
      output TEXT NOT NULL,
      
      PRIMARY KEY(problem_id, input_number),
      FOREIGN KEY(problem_id) REFERENCES problems(id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE submissions (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      problem_id INTEGER,
      user_id INTEGER,
      
			directory VARCHAR(100) NOT NULL,
      verdict VARCHAR(100) NOT NULL,
      timestamp DATETIME NOT NULL,
      runtime_error TEXT,
      
      FOREIGN KEY(user_id) REFERENCES user_account(id)
      FOREIGN KEY(problem_id) REFERENCES problems(id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE badges (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      
			title VARCHAR(100) NOT NULL,
      description VARCHAR(100) NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE received_badges (
        user_id INTEGER,
        badge_id INTEGER,
       
        PRIMARY KEY(user_id, badge_id),
        FOREIGN KEY(user_id) REFERENCES user_account(id),
        FOREIGN KEY(badge_id) REFERENCES badges(id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = register("admin", "admin", true)

	return err
}

func isAdmin(req *http.Request) bool {
	userID, ok := getUserID(req)
	if !ok {
		return false
	}

	db, err := sql.Open("sqlite3", DatabaseURL)
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.QueryRow("SELECT id FROM user_account WHERE id=? AND admin=?", userID, true).Scan(&userID)
	return err == nil
}

func isLoggedIn(req *http.Request) bool {
	_, ok := getUserID(req)
	return ok
}

func renderPage(w http.ResponseWriter, template string, data interface{}) {
	err := templates.ExecuteTemplate(w, template+".tmpl.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func errorPage(w http.ResponseWriter, statusCode int) {
	errorMessage := fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
	http.Error(w, errorMessage, statusCode)
}

func getUserID(req *http.Request) (userID int, ok bool) {
	session, _ := cookies.Get(req, "session")
	val := session.Values["user_id"]
	userID, ok = val.(int)
	return
}

func initTemplates() {
	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"showDate":     func(date time.Time) string { return date.Format("Jan 2, 2006") },
		"showDateTime": func(date time.Time) string { return date.Format(time.RFC850) },
		"showISODate":  func(date time.Time) string { return date.Format("2006-01-02") },
		"minus":        func(a, b int) int { return a - b },
		"add":          func(a, b int) int { return a + b },
		"fixNewLines": func(s string) template.HTML {
			s = template.HTMLEscapeString(s)
			s = regexp.MustCompile("\r?\n").ReplaceAllString(s, "<br>")
			return template.HTML(s)
		},
		"boldItalics": func(s string) template.HTML {
			s = template.HTMLEscapeString(s)
			imageTags := regexp.MustCompile(`&lt;img\s+src=&#34;(.*?)&#34;&gt;`)
			s = imageTags.ReplaceAllString(s, `<img src="$1" style="max-width:570px;">`)
			unescapeTags := regexp.MustCompile("&lt;(/?(b|i|pre|u|sub|sup|strike|marquee))&gt;")
			s = unescapeTags.ReplaceAllString(s, "<$1>")
			s = regexp.MustCompile("\r?\n").ReplaceAllString(s, "<br>")
			return template.HTML(s)
		},
	}).ParseGlob("./templates/*.tmpl.html"))
}
