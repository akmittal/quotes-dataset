package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Quote struct {
	Quote      string
	Author     string
	Tags       []string
	Popularity float64
	Category   string
}

func main() {
	db, err := sql.Open("sqlite3", "./quoteitup.db")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("select quotes.quote,quotes.author,quotes.tags,quotes.votes,quotecat.category from quotes, quotecat where quotes._id=quotecat.quoteid")
	if err != nil {
		log.Fatal(err)
	}
	votesrow := db.QueryRow("select max(votes) from quotes")
	var maxVotes int
	votesrow.Scan(&maxVotes)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var allRows []Quote
	for rows.Next() {
		var col1, col2, col3, col4, col5 string
		err = rows.Scan(&col1, &col2, &col3, &col4, &col5)
		if err != nil {
			log.Fatal(err)
		}

		votes, _ := strconv.Atoi(col4)
		popularity := float64(votes) / float64(maxVotes)
		quote := Quote{col1, col2, strings.Split(col3, ","), popularity, col5}
		allRows = append(allRows, quote)
		fmt.Println(col1, col2, col3, col4, col5)

	}
	b, err := json.Marshal(allRows)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./quotes.json", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
