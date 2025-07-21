package itemModel

import (
	"database/sql"
)

var db *sql.DB

type Item struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func StartDb(dbParam *sql.DB) {
	db = dbParam
}

func GetAllItem() ([]Item, error) {
	rows, err := db.Query("SELECT * FROM item")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []Item

	for rows.Next() {
		var item Item
		err := rows.Scan(&item.Id, &item.Title, &item.Description)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func PostItem(item Item) error {
	_, err := db.Exec(
		"INSERT INTO item (id, title, description) VALUES (?, ?, ?)",
		item.Id, item.Title, item.Description,
	)

	return err
}
