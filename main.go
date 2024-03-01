package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/google/uuid"
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
		return 0, fmt.Errorf("error retrieving last inserted ID: %v", err)
	}

	_, err = db.Exec("INSERT INTO players_teams (player_id, team_id) VALUES (?, ?)", req.PlayerId, id)
	if err != nil {
		return 0, fmt.Errorf("error player team relationship: %v", err)
	}

	return id, nil
}

func insertPlayer(req CreatePlayerRequest, guid string) error {
	_, err := db.Exec("INSERT INTO players (player_id, name) VALUES (?, ?)", guid, req.Name)
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

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Class struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	SpecialAbilities string `json:"special_abilities"`
}

func getNewUpgrades(w http.ResponseWriter, r *http.Request) {
	var upgrades []interface{}

	for i := 0; i < 3; i++ {
		if rand.Intn(2) == 0 {
			item, err := getRandomItem()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			upgrades = append(upgrades, item)
		} else {
			class, err := getRandomClass()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			upgrades = append(upgrades, class)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(upgrades)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getRandomItem() (Item, error) {
	var item Item
	err := db.QueryRow("SELECT ID, Name, Description FROM Items ORDER BY RANDOM() LIMIT 1").Scan(&item.ID, &item.Name, &item.Description)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func getRandomClass() (Class, error) {
	var class Class
	err := db.QueryRow("SELECT ID, Name, SpecialAbilities FROM Classes ORDER BY RANDOM() LIMIT 1").Scan(&class.ID, &class.Name, &class.SpecialAbilities)
	if err != nil {
		return Class{}, err
	}
	return class, nil
}

type Leaderboard struct {
	Name string
	Wins int
}

type CreatePlayerRequest struct {
	Name string
}

func createPlayer(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body into a struct
	var req CreatePlayerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	playerGuid := uuid.New().String()
	err = insertPlayer(req, playerGuid)
	if err != nil {
		fmt.Printf("error inserting player: %v\n", err)
		return
	}

	response := struct {
		PlayerGuid string `json:"player_guid"`
	}{
		PlayerGuid: playerGuid,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	fmt.Printf("Player %s created.\n", req.Name)
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

	if req.TeamName == "" {
		http.Error(w, "Team name cannot be empty", http.StatusBadRequest)
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
