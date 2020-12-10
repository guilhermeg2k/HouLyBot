package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Team struct {
	name string
	url  string
}

type Command struct {
	name        string
	syntax      string
	description string
}

func handleInput(db *DataBase, bot *Bot) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		command, err := reader.ReadString('\n')
		command = strings.Replace(command, "\n", "", -1)
		if err != nil {
			Log.Error(err.Error())
		}
		commandArgs := strings.Split(command, " ")
		switch commandArgs[0] {
		case "populateteams":
			err = populateTeamsWithTop30(db, commandArgs[1], commandArgs[2], commandArgs[3])
			if err != nil {
				Log.Error(err.Error())
			}
		case "commands":
			commands, err := db.getAllCommands()
			if err != nil {
				Log.Error(err.Error())
			}
			for i, command := range commands {
				fmt.Printf("%d %s Syntax: %s Description: %s\n", i, command.name, command.syntax, command.description)
			}
		case "updatecommands":
			updateCurrentCommands(db)
		case "exit":
			Log.Info("Exiting by user request")
			os.Exit(1)
		case "version":
			fmt.Printf("Version: %s ᕙ(⇀‸↼‶)ᕗ\n", VERSION)
		}
	}
}

func populateTeamsWithTop30(db *DataBase, year, month, day string) error {
	top30Teams, err := getTop30Teams(year, month, day)
	if err != nil {
		return err
	}
	for _, team := range top30Teams {
		err := db.createTeam(team.name, team.url)
		if err != nil {
			Log.Error("Failed when trying to populate the team " + team.name + " with the url " + team.url + " error: " + err.Error())
		}
	}
	return nil
}

func getTop30Teams(year, month, day string) ([]Team, error) {
	var teams []Team
	body, err := getRequestBody(fmt.Sprintf("https://www.hltv.org/ranking/teams/%s/%s/%s", year, month, day))
	if err != nil {
		return teams, err
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return teams, err
	}
	doc.Find(".ranked-team.standard-box").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find(".name").Text())
		url, _ := s.Find(".moreLink").Attr("href")
		url = "https://www.hltv.org" + strings.TrimSpace(url)
		teams = append(teams, Team{
			name: name,
			url:  url,
		})
	})
	return teams, nil
}

func updateCurrentCommands(db *DataBase) {
	var commands []Command
	commands = append(commands,
		Command{
			name:        "!team",
			syntax:      "!team <team-name>",
			description: "This command retrieves informations like roster, hltv ranking, next matches and recent results of a team.",
		},
		Command{
			name:        "!matches",
			syntax:      "!matches",
			description: "This command retrieves ongoing and upcoming matches.",
		},
		Command{
			name:        "!results",
			syntax:      "!results",
			description: "This command retrieves the most recent matches results.",
		},
		Command{
			name:        "!help",
			syntax:      "!help",
			description: "This command shows the available commands with their syntax and description.",
		},
	)
	for _, command := range commands {
		err := db.createCommand(command)
		if err != nil {
			Log.Error("Failed when trying to update the command " + command.name + "  error: " + err.Error())
		}
	}
}
