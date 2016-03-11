package judge

import (
	".././helper"
	".././notifications"
	".././problems"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
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
	Directory       string
	Verdict         string
	UvaSubmissionID int
	Runtime         float64
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

type UserSubmissions struct {
	Submissions UvaSubmissions `json:"821610"`
}

const (
	UvaNodeDirectory = `C:\Users\Sean\Desktop\uva-node`
	UvaUsername      = "CodeRanger2"
	UvaUserID        = "821610"
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
	problemList     []*problems.Problem
	problemQueue    chan *problems.Problem
	submissionList  []*Submission
	submissionQueue chan *Submission
	uvaQueue        chan *Submission
	cmd             *exec.Cmd
	stdin           io.WriteCloser
	stdout          bytes.Buffer
)

func InitQueues() {
	problemQueue = make(chan *problems.Problem)
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
			p, err := GetProblem(s.ProblemIndex)
			if err != nil {
				fmt.Println("ERR!!!!: ", err)
			}

			if p.UvaID == "" {
				go codeRangerJudge.judge(s)
			} else {
				uvaJudge.judge(s)
			}
		}
	}()

	uvaQueue = make(chan *Submission)
	go func() {
		for s := range uvaQueue {
			go uvaJudge.checkVerdict(s)
		}
	}()
	cmd := exec.Command("npm", "start")
	cmd.Dir = UvaNodeDirectory
	cmd.Stdout = &stdout
	stdin, _ = cmd.StdinPipe()

	cmd.Start()
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
					go addToSubmissionQueue(s)
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
					UpdateVerdict(s.ID, verdict)
					UpdateRuntime(s.ID, s.Runtime)
					notifications.SendMessageTo(s.UserID, s.Verdict)
				}
			}
		}
	}
}

func (UvaJudge) judge(s *Submission) {
	p, _ := GetProblem(s.ProblemIndex)

	io.WriteString(stdin, "use uva "+UvaUsername+"\n")
	str := "send " + p.UvaID + " " + s.Directory + `\Main.java` + "\n"
	io.WriteString(stdin, str)
	for !(strings.Contains(stdout.String(), "Send ok") || strings.Contains(stdout.String(), "send failed")) {
		time.Sleep(2 * time.Second)
	}

	if strings.Contains(stdout.String(), "send failed") {
		submissionQueue <- s
		return
	}

	stdout.Reset() // cleans out the stdout of the cmd to be used for another judging.

	time.Sleep(6 * time.Second)
	notgotten := true
	for notgotten {
		resp, err := http.Get("http://uhunt.felix-halim.net/api/subs-user-last/" + UvaUserID + "/1")

		if err == nil {
			defer resp.Body.Close()
			submissions := new(UvaSubmissions)
			err = json.NewDecoder(resp.Body).Decode(submissions)
			submissionID := submissions.Subs[0][0]
			if usedSubmissionID(submissionID) { // if the submission is used already that means uhunt is not updated yet. try again.
				continue
			}
			updateUvaSubmissionID(s.ID, submissionID)
			UpdateVerdict(s.ID, problems.Inqueue)
			s.UvaSubmissionID = submissionID
			uvaQueue <- s
			notgotten = false
		}
	}

}

func addToSubmissionQueue(s *Submission) {
	submissionQueue <- s
}

func (CodeRangerJudge) judge(s *Submission) {
	var err *Error

	p, _ := GetProblem(s.ProblemIndex)

	s.Verdict = problems.Compiling
	// UpdateVerdict(s.ID, Compiling)

	err = s.compile()
	if err != nil {
		s.Verdict = err.Verdict
		UpdateVerdict(s.ID, err.Verdict)
		return
	}

	s.Verdict = problems.Running
	UpdateVerdict(s.ID, problems.Running)
	t := time.Now()
	output, err := s.run(p)
	d := time.Now().Sub(t)
	UpdateRuntime(s.ID, helper.Truncate(d.Seconds(), 3))
	// fmt.Println(d)
	if err != nil {
		s.Verdict = err.Verdict
		UpdateVerdict(s.ID, err.Verdict)
		return
	}

	// s.Verdict = Judging
	// UpdateVerdict(s.ID, Judging)

	if strings.Replace(output, "\r\n", "\n", -1) != strings.Replace(p.Output, "\r\n", "\n", -1) {
		// whitespace checks..? floats? etc.
		fmt.Println(output)
		s.Verdict = problems.WrongAnswer
		UpdateVerdict(s.ID, problems.WrongAnswer)
		return
	}

	s.Verdict = problems.Accepted
	UpdateVerdict(s.ID, problems.Accepted)
}

