DROP TABLE IF EXISTS classes;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS players;

CREATE TABLE characters (
    character_id INTEGER PRIMARY KEY,
    name TEXT,			
    wins INTEGER,
    health INTEGER,
    initiative INTEGER,
    damage INTEGER,
    defense INTEGER,
    resource INTEGER,
    resource_max INTEGER,			
    lives INTEGER,
    class_id INTEGER,
    FOREIGN KEY(class_id) REFERENCES Classes(ID))
};

CREATE TABLE IF NOT EXISTS players (
    player_id VARCHAR PRIMARY KEY,
    name TEXT
);

CREATE TABLE IF NOT EXISTS players_characters (
    player_id INTEGER NOT NULL,
    character_id INTEGER NOT NULL,
    PRIMARY KEY (player_id, character_id),
    FOREIGN KEY (player_id) REFERENCES players(player_id),
    FOREIGN KEY (character_id) REFERENCES characters(character_id)
)

--CONSTANT TABLES
CREATE TABLE Classes (
    ID INTEGER PRIMARY KEY,
    Name TEXT,
    Health INTEGER,
    Initiative INTEGER,
    Damage INTEGER,
    Defense INTEGER,
    Resource_Max INTEGER,
    SpecialAbilities TEXT
);

CREATE TABLE Items (
    ID INTEGER PRIMARY KEY,
    Name TEXT,
    Description TEXT
);

--Insert initial data.
INSERT INTO Classes(ID, Name, Health, Initiative, Damage, Defense, Resource_Max, SpecialAbilities) VALUES
        (1, 'Warrior', 20, 8, 2, 3, 50, '[{"name": "Charge", "description": "Enemy skips next turn"}, {"name": "Shield up", "description": "Ready to defend the next attack."}]'),
        (2, 'Mage', 10, 12, 4, 0, 40, '[{"name": "Fireball", "description": "Launch a fiery projectile."}, {"name": "Mana shield", "description": "Temporary shield that absorbs damage."}]');
INSERT INTO Items(Name, Description) VALUES
    ('Hamburger', '+2hp'),
    ('Telescope', '+2 initiative'),
    ('Obsidian ring', '+1 damage'),        
    ('Stick pasta to body', '+1 defense'),    
    ('Full pasta armor', '+5 armor. -3 to damage');

INSERT INTO players(player_id, name) VALUES(1, 'Bot');