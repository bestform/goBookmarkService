package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	FileName string
}

func (s SqliteStorage) Read() EntryList {
	db, err := openDb(s.FileName)
	if err != nil {
		log.Fatal(err)
	}
	var entryList EntryList
	entryList.Entries = make(map[string]Entry)
	rows, err := db.Query("select title, url from entries")
	if err != nil {
		log.Fatal(err) // @todo: we should probably alter the interface to be able to return the error instead of handling it ourselves
	}
	defer rows.Close()
	for rows.Next() {
		var title string
		var url string
		rows.Scan(&title, &url)
		entryList.Entries[title] = Entry{Title: title, Url: url}
	}

	return entryList
}

func (s SqliteStorage) Write(e Entry) {
	db, err := openDb(s.FileName)
	if err != nil {
		log.Fatal(err)
	}
	tx, _ := db.Begin()
	stmt, err := tx.Prepare("replace into entries(title, url) values(?, ?)")
	if err != nil {
		log.Fatal(err) // @todo: we should probably alter the interface to be able to return the error instead of handling it ourselves
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.Title, e.Url)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func (s SqliteStorage) DeleteEntry(t string) {
	db, err := openDb(s.FileName)
	if err != nil {
		log.Fatal(err)
	}
	tx, _ := db.Begin()
	stmt, err := tx.Prepare("delete from entries where title=?")
	if err != nil {
		log.Fatal(err) // @todo: we should probably alter the interface to be able to return the error instead of handling it ourselves
	}
	defer stmt.Close()
	_, err = stmt.Exec(t)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func NewSqliteStorage(fileName string) (SqliteStorage, error) {
	if _, err := os.Stat(fileName); err == nil {
		return SqliteStorage{FileName: fileName}, nil
	}

	db, err := openDb(fileName)
	defer db.Close()
	if err != nil {
		fmt.Println("error opening db")
		return SqliteStorage{}, err
	}

	sqlStmt := "create table entries (title text not null primary key, url text);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Println("error creating table")
		return SqliteStorage{}, err
	}

	return SqliteStorage{FileName: fileName}, nil
}

func openDb(fileName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
