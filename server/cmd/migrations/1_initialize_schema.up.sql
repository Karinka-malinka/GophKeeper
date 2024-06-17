CREATE TABLE IF NOT EXISTS users (
		uuid TEXT PRIMARY KEY,
		login TEXT,
		hash_pass TEXT,
		UNIQUE (login)
	  );