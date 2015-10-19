package judge

import (
	".././dao"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func getProblemStatistics(problemID int) (submitted, accepted int) {
	db, err := sql.Open("sqlite3", dao.DatabaseURL)
	if err != nil {
		return
	}
	defer db.Close()

	db.QueryRow("SELECT COUNT(*) FROM submissions WHERE problem_id = ?", problemID).Scan(&submitted)
	db.QueryRow("SELECT COUNT(*) FROM submissions WHERE verdict = ? AND problem_id = ?", Accepted, problemID).Scan(&accepted)

	return
}
