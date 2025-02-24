package main

import (
	"database/sql"
	"fmt"
)

type Sender struct {
	Initiator string
	Email     string
}

type Recipient struct {
	Email string
}

type Mail struct {
	From    *Sender
	To      []*Recipient
	Subject string
	Body    string
}

type DataStore struct {
	db *sql.DB
}

const (
	createSenderTableSql           = "CREATE TABLE IF NOT EXISTS sender (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, initiator TEXT NOT NULL, email TEXT NOT NULL, UNIQUE (initiator, email));"
	createRecipientTableSql        = "CREATE TABLE IF NOT EXISTS recipient (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, email TEXT UNIQUE NOT NULL);"
	createMailTableSql             = "CREATE TABLE IF NOT EXISTS mail (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, sender_id INTEGER NOT NULL, subject TEXT NOT NULL, body TEXT NOT NULL, FOREIGN KEY (sender_id) REFERENCES sender(id));"
	createMailToRecipientsTableSql = "CREATE TABLE IF NOT EXISTS mail_recipient (mail_id INTEGER NOT NULL, recipient_id INTEGER NOT NULL, PRIMARY KEY (mail_id, recipient_id), FOREIGN KEY (mail_id) REFERENCES mail(id), FOREIGN KEY (recipient_id) REFERENCES recipient(id));"

	addSenderSql        = "INSERT INTO sender(initiator, email) VALUES(?, ?);"
	addRecipientSql     = "INSERT INTO recipient(email) VALUES(?);"
	addMailSql          = "INSERT INTO mail(sender_id, subject, body) VALUES(?, ?, ?);"
	addMailRecipientSql = "INSERT INTO mail_recipient(mail_id, recipient_id) VALUES(?, ?);"

	getSenderByIdSql        = "SELECT * FROM sender WHERE id = ?;"
	getSenderByInitiatorSql = "SELECT * FROM sender WHERE initiator = ?;"
	getSenderByEmailSql     = "SELECT * FROM sender WHERE email = ?;"

	getRecipientByIdSql    = "SELECT * FROM recipient WHERE id = ?;"
	getRecipientByEmailSql = "SELECT * FROM recipient WHERE email = ?;"

	getMailByIdSql          = "SELECT * FROM mail WHERE id = ?;"
	getMailsBySenderSql     = "SELECT mail.* FROM sender, mail WHERE sender.email = ? AND mail.sender_id = sender.id;"
	getMailsForRecipientSql = "SELECT mail.* FROM recipient, mail, mail_recipient WHERE recipient.email = ? AND mail_recipient.recipient_id = recipient.id AND mail_recipient.mail_id = mail.id;"
)

func MakeDataStore(db *sql.DB) (*DataStore, error) {
	if db == nil {
		return nil, fmt.Errorf("nil sql.DB supplied to MakeDataStore")
	}

	if _, err := db.Exec(createSenderTableSql); err != nil {
		return nil, err
	}

	if _, err := db.Exec(createRecipientTableSql); err != nil {
		return nil, err
	}

	if _, err := db.Exec(createMailTableSql); err != nil {
		return nil, err
	}

	if _, err := db.Exec(createMailToRecipientsTableSql); err != nil {
		return nil, err
	}

	return &DataStore{
		db,
	}, nil
}

func (datastore *DataStore) Close() {
}
