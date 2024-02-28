package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	fmt.Println("Server started")

	fmt.Println("Database initializing...")
	initDb()
	fmt.Println("Database initialised")

	http.HandleFunc("/getleaderboard", getLeaderboard)
	http.HandleFunc("/initplayer", initPlayer)
	http.HandleFunc("/createteam", createTeam)

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// Serve static files from the directory containing the Go file
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	fmt.Println("Server ready to go.")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func initDb() {
	databasePath := "/database.db"
	if _, err := os.Stat(databasePath); os.IsNotExist(err) {
		// If the database file does not exist, create it
		file, err := os.Create(databasePath)
		if err != nil {
			log.Fatalf("Error creating database file: %v", err)
		}
		file.Close()
	}
	// Open the SQLite database
	var err error
	db, err = sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Create required tables
	err = createTables()
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
}

// Function to create required tables
func createTables() error {
	if err := createTeamsTable(); err != nil {
		return err
	}
	if err := createPlayersTable(); err != nil {
		return err
	}
	return nil
}

func createTeamsTable() error {
	// Create users table
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS teams (
			teamsId INTEGER PRIMARY KEY,
        	name TEXT,
			wins INTEGER
        )
    `)
	if err != nil {
		return fmt.Errorf("error creating teams table: %v", err)
	}

	// Create other tables as needed

	return nil
}

func createPlayersTable() error {
	// Create users table
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS players (
        	guid VARCHAR PRIMARY KEY,
            name TEXT,
			teamsId INTEGER,
			FOREIGN KEY(teamsId) REFERENCES teams(teamsId)
        )
    `)
	if err != nil {
		return fmt.Errorf("error creating players table: %v", err)
	}

	// Create other tables as needed

	return nil
}

func insertTeam(name string) (int64, error) {
	result, err := db.Exec("INSERT INTO teams (name, wins) VALUES (?, 0)", name)
	if err != nil {
		return 0, fmt.Errorf("error inserting team: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error retrieving last insert ID: %v", err)
	}

	return id, nil
}

func insertPlayer(guid string, name string, teamsId int64) error {
	_, err := db.Exec("INSERT INTO players (guid, name, teamsId) VALUES (?, ?, ?)", guid, name, teamsId)
	if err != nil {
		return fmt.Errorf("error inserting team: %v", err)
	}

	return nil
}

func getLeaderboard(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
        SELECT name, wins FROM teams ORDER BY wins DESC LIMIT 5
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Initialize a slice to store the top 5 teams
	var leaderboardRecords []Leaderboard

	// Iterate through the rows and append each team to the slice
	for rows.Next() {
		var leaderboard Leaderboard
		if err := rows.Scan(&leaderboard.Name, &leaderboard.Wins); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		leaderboardRecords = append(leaderboardRecords, leaderboard)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the top 5 teams slice into JSON and send as response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(leaderboardRecords); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Leaderboard struct {
	Name string
	Wins int
}

type InitPlayerRequest struct {
	Guid       string
	TeamName   string
	PlayerName string
}

func initPlayer(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body into a struct
	var req InitPlayerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTeamId, err := insertTeam(req.TeamName)
	if err != nil {
		fmt.Fprintf(w, "error inserting team: %v", err)
		return
	}
	insertPlayer(req.Guid, req.TeamName, newTeamId)

	// Return a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Team %s with player %s created successfully", req.TeamName, req.PlayerName)
}

type CreateTeamRequest struct {
	Guid     string
	TeamName string
}

func createTeam(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body into a struct
	var req CreateTeamRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	teamId, err := insertTeam(req.TeamName)
	if err != nil {
		fmt.Fprintf(w, "error inserting team: %v", err)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Team %s[%s] created successfully", req.TeamName, teamId)
}
