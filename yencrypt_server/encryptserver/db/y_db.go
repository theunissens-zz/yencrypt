package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Import used for postgres driver
	"strings"
)

const DBName = "yotest"

type YDatabaseInterface interface {
	Setup() error
	Connect(connectToServer bool) error
	Store(id string, data string) error
	Retrieve(id string) (string, error)
	DropDB() error
}

type YDatabase struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	DB       *sql.DB
}

func (db *YDatabase) Connect(connectToServer bool) error {
	var psqlInfo string
	// If we connect to database server, we don't specify a dbname.
	// We generally do this to drop a database (Can't drop a database you are connected to)
	if connectToServer {
		psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
			db.Host, db.Port, db.User, db.Password)
	} else {
		psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			db.Host, db.Port, db.User, db.Password, db.DBName)
	}

	var err error
	db.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	return nil
}

func (db *YDatabase) Setup() error {
	err := db.createDBIfNotExists()
	err = db.createTableIfNotExists()
	return err
}

func (db *YDatabase) createDBIfNotExists() error {
	db.Connect(true)
	defer db.Close()

	exists, err := db.checkDBExists()
	if err != nil {
		return err
	}

	if !exists {
		// Create db
		dbCreationScript := db.getDBScript()
		_, err := db.DB.Exec(dbCreationScript)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("database \"%s\" already exists", db.DBName)) {
				//ignore
			} else {
				return err
			}
		}
	}
	return nil
}

func (db *YDatabase) createTableIfNotExists() error {
	db.Connect(false)
	defer db.Close()

	exists, err := db.checkTableExists()
	if err != nil {
		return err
	}

	if !exists {
		// Create table
		tableCreationScript := db.getTableScript()
		_, err = db.DB.Exec(tableCreationScript)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *YDatabase) checkDBExists() (bool, error) {
	return db.checkExists(fmt.Sprintf("SELECT datname from pg_database WHERE datname = '%s'", db.DBName))
}

func (db *YDatabase) checkTableExists() (bool, error) {
	return db.checkExists(fmt.Sprint("select tablename from pg_tables where schemaname='public' and tablename = 'cryptodata';"))
}

func (db *YDatabase) checkExists(script string) (bool, error) {
	var name string
	err := db.DB.QueryRow(script).Scan(&name)
	switch err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

func (db *YDatabase) Store(id string, data string) error {
	_, err := db.DB.Exec(fmt.Sprintf("INSERT INTO cryptodata(userid, data) VALUES('%v', '%v')", id, data))
	if err != nil {
		return err
	}
	return nil
}

func (db *YDatabase) Retrieve(id string) (string, error) {
	var data string
	err := db.DB.QueryRow(fmt.Sprintf("SELECT data FROM cryptodata WHERE userid = '%v'", id)).Scan(&data)

	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", err
	default:
		return data, nil
	}
}

func (db *YDatabase) DropDB() error {
	db.Close()
	db.Connect(true)
	_, err := db.DB.Exec(fmt.Sprintf("DROP DATABASE %s;", db.DBName))
	if err != nil {
		return err
	}
	return nil
}

func (db *YDatabase) getDBScript() string {
	return fmt.Sprintf("CREATE DATABASE %s", db.DBName)
}

func (db *YDatabase) getTableScript() string {
	script := "CREATE TABLE cryptodata(id bigserial PRIMARY KEY, userid varchar, data varchar);" +
		"CREATE UNIQUE INDEX idx_user_id ON cryptodata (userid);"
	return script
}

func (db *YDatabase) Close() {
	db.DB.Close()
}
