package judge

import (
	"bytes"
	"coderangers/achievements"
	"coderangers/helper"
	"coderangers/notifications"
	"coderangers/problems"
	"coderangers/skills"
	"coderangers/users"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var DIR string

type Judge interface {
	judge(s Submission)
}

type UvaJudge struct {
}

type CodeRangerJudge struct {
}

var uvaJudge = new(UvaJudge)
var codeRangerJudge = new(CodeRangerJudge)

type Submission struct {
	Username        string
	UserID          int
	ID              int
	ProblemIndex    int
	Directory       string `json:"-"`
	Verdict         string
	UvaSubmissionID int `json:"-"`
	Runtime         float64
	ProblemTitle    string
	Language        string
	Timestamp       time.Time
	Tries           int
}

type VerdictData struct {
	Accepted          int
	WrongAnswer       int
	CompileError      int
	RuntimeError      int
	TimeLimitExceeded int
}

type UvaSubmissions struct {
	Name  string  `json:"name"`
	Uname string  `json:"uname"`
	Subs  [][]int `json:"subs"`
}

const (
	UvaUsername = "CodeRanger2"
	UvaUserID   = "821610"
)

type UserSubmissions struct {
	Submissions UvaSubmissions `json:"821610"`
}

const (
	Java = "Java"
	C    = "C"
)

const (
	HardXP   = 50
	MediumXP = 30
	EasyXP   = 10
	Hard     = "Hard"
	Medium   = "Medium"
	Easy     = "Easy"
)

type Error struct {
	Verdict string
	Details string
}

func (e Error) Error() string {
	return e.Verdict // + ":\n" + e.Details
}

var (
	submissionQueue chan *Submission
	uvaQueue        chan *Submission
)

func InitQueues() {
	submissionQueue = make(chan *Submission)
	go func() {
		for s := range submissionQueue {
			log.Println("getting ", s.ID)
			p, err := GetProblem(s.ProblemIndex)
			if err != nil {
				log.Println(err)
			}
			log.Println("checking ", s.ID)
			if p.UvaID == "" {
				if s.Language == Java {
					codeRangerJudge.judge(s)
				} else {
					codeRangerJudge.judge(s)
				}
			} else {
				uvaJudge.judge(s)
			}
			log.Println("judged ", s.ID)
		}
		log.Println("Submission Queue Closed!!!!")
	}()

	uvaQueue = make(chan *Submission)
	go func() {
		for s := range uvaQueue {
			go uvaJudge.checkVerdict(s)
		}
		log.Println("UVa Queue Closed!!!!")
	}()

}

func (UvaJudge) checkVerdict(s *Submission) {
	prob, err := GetProblem(s.ProblemIndex)
	// fmt.Println("http://uhunt.felix-halim.net/api/subs-nums/" + UvaUserID + "/" + prob.UvaID + "/" + strconv.Itoa(s.UvaSubmissionID - 1))
	resp, err := http.Get("http://uhunt.felix-halim.net/api/subs-nums/" + UvaUserID + "/" + prob.UvaID + "/" + strconv.Itoa(s.UvaSubmissionID-1))
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		uvaQueue <- s
	} else {
		userSubmissions := new(UserSubmissions)
		json.NewDecoder(resp.Body).Decode(userSubmissions)
		submissions := userSubmissions.Submissions
		if (len(submissions.Subs)) == 0 {
			uvaQueue <- s
		}
		for i := 0; i < len(submissions.Subs); i++ {
			if submissions.Subs[i][0] == s.UvaSubmissionID {
				if submissions.Subs[i][2] == 10 {
					log.Println("sub err ", s.ID)
					submissionQueue <- s
				} else if submissions.Subs[i][2] == 20 || submissions.Subs[i][2] == 0 {
					time.Sleep(2 * time.Second)
					uvaQueue <- s
				} else {
					var verdict string
					switch submissions.Subs[i][2] {
					case 30:
						verdict = problems.CompileError
					case 35:
						verdict = problems.RestrictedFunction
					case 40:
						verdict = problems.RuntimeError
					case 45:
						verdict = problems.OutputLimitExceeded
					case 50:
						verdict = problems.TimeLimitExceeded
					case 60:
						verdict = problems.MemoryLimitExceeded
					case 70:
						verdict = problems.WrongAnswer
					case 80:
						verdict = problems.PresentationError
					case 90:
						verdict = problems.Accepted
					}
					s.Verdict = verdict
					s.Runtime = float64(submissions.Subs[i][3]) / 1000.00
					UpdateVerdict(s, verdict)
					err = UpdateRuntime(s.ID, s.Runtime)
					if err != nil {
						log.Println(err)
					}
					sendNotification(*s, prob)
				}
				break
			}
		}
	}
}

