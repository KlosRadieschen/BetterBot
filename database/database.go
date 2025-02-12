package database

import (
	"BetterScorch/secrets"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBValue struct {
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

func Insert(table string, values ...DBValue) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(getDBValueNames(values), ", "),
		strings.Repeat("?,", len(values)),
	)

	log.Println(fmt.Sprintf("Executing query: %v", query))

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

func Update(table string, keyValue *DBValue, values ...DBValue) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	setClause := strings.Join(getUpdateSetClause(values), ", ")
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s=?",
		table,
		setClause,
		keyValue.name,
	)

	log.Println(fmt.Sprintf("Executing query: %v", query))

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := append(getDBValues(values), keyValue.value)
	_, err = stmt.Exec(args)
	if err != nil {
		return err
	}

	return nil
}

func Remove(table string, keyValue *DBValue) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Executing query: DELETE FROM %v WHERE %v=%v", table, keyValue.name, keyValue.value))

	stmt, err := db.Prepare("DELETE FROM ? WHERE ?=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(table, keyValue.name, keyValue.value)
	if err != nil {
		return err
	}

	return nil
}

func getDBValues(dbVals []DBValue) []string {
	var values []string

	for _, dv := range dbVals {
		values = append(values, dv.value)
	}

	return values
}

func getDBValueNames(dbVals []DBValue) []string {
	var names []string

	for _, dv := range dbVals {
		names = append(names, dv.name)
	}

	return names
}

func getUpdateSetClause(dbVals []DBValue) []string {
	var setClause []string

	for _, dv := range dbVals {
		setClause = append(setClause, fmt.Sprintf("%s=?", dv.name))
	}

	return setClause
}
