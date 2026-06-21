package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./guppy.db")
	if err != nil {
		return err
	}

	createTables()
	return nil
}

func createTables() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS tanks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			location TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("failed to create tanks table: %v", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS water_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tank_id INTEGER NOT NULL,
			temperature REAL NOT NULL,
			pH REAL NOT NULL,
			recorded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (tank_id) REFERENCES tanks(id)
		)
	`)
	if err != nil {
		log.Fatalf("failed to create water_records table: %v", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS breeding_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tank_id INTEGER NOT NULL,
			strain TEXT NOT NULL,
			pair_date TEXT NOT NULL,
			expected_birth_date TEXT DEFAULT '',
			notes TEXT DEFAULT '',
			status TEXT DEFAULT 'breeding',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (tank_id) REFERENCES tanks(id)
		)
	`)
	if err != nil {
		log.Fatalf("failed to create breeding_records table: %v", err)
	}

	insertInitialTanks()
}

func insertInitialTanks() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM tanks").Scan(&count)
	if err != nil {
		log.Printf("error counting tanks: %v", err)
		return
	}
	if count > 0 {
		return
	}

	initialTanks := []string{"1号缸 - 白子孔雀", "2号缸 - 红孔雀", "3号缸 - 蓝孔雀", "4号缸 - 幼鱼缸"}
	for _, name := range initialTanks {
		_, err := DB.Exec("INSERT INTO tanks (name) VALUES (?)", name)
		if err != nil {
			log.Printf("error inserting initial tank: %v", err)
		}
	}
	log.Println("Initial tanks inserted")
}
