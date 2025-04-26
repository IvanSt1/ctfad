import { Database } from "bun:sqlite";

const USERS_TABLE = `
CREATE TABLE IF NOT EXISTS users (
	username 	  STRING PRIMARY KEY ON CONFLICT FAIL,
	password_hash STRING NOT NULL,
	strategy 	  STRING NOT NULL  ON CONFLICT REPLACE DEFAULT "",
	attack 		  REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.00,
	defense 	  REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.00
);
`

const CHALLENGES_TABLE = `
CREATE TABLE IF NOT EXISTS challenges (
	id 	  	 INTEGER PRIMARY KEY ON CONFLICT FAIL AUTOINCREMENT,
	attacker STRING NOT NULL,
	defender STRING NOT NULL,

	FOREIGN KEY (attacker) REFERENCES users(username),
	FOREIGN KEY (defender) REFERENCES users(username),

	UNIQUE (attacker, defender) ON CONFLICT FAIL
);
`

const DUELS_TABLE = `
CREATE TABLE IF NOT EXISTS duels (
	id 	  	 	 INTEGER PRIMARY KEY ON CONFLICT FAIL AUTOINCREMENT,
	attacker 	 STRING NOT NULL,
	defender 	 STRING NOT NULL,
	attacker_won INTEGER NOT NULL,

	FOREIGN KEY (attacker) REFERENCES users(username),
	FOREIGN KEY (defender) REFERENCES users(username)
);
`;

export function initDatabase(dbPath) {
	const db = new Database(dbPath);

	const createAll = db.transaction(() => {
		db.query(USERS_TABLE).run();
		db.query(CHALLENGES_TABLE).run();
		db.query(DUELS_TABLE).run();
	});
	createAll();

	console.log(`Using database ${dbPath}`);
	return db;
}
