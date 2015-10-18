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
	Username     string
	UserID       int
	ID           int
	ProblemIndex int
	Directory    string
	Verdict      string
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
