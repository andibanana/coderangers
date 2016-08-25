package judge

import "coderangers/problems"

func AddSamples() (err error) {
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
	err = AddProblem(p)
	if err != nil {
		return err
	}
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
	err = AddProblem(p)
	if err != nil {
		return err
	}
	//B
	p = problems.Problem{
		Index:      -1,
		Title:      "Division of Nlogonia",
		SkillID:    "B",
		Difficulty: 1,
		UvaID:      "11498",
	}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Cost Cutting"
	p.Difficulty = 1
	p.UvaID = "11727"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Save Setu"
	p.Difficulty = 2
	p.UvaID = "12403"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Celebrity jeopardy"
	p.Difficulty = 2
	p.UvaID = "1124"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Hajj-e-Akbar"
	p.Difficulty = 2
	p.UvaID = "12577"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Packing for Holiday"
	p.Difficulty = 2
	p.UvaID = "12372"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Lumberjack Sequencing"
	p.Difficulty = 2
	p.UvaID = "11942"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	//C1
	p.SkillID = "C1"

	p.Title = "Relational Operator"
	p.Difficulty = 0
	p.UvaID = "11172"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Beat the Spread!"
	p.Difficulty = 1
	p.UvaID = "10812"
	p.Tags = []string{"arithmetic", "algebra"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Automatic Answer"
	p.Difficulty = 1
	p.UvaID = "11547"
	p.Tags = []string{"arithmetic"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Ecological Premium"
	p.Difficulty = 1
	p.UvaID = "10300"
	p.Tags = []string{"arithmetic"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Summing Digits"
	p.Difficulty = 1
	p.UvaID = "11332"
	p.Tags = []string{"arithmetic", "digits"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "A Change in Thermal Unit"
	p.Difficulty = 2
	p.UvaID = "11984"
	p.Tags = []string{"arithmetic", "algebra"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Zapping"
	p.Difficulty = 2
	p.UvaID = "12468"
	p.Tags = []string{"arithmetic"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Love Calculator"
	p.Difficulty = 2
	p.UvaID = "10424"
	p.Tags = []string{"arithmetic", "digits", "strings"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Back to High School Physics"
	p.Difficulty = 0
	p.UvaID = "10071"
	p.Tags = []string{"arithmetic"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Above Average"
	p.Difficulty = 1
	p.UvaID = "10370"
	p.Tags = []string{"arithmetic"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Peter's Smokes"
	p.Difficulty = 1
	p.UvaID = "10346"
	p.Tags = []string{"arithmetic", "simulation"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	//C2
	p.SkillID = "C2"

	p.Title = "Triangle Wave"
	p.Difficulty = 1
	p.UvaID = "488"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Language Detection"
	p.Difficulty = 2
	p.UvaID = "12250"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Emoogle Balance"
	p.Difficulty = 2
	p.UvaID = "12279"
	p.Tags = []string{"arithmetic"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Horror Dash"
	p.Difficulty = 1
	p.UvaID = "11799"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Jumping Mario"
	p.Difficulty = 1
	p.UvaID = "11764"
	p.Tags = []string{"simulation"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = `A Special "Happy Birthday" Song!!!`
	p.Difficulty = 3
	p.UvaID = "12554"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Guessing Game"
	p.Difficulty = 2
	p.UvaID = "10530"
	p.Tags = []string{"simulation"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Permutation Arrays"
	p.Difficulty = 2
	p.UvaID = "482"
	p.Tags = []string{"1darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	//D1
	p.SkillID = "D1"

	p.Title = "Clock Hands"
	p.Difficulty = 1
	p.UvaID = "579"
	p.Tags = []string{"arithmetic", "algebra", "angles"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Combination Lock"
	p.Difficulty = 2
	p.UvaID = "10550"
	p.Tags = []string{"arithmetic", "angles"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Tariff Plan"
	p.Difficulty = 3
	p.UvaID = "12157"
	p.Tags = []string{"arithmetic", "algebra"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Digits"
	p.Difficulty = 3
	p.UvaID = "11687"
	p.Tags = []string{"digits"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Intersecting Lines"
	p.Difficulty = 2
	p.UvaID = "378"
	p.Tags = []string{"arithmetic", "algebra", "geometry"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Is this the easiest problem?"
	p.Difficulty = 2
	p.UvaID = "11479"
	p.Tags = []string{"arithmetic", "geometry"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Points in Figures: Rectangles"
	p.Difficulty = 2
	p.UvaID = "476"
	p.Tags = []string{"arithmetic", "geometry", "1darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Points in Figures: Rectangles and Circles"
	p.Difficulty = 2
	p.UvaID = "477"
	p.Tags = []string{"arithmetic", "geometry", "1darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Behold my quadrangle"
	p.Difficulty = 2
	p.UvaID = "11455"
	p.Tags = []string{"arithmetic", "geometry"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Perfection"
	p.Difficulty = 1
	p.UvaID = "382"
	p.Tags = []string{"arithmetic", "factors"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	//D2
	p.SkillID = "D2"

	p.Title = "Palindromes"
	p.Difficulty = 1
	p.UvaID = "401"
	p.Tags = []string{"strings", "character-mapping"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "TEX Quotes"
	p.Difficulty = 0
	p.UvaID = "272"
	p.Tags = []string{"strings", "character-mapping"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "One-Two-Three"
	p.Difficulty = 2
	p.UvaID = "12289"
	p.Tags = []string{"strings"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Hangman Judge"
	p.Difficulty = 2
	p.UvaID = "489"
	p.Tags = []string{"strings", "simulation"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "WERTYU"
	p.Difficulty = 1
	p.UvaID = "10082"
	p.Tags = []string{"strings", "character-mapping"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Soundex"
	p.Difficulty = 2
	p.UvaID = "10260"
	p.Tags = []string{"strings", "character-mapping"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Group Reverse"
	p.Difficulty = 2
	p.UvaID = "11192"
	p.Tags = []string{"strings", "1darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Newspaper"
	p.Difficulty = 2
	p.UvaID = "11340"
	p.Tags = []string{"strings", "arithmetic"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	//D3
	p.SkillID = "D3"

	p.Title = "Loansome Car Buyer"
	p.Difficulty = 3
	p.UvaID = "10114"
	p.Tags = []string{"arithmetic", "simulation"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Robot Instructions"
	p.Difficulty = 3
	p.UvaID = "12503"
	p.Tags = []string{"simulation", "interpretation", "1darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "The Snail"
	p.Difficulty = 1
	p.UvaID = "573"
	p.Tags = []string{"simulation", "arithmetic"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Die Game"
	p.Difficulty = 2
	p.UvaID = "10409"
	p.Tags = []string{"simulation", "interpretation"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Master-Mind Hints"
	p.Difficulty = 2
	p.UvaID = "340"
	p.Tags = []string{"simulation", "1darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "The Blocks Problem"
	p.Difficulty = 1
	p.UvaID = "101"
	p.Tags = []string{"simulation", "interpretation", "2darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Odd Sum"
	p.Difficulty = 1
	p.UvaID = "10783"
	p.Tags = []string{"arithmetic"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	//E
	p.SkillID = "E"

	p.Title = "Google is Feeling Lucky"
	p.Difficulty = 2
	p.UvaID = "12015"
	p.Tags = []string{"arithmetic", "1darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Greedy Gift Givers"
	p.Difficulty = 2
	p.UvaID = "119"
	p.Tags = []string{"arithmetic", "1darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Train Tracks"
	p.Difficulty = 3
	p.UvaID = "11586"
	p.Tags = []string{"pattern"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Burger Time?"
	p.Difficulty = 3
	p.UvaID = "11661"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "To Carry or not to Carry"
	p.Difficulty = 1
	p.UvaID = "10469"
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Numbering Roads"
	p.Difficulty = 2
	p.UvaID = "11723"
	p.Tags = []string{"arithmetic"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Brick Game"
	p.Difficulty = 2
	p.UvaID = "11875"
	p.Tags = []string{"arithmetic"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "The 3n + 1 problem"
	p.Difficulty = 0
	p.UvaID = "100"
	p.Tags = []string{"arithmetic", "simulation"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Primary Arithmetic"
	p.Difficulty = 1
	p.UvaID = "10035"
	p.Tags = []string{"arithmetic", "1darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Box of Bricks"
	p.Difficulty = 1
	p.UvaID = "591"
	p.Tags = []string{"arithmetic", "1darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	//F1
	p.SkillID = "F1"

	p.Title = "Average Speed"
	p.Difficulty = 3
	p.UvaID = "10281"
	p.Tags = []string{"arithmetic", "simulation"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Etruscan Warriors Never Play Chess"
	p.Difficulty = 3
	p.UvaID = "11614"
	p.Tags = []string{"arithmetic", "algebra", "pattern"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Code Refactoring"
	p.Difficulty = 3
	p.UvaID = "10879"
	p.Tags = []string{"factors"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Different Digits"
	p.Difficulty = 3
	p.UvaID = "12527"
	p.Tags = []string{"digits"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Feynman"
	p.Difficulty = 2
	p.UvaID = "12149"
	p.Tags = []string{"arithmetic", "pattern"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Pizza Cutting"
	p.Difficulty = 1
	p.UvaID = "10079"
	p.Tags = []string{"arithmetic"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Pi"
	p.Difficulty = 2
	p.UvaID = "412"
	p.Tags = []string{"arithmetic", "factors"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "LCM Cardinality"
	p.Difficulty = 3
	p.UvaID = "10892"
	p.Tags = []string{"factors"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Prime Distance"
	p.Difficulty = 3
	p.UvaID = "10140"
	p.Tags = []string{"primes", "1darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Goldbach's Conjecture"
	p.Difficulty = 1
	p.UvaID = "543"
	p.Tags = []string{"primes"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Goldbach's Conjecture (II)"
	p.Difficulty = 2
	p.UvaID = "686"
	p.Tags = []string{"arithmetic", "primes"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	//F2
	p.SkillID = "F2"

	p.Title = "Error Correction"
	p.Difficulty = 1
	p.UvaID = "541"
	p.Tags = []string{"2darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Rotated square"
	p.Difficulty = 3
	p.UvaID = "10855"
	p.Tags = []string{"2darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Spiral Tap"
	p.Difficulty = 3
	p.UvaID = "10920"
	p.Tags = []string{"arithmetic", "pattern"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Jolly Jumpers"
	p.Difficulty = 1
	p.UvaID = "10038"
	p.Tags = []string{"arithmetic", "1darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Machined Surfaces"
	p.Difficulty = 2
	p.UvaID = "414"
	p.Tags = []string{"arithmetic", "1darray", "pattern"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Mirror, Mirror"
	p.Difficulty = 3
	p.UvaID = "466"
	p.Tags = []string{"2darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Add bricks in the wall"
	p.Difficulty = 3
	p.UvaID = "11040"
	p.Tags = []string{"arithmetic", "algebra", ""}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Symmetric Matrix"
	p.Difficulty = 3
	p.UvaID = "11349"
	p.Tags = []string{"triangular-matrix"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Have Fun with Matrices"
	p.Difficulty = 3
	p.UvaID = "11360"
	p.Tags = []string{"simulation", "2darray"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	//F3
	p.SkillID = "F3"

	p.Title = "A Match Making Problem"
	p.Difficulty = 3
	p.UvaID = "12210"
	p.Tags = []string{"1darray", "greedy", "sort"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Work Reduction"
	p.Difficulty = 3
	p.UvaID = "10670"
	p.Tags = []string{"arithmetic", "1darray", "greedy", "sort"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Minimal coverage"
	p.Difficulty = 3
	p.UvaID = "10020"
	p.Tags = []string{"arithmetic", "1darray", "greedy", "sort"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "All in All"
	p.Difficulty = 1
	p.UvaID = "10340"
	p.Tags = []string{"strings", "greedy"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Dragon of Loowater"
	p.Difficulty = 2
	p.UvaID = "11292"
	p.Tags = []string{"1darray", "greedy", "sort"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Station Balance"
	p.Difficulty = 3
	p.UvaID = "410"
	p.Tags = []string{"arithmetic", "1darray", "greedy", "sort"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "The Bus Driver Problem"
	p.Difficulty = 3
	p.UvaID = "11389"
	p.Tags = []string{"arithmetic", "1darray", "greedy", "sort"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Scarecrow"
	p.Difficulty = 3
	p.UvaID = "12405"
	p.Tags = []string{"1darray", "greedy"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "Commando War"
	p.Difficulty = 3
	p.UvaID = "11729"
	p.Tags = []string{"1darray", "greedy", "sort"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	p.Title = "ShellSort"
	p.Difficulty = 2
	p.UvaID = "10152"
	p.Tags = []string{"arithmetic", "1darray", "greedy"}
	err = AddProblem(p)
	if err != nil {
		return err
	}

	return nil
}
