package main

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func(*sql.DB) error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectExec("PRAGMA foreign_keys = ON").WillReturnResult(driver.ResultNoRows)

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		t.Fatal(err)
	}

	return db, mock, func(db *sql.DB) error {
		return db.Close()
	}
}

func ensureValidDataStore(t *testing.T, datastore *DataStore) {
	assert.NotNil(t, datastore)
	if datastore != nil {
		assert.NotNil(t, datastore.db)
	}
}

func TestMakeDataStore(t *testing.T) {
	db, mock, teardown := setup(t)
	defer teardown(db)

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(driver.ResultNoRows)

	datastore, err := MakeDataStore(db)
	if !assert.NoError(t, err) {
		return
	}
	defer datastore.Close()

	ensureValidDataStore(t, datastore)
}

func TestErrorOnInvalidDb(t *testing.T) {
	datastore, err := MakeDataStore(nil)
	assert.Nil(t, datastore)
	assert.Error(t, err)
}
