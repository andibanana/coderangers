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
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var DIR string
var OS string

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

//Test darkmega12 705026
//Running CodeRanger2 821610

type UserSubmissions struct {
	Submissions UvaSubmissions `json:"821610"`
}

var UvaNodeDirectory string

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
	cmd             *exec.Cmd
	stdin           io.WriteCloser
	stdout          bytes.Buffer
)

func InitQueues() {
	OS = runtime.GOOS
	if OS == "windows" {
		UvaNodeDirectory = `C:\uva-node`
	} else {
		UvaNodeDirectory = `/root/uva-node`
	}
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
					go codeRangerJudge.judge(s)
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

	startUvaNode()
}

func restartUvaNode() {
	kill := exec.Command("killall", "-9", "node")
	log.Println(kill.Run())
	cmd.Wait()
	log.Println("uva-node restarted!")
	startUvaNode()
}

func startUvaNode() {
	cmd = exec.Command("npm", "start")
	cmd.Dir = UvaNodeDirectory
	cmd.Stdout = &stdout
	stdin, _ = cmd.StdinPipe()
	err := cmd.Start()
	if err != nil {
		log.Fatal("npm not found!")
	}
	io.WriteString(stdin, "add uva "+UvaUsername+" "+UvaUsername+"\n")
	if strings.Contains(stdout.String(), "is not recognized as an internal or external command,") ||
		strings.Contains(stdout.String(), "command not found") {
		log.Fatal("UVA NODE NOT FOUND!")
	}
	for {
		if strings.Contains(stdout.String(), "ERR!") {
			log.Fatal("UVA NODE NOT FOUND OR NPM INSTALL NOT YET RUN.")
		}
		if strings.Contains(stdout.String(), "Account added successfully") || strings.Contains(stdout.String(), "An existing account was replaced") {
			break
		}
	}
}

func (UvaJudge) checkVerdict(s *Submission) {
	// fmt.Println("checking")
	prob, err := GetProblem(s.ProblemIndex)
	// fmt.Println("http://uhunt.felix-halim.net/api/subs-nums/" + UvaUserID + "/" + prob.UvaID + "/" + strconv.Itoa(s.UvaSubmissionID - 1))
	resp, err := http.Get("http://uhunt.felix-halim.net/api/subs-nums/" + UvaUserID + "/" + prob.UvaID + "/" + strconv.Itoa(s.UvaSubmissionID-1))
	defer resp.Body.Close()
	if err != nil {
		uvaQueue <- s
	} else {
		userSubmissions := new(UserSubmissions)
		json.NewDecoder(resp.Body).Decode(userSubmissions)
		submissions := userSubmissions.Submissions
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
			log.Println(err)
		}
	}
}

func (UvaJudge) judge(s *Submission) {
	log.Println("start judge ", s.ID)
	stdout.Reset() // cleans out the stdout of the cmd to be used for another judging.
	p, _ := GetProblem(s.ProblemIndex)

	io.WriteString(stdin, "use uva "+UvaUsername+"\n")
	var language string
	if s.Language == Java {
		language = "java"
	} else {
		language = "c"
	}

	str := "send " + p.UvaID + " " + filepath.Join(s.Directory, `Main.`+language) + "\n"

	io.WriteString(stdin, str)
	timeout := time.After(30 * time.Second)
	tick := time.Tick(2 * time.Second)
	for !(strings.Contains(stdout.String(), "Send ok") || strings.Contains(stdout.String(), "send failed") ||
		strings.Contains(stdout.String(), "Login error")) {
		select {
		case <-timeout:
			log.Println("Uva-Node timedout. Here bug.")
			restartUvaNode()
			go addToSubmissionQueue(s)
			return
		case <-tick:
		}
	}

	if strings.Contains(stdout.String(), "send failed") || strings.Contains(stdout.String(), "Login error") {
		log.Println("UVA-NODE: ", stdout.String())
		restartUvaNode()
		go addToSubmissionQueue(s)
		return
	}

	time.Sleep(6 * time.Second)
	timeout = time.After(30 * time.Second)
	tick = time.Tick(2 * time.Second)
	notgotten := true
	for notgotten {
		select {
		case <-timeout:
			log.Println("Unable to get uva-id. Timeout and notgotten uva-id.")
			go addToSubmissionQueue(s)
			restartUvaNode()
			return
		case <-tick:
		}
		resp, err := http.Get("http://uhunt.felix-halim.net/api/subs-user-last/" + UvaUserID + "/1")

		if err == nil {
			defer resp.Body.Close()
			submissions := new(UvaSubmissions)
			err = json.NewDecoder(resp.Body).Decode(submissions)
			submissionID := submissions.Subs[0][0]
			used, err := usedSubmissionID(submissionID)
			if err != nil {
				log.Println(err)
				return
			}
			if used { // if the submission is used already that means uhunt is not updated yet. try again.
				continue
			}
			err = updateUvaSubmissionID(s.ID, submissionID)
			if err != nil {
				log.Println(err)
			}
			UpdateVerdict(s, problems.Inqueue)
			s.UvaSubmissionID = submissionID
			uvaQueue <- s
			notgotten = false
		}
	}
	log.Println("end judge ", s.ID)
}

func addToSubmissionQueue(s *Submission) {
	log.Println("adding ", s.ID)
	submissionQueue <- s
	log.Println("added ", s.ID)
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
	if s.Runtime != 0 {
		UpdateRuntime(s.ID, s.Runtime)
	} else {
		UpdateRuntime(s.ID, helper.Truncate(d.Seconds(), 3))
	}
	if err != nil {
		s.Verdict = err.Verdict
		UpdateVerdict(s, s.Verdict)
		sendNotification(*s, p)
		return
	}

	// s.Verdict = Judging
	// UpdateVerdict(s, Judging)

	if strings.Replace(output, "\r\n", "\n", -1) != strings.Replace(p.Output, "\r\n", "\n", -1) {
		// whitespace checks..? floats? etc.
		// fmt.Println(output)
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
		done := make(chan error)
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
			"--meta=meta.txt", "--stdin=in.txt", "--stdout=out.txt", "--run", "a.out")
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
		for _, elem := range strings.Split(meta, "\n") {
			pair := strings.Split(elem, ":")
			if pair[0] == "exitcode" && pair[1] != "0" {
				return "", &Error{problems.RuntimeError, ""}
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
