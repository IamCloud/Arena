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

func insertCharacter(req CreateCharacterRequest) (int64, error) {
	classInfo, err := getClassInfo(req.ClassId)
	if err != nil {
		return 0, fmt.Errorf("error retrieving class info: %v", err)
	}

	classInfo.Health += randRange(1, 6)
	classInfo.Initiative += randRange(1, 4)
	classInfo.Damage += randRange(1, 2)
	classInfo.Defense += 10

	result, err := db.Exec("INSERT INTO characters (name, wins, class_id, health, initiative, damage, defense, resource, resource_max, lives) VALUES (?, 0, ?, ?, ?, ?, ?, 0, ?, 3)",
		req.CharacterName, req.ClassId, classInfo.Health, classInfo.Initiative, classInfo.Damage, classInfo.Defense, classInfo.ResourceMax)
	if err != nil {
		return 0, fmt.Errorf("error inserting character: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error retrieving last inserted ID: %v", err)
	}

	_, err = db.Exec("INSERT INTO players_characters (player_id, character_id) VALUES (?, ?)", req.PlayerId, id)
	if err != nil {
		return 0, fmt.Errorf("error player character relationship: %v", err)
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
	if err := createCharacterTable(); err != nil {
		return err
	}
	if err := createPlayersTable(); err != nil {
		return err
	}
	if err := createPlayerCharRelTable(); err != nil {
		return err
	}
	return nil
}

func createCharacterTable() error {
	// Create users table
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS characters (
			character_id INTEGER PRIMARY KEY,
        	name TEXT,			
			wins INTEGER,
			health INTEGER,
			initiative INTEGER,
			damage INTEGER,
			defense INTEGER,
			resource INTEGER,
			resource_max INTEGER,			
			lives INTEGER,
			class_id INTEGER,
			FOREIGN KEY(class_id) REFERENCES Classes(ID)		
        )
    `)
	if err != nil {
		return fmt.Errorf("error creating character table: %v", err)
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

func createPlayerCharRelTable() error {
	// Create users table
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS players_characters (
        	player_id INTEGER NOT NULL,
			character_id INTEGER NOT NULL,
			PRIMARY KEY (player_id, character_id),
			FOREIGN KEY (player_id) REFERENCES players(player_id),
			FOREIGN KEY (character_id) REFERENCES characters(character_id)
        )
    `)
	if err != nil {
		return fmt.Errorf("error creating players_characters reliationship table: %v", err)
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

func getClassInfo(classId string) (ClassInfo, error) {
	var info ClassInfo
	err := db.QueryRow("SELECT Health, Initiative, Damage, Defense, Resource_max FROM Classes WHERE ID = ?", classId).Scan(&info.Health, &info.Initiative, &info.Damage, &info.Defense, &info.ResourceMax)
	if err != nil {
		return ClassInfo{}, err
	}
	return info, nil
}

func getCharacterInfo(characterId string) (CharacterInfo, error) {
	var info CharacterInfo
	err := db.QueryRow(`SELECT name, wins, health, initiative, damage, defense, resource, resource_max, lives, class_id
						FROM Characters
						WHERE character_id = ?`, characterId).Scan(&info.Name, &info.Wins, &info.Health, &info.Initiative, &info.Damage, &info.Defense, &info.Resource, &info.Resource_Max, &info.Lives, &info.ClassId)
	if err != nil {
		return CharacterInfo{}, err
	}
	return info, nil
}
