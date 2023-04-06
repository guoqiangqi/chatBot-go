package chatbot

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func CreateTable(db *sql.DB, tablename string) (bool, error) {
	createTableSQL := fmt.Sprintf(`
	CREATE TABLE %s (
		id SERIAL PRIMARY KEY,
		questiones TEXT NOT NULL,
		answers TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`, tablename)

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return false, err
	}
	fmt.Println("Table created successfully!")
	return true, nil
}

func IsTableExists(db *sql.DB, tableName string) (bool, error) {
	query := fmt.Sprintf(`
        SELECT EXISTS (
            SELECT FROM information_schema.tables
            WHERE table_schema = 'public'
            AND table_name = '%s'
        );
    `, tableName)

	var exists bool
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func InsertData(db *sql.DB, tableName string, question string, answer string) (bool, error) {
	insertSQL := fmt.Sprintf("INSERT INTO %s (questiones, answers) VALUES ($1, $2)", tableName)
	_, err := db.Exec(insertSQL, question, answer)
	if err != nil {
		return false, err
	}
	fmt.Println("Data inserted successfully!")
	return true, nil
}

func TableInit(db *sql.DB, tableName string) (bool, error) {
	exists, err := IsTableExists(db, tableName)
	if err != nil {
		return false, err
	}

	if !exists {
		return CreateTable(db, tableName)
	} else {
		fmt.Println("Table is already existed.")
	}
	return true, nil
}
