package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/lib/pq" // postgres
)

// Postgres contains objects for database communication.
type Postgres struct {
	*sql.DB
}

// NewPostgres create new database instance.
func NewPostgres(config *PostgresConfig) (*Postgres, error) {
	psqlInfo := fmt.Sprintf("host=%v port=%v user=%v "+
		"password=%v dbname=%v sslmode=%v",
		config.Host, config.Port, config.User, config.Password, config.DbName, config.SSL)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}

//Exec run sql query
func (db *Postgres) Exec(query string, args ...interface{}) (bool, error) {
	_, err := db.DB.Exec(query, args)
	if err != nil {
		return false, err
	}

	return true, nil
}

//QueryRow run sql query and expect return
func (db *Postgres) QueryRow(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.DB.Query(query, args)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

//Shutdown close connection
func (db *Postgres) Shutdown() {
	db.DB.Close()
}

//RunFileQuery Run sql schema
func (db *Postgres) RunFileQuery(file string) error {
	query, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	if _, err := db.DB.Exec(string(query)); err != nil {
		return err
	}

	return nil
}
