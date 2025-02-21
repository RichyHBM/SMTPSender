package main

import (
	"database/sql"
	"fmt"
)

type emptyFuncReturnErr func() error

func nilErrsOrFail(errorReturnFuncs ...emptyFuncReturnErr) error {
	for _, elem := range errorReturnFuncs {
		if err := elem(); err != nil {
			return err
		}
	}
	return nil
}

func prepareStatement(db *sql.DB, statement *sql.Stmt, query string) emptyFuncReturnErr {
	return func() error {
		var err error = nil
		statement, err = db.Prepare(query)
		return err
	}
}

func execQueryIgnoreResult(db *sql.DB, query string) emptyFuncReturnErr {
	return func() error {
		if _, err := db.Exec(query); err != nil {
			return err
		}
		return nil
	}
}

type DataStore struct {
	db *sql.DB

	addSenderStmt        *sql.Stmt
	addRecipientStmt     *sql.Stmt
	addMailStmt          *sql.Stmt
	addMailRecipientStmt *sql.Stmt

	getSenderByIdStmt        *sql.Stmt
	getSenderByInitiatorStmt *sql.Stmt
	getSenderByEmailStmt     *sql.Stmt

	getRecipientByIdStmt    *sql.Stmt
	getRecipientByEmailStmt *sql.Stmt

	getMailByIdStmt          *sql.Stmt
	getMailsBySenderStmt     *sql.Stmt
	getMailsForRecipientStmt *sql.Stmt
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

	if err := nilErrsOrFail(
		execQueryIgnoreResult(db, createSenderTableSql),
		execQueryIgnoreResult(db, createRecipientTableSql),
		execQueryIgnoreResult(db, createMailTableSql),
		execQueryIgnoreResult(db, createMailToRecipientsTableSql)); err != nil {
		return nil, err
	}

	var addSenderStmt,
		addRecipientStmt,
		addMailStmt,
		addMailRecipientStmt,
		getSenderByIdStmt,
		getSenderByInitiatorStmt,
		getSenderByEmailStmt,
		getRecipientByIdStmt,
		getRecipientByEmailStmt,
		getMailByIdStmt,
		getMailsBySenderStmt,
		getMailsForRecipientStmt *sql.Stmt

	if err := nilErrsOrFail(
		prepareStatement(db, addSenderStmt, addSenderSql),
		prepareStatement(db, addRecipientStmt, addRecipientSql),
		prepareStatement(db, addMailStmt, addMailSql),
		prepareStatement(db, addMailRecipientStmt, addMailRecipientSql),
		prepareStatement(db, getSenderByIdStmt, getSenderByIdSql),
		prepareStatement(db, getSenderByInitiatorStmt, getSenderByInitiatorSql),
		prepareStatement(db, getSenderByEmailStmt, getSenderByEmailSql),
		prepareStatement(db, getRecipientByIdStmt, getRecipientByIdSql),
		prepareStatement(db, getRecipientByEmailStmt, getRecipientByEmailSql),
		prepareStatement(db, getMailByIdStmt, getMailByIdSql),
		prepareStatement(db, getMailsBySenderStmt, getMailsBySenderSql),
		prepareStatement(db, getMailsForRecipientStmt, getMailsForRecipientSql)); err != nil {
		return nil, err
	}

	return &DataStore{
		db,
		addSenderStmt,
		addRecipientStmt,
		addMailStmt,
		addMailRecipientStmt,
		getSenderByIdStmt,
		getSenderByInitiatorStmt,
		getSenderByEmailStmt,
		getRecipientByIdStmt,
		getRecipientByEmailStmt,
		getMailByIdStmt,
		getMailsBySenderStmt,
		getMailsForRecipientStmt,
	}, nil
}

func (datastore *DataStore) Close() {
	datastore.addSenderStmt.Close()
	datastore.addRecipientStmt.Close()
	datastore.addMailStmt.Close()
	datastore.addMailRecipientStmt.Close()
	datastore.getSenderByIdStmt.Close()
	datastore.getSenderByInitiatorStmt.Close()
	datastore.getSenderByEmailStmt.Close()
	datastore.getRecipientByIdStmt.Close()
	datastore.getRecipientByEmailStmt.Close()
	datastore.getMailByIdStmt.Close()
	datastore.getMailsBySenderStmt.Close()
	datastore.getMailsForRecipientStmt.Close()
}