func ResendNotification(submissionID int) {
	sub, err := GetSubmission(submissionID)
	if err != nil {
		log.Println(err)
		return
	}
	prob, err := GetProblem(sub.ProblemIndex)
	if err != nil {
		log.Println(err)
		return
	}
	sendNotification(sub, prob)
}

func sendNotification(s Submission, prob problems.Problem) {
	var relatedProblems []problems.Problem
	var newAchievements []achievements.Achievement
	var err error
	if s.Verdict == problems.Accepted {
		newAchievements, err = achievements.CheckNewAchievementsInSkill(s.UserID, s.ID, prob.SkillID)
		if err != nil {
			log.Println(err)
		}
		relatedProblems, err = GetUnsolvedUnlockedProblem(s.UserID)
		if err != nil {
			log.Println(err)
		}
	} else {
		relatedProblems, err = GetRelatedProblems(s.UserID, s.ProblemIndex)
		if err != nil {
			log.Println(err)
		}
		if len(relatedProblems) == 0 {
			relatedProblems, err = GetUnsolvedUnlockedProblem(s.UserID)
			if err != nil {
				log.Println(err)
			}
			for i, element := range relatedProblems {
				if element.Index == s.ProblemIndex {
					relatedProblems = append(relatedProblems[:i], relatedProblems[i+1:]...)
					break
				}
			}
		}
	}
	user, err := users.GetUserData(s.UserID)
	if err != nil {
		log.Println(err)
	}
	skill, err := skills.GetUserDataOnSkill(s.UserID, prob.SkillID)
	if err != nil {
		log.Println(err)
	}
	problemList, err := skills.GetProblemsInSkill(prob.SkillID)
	if err != nil {
		log.Println(err)
	}
	firstTime, err := firstTimeSolved(s.UserID, prob.Index)
	if err != nil {
		log.Println(err)
	}
	skill.NumberOfProblems = len(problemList)
	data := struct {
		Submission      Submission
		Problem         problems.Problem
		User            users.UserData
		Skill           skills.Skill
		RelatedProblems []problems.Problem
		NewAchievements []achievements.Achievement
		FirstTime       bool
	}{
		s,
		prob,
		user,
		skill,
		relatedProblems,
		newAchievements,
		firstTime,
	}
	message, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	} else {
		notifications.SendMessageTo(s.UserID, string(message), notifications.Notifications)
		err = notifications.AddNotification(s.ID, s.UserID)
		if err != nil {
			// log.Println(err)
		}
	}
}

func (UvaJudge) judge(s *Submission) {
	log.Println("start judge ", s.ID)
	p, _ := GetProblem(s.ProblemIndex)
	var language string
	if s.Language == Java {
		language = "java"
	} else {
		language = "c"
	}

	bytes, err := ioutil.ReadFile(filepath.Join(s.Directory, `Main.`+language))
	if err != nil {
		log.Println(err)
		go addToSubmissionQueue(s)
		return
	}
	code := string(bytes)
	msg, err := submitUVA(UvaUsername, UvaUsername, p.UvaID, languageToUvaID(s.Language), code)
	if err != nil {
		log.Println(err)
		go addToSubmissionQueue(s)
		return
	}

	// if you want to get the submission ID:
	var submissionID int
	n, err := fmt.Sscanf(msg, "Submission received with ID %d", &submissionID)
	if n != 1 || err != nil {
		log.Println(msg)
		log.Println(err)
		go addToSubmissionQueue(s)
		return
	}
	err = updateUvaSubmissionID(s.ID, submissionID)
	if err != nil {
		log.Println(err)
		go addToSubmissionQueue(s)
		return
	}
	s.Verdict = problems.Inqueue
	UpdateVerdict(s, problems.Inqueue)
	s.UvaSubmissionID = submissionID
	uvaQueue <- s
	log.Println("end judge ", s.ID)
}

