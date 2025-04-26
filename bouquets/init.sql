CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(32) unique,
    password VARCHAR(64),
    preferences VARCHAR(64),
    current_balance INTEGER DEFAULT 100,
    priveleged BOOLEAN
);

CREATE TABLE IF NOT EXISTS bouquet(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(32),
    owner VARCHAR(32),
    cost INTEGER,
    description VARCHAR(64),
    sent_to VARCHAR(32), --THINK ABOUT IT
    FOREIGN KEY (owner)
        REFERENCES users(username)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS flowers(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(32),
    color CHAR,
    cost INTEGER
);

CREATE TABLE IF NOT EXISTS flowers_bouquet(
    bouquet_id INTEGER,
    flower_id INTEGER,
    FOREIGN KEY (bouquet_id) 
        REFERENCES bouquet(id)
        ON DELETE CASCADE,
     FOREIGN KEY (flower_id) 
        REFERENCES flowers(id)
        ON DELETE CASCADE
);

INSERT OR REPLACE INTO flowers VALUES (1, "rose", "R", 100), (2, "aster", "W", 100), (3, "chamomile", "Y", 100), (4, "peony", "P", 100),  (5, "tulip", "P", 100);
