package core

import "database/sql"

type AppState struct {
	DB      *sql.DB
	Stores  Stores
	Modules map[string]Module
	Order   []string
	Active  string
}
