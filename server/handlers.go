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
        SELECT p.name, c.name, c.wins 
		FROM characters c
		INNER JOIN players_characters pc ON c.character_id = pc.character_id
		INNER JOIN players p ON p.player_id = pc.player_id
		WHERE p.player_id != '1'
		ORDER BY c.wins DESC 
		LIMIT 10
    `)
	if err != nil {
		fmt.Printf("error getting leaderboard top 10: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Initialize a slice to store the top 5 teams
	var leaderboardRecords []Leaderboard

	// Iterate through the rows and append each team to the slice
	for rows.Next() {
		var leaderboard Leaderboard
		if err := rows.Scan(&leaderboard.PlayerName, &leaderboard.CharacterName, &leaderboard.Wins); err != nil {
			fmt.Printf("error filling leaderboard rows: %v\n", err)
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

	req.Wins = "0"
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

	opponentChar, err := getOpponentInfo(req.CharacterId, playerChar.Wins)
	if err != nil {
		fmt.Printf("error finding opponent info: %v\n", err)
		return
	}

	fmt.Printf("Opponent found: %s\n", opponentChar.Name)

	// Simulate the fight on the server
	fightResult := simulateFightLogic(&playerChar, &opponentChar)

	// Send all fight event data in a single response
	//fightEventData := prepareFightEventData(fightResult, charInfo, opponentInfo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fightResult)
}

func simulateFightLogic(playerChar, opponentChar *Character) []FightEvent {
	var fightEvents []FightEvent
	attacker, defender := determineInitialAttacker(playerChar, opponentChar)

	fightEvents = append(fightEvents,
		charactersUpdateEvent(playerChar, opponentChar),
		initiativeEvent(attacker.Name),
	)

	for {
		if fightEnded(playerChar, opponentChar) {
			winner := getWinner(playerChar, opponentChar)
			if winner == playerChar {
				playerChar.Wins += 1
				playerChar.incrWins()
			}
			fightEvents = append(fightEvents, combatEndEvent(winner))
			break
		}

		attacker.Attack(&fightEvents, defender)
		fightEvents = append(fightEvents, charactersUpdateEvent(playerChar, opponentChar))

		attacker, defender = swapAttackerAndDefender(attacker, defender)
	}

	return fightEvents
}

// Helper functions for readability and potential reusability:

func determineInitialAttacker(player, opponent *Character) (*Character, *Character) {
	if player.Initiative > opponent.Initiative {
		return player, opponent
	}
	return opponent, player
}

func initiativeEvent(startingCharacterName string) FightEvent {
	initEvent := InitiativeEvent{StartingCharacterName: startingCharacterName}
	return createEvent(EV_TP_INIT, initEvent)
}

func fightEnded(player, opponent *Character) bool {
	return player.Health <= 0 || opponent.Health <= 0
}

func getWinner(player, opponent *Character) *Character {
	if player.Health > 0 {
		return player
	}
	return opponent
}

func combatEndEvent(winner *Character) FightEvent {
	endEvent := CombatEndEvent{Winner: *winner}
	return createEvent(EV_TP_END, endEvent)
}

func charactersUpdateEvent(player, opponent *Character) FightEvent {
	updateCharsEvent := UpdateCharactersEvent{*player, *opponent}
	return createEvent(EV_TP_UPD_CHARS, updateCharsEvent)
}

func swapAttackerAndDefender(attacker, defender *Character) (*Character, *Character) {
	return defender, attacker
}

// Generic event creation for potential reuse:

func createEvent(eventType string, data interface{}) FightEvent {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		// Handle error appropriately
	}
	return FightEvent{Type: eventType, Data: string(jsonData)}
}
