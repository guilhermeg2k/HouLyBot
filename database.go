package main

import (
	"database/sql"
)

type DataBase struct {
	DB *sql.DB
}

func (db *DataBase) setupDB() error {
	database, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return err
	}
	createTableStatement := `
		CREATE TABLE IF NOT EXISTS teams (
			id INTEGER PRIMARY KEY, 
			name TEXT UNIQUE NOT NULL, 
			url TEXT NOT NULL UNIQUE
		);
	`
	_, err = database.Exec(createTableStatement)
	if err != nil {
		return err
	}
	createTableStatement = `
		CREATE TABLE IF NOT EXISTS commands (
			id INTEGER PRIMARY KEY, 
			command TEXT UNIQUE NOT NULL UNIQUE,
			syntax TEXT NOT NULL, 
			description TEXT NOT NULL
		);
	`
	_, err = database.Exec(createTableStatement)
	if err != nil {
		return err
	}
	createTableStatement = `
		CREATE TABLE IF NOT EXISTS logs(
			id INTEGER PRIMARY KEY,
			type INTEGER NOT NULL,
			file TEXT NOT NULL,
			time TEXT NOT NULL,
			log TEXT NOT NULL
		);
	`
	_, err = database.Exec(createTableStatement)
	if err != nil {
		return err
	}
	db.DB = database
	return nil
}

func (db *DataBase) loadAllTeams() ([]Team, error) {
	var teams []Team
	allTeams, err := db.DB.Query("SELECT name, url FROM teams")
	if err != nil {
		return teams, err
	}
	for allTeams.Next() {
		var team Team
		err = allTeams.Scan(&team.name, &team.url)
		if err != nil {
			Log.Error("Failed to load the team " + team.name + " with the url " + team.url)
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (db *DataBase) createTeam(name, url string) error {
	createTeam, err := db.DB.Prepare("INSERT INTO teams(name, url) VALUES (?,?)")
	if err != nil {
		return err
	}
	_, err = createTeam.Exec(name, url)
	if err != nil {
		return err
	}
	return nil
}

func (db *DataBase) createLog(log LogData) error {
	createTeam, err := db.DB.Prepare("INSERT INTO logs(type, file, time, log) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = createTeam.Exec(log.logType, log.file, log.time, log.log)
	if err != nil {
		return err
	}
	return nil
}
