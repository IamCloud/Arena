package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	fmt.Println("Server started")

	fmt.Println("Database initializing...")
	initDb()
	fmt.Println("Database initialised")

	http.HandleFunc("/getleaderboard", getLeaderboard)
	http.HandleFunc("/createplayer", createPlayer)
	http.HandleFunc("/createteam", createTeam)
	http.HandleFunc("/getnewupgrades", getNewUpgrades)

	rootDir := "./.." // Assuming your main folder is named Arena

	// Construct the path to the client directory
	clientDir := filepath.Join(rootDir, "Client")

	// Serve static files from the client directory
	fs := http.FileServer(http.Dir(clientDir))
	http.Handle("/", fs)

	fmt.Println("Server ready to go.")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
