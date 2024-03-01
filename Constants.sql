DROP TABLE IF EXISTS Classes;
DROP TABLE IF EXISTS Items;
DROP TABLE IF EXISTS Rarity;
-- Create Classes table
CREATE TABLE Classes (
    ID INTEGER PRIMARY KEY,
    Name TEXT,
    SpecialAbilities TEXT
);
INSERT INTO Classes(Name, SpecialAbilities) VALUES
        ('Warrior', '[{"name": "Charge", "description": "Dash forward and stuns an enemy."}, {"name": "Whirlwind", "description": "Spin around and damages all enemies."}]'),
        ('Mage', '[{"name": "Fireball", "description": "Launch a fiery projectile."}, {"name": "Frost Nova", "description": "Slows all enemies."}]'),
        ('Priest', '[{"name": "Heal", "description": "Heal the weakest ally."}, {"name": "Smite", "description": "Smite down a weakest opponent."}]');

-- Create Items table
CREATE TABLE Items (
    ID INTEGER PRIMARY KEY,
    Name TEXT,
    Description TEXT
);
INSERT INTO Items(Name, Description) VALUES
    ('Obsidian Ring', 'Doubles physical damage.'),
    ('Magical Gloves', 'Doubles magical damage.'),
    ('Pasta Armor', '+2 armor.'),    
    ('Heavy Armor', '+5 armor. -30% physical damage. -60% magical damage.');

-- Create Rarity table
CREATE TABLE Rarity (
    ID INTEGER PRIMARY KEY,
    Name TEXT
);
INSERT INTO Rarity(ID, Name) VALUES
    (0, 'Common'),
    (1, 'Uncommon'),
    (2, 'Rare'),    
    (3, 'Epic'),    
    (4, 'Legendary');
