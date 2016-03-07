package judge

import (
	".././dao"
	".././problems"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func getProblemStatistics(problemID int) (submitted int, verdictData VerdictData) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()

	db.QueryRow("SELECT COUNT(*) FROM submissions WHERE problem_id = ?", problemID).Scan(&submitted)
	row, err := db.Query("SELECT COUNT(*), verdict FROM submissions WHERE problem_id = ? GROUP BY verdict ORDER BY verdict", problemID)

	for row.Next() {
		var count int
		var verdict string
		row.Scan(&count, &verdict)
		switch verdict {
		case problems.Accepted:
			verdictData.Accepted = count
		case problems.WrongAnswer:
			verdictData.WrongAnswer = count
		case problems.CompileError:
			verdictData.CompileError = count
		case problems.RuntimeError:
			verdictData.RuntimeError = count
		case problems.TimeLimitExceeded:
			verdictData.TimeLimitExceeded = count

		}
	}

	return
}
