package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func initDb(needReset bool) {
	err := openDb()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if needReset {
		err = resetDb()
		if err != nil {
			log.Fatalf("Error resetting database: %v", err)
		}
		err = createBots()
		if err != nil {
			log.Fatalf("Error creating bots: %v", err)
		}
	}
}

func openDb() error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Construct the path to the database file in the project folder
	databasePath := filepath.Join(cwd, "database.db")

	// Open the SQLite database
	db, err = sql.Open("sqlite3", databasePath)
	if err != nil {
		return err
	}

	return nil
}

func deleteSqliteDatabase() error {
	// Check if the file exists
	if _, err := os.Stat("database.db"); os.IsNotExist(err) {
		return nil // Database file doesn't exist, nothing to delete
	} else if err != nil {
		return err // Handle other errors during stat
	}

	// Delete the file if it exists
	err := os.Remove("database.db")
	if err != nil {
		return err // Handle deletion errors
	}
	return nil
}

func resetDb() error {
	err := deleteSqliteDatabase()
	if err != nil {
		return err
	}

	err = openDb()
	if err != nil {
		return err
	}

	// Read the script file
	scriptBytes, err := os.ReadFile("DbInit.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(scriptBytes))
	if err != nil {
		return err
	}

	return nil
}

func createBots() error {
	for i := 0; i < 100; i++ {
		err := createBot(1, i)
		if err != nil {
			return err
		}
		err = createBot(2, i)
		if err != nil {
			return err
		}
	}

	return nil
}
func createBot(class_id int, wins int) error {
	var req CreateCharacterRequest
	req.PlayerId = "1"
	req.ClassId = strconv.Itoa(class_id)
	req.CharacterName = "Bot (" + strconv.Itoa(wins) + ")"
	req.Wins = strconv.Itoa(wins)
	_, err := insertCharacter(req)
	if err != nil {
		return err
	}

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

	result, err := db.Exec("INSERT INTO characters (name, wins, class_id, health, initiative, damage, defense, resource, resource_max, lives) VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?, 3)",
		req.CharacterName, req.Wins, req.ClassId, classInfo.Health, classInfo.Initiative, classInfo.Damage, classInfo.Defense, classInfo.ResourceMax)
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
	err := db.QueryRow(`SELECT character_id, name, wins, health, initiative, damage, defense, resource, resource_max, lives, class_id
						FROM characters
						WHERE character_id = ?`, characterId).Scan(&info.CharacterId, &info.Name, &info.Wins, &info.Health, &info.Initiative, &info.Damage, &info.Defense, &info.Resource, &info.ResourceMax, &info.Lives, &info.ClassId)

	info.HealthMax = info.Health
	if err != nil {
		return Character{}, err
	}
	return info, nil
}

func getOpponentInfo(playerCharID string, wins int) (Character, error) {
	var opponentCharInfo Character
	err := db.QueryRow(`SELECT character_id, name, wins, health, initiative, damage, defense, resource, resource_max, lives, class_id
						FROM characters
						WHERE character_id != ? AND wins = ?
						ORDER BY RANDOM()
						LIMIT 1`, playerCharID, wins).Scan(&opponentCharInfo.CharacterId, &opponentCharInfo.Name, &opponentCharInfo.Wins, &opponentCharInfo.Health, &opponentCharInfo.Initiative, &opponentCharInfo.Damage, &opponentCharInfo.Defense, &opponentCharInfo.Resource, &opponentCharInfo.ResourceMax, &opponentCharInfo.Lives, &opponentCharInfo.ClassId)
	opponentCharInfo.HealthMax = opponentCharInfo.Health
	if err != nil {
		return Character{}, err
	}
	return opponentCharInfo, nil
}

func AddUpgardeToCharacter(upgradeId string, characterId string) {
	switch upgradeId {
	case "1":
		incrHp(characterId, 2)
	case "2":
		incrInit(characterId, 2)
	case "3":
		incrDmg(characterId, 1)
	case "4":
		incrDef(characterId, 1)
	case "5":
		incrDef(characterId, 5)
		incrDmg(characterId, -3)
	}
}

func incrHp(characterId string, value int) {
	_, err := db.Exec("UPDATE characters SET health = health + ? WHERE character_id = ?", value, characterId)
	if err != nil {
		fmt.Println("error increasing character health.", err)
		return
	}
}
func incrInit(characterId string, value int) {
	_, err := db.Exec("UPDATE characters SET initiative = initiative + ? WHERE character_id = ?", value, characterId)
	if err != nil {
		fmt.Println("error increasing character initiative.", err)
		return
	}
}

func incrDmg(characterId string, value int) {
	_, err := db.Exec("UPDATE characters SET damage = damage + ? WHERE character_id = ?", value, characterId)
	if err != nil {
		fmt.Println("error increasing character damage.", err)
		return
	}
}

func incrDef(characterId string, value int) {
	_, err := db.Exec("UPDATE characters SET defense = defense + ? WHERE character_id = ?", value, characterId)
	if err != nil {
		fmt.Println("error increasing character defense.", err)
		return
	}
}
