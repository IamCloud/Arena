package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	fmt.Println("Server started")

	fmt.Println("Database initializing...")

	needCreateBots := flag.Bool("bots", false, "Need to create bots")
	flag.Parse()

	fmt.Println("Need to create bots: ", *needCreateBots)
	initDb(*needCreateBots)
	fmt.Println("Database initialised")

	http.HandleFunc("/getleaderboard", getLeaderboard)
	http.HandleFunc("/createplayer", createPlayer)
	http.HandleFunc("/createcharacter", createCharacter)
	http.HandleFunc("/getnewupgrades", getNewUpgrades)
	http.HandleFunc("/simulatefight", simulateFight)

	rootDir := "./.." // Assuming your main folder is named Arena

	// Construct the path to the client directory
	clientDir := filepath.Join(rootDir, "Client")

	// Serve static files from the client directory
	fs := http.FileServer(http.Dir(clientDir))
	http.Handle("/", fs)

	fmt.Println("Server ready to go.")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
