package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func getLeaderboard(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
        SELECT name, wins FROM characters ORDER BY wins DESC LIMIT 5
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

func getNewUpgrades(w http.ResponseWriter, r *http.Request) {
	var upgrades []interface{}

	for i := 0; i < 3; i++ {
		item, err := getRandomItem()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		upgrades = append(upgrades, item)
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(upgrades)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

func createCharacter(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body into a struct
	var req CreateCharacterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.CharacterName == "" {
		http.Error(w, "Character name cannot be empty", http.StatusBadRequest)
		return
	}

	characterId, err := insertCharacter(req)
	if err != nil {
		fmt.Printf("error inserting character: %v\n", err)
		return
	}

	response := struct {
		CharacterID int64 `json:"character_id"`
	}{
		CharacterID: characterId,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	fmt.Printf("Character %s[%s] created successfully for player of id %s\n", req.CharacterName, strconv.FormatInt(characterId, 10), req.PlayerId)
}

func simulateFight(w http.ResponseWriter, r *http.Request) {
	var req SimulateFightRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.CharacterId == "" {
		http.Error(w, "CharacterId cannot be empty", http.StatusBadRequest)
		return
	}

	playerChar, err := getCharacterInfo(req.CharacterId)
	if err != nil {
		fmt.Printf("error getting character info: %v\n", err)
		return
	}
	fmt.Printf("%+v\n", playerChar)
	fmt.Printf("Fight of character %s starting. Searching for target...\n", playerChar.Name)

	targetChar, err := getOpponentInfo(req.CharacterId, playerChar.Wins)
	if err != nil {
		fmt.Printf("error finding target info: %v\n", err)
		return
	}

	fmt.Printf("Target found: %s\n", targetChar.Name)

	// Simulate the fight on the server
	fightResult := simulateFightLogic(&playerChar, &targetChar)

	// Send all fight event data in a single response
	//fightEventData := prepareFightEventData(fightResult, charInfo, opponentInfo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fightResult)
}

func simulateFightLogic(player *Character, target *Character) FightData {
	var FightData FightData
	FightData.StartPlayerInfo = *player
	FightData.StartOpponentInfo = *target

	// Determine the attacker based on initiative
	var currentAttacker *Character
	var currentTarget *Character
	if player.Initiative > target.Initiative {
		currentAttacker = player
		currentTarget = target
	} else {
		currentAttacker = target
		currentTarget = player
	}

	for {
		// Check for win condition (either character's health <= 0)

		if player.Health <= 0 || target.Health <= 0 {
			fmt.Printf("Combat ended!\n")
			break
		}

		// Attack logic based on current attacker
		currentAttacker.Attack(&FightData.Events, currentTarget)

		// Switch for next round
		if currentAttacker == target {
			currentAttacker = player
			currentTarget = target
		} else {
			currentAttacker = target
			currentTarget = player
		}
	}
	return FightData
}
