package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func initDb(needCreateBots bool) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	// Construct the path to the database file in the project folder
	databasePath := filepath.Join(cwd, "database.db")

	// Open the SQLite database
	db, err = sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if needCreateBots {
		err = createBots()
		if err != nil {
			log.Fatalf("Error creating bots: %v", err)
		}
	}
}

func createBots() error {
	var reqWarr CreateCharacterRequest
	reqWarr.PlayerId = "1"
	reqWarr.ClassId = "1"
	reqWarr.CharacterName = "Warrior bot"
	_, err := insertCharacter(reqWarr)
	if err != nil {
		return fmt.Errorf("error inserting bot character 1: %v", err)
	}
	var reqMage CreateCharacterRequest
	reqMage.PlayerId = "1"
	reqMage.ClassId = "2"
	reqMage.CharacterName = "Mage bot"
	_, err = insertCharacter(reqMage)
	if err != nil {
		return fmt.Errorf("error inserting bot character 2: %v", err)
	}
	fmt.Println("Bots created successfully.")

	return nil
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

func getRandomItem() (Item, error) {
	var item Item
	err := db.QueryRow("SELECT ID, Name, Description FROM items ORDER BY RANDOM() LIMIT 1").Scan(&item.ID, &item.Name, &item.Description)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func getClassInfo(classId string) (ClassInfo, error) {
	var info ClassInfo
	err := db.QueryRow("SELECT Health, Initiative, Damage, Defense, Resource_max FROM classes WHERE ID = ?", classId).Scan(&info.Health, &info.Initiative, &info.Damage, &info.Defense, &info.ResourceMax)
	if err != nil {
		return ClassInfo{}, err
	}
	return info, nil
}

func getCharacterInfo(characterId string) (Character, error) {
	var info Character
	err := db.QueryRow(`SELECT name, wins, health, initiative, damage, defense, resource, resource_max, lives, class_id
						FROM characters
						WHERE character_id = ?`, characterId).Scan(&info.Name, &info.Wins, &info.Health, &info.Initiative, &info.Damage, &info.Defense, &info.Resource, &info.ResourceMax, &info.Lives, &info.ClassId)

	info.HealthMax = info.Health
	if err != nil {
		return Character{}, err
	}
	return info, nil
}

func getOpponentInfo(playerCharID string, wins int) (Character, error) {
	var opponentCharInfo Character
	err := db.QueryRow(`SELECT name, wins, health, initiative, damage, defense, resource, resource_max, lives, class_id
						FROM characters
						WHERE character_id != ? AND wins = ?
						ORDER BY RANDOM()
						LIMIT 1`, playerCharID, wins).Scan(&opponentCharInfo.Name, &opponentCharInfo.Wins, &opponentCharInfo.Health, &opponentCharInfo.Initiative, &opponentCharInfo.Damage, &opponentCharInfo.Defense, &opponentCharInfo.Resource, &opponentCharInfo.ResourceMax, &opponentCharInfo.Lives, &opponentCharInfo.ClassId)
	opponentCharInfo.HealthMax = opponentCharInfo.Health
	if err != nil {
		return Character{}, err
	}
	return opponentCharInfo, nil
}
