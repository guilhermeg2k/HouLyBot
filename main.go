package main

import (
	"fmt"
)

const TIMEZONE int = -4
const VERSION string = "v0.0.7"
const GITHUB string = "github.com/guilhermeg2k/HouLyTVBot"

func main() {
	var db DataBase
	var bot Bot
	err := db.setupDB()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	setupLogger(&db)
	err = bot.setupBot(&db)
	if err != nil {
		Log.Error(err.Error())
	}
	handleInput(&db, &bot)
}
