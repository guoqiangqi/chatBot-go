package main

import (
	chatbot "chatbot/utils"
	"fmt"
	"database/sql"
)

const (
	host      = ""
	port      = 5432
	user      = ""
	password  = ""
	dbname    = "postgresdb"
	tablename = "chatbot"
)

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return
	}
	defer db.Close()

	_, err = chatbot.TableInit(db, tablename)
	if err != nil {
		fmt.Println(err)
	}

	_, err = chatbot.InsertData(db, tablename, "Hello!", "What can i help you?") 
	if err != nil {
		fmt.Println(err)
	}

}
