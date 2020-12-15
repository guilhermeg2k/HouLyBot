package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Cli struct {
	bot *Bot
}

func newCli(bot *Bot) (Cli, error) {
	var cli Cli
	if bot == nil {
		return cli, errors.New("Bot can't be nil")
	}
	cli.bot = bot
	return cli, nil
}

func (cli *Cli) handleInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		command, err := reader.ReadString('\n')
		command = strings.Replace(command, "\n", "", -1)
		if err != nil {
			logger.Error(err.Error())
		}
		commandArgs := strings.Split(command, " ")
		switch commandArgs[0] {
		case "populateteams":
			err = cli.populateTeamsWithTop30(commandArgs[1], commandArgs[2], commandArgs[3])
			if err != nil {
				logger.Error(err.Error())
			}
		case "commands":
			err = cli.showCommands()
			if err != nil {
				logger.Error(err.Error())
			}
		case "updatecommands":
			err = cli.bot.db.updateCommands()
			if err != nil {
				logger.Error(err.Error())
			}
		case "logs":
			err = cli.showLogs(commandArgs[1:])
			if err != nil {
				logger.Error(err.Error())
			}
		case "help":
			err = cli.showCliCommands()
			if err != nil {
				logger.Error(err.Error())
			}
		case "version":
			fmt.Printf("Version: %s ᕙ(⇀‸↼‶)ᕗ\n", VERSION)
		case "exit":
			logger.Info("Exiting by user request")
			os.Exit(1)
		}
	}
}

func (cli *Cli) populateTeamsWithTop30(year, month, day string) error {
	top30Teams, err := cli.bot.getTop30Teams(year, month, day)
	if err != nil {
		return err
	}
	for _, team := range top30Teams {
		err := cli.bot.db.createTeam(team)
		if err != nil {
			logger.Error("Failed when trying to populate the team " + team.name + " with the url " + team.url + " error: " + err.Error())
		}
	}
	return nil
}

func (cli *Cli) showCommands() error {
	commands, err := cli.bot.db.getBotCommands()
	if err != nil {
		return err
	}
	for i, command := range commands {
		fmt.Printf("%d %s Syntax: %s Description: %s\n", i, command.name, command.syntax, command.description)
	}
	return nil
}

func (cli *Cli) showCliCommands() error {
	commands, err := cli.bot.db.getCliCommands()
	if err != nil {
		return err
	}
	for _, command := range commands {
		fmt.Printf("%s - Syntax: %s\nDescription: %s\n\n", command.name, command.syntax, strings.ReplaceAll(command.description, "\\n", "\n"))
	}
	return nil
}

func (cli *Cli) showLogs(args []string) error {
	var logs []Log
	if len(args) > 0 {
		var logType string
		switch args[0] {
		case "a":
			logType = "-1"
		case "i":
			logType = "0"
		case "w":
			logType = "1"
		case "e":
			logType = "2"
		}
		if len(args) > 1 {
			if strings.EqualFold(logType, "-1") {
				_logs, err := cli.bot.db.getLogsWithLimit(args[1])
				if err != nil {
					return err
				}
				logs = _logs
			} else {
				_logs, err := cli.bot.db.getLogsByTypeWithLimit(logType, args[1])
				if err != nil {
					return err
				}
				logs = _logs
			}
		} else {
			if strings.EqualFold(logType, "-1") {
				_logs, err := cli.bot.db.getAllLogs()
				if err != nil {
					return err
				}
				logs = _logs
			} else {
				_logs, err := cli.bot.db.getLogsByType(logType)
				if err != nil {
					return err
				}
				logs = _logs
			}
		}
	} else {
		_logs, err := cli.bot.db.getAllLogs()
		if err != nil {
			return err
		}
		logs = _logs
	}

	for _, log := range logs {
		fmt.Printf("[%s][%s]: %s\n", log.time, log.file, log.text)
	}
	return nil
}
