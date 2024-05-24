package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	HIT_THRESHOLD   = 20
	EV_TP_ATK       = "atk"
	EV_TP_INIT      = "init"
	EV_TP_END       = "end"
	EV_TP_UPD_CHARS = "upd"
	EV_TP_DEAD      = "dead"
)

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

type Character struct {
	CharacterId int
	Name        string
	Wins        int
	Health      int
	HealthMax   int
	Initiative  int
	Damage      int
	Defense     int
	Resource    int
	ResourceMax int
	Lives       int
	ClassId     int
}

func (c *Character) incrWins() error {
	_, err := db.Exec("UPDATE characters SET wins = wins + 1 WHERE character_id = ?", c.CharacterId)
	if err != nil {
		return fmt.Errorf("error incrementing character wins: %w", err)
	}
	return nil
}

func (c *Character) decrLives() error {
	_, err := db.Exec("UPDATE characters SET lives = lives - 1 WHERE character_id = ?", c.CharacterId)
	if err != nil {
		return fmt.Errorf("error decrementing character lives: %w", err)
	}
	return nil
}

func (a *Character) Attack(events *[]FightEvent, d *Character) {
	var desc strings.Builder
	var damageDealt int

	roll := randRange(1, 20)
	isMiss := roll+d.Defense > HIT_THRESHOLD
	if !isMiss {
		damageDealt = a.Damage
		d.Health -= a.Damage
		if d.Health < 0 {
			d.Health = 0
		}
		desc.WriteString(fmt.Sprintf("%s rolls a %d + (%d), fails to defend and receives %d damage!", d.Name, roll, d.Defense, a.Damage))
	} else {
		desc.WriteString(fmt.Sprintf("%s rolls a %d + (%d) and defends the attack!", d.Name, roll, d.Defense))
	}

	event := AttackEvent{AttackerName: a.Name, DefenderName: d.Name, Success: !isMiss, Damage: damageDealt}
	jsonData, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	fightEvent := FightEvent{Type: EV_TP_ATK, Data: string(jsonData)}
	*events = append(*events, fightEvent)
}

type Player struct {
}
type Classes struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	Health           int        `json:"health"`
	Initiative       int        `json:"initiative"`
	Damage           int        `json:"damage"`
	Defense          int        `json:"defense"`
	ResourceMax      int        `json:"resource_max"`
	SpecialAbilities []struct { // Consider a struct for complex abilities
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"special_abilities"`
}

type FightData struct {
	Events []FightEvent
}
type FightEvent struct {
	Type string
	Data string
}
type InitiativeEvent struct {
	StartingCharacterName string
}
type CombatEndEvent struct {
	Winner Character
}
type AttackEvent struct {
	AttackerName string
	DefenderName string
	Success      bool
	Damage       int
}
type UpdateCharactersEvent struct {
	Player   Character
	Opponent Character
}

type Leaderboard struct {
	PlayerName    string
	CharacterName string
	Wins          int
}

type CreatePlayerRequest struct {
	Name string
}
type UpgradeRequest struct {
	CharacterId string
	UpgradeId   string
}
type CreateCharacterRequest struct {
	PlayerId      string
	CharacterName string
	ClassId       string
	Wins          string
}

type SimulateFightRequest struct {
	CharacterId string
}
