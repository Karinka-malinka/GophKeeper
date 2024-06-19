CREATE TABLE IF NOT EXISTS users (
		uuid TEXT PRIMARY KEY,
		login TEXT,
		hash_pass TEXT,
		UNIQUE (login)
	  );

CREATE TABLE IF NOT EXISTS logindata (
		uuid TEXT PRIMARY KEY,
		created TIMESTAMP,
		login BYTEA,
		password BYTEA,
		meta BYTEA,
		user_id TEXT,
		FOREIGN KEY (user_id) REFERENCES users(uuid)
	  );

CREATE TABLE IF NOT EXISTS textdata (
		uuid TEXT PRIMARY KEY,
		textdata BYTEA,
		meta BYTEA,
		user_id TEXT,
		FOREIGN KEY (user_id) REFERENCES users(uuid)
	  );

CREATE TABLE IF NOT EXISTS filedata (
		uuid TEXT PRIMARY KEY,
		filedata BYTEA,
		meta BYTEA,
		user_id TEXT,
		FOREIGN KEY (user_id) REFERENCES users(uuid)
	  );

CREATE TABLE IF NOT EXISTS bankcard (
		numbercard BYTEA PRIMARY KEY,
		term BYTEA,
		ccv BYTEA,
		meta BYTEA,
		user_id TEXT,
		FOREIGN KEY (user_id) REFERENCES users(uuid)
	  );