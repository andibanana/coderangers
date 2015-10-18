package judge

import (
	".././data"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

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

type Submission struct {
	Username       string
	UserID         int
	ID             int
	ProblemIndex   int
	Directory      string
	Verdict        string
	DailyChallenge bool
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
	problemList     []*Problem
	problemQueue    chan *Problem
	submissionList  []*Submission
	submissionQueue chan *Submission
)

func InitQueues() {
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

func (s *Submission) judge() {
	var err *Error

	p, _ := GetProblem(s.ProblemIndex)

	s.Verdict = Compiling
	UpdateVerdict(s.ID, Compiling)

	err = s.compile()
	if err != nil {
		s.Verdict = err.Verdict
		UpdateVerdict(s.ID, err.Verdict)
		return
	}

	s.Verdict = Running
	UpdateVerdict(s.ID, Running)
	t := time.Now()
	output, err := s.run(p)
	d := time.Now().Sub(t)
	fmt.Println(d)
	if err != nil {
		s.Verdict = err.Verdict
		UpdateVerdict(s.ID, err.Verdict)
		return
	}

	s.Verdict = Judging
	UpdateVerdict(s.ID, Judging)

	if strings.Replace(output, "\r\n", "\n", -1) != strings.Replace(p.Output, "\r\n", "\n", -1) {
		// whitespace checks..? floats? etc.
		fmt.Println(output)
		s.Verdict = WrongAnswer
		UpdateVerdict(s.ID, WrongAnswer)
		return
	}

	s.Verdict = Accepted
	if !acceptedAlready(s.UserID, s.ProblemIndex) {
		data.IncrementCount(s.UserID, data.Accepted)
		multiplier := 1
		if s.DailyChallenge {
			multiplier = 2
			data.IncrementCount(s.UserID, data.DailyChallenge)
		}
		switch {
		case 1 <= p.Difficulty && p.Difficulty <= 3:
			data.AddExperienceAndCoins(s.UserID, EasyXP*multiplier, EasyXP/10*multiplier)
		case 4 <= p.Difficulty && p.Difficulty <= 8:
			data.AddExperienceAndCoins(s.UserID, MediumXP*multiplier, MediumXP/10*multiplier)
		case 9 <= p.Difficulty && p.Difficulty <= 10:
			data.AddExperienceAndCoins(s.UserID, HardXP*multiplier, HardXP/10*multiplier)
		}
	} else if s.DailyChallenge && !acceptedAlreadyAndDailyChallenge(s.UserID, s.ProblemIndex) {
		data.IncrementCount(s.UserID, data.DailyChallenge)
		switch {
		case 1 <= p.Difficulty && p.Difficulty <= 3:
			data.AddExperienceAndCoins(s.UserID, EasyXP, EasyXP/10)
		case 4 <= p.Difficulty && p.Difficulty <= 8:
			data.AddExperienceAndCoins(s.UserID, MediumXP, MediumXP/10)
		case 9 <= p.Difficulty && p.Difficulty <= 10:
			data.AddExperienceAndCoins(s.UserID, HardXP, HardXP/10)
		}
	}
	UpdateVerdict(s.ID, Accepted)
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

func AddSamples() {
	p := Problem{
		Index: -1,
		Title: "Hashmat the Brave Warrior",
		Description: `Hashmat is a brave warrior who with his group of
                    young soldiers moves from one place to another to Fight against his opponents.
                    Before Fighting he just calculates one thing, 
                    the difference between his soldier number and the opponent's soldier number.
                    From this difference he decides whether to Fight or not.
                    Hashmat's soldier number is never greater than his opponent.`,
		Category:     "Math",
		Difficulty:   1,
		Hint:         "Subtract",
		Input:        "10 12\n10 14\n100 200\n4294967295 4294967294\n",
		Output:       "2\n4\n100\n1\n",
		SampleInput:  "10 12\n10 14\n100 200\n4294967295 4294967294\n",
		SampleOutput: "2\n4\n100\n1\n",
		TimeLimit:    5,
		MemoryLimit:  200,
	}
	AddProblem(p)
	p = Problem{
		Index: -1,
		Title: "Hashmat the Brave Warrior",
		Description: `Hashmat is a brave warrior who with his group of
                    young soldiers moves from one place to another to Fight against his opponents.
                    Before Fighting he just calculates one thing, 
                    the difference between his soldier number and the opponent's soldier number.
                    From this difference he decides whether to Fight or not.
                    Hashmat's soldier number is never greater than his opponent.`,
		Category:     "Math",
		Difficulty:   5,
		Hint:         "Subtract",
		Input:        "10 12\n10 14\n100 200\n4294967295 4294967294\n",
		Output:       "2\n4\n100\n1\n",
		SampleInput:  "10 12\n10 14\n100 200\n4294967295 4294967294\n",
		SampleOutput: "2\n4\n100\n1\n",
		TimeLimit:    5,
		MemoryLimit:  200,
	}
	AddProblem(p)
	p = Problem{
		Index: -1,
		Title: "Hashmat the Brave Warrior",
		Description: `Hashmat is a brave warrior who with his group of
                    young soldiers moves from one place to another to Fight against his opponents.
                    Before Fighting he just calculates one thing, 
                    the difference between his soldier number and the opponent's soldier number.
                    From this difference he decides whether to Fight or not.
                    Hashmat's soldier number is never greater than his opponent.`,
		Category:     "Math",
		Difficulty:   9,
		Hint:         "Subtract",
		Input:        "10 12\n10 14\n100 200\n4294967295 4294967294\n",
		Output:       "2\n4\n100\n1\n",
		SampleInput:  "10 12\n10 14\n100 200\n4294967295 4294967294\n",
		SampleOutput: "2\n4\n100\n1\n",
		TimeLimit:    5,
		MemoryLimit:  200,
	}
	AddProblem(p)

	AddDailyChallenge(time.Date(2015, time.October, 18, 0, 0, 0, 0, time.Local), Easy, 1)
	AddDailyChallenge(time.Date(2015, time.October, 18, 0, 0, 0, 0, time.Local), Medium, 2)
	AddDailyChallenge(time.Date(2015, time.October, 18, 0, 0, 0, 0, time.Local), Hard, 3)
	AddDailyChallenge(time.Date(2015, time.October, 19, 0, 0, 0, 0, time.Local), Easy, 1)
	AddDailyChallenge(time.Date(2015, time.October, 19, 0, 0, 0, 0, time.Local), Medium, 2)
	AddDailyChallenge(time.Date(2015, time.October, 19, 0, 0, 0, 0, time.Local), Hard, 3)
	AddDailyChallenge(time.Date(2015, time.October, 20, 0, 0, 0, 0, time.Local), Easy, 1)
	AddDailyChallenge(time.Date(2015, time.October, 20, 0, 0, 0, 0, time.Local), Medium, 2)
	AddDailyChallenge(time.Date(2015, time.October, 20, 0, 0, 0, 0, time.Local), Hard, 3)
	AddDailyChallenge(time.Date(2015, time.October, 21, 0, 0, 0, 0, time.Local), Easy, 1)
	AddDailyChallenge(time.Date(2015, time.October, 21, 0, 0, 0, 0, time.Local), Medium, 2)
	AddDailyChallenge(time.Date(2015, time.October, 21, 0, 0, 0, 0, time.Local), Hard, 3)

}