func addToSubmissionQueue(s *Submission) {
	submissionQueue <- s
}

func (CodeRangerJudge) judge(s *Submission) {
	var err *Error

	p, er := GetProblem(s.ProblemIndex)
	if er != nil {
		log.Println(er)
	}

	s.Verdict = problems.Compiling
	// UpdateVerdict(s, Compiling)

	err = s.compile()
	if err != nil {
		s.Verdict = err.Verdict
		UpdateVerdict(s, s.Verdict)
		sendNotification(*s, p)
		return
	}

	s.Verdict = problems.Running
	UpdateVerdict(s, problems.Running)
	t := time.Now()
	output, err := s.run(p)
	d := time.Now().Sub(t)
	if s.Runtime == 0 {
		s.Runtime = helper.Truncate(d.Seconds(), 3)
	}
	UpdateRuntime(s.ID, s.Runtime)
	if err != nil {
		s.Verdict = err.Verdict
		UpdateVerdict(s, s.Verdict)
		sendNotification(*s, p)
		if s.Verdict == problems.RuntimeError {
			UpdateRutimeError(s.ID, err.Details)
		}
		return
	}

	// s.Verdict = Judging
	// UpdateVerdict(s, Judging)

	if strings.Replace(output, "\r\n", "\n", -1) != strings.Replace(p.Output, "\r\n", "\n", -1) {
		// whitespace checks..? floats? etc.
		s.Verdict = problems.WrongAnswer
		UpdateVerdict(s, problems.WrongAnswer)
		sendNotification(*s, p)
		return
	}

	s.Verdict = problems.Accepted
	UpdateVerdict(s, problems.Accepted)
	sendNotification(*s, p)
}

func UpdateVerdict(s *Submission, verdict string) {
	err := UpdateVerdictInDB(s.ID, verdict)
	if err != nil {
		log.Println(err)
	}
	message, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	} else {
		notifications.SendMessageTo(s.UserID, string(message), notifications.Submissions)
	}
}

func (s Submission) compile() *Error {
	var stderr bytes.Buffer
	var cmd *exec.Cmd
	switch s.Language {
	case Java:
		cmd = exec.Command("javac", "Main.java")
		cmd.Dir = s.Directory
		cmd.Stderr = &stderr
	case C:
		cmd = exec.Command("gcc", "Main.c")
		cmd.Dir = s.Directory
		cmd.Stderr = &stderr
	}
	err := cmd.Run()
	if err != nil {
		// fmt.Println(stderr.String())
		return &Error{problems.CompileError, stderr.String()}
	}

	return nil
}

func (s *Submission) run(p problems.Problem) (string, *Error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	switch s.Language {
	case Java:
		cmd := exec.Command("java", "-Djava.security.manager", "Main") // "-Xmx20m"
		cmd.Dir = s.Directory
		cmd.Stdin = strings.NewReader(p.Input)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		cmd.Start()
		timeout := time.After(time.Duration(p.TimeLimit) * time.Second)
		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()
		select {
		case <-timeout:
			cmd.Process.Kill()
			return "", &Error{problems.TimeLimitExceeded, ""}
		case err := <-done:
			if err != nil {
				// fmt.Println(stderr.String())
				return "", &Error{problems.RuntimeError, stderr.String()}
			}
		}
	case C:
		cmd := exec.Command("isolate", "--init")
		cmd.Stdout = &stdout
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		for stdout.String() == "" {

		}
		var dir = filepath.Join(strings.Replace(stdout.String(), "\n", "", -1), "box")
		cmd = exec.Command("mv", filepath.Join(s.Directory, "a.out"), dir)

		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		cmd.Wait()
		ioutil.WriteFile(filepath.Join(dir, "in.txt"), []byte(p.Input), 0600)
		cmd = exec.Command("isolate", "--time="+fmt.Sprintf("%d", p.TimeLimit), "--mem=262144",
			"--meta=meta.txt", "--stdin=in.txt", "--stdout=out.txt", "--stderr=stderr.txt", "--run", "a.out")
		stdout.Reset()
		cmd.Dir = dir
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		cmd.Wait()
		bytes, err := ioutil.ReadFile(filepath.Join(dir, "meta.txt"))
		if err != nil {
			log.Println(err)
		}
		meta := string(bytes)
		bytes, err = ioutil.ReadFile(filepath.Join(dir, "stderr.txt"))
		if err != nil {
			log.Println(err)
		}
		solutionStderr := string(bytes)
		for _, elem := range strings.Split(meta, "\n") {
			pair := strings.Split(elem, ":")
			if pair[0] == "exitsig" && pair[1] != "0" {
				return "", &Error{problems.RuntimeError, solutionStderr}
			}
			if pair[0] == "time-wall" {
				s.Runtime, err = strconv.ParseFloat(pair[1], 64)
				if err != nil {
					log.Println(err)
				}
			}
		}
		if strings.Contains(stdout.String(), "Time limit exceeded") || strings.Contains(stderr.String(), "Time limit exceeded") {
			return "", &Error{problems.TimeLimitExceeded, ""}
		}
		bytes, err = ioutil.ReadFile(filepath.Join(dir, "out.txt"))
		if err != nil {
			log.Println(err)
		}
		out := string(bytes)
		cmd = exec.Command("isolate", "--cleanup")
		err = cmd.Start()
		if err != nil {
			log.Println(err)
		}
		cmd.Wait()
		return out, nil

	}

	return stdout.String(), nil
}

