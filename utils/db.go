package chatbot

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	dbHost      = os.Getenv("PGSQL_HOST")
	dbPort      = os.Getenv("PGSQL_PORT")
	dbUser      = os.Getenv("PGSQL_USER")
	dbPassword  = os.Getenv("PGSQL_PASSWORD")
	dbName      = os.Getenv("PGSQL_DBNAME")
	dbTablename = os.Getenv("PGSQL_TABLENAME")
)

func WriteQAToDB(question string, answer string) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = TableInit(db, dbTablename)
	if err != nil {
		return err
	}

	_, err = InsertData(db, dbTablename, question, answer)
	return err
}

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
	log.Println("Table created successfully!")
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
	log.Println("Data inserted successfully!")
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
		log.Println("Table is already existed.")
	}
	return true, nil
}
