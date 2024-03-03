package main

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Class struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	SpecialAbilities string `json:"special_abilities"`
}

type Leaderboard struct {
	Name string
	Wins int
}

type CreatePlayerRequest struct {
	Name string
}

type CreateTeamRequest struct {
	PlayerId string
	TeamName string
}
