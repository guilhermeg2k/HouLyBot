package main

import (
	"log"
)

func main() {
	var db DataBase
	var bot Bot
	err := db.setupDB()
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = bot.setupBot(&db)
	if err != nil {
		log.Fatalln(err.Error())
	}
	handleInput(&db, &bot)
}
