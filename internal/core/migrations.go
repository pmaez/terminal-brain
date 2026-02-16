package core

import "database/sql"

func Migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS stack_tasks (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			position INTEGER NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS dump_entries (
			id INTEGER PRIMARY KEY,
			body TEXT NOT NULL,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS duck_sessions (
			id INTEGER PRIMARY KEY,
			q1 TEXT NOT NULL,
			q2 TEXT NOT NULL,
			q3 TEXT NOT NULL,
			q4 TEXT NOT NULL,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS commands (
			id INTEGER PRIMARY KEY,
			command TEXT NOT NULL,
			description TEXT NOT NULL,
			uses INTEGER NOT NULL DEFAULT 0,
			last_used TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS time_logs (
			id INTEGER PRIMARY KEY,
			label TEXT NOT NULL,
			start TEXT NOT NULL,
			end TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		);`,
	}

	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}
