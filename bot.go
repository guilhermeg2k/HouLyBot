package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

type TeamInfo struct {
	name    string
	country string
	roster  []string
	ranking string
}

type Match struct {
	firstTeam  string
	secondTeam string
	score      []string
	date       string
}

func setupDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return database, err
	}
	createTeamsTable, err := database.Prepare("CREATE TABLE IF NOT EXISTS teams (id INTEGER PRIMARY KEY, name TEXT, url TEXT)")
	if err != nil {
		return database, err
	}
	createTeamsTable.Exec()
	return database, nil
}
func main() {
	bot, err := discordgo.New("Bot " + os.Getenv("HOULY_TOKEN"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = bot.Open()
	if err != nil {
		log.Fatalf(err.Error())
	}
	bot.AddHandler(messageCreate)
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	var s string
	fmt.Scanf("%s", &s)
	bot.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	command := strings.Split(m.Content, " ")
	switch command[0] {
	case "!team":
		displayTeam, _ := displayTeam(command[1])
		s.ChannelMessageSend(m.ChannelID, displayTeam)
	}
}

func displayTeam(teamName string) (string, error) {
	var display string
	db, _ := setupDB()
	url, err := getTeamURL(db, teamName)
	if err != nil {
		return "", err
	}
	teamInfo, err := getTeamInfo(url)
	teamMatches, err := getTeamMatches(url)
	display = fmt.Sprintf("**%s %s**\n%s [ ", teamInfo.name, teamInfo.ranking, teamInfo.country)
	for _, player := range teamInfo.roster {
		display += player + " "
	}
	display += " ]\n"
	display += "**Next Matches**\n```"
	if teamMatches[0].score[0] != "-" {
		display += "No upcoming matches, check back later."
	} else {
		i := 0
		for {
			display += fmt.Sprintf("[%s] %s x %s\n", teamMatches[i].date, teamMatches[i].firstTeam, teamMatches[i].secondTeam)
			i++
			if i > len(teamMatches)-1 || teamMatches[i].score[0] != "-" {
				break
			}
		}
		if i < len(teamMatches)-1 {
			display += "```**Recent results**\n```"
			for i < len(teamMatches)-1 {
				display += fmt.Sprintf("[%s] %s %s:%s %s\n",
					teamMatches[i].date,
					teamMatches[i].firstTeam,
					teamMatches[i].score[0],
					teamMatches[i].score[1],
					teamMatches[i].secondTeam,
				)
				i++
			}
		}
	}
	return display + "```", nil
}

func getTeamURL(database *sql.DB, teamName string) (string, error) {
	var url string
	err := database.QueryRow("SELECT url FROM teams where name = ?", teamName).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func getTeamInfo(url string) (TeamInfo, error) {
	var teamInfo TeamInfo
	var roster []string
	body, err := getRequestBody(url)
	if err != nil {
		return teamInfo, err
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return teamInfo, err
	}

	teamInfo.name = strings.TrimSpace(doc.Find(".profile-team-name").Text())
	teamInfo.country = strings.TrimSpace(doc.Find(".team-country").Text())
	teamInfo.ranking = strings.TrimSpace(doc.Find(".ranking-info").First().Find(".wrap").Find(".value").Text())
	doc.Find(".bodyshot-team-bg").Find(".playerFlagName").Each(func(i int, s *goquery.Selection) {
		roster = append(roster, strings.TrimSpace(s.Text()))
	})
	teamInfo.roster = roster
	return teamInfo, nil
}

func getTeamMatches(url string) ([]Match, error) {
	var matches []Match

	body, err := getRequestBody(url + "#tab-matchesBox")
	if err != nil {
		return matches, err
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return matches, err
	}

	doc.Find(".team-row").Each(func(i int, s *goquery.Selection) {
		var score []string
		team1 := strings.TrimSpace(s.Find(".team-1").Text())
		team2 := strings.TrimSpace(s.Find(".team-2").Text())
		date := strings.TrimSpace(s.Find(".date-cell").Text())
		s.Find(".score").Each(func(i int, s *goquery.Selection) {
			score = append(score, strings.TrimSpace(s.Text()))
		})
		matches = append(matches, Match{
			firstTeam:  team1,
			secondTeam: team2,
			date:       date,
			score:      score,
		})
	})
	return matches, nil
}

func getRequestBody(url string) (io.ReadCloser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+"#tab-matchesBox", nil)
	if err != nil {
		return req.Body, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 Gecko/20100101 Firefox/84.0")
	res, err := client.Do(req)
	if err != nil {
		return res.Body, err
	}
	return res.Body, nil
}
