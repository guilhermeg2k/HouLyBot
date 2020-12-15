package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

type Bot struct {
	teams          []Team
	db             *DataBase
	discordSession *discordgo.Session
}
type TeamInfo struct {
	name    string
	country string
	roster  []string
	ranking string
}

type Match struct {
	firstTeam  MatchTeam
	secondTeam MatchTeam
	date       string
}

type MatchTeam struct {
	name  string
	score string
}

func newBot() (Bot, error) {
	var bot Bot

	db, err := newDataBase()
	if err != nil {
		return bot, err
	}
	bot.db = &db

	bot.teams, err = db.getAllTeams()
	if err != nil {
		return bot, err
	}

	discordSession, err := discordgo.New("Bot " + os.Getenv("HOULY_TOKEN"))
	if err != nil {
		logger.FatalError(err.Error())
	}
	bot.discordSession = discordSession

	go bot.startBot()
	return bot, nil
}

func (bot *Bot) startBot() {
	err := bot.discordSession.Open()
	if err != nil {
		logger.FatalError(err.Error())
	}
	bot.discordSession.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		bot.onMessageCreate(m)
	})
	logger.Info("Bot successfully started")
}

func (bot *Bot) onMessageCreate(m *discordgo.MessageCreate) {
	if m.Author.ID == bot.discordSession.State.User.ID {
		return
	}
	command := strings.Split(m.Content, " ")
	switch command[0] {
	case "!team":
		teamName := strings.Join(command[1:], " ")
		teamText, err := bot.teamText(teamName)
		if err != nil {
			logger.Error(err.Error())
		}
		err = bot.sendMessageToChannel(m.ChannelID, teamText)
		if err != nil {
			logger.Error(err.Error())
		}
	case "!matches":
		todayMatchesText, err := bot.todayMatchesText()
		if err != nil {
			logger.Error(err.Error())
		}
		err = bot.sendMessageToChannel(m.ChannelID, todayMatchesText)
		if err != nil {
			logger.Error(err.Error())
		}
	case "!results":
		recentResultsText, err := bot.recentResultsText()
		if len(recentResultsText) > 2000 {
			fmt.Println(recentResultsText)
		}
		if err != nil {
			logger.Error(err.Error())
		}
		err = bot.sendMessageToChannel(m.ChannelID, recentResultsText)
		if err != nil {
			logger.Error(err.Error())
		}
	case "!commands":
		commandsText, err := bot.commandsText()
		if err != nil {
			logger.Error(err.Error())
		}
		err = bot.sendMessageToChannel(m.ChannelID, commandsText)
		if err != nil {
			logger.Error(err.Error())
		}
	case "!about":
		err := bot.sendMessageToChannel(m.ChannelID, bot.aboutText())
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func (bot *Bot) recentResultsText() (string, error) {
	var resultsDisplay string
	matches, err := bot.getRecentResults()
	if err != nil {
		return "", err
	}
	resultsDisplay = "**Recent Results**\n```"
	for _, match := range matches {
		resultsDisplay += fmt.Sprintf(
			"%s [%s] x [%s] %s\n",
			match.firstTeam.name,
			match.firstTeam.score,
			match.secondTeam.score,
			match.secondTeam.name,
		)
	}
	resultsDisplay += "```"
	return resultsDisplay, nil
}

func (bot *Bot) getRecentResults() ([]Match, error) {
	var matches []Match
	body, err := getRequestBody("https://www.hltv.org")
	if err != nil {
		return matches, err
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return matches, err
	}
	doc.Find(".rightCol .result-box").Each(func(i int, s *goquery.Selection) {
		var match Match
		s.Find(".team").EachWithBreak(func(i int, ss *goquery.Selection) bool {
			if i == 0 {
				match.firstTeam.name = ss.Text()
			}
			if i == 1 {
				match.secondTeam.name = ss.Text()
				return false
			}
			return true
		})
		s.Find(".twoRowExtraRow").EachWithBreak(func(i int, ss *goquery.Selection) bool {
			if i == 0 {
				match.firstTeam.score = ss.Text()
			}
			if i == 1 {
				match.secondTeam.score = ss.Text()
				return false
			}
			return true
		})
		matches = append(matches, match)
	})
	return matches, nil
}

func (bot *Bot) todayMatchesText() (string, error) {
	var matchesDisplay string
	matches, err := bot.getTodayMatches()
	if err != nil {
		return "", err
	}
	matchesDisplay = "**Today's Matches**\n```"
	for _, match := range matches {
		matchesDisplay += fmt.Sprintf("%s %s x %s\n", match.date, match.firstTeam.name, match.secondTeam.name)
	}
	matchesDisplay += "```"
	return matchesDisplay, nil
}

func (bot *Bot) getTodayMatches() ([]Match, error) {
	var matches []Match
	body, err := getRequestBody("https://www.hltv.org/")
	if err != nil {
		return matches, err
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return matches, err
	}
	doc.Find(".rightCol aside .hotmatch-box").Each(func(i int, s *goquery.Selection) {
		var match Match
		s.Find(".team").EachWithBreak(func(i int, ss *goquery.Selection) bool {
			if i == 0 {
				match.firstTeam.name = ss.Text()
			}
			if i == 1 {
				match.secondTeam.name = ss.Text()
				return false
			}
			return true
		})
		matchTime := s.Find(".middleExtra").Text()
		if matchTime != "" {
			match.date = convertTimeZone(matchTime)
		} else {
			match.date = "LIVE"
		}
		if match.firstTeam.name == "" {
			return
		}
		matches = append(matches, match)
	})
	return matches, nil
}

func (bot *Bot) teamText(teamName string) (string, error) {
	var display string
	url := bot.getTeamUrl(teamName)
	if url == "" {
		return "", errors.New("Team url not founded")
	}
	teamInfo, err := bot.getTeamInfo(url)
	if err != nil {
		return "", err
	}
	teamMatches, err := bot.getTeamMatches(url)
	if err != nil {
		return "", err
	}
	display = fmt.Sprintf("**%s %s**\n%s [ ", teamInfo.name, teamInfo.ranking, teamInfo.country)
	for _, player := range teamInfo.roster {
		display += player + " "
	}
	display += " ]\n"
	display += "**Next Matches**\n```"
	if teamMatches[0].firstTeam.score != "-" {
		display += "No upcoming matches, check back later."
	}
	i := 0
	for {
		if i > len(teamMatches)-1 || teamMatches[i].firstTeam.score != "-" {
			break
		}
		display += fmt.Sprintf("[%s] %s x %s\n", teamMatches[i].date, teamMatches[i].firstTeam.name, teamMatches[i].secondTeam.name)
		i++
	}
	//TODO: Handle when you don't have any recent match, which is rare
	display += "```**Recent results**\n```"
	for i < len(teamMatches)-1 {
		display += fmt.Sprintf("[%s] %s %s:%s %s\n",
			teamMatches[i].date,
			teamMatches[i].firstTeam.name,
			teamMatches[i].firstTeam.score,
			teamMatches[i].secondTeam.score,
			teamMatches[i].secondTeam.name,
		)
		i++
	}
	return display + "```", nil
}

func (bot *Bot) getTeamInfo(url string) (TeamInfo, error) {
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

func (bot *Bot) getTeamMatches(url string) ([]Match, error) {
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
		var match Match
		match.firstTeam.name = strings.TrimSpace(s.Find(".team-1").Text())
		match.secondTeam.name = strings.TrimSpace(s.Find(".team-2").Text())
		match.date = strings.TrimSpace(s.Find(".date-cell").Text())
		s.Find(".score").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if i == 0 {
				match.firstTeam.score = s.Text()
			}
			if i == 1 {
				match.secondTeam.score = s.Text()
				return false
			}
			return true
		})
		if strings.Contains(match.date, ":") {
			match.date = convertTimeZone(match.date)
		}
		matches = append(matches, match)
	})
	return matches, nil
}

