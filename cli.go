package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Team struct {
	name string
	url  string
}

func handleInput(db *DataBase, bot *Bot) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		command, err := reader.ReadString('\n')
		command = strings.Replace(command, "\n", "", -1)
		if err != nil {
			log.Println(err.Error())
		}
		commandArgs := strings.Split(command, " ")
		switch commandArgs[0] {
		case "populateteams":
			err = populateTeamsWithTop30(db, commandArgs[1], commandArgs[2], commandArgs[3])
			if err != nil {
				log.Println(err.Error())
			}
		case "exit":
			log.Fatalln("Exiting by user request")
		}
	}
}

func populateTeamsWithTop30(db *DataBase, year, month, day string) error {
	top30Teams, err := getTop30Teams(year, month, day)
	if err != nil {
		return err
	}
	for _, team := range top30Teams {
		db.createTeam(team.name, team.url)
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
