package main

import (
	"database/sql"
)

type DataBase struct {
	conn *sql.DB
}

type Team struct {
	id   uint
	name string
	url  string
}

type Log struct {
	id      uint
	logType uint
	file    string
	time    string
	text    string
}
type Command struct {
	id          uint
	name        string
	syntax      string
	description string
}

func newDataBase() (DataBase, error) {
	var db DataBase
	conn, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return db, err
	}
	db.conn = conn
	createTableStatement := `
		CREATE TABLE IF NOT EXISTS teams (
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			url TEXT NOT NULL UNIQUE
		);
	`
	_, err = conn.Exec(createTableStatement)
	if err != nil {
		return db, err
	}
	createTableStatement = `
		CREATE TABLE IF NOT EXISTS commands (
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE NOT NULL UNIQUE,
			syntax TEXT NOT NULL,
			description TEXT NOT NULL
		);
	`
	_, err = conn.Exec(createTableStatement)
	if err != nil {
		return db, err
	}
	createTableStatement = `
		CREATE TABLE IF NOT EXISTS logs(
			id INTEGER PRIMARY KEY,
			type INTEGER NOT NULL,
			file TEXT NOT NULL,
			time TEXT NOT NULL,
			text TEXT NOT NULL
		);
	`
	_, err = conn.Exec(createTableStatement)
	if err != nil {
		return db, err
	}
	return db, nil
}

func (db *DataBase) createTeam(team Team) error {
	createTeam, err := db.conn.Prepare("INSERT INTO teams(name, url) VALUES (?,?)")
	if err != nil {
		return err
	}
	_, err = createTeam.Exec(team.name, team.url)
	if err != nil {
		return err
	}
	return nil
}

func (db *DataBase) getAllTeams() ([]Team, error) {
	var teams []Team
	allTeams, err := db.conn.Query("SELECT * FROM teams")
	if err != nil {
		return teams, err
	}
	for allTeams.Next() {
		var team Team
		err = allTeams.Scan(&team.id, &team.name, &team.url)
		if err != nil {
			logger.Error(err.Error())
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (db *DataBase) createCommand(command Command) error {
	createCommand, err := db.conn.Prepare(`
		INSERT INTO commands(
			name,
			syntax,
			description
		) VALUES (?,?,?)
	`)
	if err != nil {
		return err
	}
	_, err = createCommand.Exec(command.name, command.syntax, command.description)
	if err != nil {
		return err
	}
	return nil
}

func (db *DataBase) getAllCommands() ([]Command, error) {
	var commands []Command
	commandsQuery, err := db.conn.Query("select * from commands order by name;")
	if err != nil {
		return commands, err
	}
	for commandsQuery.Next() {
		var command Command
		err = commandsQuery.Scan(&command.id, &command.name, &command.syntax, &command.description)
		if err != nil {
			logger.Error(err.Error())
		}
		commands = append(commands, command)
	}
	return commands, nil
}

func (db *DataBase) createLog(log Log) error {
	createTeam, err := db.conn.Prepare("INSERT INTO logs(type, file, time, log) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = createTeam.Exec(log.logType, log.file, log.time, log.text)
	if err != nil {
		return err
	}
	return nil
}

func (db *DataBase) getAllLogs() ([]Log, error) {
	var logs []Log
	allTeams, err := db.conn.Query("SELECT * FROM logs")
	if err != nil {
		return logs, err
	}
	for allTeams.Next() {
		var log Log
		err = allTeams.Scan(&log.id, &log.logType, &log.file, &log.time, &log.text)
		if err != nil {
			logger.Error(err.Error())
		}
		logs = append(logs, log)
	}
	return logs, nil
}
