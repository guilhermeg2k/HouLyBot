package main

import "log"

const TIMEZONE int = -4
const VERSION string = "v0.0.7"
const GITHUB string = "github.com/guilhermeg2k/HouLyTVBot"

var logger *Logger

func main() {
	bot, err := newBot()
	if err != nil {
		log.Fatalln(err.Error())
	}
	setupLogger(bot.db)
	handleInput(&bot)
}
