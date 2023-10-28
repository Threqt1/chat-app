package main

/**
TODO:
add notification system into websockt
*/

import (
	"chat-app/pkg/database"
	"chat-app/pkg/httpserver"
	"chat-app/pkg/ws"
	"chat-app/repository"
	"flag"
	"log"
	"sync"

	"github.com/joho/godotenv"
)

// init initializes the environment variables.
func init() {
	//Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to load .env", err)
	}
}

func main() {
	var db repository.DatabaseService
	db, err := database.NewRedisDatabaseProvider(database.TEST_DATABASE)
	if err != nil {
		log.Fatal("Failed to connect to Redis database", err)
		return
	}

	serverType := flag.String("server", "", "ws,http")
	flag.Parse()

	switch *serverType {
	case "http":
		var hp repository.HTTPService
		hp, err = httpserver.NewHTTPProvider(db)
		if err != nil {
			log.Fatal("Failed to start HTTP server", err)
		}
		log.Print("Starting HTTP server on :8080")
		hp.Start()
	case "ws":
		var websocket repository.WebsocketService
		websocket = ws.NewWebsocketProvider(db)
		log.Print("Starting Websocket server on :8081")
		websocket.Start()
	default:
		var wg sync.WaitGroup

		var hp repository.HTTPService
		hp, err = httpserver.NewHTTPProvider(db)
		if err != nil {
			log.Fatal("Failed to start HTTP server", err)
		}

		var websocket repository.WebsocketService
		websocket = ws.NewWebsocketProvider(db)

		wg.Add(2)

		log.Print("Starting HTTP server on :8080")
		go func() {
			defer wg.Done()
			hp.Start()
		}()
		log.Print("Starting Websocket server on :8081")
		go func() {
			defer wg.Done()
			websocket.Start()
		}()

		wg.Wait()
	}
}
