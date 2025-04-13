package database

import (
	"BetterScorch/secrets"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
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

	slog.Info("Database interaction", "query", query)

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

func Get(table string, fields []string, whereValues ...*DBValue) ([][]string, error) {
	if len(whereValues) == 0 {
		return nil, fmt.Errorf("at least one where value must be provided")
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Build WHERE clause and args
	whereClauses := make([]string, len(whereValues))
	args := make([]any, len(whereValues))

	for i, wv := range whereValues {
		whereClauses[i] = fmt.Sprintf("%s = ?", wv.Name)
		args[i] = wv.Value
	}

	query := fmt.Sprintf(
		"SELECT %s FROM `%s` WHERE %s",
		strings.Join(fields, ", "),
		table,
		strings.Join(whereClauses, " AND "), // 🧠 spacing fixed, no more SQL-pocalypse
	)

	slog.Info("Database interaction", "query", query, "args", args)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var results [][]string

	for rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch columns: %w", err)
		}

		// Prepare string slices for scan
		row := make([]string, len(cols))
		rowPtrs := make([]any, len(cols))
		for i := range row {
			rowPtrs[i] = &row[i]
		}

		if err := rows.Scan(rowPtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func GetAll(table string) ([][]string, error) {
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT * FROM `%s`", table)
	// slog.Info("Database interaction", "query", query)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results [][]string

	for rows.Next() {
		columns, _ := rows.Columns()
		row := make([]string, len(columns))
		rowPtrs := make([]any, len(columns))

		for i := range row {
			rowPtrs[i] = &row[i]
		}

		if err := rows.Scan(rowPtrs...); err != nil {
			return nil, err
		}

		results = append(results, row)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func Update(table string, keyValues []*DBValue, values ...*DBValue) error {
	if len(keyValues) == 0 {
		return fmt.Errorf("at least one key value must be provided")
	}

	if err := db.Ping(); err != nil {
		return err
	}

	setClauses := make([]string, len(values))
	for i, v := range values {
		setClauses[i] = fmt.Sprintf("%s = ?", v.Name)
	}
	setClause := strings.Join(setClauses, ", ")

	whereClauses := make([]string, len(keyValues))
	for i, kv := range keyValues {
		whereClauses[i] = fmt.Sprintf("%s = ?", kv.Name)
	}
	whereClause := strings.Join(whereClauses, " AND ")

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		table,
		setClause,
		whereClause,
	)

	// Proper order: SET values first, then WHERE key values
	args := make([]any, 0, len(values)+len(keyValues))
	for _, v := range values {
		args = append(args, v.Value)
	}
	for _, kv := range keyValues {
		args = append(args, kv.Value)
	}

	slog.Info("Database interaction", "query", query)

	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	return nil
}

func Remove(table string, conditions ...*DBValue) (int, error) {
	err := db.Ping()
	if err != nil {
		return -1, err
	}

	// Build WHERE clause
	whereConditions := make([]string, len(conditions))
	values := make([]interface{}, len(conditions))
	for i, condition := range conditions {
		whereConditions[i] = fmt.Sprintf("`%s` = ?", condition.Name)
		values[i] = condition.Value
	}
	whereClause := strings.Join(whereConditions, " AND ")

	query := fmt.Sprintf("DELETE FROM `%s` WHERE %s", table, whereClause)
	slog.Info("Database interaction", "query", query)

	stmt, err := db.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(values...)
	if err != nil {
		return -1, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	return int(affected), nil
}

func InsertOrUpdate(table string, keys []*DBValue, value *DBValue) error {
	allColumns := append(keys, value)

	columnNames := make([]string, len(allColumns))
	placeholders := make([]string, len(allColumns))
	args := make([]any, len(allColumns))

	for i, col := range allColumns {
		columnNames[i] = col.Name
		placeholders[i] = "?"
		args[i] = col.Value
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (%s) VALUES (%s)
		ON DUPLICATE KEY UPDATE %s=?`,
		table,
		strings.Join(columnNames, ", "),
		strings.Join(placeholders, ", "),
		value.Name,
	)

	args = append(args, value.Value)

	slog.Info("Database interaction", "query", query, "args", args)
	_, err := db.Exec(query, args...)
	return err
}

func joinNames(values []*DBValue) string {
	names := make([]string, len(values))
	for i, v := range values {
		names[i] = v.Name
	}
	return strings.Join(names, ", ")
}

func joinPlaceholders(n int) string {
	p := make([]string, n)
	for i := 0; i < n; i++ {
		p[i] = "?"
	}
	return strings.Join(p, ", ")
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

func anysToStrings(anySlice []any) []string {
	result := make([]string, len(anySlice))
	for i, val := range anySlice {
		str, _ := val.(string)
		result[i] = str
	}
	return result
}
