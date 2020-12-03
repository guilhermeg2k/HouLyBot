package main

import (
	"fmt"
)

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
