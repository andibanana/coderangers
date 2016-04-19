package dao

import (
	"coderangers/cookies"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

const SQLiteDatabaseURL = "file:database.sqlite?cache=shared&mode=rwc"
const MySQLDatabaseURL = "root:p@ssword@tcp(127.0.0.1:3306)/"
const MySQLDB = "coderangers"
const MySQL = false

var db *sql.DB

func initDB() {
	var err error
	if MySQL {
		db, err = sql.Open("mysql", MySQLDatabaseURL+MySQLDB)
	} else {
		db, err = sql.Open("sqlite3", SQLiteDatabaseURL)
	}
	if err != nil {
		log.Fatal("DB CONNECTION NOT MADE")
	}
}

func Open() (*sql.DB, error) {
	if db == nil {
		initDB()
	}
	return db, db.Ping()
}

func IsAdmin(req *http.Request) bool {
	userID, ok := cookies.GetUserID(req)
	if !ok {
		return false
	}

	db, err := Open()
	if err != nil {
		return false
	}

	err = db.QueryRow("SELECT id FROM user_account WHERE id=? AND admin=?", userID, true).Scan(&userID)
	return err == nil
}

func CheckRealUser(req *http.Request) bool {
	userID, ok := cookies.GetUserID(req)
	if !ok {
		return false
	}
	db, err := Open()
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.QueryRow("SELECT id FROM user_account WHERE id=?", userID).Scan(&userID)
	return err == nil
}

func CreateDB() error {
	var AUTOINCREMENT = "AUTOINCREMENT"
	if MySQL {
		AUTOINCREMENT = "AUTO_INCREMENT"
		db, err := sql.Open("mysql", MySQLDatabaseURL)
		if err != nil {
			return err
		}

		_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + MySQLDB)
		if err != nil {
			return err
		}
	}

	db, err := Open()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE user_account (
			id INTEGER PRIMARY KEY ` + AUTOINCREMENT + `,
					
			username VARCHAR(50) UNIQUE NOT NULL,
			hashed_password CHARACTER(60) NOT NULL,
      email VARCHAR(50) UNIQUE NOT NULL,
			admin BOOLEAN NOT NULL DEFAULT FALSE,
      date_joined DATE NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE skills (
      id VARCHAR(20) PRIMARY KEY,
      
			title VARCHAR(100) NOT NULL,
      description TEXT NOT NULL,
      number_of_problems_to_unlock INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE problems (
      id INTEGER PRIMARY KEY ` + AUTOINCREMENT + `,
			
			title VARCHAR(100) NOT NULL,
      description TEXT NOT NULL,
      skill_id VARCHAR(20) NOT NULL,
      uva_id VARCHAR(100) NOT NULL,
      difficulty INTEGER,
      time_limit INTEGER,
      memory_limit INTEGER,
      sample_input TEXT,
      sample_output TEXT,
      
      FOREIGN KEY(skill_id) REFERENCES skills(id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE tags (
      problem_id INTEGER NOT NULL,
      tag TEXT NOT NULL,
      
      FOREIGN KEY(problem_id) REFERENCES problems(id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE prerequisites (
      skill_id VARCHAR(20) NOT NULL,
      prerequisite_id VARCHAR(20) NOT NULL,
      
      FOREIGN KEY(skill_id) REFERENCES skills(id),
      FOREIGN KEY(prerequisite_id) REFERENCES skills(id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE inputoutput (
      problem_id INTEGER,
      input_number INTEGER,
    
			input TEXT NOT NULL,
      output TEXT NOT NULL,
      
      PRIMARY KEY(problem_id, input_number),
      FOREIGN KEY(problem_id) REFERENCES problems(id)
		)
	`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE submissions (
      id INTEGER PRIMARY KEY ` + AUTOINCREMENT + `,
      problem_id INTEGER,
      user_id INTEGER,
      
      uva_submission_id INTEGER,
			directory VARCHAR(100) NOT NULL,
      verdict VARCHAR(20) NOT NULL,
      language VARCHAR(5),
      timestamp DATETIME NOT NULL,
      runtime_error TEXT,
      runtime NUMERIC,
      
      FOREIGN KEY(user_id) REFERENCES user_account(id),
      FOREIGN KEY(problem_id) REFERENCES problems(id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE badges (
      id INTEGER PRIMARY KEY ` + AUTOINCREMENT + `,
      
			title VARCHAR(100) NOT NULL,
      description VARCHAR(100) NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE received_badges (
        user_id INTEGER,
        badge_id INTEGER,
       
        PRIMARY KEY(user_id, badge_id),
        FOREIGN KEY(user_id) REFERENCES user_account(id),
        FOREIGN KEY(badge_id) REFERENCES badges(id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE viewed_problems (
        user_id INTEGER,
        problem_id INTEGER,
       
        PRIMARY KEY(user_id, problem_id),
        FOREIGN KEY(user_id) REFERENCES user_account(id),
        FOREIGN KEY(problem_id) REFERENCES problems(id)
		)
	`)
	if err != nil {
		return err
	}

	// _, err = db.Exec(`
	// CREATE TABLE user_data (
	// user_id INTEGER PRIMARY KEY,

	// submitted_count INTEGER,
	// accepted_count INTEGER,
	// attempted_count INTEGER,
	// viewed_problems_count INTEGER,
	// experience INTEGER,

	// FOREIGN KEY(user_id) REFERENCES user_account(id)
	// )
	// `)

	// if err != nil {
	// return err
	// }

	return err
}

func AddTables() (err error) {
	db, err := Open()

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS notifications (
      submission_id INTEGER,
      user_id INTEGER,
      viewed boolean,
      
      PRIMARY KEY(submission_id, user_id),
      FOREIGN KEY(submission_id) REFERENCES submissions(id),
      FOREIGN KEY(user_id) REFERENCES user_account(id)
    )
  `)
	if err != nil {
		return err
	}

	return
}
