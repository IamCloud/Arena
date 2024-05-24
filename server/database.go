package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const databaseFileName = "database.db"

func initDb(needReset bool) {
	if err := openDb(); err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if needReset {
		if err := resetDb(); err != nil {
			log.Fatalf("Error resetting database: %v", err)
		}
		if err := createBots(); err != nil {
			log.Fatalf("Error creating bots: %v", err)
		}
	}
}

func openDb() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	databasePath := filepath.Join(cwd, databaseFileName)

	db, err = sql.Open("sqlite3", databasePath)
	if err != nil {
		return err
	}

	return nil
}

func deleteSqliteDatabase() error {
	if _, err := os.Stat(databaseFileName); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	return os.Remove(databaseFileName)
}

func resetDb() error {
	if err := deleteSqliteDatabase(); err != nil {
		return err
	}

	if err := openDb(); err != nil {
		return err
	}

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
		if err := createBot(1, i); err != nil {
			return err
		}
		if err := createBot(2, i); err != nil {
			return err
		}
	}
	return nil
}

func createBot(classID, wins int) error {
	req := CreateCharacterRequest{
		PlayerId:      "1",
		ClassId:       strconv.Itoa(classID),
		CharacterName: "Bot (" + strconv.Itoa(wins) + ")",
		Wins:          strconv.Itoa(wins),
	}

	if _, err := insertCharacter(req); err != nil {
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
		return 0, fmt.Errorf("error inserting player character relationship: %v", err)
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
	if err != nil {
		return Character{}, err
	}
	info.HealthMax = info.Health
	return info, nil
}

func getOpponentInfo(playerCharID string, wins int) (Character, error) {
	var opponentCharInfo Character
	err := db.QueryRow(`SELECT character_id, name, wins, health, initiative, damage, defense, resource, resource_max, lives, class_id
                        FROM characters
                        WHERE character_id != ? AND wins = ?
                        ORDER BY RANDOM()
                        LIMIT 1`, playerCharID, wins).Scan(&opponentCharInfo.CharacterId, &opponentCharInfo.Name, &opponentCharInfo.Wins, &opponentCharInfo.Health, &opponentCharInfo.Initiative, &opponentCharInfo.Damage, &opponentCharInfo.Defense, &opponentCharInfo.Resource, &opponentCharInfo.ResourceMax, &opponentCharInfo.Lives, &opponentCharInfo.ClassId)
	if err != nil {
		return Character{}, err
	}
	opponentCharInfo.HealthMax = opponentCharInfo.Health
	return opponentCharInfo, nil
}

func AddUpgradeToCharacter(upgradeId, characterId string) {
	switch upgradeId {
	case "1":
		incrStat("health", characterId, 2)
	case "2":
		incrStat("initiative", characterId, 2)
	case "3":
		incrStat("damage", characterId, 1)
	case "4":
		incrStat("defense", characterId, 1)
	case "5":
		incrStat("defense", characterId, 5)
		incrStat("damage", characterId, -3)
	}
}

func incrStat(stat, characterId string, value int) {
	_, err := db.Exec(fmt.Sprintf("UPDATE characters SET %s = %s + ? WHERE character_id = ?", stat, stat), value, characterId)
	if err != nil {
		fmt.Printf("error increasing character %s: %v\n", stat, err)
	}
}
