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
	createTeamsTable, err := database.Prepare(`
		CREATE TABLE IF NOT EXISTS teams (id INTEGER PRIMARY KEY, name TEXT UNIQUE, url TEXT);
		CREATE TABLE IF NOT EXISTS commands (id INTEGER PRIMARY KEY, command TEXT UNIQUE, syntax TEXT, description TEXT);
	`)
	if err != nil {
		return err
	}
	createTeamsTable.Exec()
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
		allTeams.Scan(&team.name, &team.url)
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

func (db *DataBase) getTeamURL(teamName string) (string, error) {
	var url string
	err := db.DB.QueryRow("SELECT url FROM teams where name = ?", teamName).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}
