package main

import (
	"fmt"
	"strconv"
)

const (
	HIT_THRESHOLD = 15
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
	Name        string
	Wins        int
	Health      int
	Initiative  int
	Damage      int
	Defense     int
	Resource    int
	ResourceMax int
	Lives       int
	ClassId     int
}

func (c *Character) Attack(events *[]FightEvent, target *Character) {
	fmt.Printf("%s attacks %s!\n", c.Name, target.Name)

	isMiss := randRange(1, 20)+target.Defense > HIT_THRESHOLD
	if !isMiss {
		(*target).Health -= c.Damage
		fmt.Printf("%s hits %s for %s!\n", c.Name, target.Name, strconv.Itoa(c.Damage))
		fmt.Printf("Target %s current health: %s!\n", target.Name, strconv.Itoa(target.Health))
	} else {
		fmt.Printf("%s misses %s\n", c.Name, target.Name)
	}

	event := FightEvent{Defender: target.Name, Hit: !isMiss, DmgReceived: c.Damage}
	*events = append(*events, event)
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

type SpecialAbility struct {
	Name        string
	Description string
}

type FightData struct {
	StartPlayerInfo   Character
	StartOpponentInfo Character
	Events            []FightEvent
}
type FightEvent struct {
	Defender    string
	Hit         bool
	DmgReceived int
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
