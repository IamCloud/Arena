package main

import (
	"encoding/json"
	"fmt"
	"strconv"
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

func (winner *Character) incrWins() {
	_, err := db.Exec("UPDATE characters SET wins = wins + 1 WHERE character_id = ?", winner.CharacterId)
	if err != nil {
		fmt.Println("error incrementing character wins.", err)
		return
	}
}

func (loser *Character) decrLives() {
	_, err := db.Exec("UPDATE characters SET lives = lives - 1 WHERE character_id = ?", loser.CharacterId)
	if err != nil {
		fmt.Println("error decrementing character lives.", err)
		return
	}
}

func (attacker *Character) Attack(events *[]FightEvent, defender *Character) {
	var desc strings.Builder

	var damageDealt int = 0

	roll := randRange(1, 20)
	isMiss := roll+defender.Defense > HIT_THRESHOLD
	if !isMiss {
		damageDealt = attacker.Damage
		(*defender).Health -= attacker.Damage
		if (*defender).Health < 0 {
			(*defender).Health = 0
		}
		desc.WriteString(fmt.Sprintf("%s rolls a %s + (%s), fails to defend and receives %s damage!", defender.Name, strconv.Itoa(roll), strconv.Itoa(defender.Defense), strconv.Itoa(attacker.Damage)))
	} else {
		desc.WriteString(fmt.Sprintf("%s rolls a %s + (%s) and defends the attack !", defender.Name, strconv.Itoa(roll), strconv.Itoa(defender.Defense)))
	}

	event := AttackEvent{AttackerName: attacker.Name, DefenderName: defender.Name, Success: !isMiss, Damage: damageDealt}
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

type CreateCharacterRequest struct {
	PlayerId      string
	CharacterName string
	ClassId       string
	Wins          string
}

type SimulateFightRequest struct {
	CharacterId string
}