func (s Submission) compile() *Error {
	var stderr bytes.Buffer

	cmd := exec.Command("javac", "Main.java")
	cmd.Dir = s.Directory
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		return &Error{problems.CompileError, stderr.String()}
	}

	return nil
}

func (s Submission) run(p problems.Problem) (string, *Error) {
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
		return "", &Error{problems.TimeLimitExceeded, ""}
	case err := <-done:
		if err != nil {
			fmt.Println(stderr.String())
			return "", &Error{problems.RuntimeError, stderr.String()}
		}
	}

	return stdout.String(), nil
}

func AddSamples() {
	p := problems.Problem{
		Index: -1,
		Title: "Hashmat the Brave Warrior",
		Description: "Hashmat is a brave warrior who with his group of " +
			"young soldiers moves from one place to another to Fight against his opponents. " +
			"Before Fighting he just calculates one thing, " +
			"the difference between his soldier number and the opponent's soldier number. " +
			"From this difference he decides whether to Fight or not. " +
			"Hashmat's soldier number is never greater than his opponent. ",
		SkillID:      "1",
		Difficulty:   1,
		UvaID:        "10055",
		Input:        "10 12\n10 14\n100 200\n4294967295 4294967294\n",
		Output:       "2\n4\n100\n1\n",
		SampleInput:  "10 12\n10 14\n100 200\n4294967295 4294967294\n",
		SampleOutput: "2\n4\n100\n1\n",
		TimeLimit:    5,
		MemoryLimit:  200,
		Tags:         []string{"Subtract", "Math"},
	}
	AddProblem(p)
	p = problems.Problem{
		Index: -1,
		Title: "Relational Operator",
		Description: "Some operators checks about the relationship between " +
			"two values and these operators are called relational operators. " +
			"Given two numerical values your job is just to and out the relationship " +
			"between them that is (i) First one is greater than the second (ii) " +
			"First one is less than the second or (iii) First and second one is equal." +
			"For each line of input produce one line of output. " +
			"This line contains any one of the relational operators '>', '<' or '=', " +
			"which indicates the relation that is appropriate for the given two numbers.",
		SkillID:      "2",
		Difficulty:   5,
		UvaID:        "11172",
		Input:        "3\n10 20\n20 10\n10 10\n",
		Output:       "<\n>\n=\n",
		SampleInput:  "3\n10 20\n20 10\n10 10\n",
		SampleOutput: "<\n>\n=\n",
		TimeLimit:    5,
		MemoryLimit:  200,
		Tags:         []string{"Relational", "Math"},
	}
	AddProblem(p)
	p = problems.Problem{
		Index: -1,
		Title: "Big Mod",
		Description: "Calculate R : B^P mod M\n" +
			"for large values of B, P, and M using an efficient algorithm. (That's right, this problem has a time dependency !!!.)\n" +
			"Three integer values (in the order B, P, M) will be read one number per line. " +
			"B and P are integers in the range 0 to 2147483647 inclusive. M is an integer in the range 1 to 46340 inclusive. ",
		SkillID:      "1",
		Difficulty:   9,
		UvaID:        "374",
		Input:        "3\n18132\n17\n\n17\n1765\n3\n\n2374859\n3029382\n36123\n",
		Output:       "13\n2\n13195\n",
		SampleInput:  "3\n18132\n17\n\n17\n1765\n3\n\n2374859\n3029382\n36123\n",
		SampleOutput: "13\n2\n13195\n",
		TimeLimit:    5,
		MemoryLimit:  200,
		Tags:         []string{"Modulo", "Math"},
	}
	AddProblem(p)
	p.Title = "Clock Hands"
	p.Tags = nil
	AddProblem(p)
	p.Title = "Y3K Problem"
	AddProblem(p)
	p.Title = "Cancer or Scorpio"
	AddProblem(p)
	p.Title = "Amazing"
	AddProblem(p)
	p.Title = "All Integer Average"
	AddProblem(p)
	p.Title = "Mobile Casanova"
	AddProblem(p)
	p.Title = "Horror Dash"
	AddProblem(p)
	p.Title = "Hangman Judge"
	AddProblem(p)

}
