package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// Call of Cthulhu modules
func init() {
}

func createCharacter() {
	conn, err := sql.Open("sqlite3", "file:locked.sqlite?cache=shared")
	if err != nil {
		fmt.Println("Initialize database failed! ", err)
	}
	conn.Exec("select * from character;")
}
