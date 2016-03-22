package problems

const (
	Received            = "received"
	Compiling           = "compiling"
	Running             = "running"
	Judging             = "judging"
	Inqueue             = "inqueue"
	Accepted            = "accepted"
	PresentationError   = "presentation error"
	WrongAnswer         = "wrong answer"
	CompileError        = "compile error"
	RuntimeError        = "runtime error"
	TimeLimitExceeded   = "time limit exceeded"
	MemoryLimitExceeded = "memory limit exceeded"
	OutputLimitExceeded = "output limit exceeded"
	SubmissionError     = "submission error"
	RestrictedFunction  = "restricted function"
	CantBeJudged        = "can't be judged"
)

type Problem struct {
	Index        int
	Title        string
	Description  string
	Difficulty   int
	SkillID      string
	SampleInput  string
	SampleOutput string
	UvaID        string
	Input        string `json:"-"`
	Output       string `json:"-"`
	TimeLimit    int
	MemoryLimit  int
	Solved       bool
	Tags         []string `json:"-"`
}
