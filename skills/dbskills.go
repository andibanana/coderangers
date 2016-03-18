package skills

func AddSamples() {
	s := Skill{
		ID:                       "A",
		Title:                    "Introduction to Competitive Programming",
		Description:              "Trivial problems focused on familiarizing yourself with the software",
		NumberOfProblemsToUnlock: 2,
		//Prerequisites:			  []string{"a", "b", "c"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "B",
		Title:                    "Ad Hoc 101",
		Description:              "Problems that can be solved with basic programming skills... I hope...",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"A"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "C1",
		Title:                    "Simple Math",
		Description:              "Problems involving basic math problems such as multiplication and fractions",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"B"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "C2",
		Title:                    "Garbage in, Garbage out",
		Description:              "Memory is cheap but not infinite, plus we need to cut down on defense spending where we can if we wanna keep the free coffee at the mess hall",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"B"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "D1",
		Title:                    "More Math",
		Description:              "A lot of people don't like math. I intend to change that",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"C1"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "D2",
		Title:                    "Text Twist",
		Description:              "'RACE CAR' read backwards is actually 'RACE CAR'... who knew?",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"C2"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "D3",
		Title:                    "Try Try Again",
		Description:              "If you keep hitting the compile button, it's bound to work eventually right?",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"C2"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "E",
		Title:                    "Back to Basics",
		Description:              "I hope you still know how to pitch a tent",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"C1", "D2", "D3"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "F1",
		Title:                    "Even More Math",
		Description:              "As if there wasn't enough numbers already, they added letters",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"D1", "E"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "F2",
		Title:                    "Know your Data Structures I",
		Description:              "There is one rule in this organization... actually a lot more but this one is important; Keep your data organized or die",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"E"},
	}
	addSkill(s)

	s = Skill{
		ID:                       "F3",
		Title:                    "Greed is Good",
		Description:              "Follow the money, and hopefully it leads to more money",
		NumberOfProblemsToUnlock: 3,
		Prerequisites:            []string{"D3"},
	}
	addSkill(s)
}
