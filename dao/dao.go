package dao

import (
	".././cookies"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

const DatabaseURL = "file:database.sqlite?cache=shared&mode=rwc"

func IsAdmin(req *http.Request) bool {
	userID, ok := cookies.GetUserID(req)
	if !ok {
		return false
	}

	db, err := sql.Open("sqlite3", DatabaseURL)
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.QueryRow("SELECT id FROM user_account WHERE id=? AND admin=?", userID, true).Scan(&userID)
	return err == nil
}

func CreateDB() error {
	db, err := sql.Open("sqlite3", DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE user_account (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
					
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
		CREATE TABLE problems (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
			
			title VARCHAR(100) NOT NULL,
      description VARCHAR(200) NOT NULL,
      skill_id STRING NOT NULL,
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
		CREATE TABLE skills (
      id STRING PRIMARY KEY,
      
			title VARCHAR(100) NOT NULL,
      description VARCHAR(200) NOT NULL,
      number_of_problems_to_unlock INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE prerequisites (
      skill_id STRING NOT NULL,
      prerequisite_id STRING NOT NULL,
      
      FOREIGN KEY(skill_id) REFERENCES skills(id)
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
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      problem_id INTEGER,
      user_id INTEGER,
      
      uva_submission_id INTEGER,
			directory VARCHAR(100) NOT NULL,
      verdict VARCHAR(100) NOT NULL,
      timestamp DATETIME NOT NULL,
      runtime_error TEXT,
      runtime NUMERIC,
      
      FOREIGN KEY(user_id) REFERENCES user_account(id)
      FOREIGN KEY(problem_id) REFERENCES problems(id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE badges (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      
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

	_, err = db.Exec(`
		CREATE TABLE user_data (
        user_id INTEGER PRIMARY KEY,
        
        submitted_count INTEGER,
        accepted_count INTEGER,
        attempted_count INTEGER,
        viewed_problems_count INTEGER,
        experience INTEGER,
        
        FOREIGN KEY(user_id) REFERENCES user_account(id)
		)
	`)

	if err != nil {
		return err
	}

	return err
}
