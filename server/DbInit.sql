DROP TABLE IF EXISTS Classes;
DROP TABLE IF EXISTS Items;
-- Create Classes table
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
INSERT INTO Classes(ID, Name, Health, Initiative, Damage, Defense, Resource_Max, SpecialAbilities) VALUES
        (1, 'Warrior', 20, 8, 2, 3, 50, '[{"name": "Charge", "description": "Enemy skips next turn"}, {"name": "Shield up", "description": "Ready to defend the next attack."}]'),
        (2, 'Mage', 10, 12, 4, 0, 40, '[{"name": "Fireball", "description": "Launch a fiery projectile."}, {"name": "Mana shield", "description": "Temporary shield that absorbs damage."}]');

-- Create Items table
CREATE TABLE Items (
    ID INTEGER PRIMARY KEY,
    Name TEXT,
    Description TEXT
);
INSERT INTO Items(Name, Description) VALUES
    ('Hamburger', '+2hp'),
    ('Telescope', '+2 initiative'),
    ('Obsidian ring', '+1 damage'),        
    ('Stick pasta to body', '+1 defense'),    
    ('Full pasta armor', '+5 armor. -3 to damage');
