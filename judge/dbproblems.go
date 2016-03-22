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
		//Tags:         []string{"Subtract", "Math"},
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
		//Tags:         []string{"Subtract", "Math"},
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

	p.Title = "Back to High School Physics"
	p.Difficulty = 0
	p.UvaID = "10071"
	AddProblem(p)

	p.Title = "Above Average"
	p.Difficulty = 1
	p.UvaID = "10370"
	AddProblem(p)

	p.Title = "Peter's Smokes"
	p.Difficulty = 1
	p.UvaID = "10346"
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

	p.Title = "Permutation Arrays"
	p.Difficulty = 2
	p.UvaID = "482"
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

	p.Title = "Perfection"
	p.Difficulty = 1
	p.UvaID = "382"
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

	p.Title = "Soundex"
	p.Difficulty = 2
	p.UvaID = "10260"
	AddProblem(p)

	p.Title = "Group Reverse"
	p.Difficulty = 2
	p.UvaID = "11192"
	AddProblem(p)

	p.Title = "Newspaper"
	p.Difficulty = 2
	p.UvaID = "11340"
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

	p.Title = "The Blocks Problem"
	p.Difficulty = 1
	p.UvaID = "101"
	AddProblem(p)

	p.Title = "Odd Sum"
	p.Difficulty = 1
	p.UvaID = "10783"
	AddProblem(p)

	//E
	p.SkillID = "E"

	p.Title = "Google is Feeling Lucky"
	p.Difficulty = 2
	p.UvaID = "12015"
	AddProblem(p)

	p.Title = "Greedy Gift Givers"
	p.Difficulty = 2
	p.UvaID = "119"
	AddProblem(p)

	p.Title = "Train Tracks"
	p.Difficulty = 3
	p.UvaID = "11586"
	AddProblem(p)

	p.Title = "Burger Time?"
	p.Difficulty = 3
	p.UvaID = "11661"
	AddProblem(p)

	p.Title = "To Carry or not to Carry"
	p.Difficulty = 1
	p.UvaID = "10469"
	AddProblem(p)

	p.Title = "Numbering Roads"
	p.Difficulty = 2
	p.UvaID = "11723"
	AddProblem(p)

	p.Title = "Brick Game"
	p.Difficulty = 2
	p.UvaID = "11875"
	AddProblem(p)

	p.Title = "The 3n + 1 problem"
	p.Difficulty = 0
	p.UvaID = "100"
	AddProblem(p)

	p.Title = "Primary Arithmetic"
	p.Difficulty = 1
	p.UvaID = "10035"
	AddProblem(p)

	p.Title = "Box of Bricks"
	p.Difficulty = 1
	p.UvaID = "591"
	AddProblem(p)

	//F1
	p.SkillID = "F1"

	p.Title = "Average Speed"
	p.Difficulty = 3
	p.UvaID = "10281"
	AddProblem(p)

	p.Title = "Etruscan Warriors Never Play Chess"
	p.Difficulty = 3
	p.UvaID = "11614"
	AddProblem(p)

	p.Title = "Code Refactoring"
	p.Difficulty = 3
	p.UvaID = "10879"
	AddProblem(p)

	p.Title = "Different Digits"
	p.Difficulty = 3
	p.UvaID = "12527"
	AddProblem(p)

	p.Title = "Feynman"
	p.Difficulty = 2
	p.UvaID = "12149"
	AddProblem(p)

	p.Title = "Pizza Cutting"
	p.Difficulty = 1
	p.UvaID = "10079"
	AddProblem(p)

	p.Title = "Pi"
	p.Difficulty = 2
	p.UvaID = "412"
	AddProblem(p)

	p.Title = "LCM Cardinality"
	p.Difficulty = 3
	p.UvaID = "10892"
	AddProblem(p)

	p.Title = "Prime Distance"
	p.Difficulty = 3
	p.UvaID = "10140"
	AddProblem(p)

	p.Title = "Goldbach's Conjecture"
	p.Difficulty = 1
	p.UvaID = "543"
	AddProblem(p)

	p.Title = "Goldbach's Conjecture (II)"
	p.Difficulty = 2
	p.UvaID = "686"
	AddProblem(p)

	//F2
	p.SkillID = "F2"

	p.Title = "Error Correction"
	p.Difficulty = 1
	p.UvaID = "541"
	AddProblem(p)

	p.Title = "Rotated square"
	p.Difficulty = 3
	p.UvaID = "10855"
	AddProblem(p)

	p.Title = "Spiral Tap"
	p.Difficulty = 3
	p.UvaID = "10920"
	AddProblem(p)

	p.Title = "Jolly Jumpers"
	p.Difficulty = 1
	p.UvaID = "10038"
	AddProblem(p)

	p.Title = "Machined Surfaces"
	p.Difficulty = 2
	p.UvaID = "414"
	AddProblem(p)

	p.Title = "Mirror, Mirror"
	p.Difficulty = 3
	p.UvaID = "466"
	AddProblem(p)

	p.Title = "Add bricks in the wall"
	p.Difficulty = 3
	p.UvaID = "11040"
	AddProblem(p)

	p.Title = "Symmetric Matrix"
	p.Difficulty = 3
	p.UvaID = "11349"
	AddProblem(p)

	p.Title = "Have Fun with Matrices"
	p.Difficulty = 3
	p.UvaID = "11360"
	AddProblem(p)

	//F3
	p.SkillID = "F3"

	p.Title = "A Match Making Problem"
	p.Difficulty = 3
	p.UvaID = "12210"
	AddProblem(p)

	p.Title = "Work Reduction"
	p.Difficulty = 3
	p.UvaID = "10670"
	AddProblem(p)

	p.Title = "Minimal coverage"
	p.Difficulty = 3
	p.UvaID = "10020"
	AddProblem(p)

	p.Title = "All in All"
	p.Difficulty = 1
	p.UvaID = "10340"
	AddProblem(p)

	p.Title = "Dragon of Loowater"
	p.Difficulty = 2
	p.UvaID = "11292"
	AddProblem(p)

	p.Title = "Station Balance"
	p.Difficulty = 3
	p.UvaID = "410"
	AddProblem(p)

	p.Title = "The Bus Driver Problem"
	p.Difficulty = 3
	p.UvaID = "11389"
	AddProblem(p)

	p.Title = "Scarecrow"
	p.Difficulty = 3
	p.UvaID = "12405"
	AddProblem(p)

	p.Title = "Commando War"
	p.Difficulty = 3
	p.UvaID = "11729"
	AddProblem(p)

	p.Title = "ShellSort"
	p.Difficulty = 2
	p.UvaID = "10152"
	AddProblem(p)

}
