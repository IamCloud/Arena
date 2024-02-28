package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

type Player struct {
	id string

}
type Team struct {
	Name       string
	PlayerName string
	Wins       int
}

var teams []Team

func main() {
	fmt.Println("Server started...")

	http.HandleFunc("/leaderboard", getLeaderboard)
	http.HandleFunc("/initplayer", initPlayer)

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// Serve static files from the directory containing the Go file
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func getLeaderboard(w http.ResponseWriter, r *http.Request) {
	// Check if teams slice is empty
	if len(teams) == 0 {
		// Return an empty array as response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	// Encode teams slice into JSON and send as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func initPlayer(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body into a struct
	var newTeam Team
	err := json.NewDecoder(r.Body).Decode(&newTeam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, team := range teams {
		if team.Name == newTeam.Name {
			http.Error(w, "Team name already exists", http.StatusBadRequest)
			return
		}
		if team.PlayerName == newTeam.PlayerName {
			http.Error(w, "Player name already exists", http.StatusBadRequest)
			return
		}
	}

	teams = append(teams, newTeam)

	// Return a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Team %s with player %s created successfully", newTeam.Name, newTeam.PlayerName)
}
