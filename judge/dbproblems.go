package judge

import ".././problems"

func AddSamples() {
	p := problems.Problem{
		Index: -1,
		Title: "Hello, world!",
		Description: "You've just been born into the world and " +
			"there's a lot of people around you. Your job is to call them all " +
			"by their name and saying hello before it. Given a <name> which consists " +
			"alphaneumeric characters and no spaces, print Hello, <name> in a line.",
		SkillID:      "A",
		Difficulty:   1,
		Input:        "Sean\nMatthew\nJM\nKiel\n",
		Output:       "Hello, Sean\nHello, Matthew\nHello, JM\nHello, Kiel\n",
		SampleInput:  "Sean\nMatthew\n",
		SampleOutput: "Hello, Sean\nHello, Matthew\n",
		TimeLimit:    2,
		MemoryLimit:  200,
		Tags:         []string{"Subtract", "Math"},
	}
	AddProblem(p)
	p = problems.Problem{
		Index: -1,
		Title: "Caveman",
		Description: "You are in a cave and there's echo everywhere. Each time you say something " +
			"an echo is repeated three times. For each <line> of input, output the <line> and append EcHO ECHO ECHo.",
		SkillID:      "A",
		Difficulty:   1,
		Input:        "Hello\nBye\nEchooooo\nNooooooooooooooooo\n",
		Output:       "Hello EcHO ECHO ECHo\nBye EcHO ECHO ECHo\nEchooooo EcHO ECHO ECHo\nNooooooooooooooooo EcHO ECHO ECHo\n",
		SampleInput:  "Hello\nBye\n",
		SampleOutput: "Hello EcHO ECHO ECHo\nBye EcHO ECHO ECHo\n",
		TimeLimit:    2,
		MemoryLimit:  200,
		Tags:         []string{"Subtract", "Math"},
	}
	AddProblem(p)
	//B
	p = problems.Problem{
		Index:      -1,
		Title:      "Division of Nlogonia",
		SkillID:    "B",
		Difficulty: 1,
		UvaID:      "11498",
	}
	AddProblem(p)

	p.Title = "Cost Cutting"
	p.Difficulty = 1
	p.UvaID = "11727"
	AddProblem(p)

	p.Title = "Save Setu"
	p.Difficulty = 2
	p.UvaID = "12403"
	AddProblem(p)

	p.Title = "Celebrity jeopardy"
	p.Difficulty = 2
	p.UvaID = "1124"
	AddProblem(p)

	p.Title = "Hajj-e-Akbar"
	p.Difficulty = 2
	p.UvaID = "12577"
	AddProblem(p)

	p.Title = "Packing for Holiday"
	p.Difficulty = 2
	p.UvaID = "12372"
	AddProblem(p)

	p.Title = "Lumberjack Sequencing"
	p.Difficulty = 2
	p.UvaID = "11942"
	AddProblem(p)

	//C1
	p.SkillID = "C1"

	p.Title = "Relational Operator"
	p.Difficulty = 0
	p.UvaID = "11172"
	AddProblem(p)

	p.Title = "Beat the Spread!"
	p.Difficulty = 1
	p.UvaID = "10812"
	AddProblem(p)

	p.Title = "Automatic Answer"
	p.Difficulty = 1
	p.UvaID = "11547"
	AddProblem(p)

	p.Title = "Ecological Premium"
	p.Difficulty = 1
	p.UvaID = "10300"
	AddProblem(p)

	p.Title = "Summing Digits"
	p.Difficulty = 1
	p.UvaID = "11332"
	AddProblem(p)

	p.Title = "A Change in Thermal Unit"
	p.Difficulty = 2
	p.UvaID = "11984"
	AddProblem(p)

	p.Title = "Zapping"
	p.Difficulty = 2
	p.UvaID = "12468"
	AddProblem(p)

	p.Title = "Love Calculator"
	p.Difficulty = 2
	p.UvaID = "10424"
	AddProblem(p)

	//C2
	p.SkillID = "C2"

	p.Title = "Triangle Wave"
	p.Difficulty = 1
	p.UvaID = "488"
	AddProblem(p)

	p.Title = "Language Detection"
	p.Difficulty = 2
	p.UvaID = "12250"
	AddProblem(p)

	p.Title = "Emoogle Balance"
	p.Difficulty = 2
	p.UvaID = "12279"
	AddProblem(p)

	p.Title = "Horror Dash"
	p.Difficulty = 1
	p.UvaID = "11799"
	AddProblem(p)

	p.Title = "Jumping Mario"
	p.Difficulty = 1
	p.UvaID = "11764"
	AddProblem(p)

	p.Title = `A Special "Happy Birthday" Song!!!`
	p.Difficulty = 3
	p.UvaID = "12554"
	AddProblem(p)

	p.Title = "Guessing Game"
	p.Difficulty = 2
	p.UvaID = "10530"
	AddProblem(p)

	//D1
	p.SkillID = "D1"

	p.Title = "Clock Hands"
	p.Difficulty = 1
	p.UvaID = "579"
	AddProblem(p)

	p.Title = "Combination Lock"
	p.Difficulty = 2
	p.UvaID = "10550"
	AddProblem(p)

	p.Title = "Tariff Plan"
	p.Difficulty = 3
	p.UvaID = "12157"
	AddProblem(p)

	p.Title = "Digits"
	p.Difficulty = 3
	p.UvaID = "11687"
	AddProblem(p)

	p.Title = "Intersecting Lines"
	p.Difficulty = 2
	p.UvaID = "378"
	AddProblem(p)

	p.Title = "Is this the easiest problem?"
	p.Difficulty = 2
	p.UvaID = "11479"
	AddProblem(p)

	p.Title = "Points in Figures: Rectangles"
	p.Difficulty = 2
	p.UvaID = "476"
	AddProblem(p)

	p.Title = "Points in Figures: Rectangles and Circles"
	p.Difficulty = 2
	p.UvaID = "477"
	AddProblem(p)

	p.Title = "Behold my quadrangle"
	p.Difficulty = 2
	p.UvaID = "11455"
	AddProblem(p)

	//D2
	p.SkillID = "D2"

	p.Title = "Palindromes"
	p.Difficulty = 1
	p.UvaID = "401"
	AddProblem(p)

	p.Title = "TEX Quotes"
	p.Difficulty = 0
	p.UvaID = "272"
	AddProblem(p)

	p.Title = "One-Two-Three"
	p.Difficulty = 2
	p.UvaID = "12289"
	AddProblem(p)

	p.Title = "Hangman Judge"
	p.Difficulty = 2
	p.UvaID = "489"
	AddProblem(p)

	p.Title = "WERTYU"
	p.Difficulty = 1
	p.UvaID = "10082"
	AddProblem(p)

	//D3
	p.SkillID = "D3"

	p.Title = "Loansome Car Buyer"
	p.Difficulty = 3
	p.UvaID = "10114"
	AddProblem(p)

	p.Title = "Robot Instructions"
	p.Difficulty = 3
	p.UvaID = "12503"
	AddProblem(p)

	p.Title = "The Snail"
	p.Difficulty = 1
	p.UvaID = "573"
	AddProblem(p)

	p.Title = "Die Game"
	p.Difficulty = 2
	p.UvaID = "10409"
	AddProblem(p)

	p.Title = "Master-Mind Hints"
	p.Difficulty = 2
	p.UvaID = "340"
	AddProblem(p)

	//E
	p.SkillID = "E"

	p.Title = "Google is Feeling Lucky"
	p.Difficulty = 1
	p.UvaID = "12015"
	AddProblem(p)

	p.Title = "Greedy Gift Givers"
	p.Difficulty = 1
	p.UvaID = "119"
	AddProblem(p)

	p.Title = "Train Tracks"
	p.Difficulty = 1
	p.UvaID = "11586"
	AddProblem(p)

	p.Title = "Burger Time?"
	p.Difficulty = 1
	p.UvaID = "11661"
	AddProblem(p)

}
