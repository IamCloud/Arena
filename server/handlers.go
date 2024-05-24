package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		handleError(w, "error getting leaderboard top 10", err)
		return
	}
	defer rows.Close()

	var leaderboardRecords []Leaderboard
	for rows.Next() {
		var record Leaderboard
		if err := rows.Scan(&record.PlayerName, &record.CharacterName, &record.Wins); err != nil {
			handleError(w, "error filling leaderboard rows", err)
			return
		}
		leaderboardRecords = append(leaderboardRecords, record)
	}

	if err := rows.Err(); err != nil {
		handleError(w, "row iteration error", err)
		return
	}

	writeJSON(w, http.StatusOK, leaderboardRecords)
}

func getNewUpgrades(w http.ResponseWriter, r *http.Request) {
	var upgrades []interface{}
	for i := 0; i < 3; i++ {
		item, err := getRandomItem()
		if err != nil {
			handleError(w, "error getting random item", err)
			return
		}
		upgrades = append(upgrades, item)
	}
	writeJSON(w, http.StatusOK, upgrades)
}

func createPlayer(w http.ResponseWriter, r *http.Request) {
	var req CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, "invalid JSON", err)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	playerGuid := uuid.New().String()
	if err := insertPlayer(req, playerGuid); err != nil {
		handleError(w, "error inserting player", err)
		return
	}

	response := struct {
		PlayerGuid string `json:"player_guid"`
	}{PlayerGuid: playerGuid}
	writeJSON(w, http.StatusOK, response)
}

func chooseUpgrade(w http.ResponseWriter, r *http.Request) {
	var req UpgradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, "invalid JSON", err)
		return
	}

	// TODO: Add verification that the number of upgrades fits the number of upgrades to avoid exploits.
	AddUpgradeToCharacter(req.UpgradeId, req.CharacterId)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Upgrade added successfully"))
}

func createCharacter(w http.ResponseWriter, r *http.Request) {
	var req CreateCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, "invalid JSON", err)
		return
	}

	if req.CharacterName == "" {
		http.Error(w, "Character name cannot be empty", http.StatusBadRequest)
		return
	}

	req.Wins = "0"
	characterId, err := insertCharacter(req)
	if err != nil {
		handleError(w, "error inserting character", err)
		return
	}

	response := struct {
		CharacterID int64 `json:"character_id"`
	}{CharacterID: characterId}
	writeJSON(w, http.StatusOK, response)
}

func simulateFight(w http.ResponseWriter, r *http.Request) {
	var req SimulateFightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, "invalid JSON", err)
		return
	}

	if req.CharacterId == "" {
		http.Error(w, "CharacterId cannot be empty", http.StatusBadRequest)
		return
	}

	playerChar, err := getCharacterInfo(req.CharacterId)
	if err != nil {
		handleError(w, "error getting character info", err)
		return
	}

	opponentChar, err := getOpponentInfo(req.CharacterId, playerChar.Wins)
	if err != nil {
		handleError(w, "error finding opponent info", err)
		return
	}

	fightResult := simulateFightLogic(&playerChar, &opponentChar)
	writeJSON(w, http.StatusOK, fightResult)
}

func handleError(w http.ResponseWriter, message string, err error) {
	fmt.Printf("%s: %v\n", message, err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
				fmt.Println("Player wins")
			} else {
				playerChar.Lives -= 1
				playerChar.decrLives()
				fmt.Println("Player loses")

				if playerChar.Lives < 1 {
					fightEvents = append(fightEvents, playerCharDeadEvent())
				}
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

func playerCharDeadEvent() FightEvent {
	return FightEvent{Type: EV_TP_DEAD, Data: "{}"}
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
