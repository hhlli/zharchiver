package main

import (
	"database/sql"
	"fmt"
	"log"
	"zharchiver/services"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./data/db/zharchiver.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	url := "https://x.com/elonmusk/status/1585841080431321088"
	
	ans, err := services.ProcessTwitterTask(db, url, "TEST")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Success! ID: %s, Title: %s\n", ans.AnswerID, ans.Title)
}