func (bot *Bot) getTeamUrl(teamName string) string {
	for _, team := range bot.teams {
		if strings.EqualFold(teamName, team.name) {
			return team.url
		}
	}
	return ""
}

func (bot *Bot) commandsText() (string, error) {
	commands, err := bot.db.getBotCommands()
	if err != nil {
		logger.Error(err.Error())
	}
	commandsDisplay := "Available commands \n"
	for _, command := range commands {
		commandsDisplay += fmt.Sprintf("```%s Syntax: %s Description: %s```", command.name, command.syntax, command.description)
	}
	return commandsDisplay, nil
}

func (bot *Bot) aboutText() string {
	about := fmt.Sprintf(
		"```HoulyTVBot version: %s ᕙ(⇀‸↼‶)ᕗ\nCode available on github: %s ```",
		VERSION,
		GITHUB,
	)
	return about
}

func (bot *Bot) sendMessageToChannel(channelId string, content string) error {
	if len(content) > 2000 {
		return errors.New("Content string must be 2000 or fewer in length.")
	}
	_, err := bot.discordSession.ChannelMessageSend(channelId, content)
	if err != nil {
		return err
	}
	return nil
}

func (bot *Bot) getTop30Teams(year, month, day string) ([]Team, error) {
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

func getRequestBody(url string) (io.ReadCloser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 Gecko/20100101 Firefox/84.0")
	res, err := client.Do(req)
	if err != nil {
		return res.Body, err
	}
	return res.Body, nil
}

func convertTimeZone(time string) string {
	hourAndMinutes := strings.Split(time, ":")
	hour, err := strconv.Atoi(hourAndMinutes[0])
	if err != nil {
		logger.Error("Failed to parse hour to int")
		return ""
	}
	if TIMEZONE+hour == 0 {
		return fmt.Sprintf("00:%s", hourAndMinutes[1])
	}
	if TIMEZONE+hour < 0 {
		return fmt.Sprintf("%d:%s", TIMEZONE+hour+24, hourAndMinutes[1])
	}
	return fmt.Sprintf("%d:%s", TIMEZONE+hour, hourAndMinutes[1])
}
