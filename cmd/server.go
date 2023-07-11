package main

import (
	"encoding/json"
	"log"
	"mod36a41/pkg/api"
	"mod36a41/pkg/rss"
	"mod36a41/pkg/storage"
	"mod36a41/pkg/storage/postgresdb"
	"net/http"
	"os"
	"time"
)

// configuration
type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

type server struct {
	db  storage.DBinterface
	api *api.API
}

func main() {
	// reads config.json
	b, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	// decoding
	var c config
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Fatal(err)
	}

	var srv server

	pgPass := os.Getenv("pgPass") // pass
	if pgPass == "" {
		os.Exit(1)
	}
	// connects to postgresdb
	pgConn := "postgres://postgres:" + pgPass + "@localhost:5432/news"

	// initializes postgresdb
	db, err := postgresdb.New(pgConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// initializes server storage
	srv.db = db

	//registers API handlers
	srv.api = api.New(srv.db)

	// launches news parser in individual goroutines for each RSS-feed
	chNews := make(chan []storage.Post)
	chErrors := make(chan error)
	for _, url := range c.URLS {
		go parse(url, chNews, chErrors, c.Period)
	}

	// writes news stream to db
	go func() {
		for posts := range chNews {
			srv.db.StoreNews(posts)
		}
	}()

	// error chan handling
	go func() {
		for err := range chErrors {
			log.Println("RSS feed parser error:", err)
		}
	}()

	http.ListenAndServe(":80", srv.api.Router())
}

// parser
func parse(url string, posts chan<- []storage.Post, errors chan<- error, period int) {
	for {
		news, err := rss.ParseFeed(url)
		if err != nil {
			errors <- err
		}
		posts <- news
		time.Sleep(time.Duration(period) * time.Minute)
	}
}
