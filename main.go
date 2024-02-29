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
	"strconv"

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
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	// Construct the path to the database file in the project folder
	databasePath := filepath.Join(cwd, "database.db")

	// If the database file does not exist, create it
	if _, err := os.Stat(databasePath); os.IsNotExist(err) {
		file, err := os.Create(databasePath)
		if err != nil {
			log.Fatalf("Error creating database file: %v", err)
		}
		file.Close()
	}

	// Open the SQLite database
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
	if err := createPlayerTeamsRelTable(); err != nil {
		return err
	}
	return nil
}

func createTeamsTable() error {
	// Create users table
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS teams (
			team_id INTEGER PRIMARY KEY,
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
        	player_id VARCHAR PRIMARY KEY,
            name TEXT
        )
    `)
	if err != nil {
		return fmt.Errorf("error creating players table: %v", err)
	}
	return nil
}

func createPlayerTeamsRelTable() error {
	// Create users table
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS players_teams (
        	player_id INTEGER NOT NULL,
			team_id INTEGER NOT NULL,
			PRIMARY KEY (player_id, team_id),
			FOREIGN KEY (player_id) REFERENCES players(player_id),
			FOREIGN KEY (team_id) REFERENCES teams(team_id)
        )
    `)
	if err != nil {
		return fmt.Errorf("error creating players_teams reliationship table: %v", err)
	}
	return nil
}

func insertTeam(req CreateTeamRequest) (int64, error) {
	result, err := db.Exec("INSERT INTO teams (name, wins) VALUES (?, 0)", req.TeamName)
	if err != nil {
		return 0, fmt.Errorf("error inserting team: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error retrieving last insert ID: %v", err)
	}

	_, err = db.Exec("INSERT INTO players_teams (player_id, team_id) VALUES (?, ?)", req.PlayerId, id)
	if err != nil {
		return 0, fmt.Errorf("error player team relationship: %v", err)
	}

	return id, nil
}

func insertPlayer(req CreatePlayerRequest) error {
	_, err := db.Exec("INSERT INTO players (player_id, name) VALUES (?, ?)", req.PlayerId, req.PlayerName)
	if err != nil {
		return fmt.Errorf("error inserting player: %v", err)
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

type CreatePlayerRequest struct {
	PlayerId   string
	PlayerName string
}

func initPlayer(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body into a struct
	var req CreatePlayerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	insertPlayer(req)

	// Return a success response
	w.WriteHeader(http.StatusOK)
	fmt.Printf("Player %s created successfully\n", req.PlayerName)
}

type CreateTeamRequest struct {
	PlayerId string
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

	teamId, err := insertTeam(req)
	if err != nil {
		fmt.Printf("error inserting team: %v\n", err)
		return
	}

	response := struct {
		TeamID int64 `json:"team_id"`
	}{
		TeamID: teamId,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	fmt.Printf("Team %s[%s] created successfully for player of id %s\n", req.TeamName, strconv.FormatInt(teamId, 10), req.PlayerId)
}
