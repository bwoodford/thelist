package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var CurrentTime = func() time.Time {
	return time.Now().UTC()
}

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		created_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		modified_date DATETIME,
		completed_date DATETIME,
		is_active BOOLEAN NOT NULL DEFAULT 1
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// InsertItem adds a new item to the database
func InsertItem(db *sql.DB, item *Item) (int64, error) {
	stmt := `
		INSERT INTO items (
			title,
			description,
			created_date,
			modified_date,
			completed_date,
			is_active
		) VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(stmt,
		item.Title,
		item.Description,
		item.CreatedDate,
		item.ModifiedDate,
		item.CompletedDate,
		item.IsActive)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func GetItem(db *sql.DB, id int64) (*Item, error) {
	stmt := `SELECT * FROM items WHERE id = ?`
	var item Item
	err := db.QueryRow(stmt, id).Scan(&item.ID, &item.Title, &item.Description, &item.CreatedDate, &item.ModifiedDate, &item.CompletedDate, &item.IsActive)
	return &item, err
}

// GetItems retrieves all items from the database
func GetItems(db *sql.DB) ([]Item, error) {
	rows, err := db.Query(`SELECT 
				id, 
				title, 
				description, 
				created_date, 
				modified_date, 
				completed_date, 
				is_active FROM items`)

	defer rows.Close()
	if err != nil {
		return nil, err
	}

	items := []Item{}
	for rows.Next() {
		var item Item
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Description,
			&item.CreatedDate,
			&item.ModifiedDate,
			&item.CompletedDate,
			&item.IsActive,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func UpdateItem(db *sql.DB, id int, item *Item) error {
	stmt := `
		UPDATE items 
		SET title = ?, 
			description = ?, 
			created_date = ?,
			modified_date = ?, 
			completed_date = ?, 
			is_active = ? 
		WHERE id = ?`

	_, err := db.Exec(
		stmt,
		item.Title,
		item.Description,
		item.CreatedDate,
		item.ModifiedDate,
		item.CompletedDate,
		item.IsActive,
		id)

	return err
}

func DeleteItem(db *sql.DB, id int) error {
	stmt := `
		UPDATE items 
		SET is_active = ? 
		WHERE id = ?`

	_, err := db.Exec(
		stmt,
		false,
		id)

	return err
}
