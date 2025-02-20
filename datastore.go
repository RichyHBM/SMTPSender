package main

import "database/sql"

type Sender struct {
	Initiator string
	Email string
}

type Recipient struct {
	Email string
}

type Mail struct {
	From    Sender
	To      []Recipient
	Subject string
	Body    string
}

type DataStore struct {
	db *sql.DB
}

func MakeDataStore(db *sql.DB) (*DataStore, error) {
	return &DataStore{
		db,
	}, nil
}