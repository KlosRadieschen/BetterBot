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
	Name  string
	Value string
}

var db *sql.DB

func Connect() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", secrets.DBUser, secrets.DBPassword, secrets.DBAddress, secrets.DBName)
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

func Insert(table string, values ...*DBValue) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	fmt.Println(values)

	query := fmt.Sprintf(
		"INSERT INTO `%s`(%s) VALUES (%s)",
		table,
		strings.Join(getDBValueNames(values), ", "),
		strings.TrimSuffix(strings.Repeat("?,", len(values)), ","),
	)

	log.Println(fmt.Sprintf("Executing query: %v", query))

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(stringsToAnys(getDBValues(values))...)
	if err != nil {
		return err
	}

	return nil
}

func Update(table string, keyValue *DBValue, values ...*DBValue) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	setClause := strings.Join(getUpdateSetClause(values), ", ")
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s=?",
		table,
		setClause,
		keyValue.Name,
	)

	log.Println(fmt.Sprintf("Executing query: %v", query))

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := append(getDBValues(values), keyValue.Value)
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

	log.Println(fmt.Sprintf("Executing query: DELETE FROM %v WHERE %v=%v", table, keyValue.Name, keyValue.Value))

	stmt, err := db.Prepare("DELETE FROM ? WHERE ?=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(table, keyValue.Name, keyValue.Value)
	if err != nil {
		return err
	}

	return nil
}

func getDBValues(dbVals []*DBValue) []string {
	var values []string

	for _, dv := range dbVals {
		values = append(values, dv.Value)
	}

	return values
}

func getDBValueNames(dbVals []*DBValue) []string {
	var names []string

	for _, dv := range dbVals {
		names = append(names, dv.Name)
	}

	return names
}

func getUpdateSetClause(dbVals []*DBValue) []string {
	var setClause []string

	for _, dv := range dbVals {
		setClause = append(setClause, fmt.Sprintf("%s=?", dv.Name))
	}

	return setClause
}

func stringsToAnys(strings []string) []interface{} {
	result := make([]interface{}, len(strings))
	for i, v := range strings {
		result[i] = v
	}
	return result
}
