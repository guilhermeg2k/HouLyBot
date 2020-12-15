package main

import (
	"fmt"
	"log"
)

const TIMEZONE int = -4
const VERSION string = "v0.0.8"
const GITHUB string = "github.com/guilhermeg2k/HouLyTVBot"

var logger *Logger

func main() {

	bot, err := newBot()
	if err != nil {
		log.Fatalln(err.Error())
	}

	setupLogger(bot.db)

	cli, err := newCli(&bot)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Printf("\\ (•◡•) / HoulyTV Bot - <%s>\n", VERSION)
	cli.handleInput()
}
