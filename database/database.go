package database

import (
	"BetterScorch/secrets"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type dbValue struct {
	name  string
	value string
}

var db *sql.DB

func Connect() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", secrets.DBAddress, secrets.DBPassword, secrets.DBAddress, secrets.DBName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(3)                // Maximum number of open connections to the database.
	db.SetMaxIdleConns(1)                // Maximum number of idle connections.
	db.SetConnMaxLifetime(time.Hour * 2) // Connections are recycled after two hours.
}

func Insert(table string, values ...dbValue) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(getDBValueNames(values), ", "),
		strings.Repeat("?,", len(values)))

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(getDBValues(values))
	if err != nil {
		return err
	}

	return nil
}

func getDBValues(dbVals []dbValue) []string {
	var values []string

	for _, dv := range dbVals {
		values = append(values, dv.value)
	}

	return values
}

func getDBValueNames(dbVals []dbValue) []string {
	var names []string

	for _, dv := range dbVals {
		names = append(names, dv.name)
	}

	return names
}
