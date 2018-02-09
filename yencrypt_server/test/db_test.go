package test

import (
	. "github.com/yencrypt/yencrypt_server/encryptserver/db"
	"strings"
	"testing"
)

const testDbName = "yotitest_test"

// This is an integration test for the database persistence layer

func TestDBConnection_success(t *testing.T) {
	_, err := connectToDB(true)
	if err != nil {
		t.Errorf("An error occurred trying to connect to database server: %v", err)
	}
}

func TestDBCreation_success(t *testing.T) {
	dbConnToServer := YDatabase{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   testDbName,
	}
	err := dbConnToServer.Setup()
	if err != nil {
		t.Errorf("An error occurred trying to run setup scripts: %v", err)
	}

	err = dbConnToServer.Setup()
	if err != nil {
		t.Errorf("No error expected, database setup: %v", err)
	}

	dbConnToDB, err := connectToDB(false)
	if err != nil {
		t.Errorf("An error occurred trying to connect to database: %v", err)
	}

	dbConnToDB.Close()

	err = dropDB()
	if err != nil {
		t.Errorf("An error occurred trying to drop database: %v", err)
	}
}

func TestTableCreation_success(t *testing.T) {
	dbConnToServer, err := connectToDB(true)
	if err != nil {
		t.Errorf("An error occurred trying to connect to database server: %v", err)
	}

	err = dbConnToServer.Setup()
	if err != nil {
		t.Errorf("An error occurred trying to run setup scripts: %v", err)
	}

	err = dropDB()
	if err != nil {
		t.Errorf("An error occurred trying to drop database: %v", err)
	}
}

func TestPersistenceAndRetrieval_success(t *testing.T) {
	id := "1"
	plainText := "someData,moreData"

	dbConnToServer, err := connectToDB(true)
	if err != nil {
		t.Errorf("An error occurred trying to connect to database server: %v", err)
	}

	err = dbConnToServer.Setup()
	if err != nil {
		t.Errorf("An error occurred trying to run setup scripts: %v", err)
	}

	dbConnToDB, err := connectToDB(false)
	if err != nil {
		t.Errorf("An error occurred trying to connect to database: %v", err)
	}

	err = dbConnToDB.Store(id, plainText)
	if err != nil {
		t.Errorf("An error occurred trying to persist to database: %v", err)
	}

	retrievedData, err := dbConnToDB.Retrieve(id)
	if err != nil {
		t.Errorf("An error occurred trying to retrieve from database: %v", err)
	}

	dbConnToDB.Close()

	if strings.Compare(plainText, retrievedData) != 0 {
		t.Error("Retrieved data does not match stored data")
	}

	err = dropDB()
	if err != nil {
		t.Errorf("An error occurred trying to drop database: %v", err)
	}
}

func connectToDB(connectToServer bool) (*YDatabase, error) {
	db := YDatabase{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   testDbName,
	}
	err := db.Connect(connectToServer)

	return &db, err

}

func dropDB() error {
	conn, err := connectToDB(true)
	if err != nil {
		return err
	}
	err = conn.DropDB()
	if err != nil {
		return err
	}
	return nil
}
