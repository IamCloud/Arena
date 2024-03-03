package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

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
