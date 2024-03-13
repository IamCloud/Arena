package main

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ClassInfo struct {
	Health      int
	Initiative  int
	Damage      int
	Defense     int
	ResourceMax int
}

type CharacterInfo struct {
	Name         string
	Wins         int
	Health       int
	Initiative   int
	Damage       int
	Defense      int
	Resource     int
	Resource_Max int
	Lives        int
	ClassId      int
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

type CreateCharacterRequest struct {
	PlayerId      string
	CharacterName string
	ClassId       string
}

type SimulateFightRequest struct {
	CharacterId string
}
