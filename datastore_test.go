package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*sql.DB, func(*sql.DB) error) {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		t.Fatal(err)
	}

	return db, func(db *sql.DB) error {
		return db.Close()
	}
}

func ensureValidDataStore(t *testing.T, datastore *DataStore) {
	assert.NotNil(t, datastore)
	assert.NotNil(t, datastore.db)
	assert.NotNil(t, datastore.addSenderStmt)
	assert.NotNil(t, datastore.addRecipientStmt)
	assert.NotNil(t, datastore.addMailStmt)
	assert.NotNil(t, datastore.addMailRecipientStmt)
	assert.NotNil(t, datastore.getSenderByIdStmt)
	assert.NotNil(t, datastore.getSenderByInitiatorStmt)
	assert.NotNil(t, datastore.getSenderByEmailStmt)
	assert.NotNil(t, datastore.getRecipientByIdStmt)
	assert.NotNil(t, datastore.getRecipientByEmailStmt)
	assert.NotNil(t, datastore.getMailByIdStmt)
	assert.NotNil(t, datastore.getMailsBySenderStmt)
	assert.NotNil(t, datastore.getMailsForRecipientStmt)
}

func TestMakeDataStore(t *testing.T) {
	db, teardown := setup(t)
	defer teardown(db)

	datastore, err := MakeDataStore(db)
	assert.NoError(t, err)
	defer datastore.Close()
	
	ensureValidDataStore(t, datastore)
}

func TestErrorOnInvalidDb(t *testing.T) {
	datastore, err := MakeDataStore(nil)
	assert.Nil(t, datastore)
	assert.Error(t, err)
}
