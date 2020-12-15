package main

import (
	"database/sql"
	"fmt"
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
	commandType int
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
			type INTEGER NOT NULL,
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
			type,
			syntax,
			description
		) VALUES (?,?,?,?)
	`)
	if err != nil {
		return err
	}
	_, err = createCommand.Exec(command.name, command.commandType, command.syntax, command.description)
	if err != nil {
		return err
	}
	return nil
}

func (db *DataBase) getCommands(commandType string) ([]Command, error) {
	var commands []Command
	commandsQuery, err := db.conn.Query(
		fmt.Sprintf("select * from commands where type = %s order by name;", commandType),
	)
	if err != nil {
		return commands, err
	}
	for commandsQuery.Next() {
		var command Command
		err = commandsQuery.Scan(&command.id, &command.commandType, &command.name, &command.syntax, &command.description)
		if err != nil {
			logger.Error(err.Error())
		}
		commands = append(commands, command)
	}

	return commands, nil
}

func (db *DataBase) deleteAllCommands() error {
	_, err := db.conn.Exec("delete from commands;")
	if err != nil {
		return err
	}

	return nil
}

func (db *DataBase) getBotCommands() ([]Command, error) {
	return db.getCommands("0")
}

func (db *DataBase) getCliCommands() ([]Command, error) {
	return db.getCommands("1")
}

func (db *DataBase) createLog(log Log) error {
	createTeam, err := db.conn.Prepare("INSERT INTO logs(type, file, time, text) VALUES (?,?,?,?)")
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
	allLogs, err := db.conn.Query("SELECT * FROM logs")
	if err != nil {
		return logs, err
	}
	for allLogs.Next() {
		var log Log
		err = allLogs.Scan(&log.id, &log.logType, &log.file, &log.time, &log.text)
		if err != nil {
			logger.Error(err.Error())
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (db *DataBase) getLogsByType(logType string) ([]Log, error) {
	var logs []Log
	logsQuery, err := db.conn.Query(fmt.Sprintf("SELECT * FROM logs WHERE type = %s", logType))
	if err != nil {
		return logs, err
	}
	for logsQuery.Next() {
		var log Log
		err = logsQuery.Scan(&log.id, &log.logType, &log.file, &log.time, &log.text)
		if err != nil {
			logger.Error(err.Error())
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (db *DataBase) getLogsWithLimit(limit string) ([]Log, error) {
	var logs []Log
	logsQuery, err := db.conn.Query(fmt.Sprintf("SELECT * FROM logs LIMIT  %s", limit))
	if err != nil {
		return logs, err
	}
	for logsQuery.Next() {
		var log Log
		err = logsQuery.Scan(&log.id, &log.logType, &log.file, &log.time, &log.text)
		if err != nil {
			logger.Error(err.Error())
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (db *DataBase) getLogsByTypeWithLimit(logType, limit string) ([]Log, error) {
	var logs []Log
	logsQuery, err := db.conn.Query(
		fmt.Sprintf(
			"SELECT * FROM logs WHERE type = %s LIMIT %s;",
			logType,
			limit,
		),
	)
	if err != nil {
		return logs, err
	}
	for logsQuery.Next() {
		var log Log
		err = logsQuery.Scan(&log.id, &log.logType, &log.file, &log.time, &log.text)
		if err != nil {
			logger.Error(err.Error())
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (db *DataBase) updateCommands() error {
	var commands []Command
	err := db.deleteAllCommands()
	if err != nil {
		return err
	}
	commands = append(commands,
		Command{
			name:        "!team",
			commandType: 0,
			syntax:      "!team <team-name>",
			description: "This command retrieves informations like roster, hltv ranking, next matches and recent results of a team.",
		},
		Command{
			name:        "!matches",
			commandType: 0,
			syntax:      "!matches",
			description: "This command retrieves ongoing and upcoming matches.",
		},
		Command{
			name:        "!results",
			commandType: 0,
			syntax:      "!results",
			description: "This command retrieves the most recent matches results.",
		},
		Command{
			name:        "!commands",
			commandType: 0,
			syntax:      "!commands",
			description: "This command shows the available commands with their syntax and description.",
		},
		Command{
			name:        "populateteams",
			commandType: 1,
			syntax:      "populateteams <year> <month> <day>",
			description: "This command populate the bot with the top 30 teams on the hltv ranking with the publish date of <year>/<month>/<day>\nExample: populateteams 2020 november 30",
		},
		Command{
			name:        "commands",
			commandType: 1,
			syntax:      "commands",
			description: "This command shows the available bot commands with their syntax and description.",
		},
		Command{
			name:        "updatecommands",
			commandType: 1,
			syntax:      "updatecommands",
			description: "This command update the current commands",
		},
		Command{
			name:        "help",
			commandType: 1,
			syntax:      "help",
			description: "This command shows the available cli commands with their syntax and description.",
		},
		Command{
			name:        "logs",
			commandType: 1,
			syntax:      "logs <type> <limit>",
			description: "This command shows the bot logs, <type> and <limit> are optional.\nThe available types are: a: all, i: info, w: warnings: e: errors\nExamples:\n> logs>\n> logs w\n> logs i 5",
		},
		Command{
			name:        "version",
			commandType: 1,
			syntax:      "version",
			description: "This command shows the current bot version",
		},
		Command{
			name:        "exit",
			commandType: 1,
			syntax:      "exit",
			description: "This command exits closing the bot",
		},
	)
	for _, command := range commands {
		err := db.createCommand(command)
		if err != nil {
			return err
		}
	}
	return nil
}