func ResendReceivedAndCheckInqueue() (err error) {
	subs, err := getSubmissionsReceivedAndInqueue()
	if err != nil {
		return err
	}
	for _, sub := range subs {
		resub := sub
		if sub.Verdict == problems.Inqueue {
			uvaQueue <- &resub
		} else if sub.Verdict == problems.Received ||
			sub.Verdict == problems.Compiling || sub.Verdict == problems.Running {
			submissionQueue <- &resub
		}
	}
	return
}

type Language int

const (
	AnsiC   = "1"
	JavaUva = "2"
	Cpp     = "3"
	Pascal  = "4"
	Cpp11   = "5"
	Python3 = "6"
)

func languageToUvaID(lang string) string {
	switch lang {
	case Java:
		return JavaUva
	case C:
		return AnsiC
	}
	return "-1"
}

func submitUVA(username, password, problemID, language, code string) (string, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get("https://uva.onlinejudge.org/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(?s)<form[^>]*id="mod_loginform"[^>]*>(.*?)</form>`)
	loginForm := re.FindString(string(body))
	if loginForm == "" {
		return "", errors.New("uva-go: could not find UVA login form")
	}

	loginVals := url.Values{}
	re = regexp.MustCompile(`name="([^"]*)"[^>]*value="([^"]*)"`)
	for _, match := range re.FindAllStringSubmatch(loginForm, -1) {
		name := match[1]
		value := match[2]
		loginVals.Set(name, value)
	}

	re = regexp.MustCompile(`action="(.*?)"`)
	match := re.FindStringSubmatch(loginForm)
	if len(match) != 2 {
		return "", errors.New("uva-go: could not find UVA login URL")
	}
	loginURL := html.UnescapeString(match[1])

	loginVals.Set("username", username)
	loginVals.Set("passwd", password)
	resp, err = client.PostForm(loginURL, loginVals)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	redirectURL, err := resp.Location()
	if err != nil {
		return "", errors.New("uva-go: login failed")
	}

	mpart := &bytes.Buffer{}
	w := multipart.NewWriter(mpart)
	err = w.WriteField("problemid", "")
	if err != nil {
		return "", err
	}
	err = w.WriteField("category", "")
	if err != nil {
		return "", err
	}
	err = w.WriteField("localid", problemID)
	if err != nil {
		return "", err
	}
	err = w.WriteField("language", language)
	if err != nil {
		return "", err
	}
	err = w.WriteField("code", code)
	if err != nil {
		return "", err
	}
	_, err = w.CreateFormFile("codeupl", "")
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}

	submitURL := "https://uva.onlinejudge.org/index.php?option=com_onlinejudge&Itemid=25&page=save_submission"
	req, err := http.NewRequest("POST", submitURL, mpart)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", w.FormDataContentType())

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	redirectURL, err = resp.Location()
	if err != nil {
		return "", errors.New("uva-go: submission failed")
	}

	msg := redirectURL.Query().Get("mosmsg")

	return msg, nil
}
