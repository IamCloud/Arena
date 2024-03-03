package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

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
