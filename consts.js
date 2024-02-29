const STORED_PLAYERID = "playerid";
const STORED_TEAMID = "teamid";
const MAX_TEAM_SIZE = 5;

const Abilities = {
    Heal: {
        name: "Heal",
        description: "Heals an ally or itself for x health."
    },
    Whirlwind: {
        name: "Whirlwind",
        description: "Spins around and attacks all enemies."
    },
}
const Heroes = {
    Warrior: {
        name: "Warrior",
        health: 200,
        damage: 5,
        abilities: {
            0: Abilities.Whirlwind
        }
    },
    Archer: {
        name: "Archer",
        health: 100,
        damage: 10
    },
    Priest: {
        name: "Priest",
        health: 125,
        damage: 1,
        abilities: {
            0: Abilities.Heal
        }
    },
}