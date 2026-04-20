package database_folder

import (
	"database/sql"
	"os"
)

func CreateDB() (*DB, error) {
	db, err := sql.Open("sqlite3", "database_folder/db.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := &DB{Db: db}

	if err := generateTables(database); err != nil {
		return nil, err
	}
	if err := addBaseCategory(database); err != nil {
		return nil, err
	}

	return database, nil
}

func generateTables(db *DB) error {
	queryBytes, err := os.ReadFile("database_folder/command.sql")
	if err != nil {
		return err
	}

	_, err = db.Db.Exec(string(queryBytes))
	return err
}

func addBaseCategory(db *DB) error {
	queryBates, err := os.ReadFile("database_folder/baseCategory.sql")
	if err != nil {
		return err
	}
	_, err = db.Db.Exec(string(queryBates))
	if err != nil {
		return err
	}

	return nil
}
