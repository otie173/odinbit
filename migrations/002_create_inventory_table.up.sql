CREATE TABLE IF NOT EXISTS Inventory (
    player_id INTEGER UNIQUE NOT NULL,
    inventory_data BLOB NOT NULL DEFAULT 0,
    FOREIGN KEY(player_id) REFERENCES Player(player_id)
);
