package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"todo-list/internal/config"
)

type DB struct {
	conn   *pgxpool.Pool
	config config.DB
}

func connect(dsn string, ctx context.Context) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, err
	}

	return dbPool, nil
}

func NewDB(config config.DB, ctx context.Context) (*DB, error) {
	db := &DB{
		config: config,
		conn:   nil,
	}

	dsn := db.GetDSN()

	conn, err := connect(dsn, ctx)
	if err != nil {
		return nil, err
	}

	db.conn = conn

	return db, nil
}

func (db *DB) Connection() *pgxpool.Pool {
	return db.conn
}

func (db *DB) Close() error {
	db.conn.Close()
	return nil
}

func (db *DB) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db.config.User, db.config.Password, db.config.Host, db.config.Port, db.config.Database)
}

func (db *DB) Insert(table string, columns []string, values []any) (pgx.Rows, error) {
	if len(columns) != len(values) {
		return nil, fmt.Errorf("columns and values length mismatch")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", table, strings.Join(columns, ", "))
	args := []any{}

	query += "("
	for i := range values {
		query += fmt.Sprintf("$%d, ", i+1)
	}

	query = query[0 : len(query)-2]
	query += ") RETURNING *"

	args = append(args, values...)

	rows, err := db.conn.Query(context.Background(), query, args...)

	if err != nil {
		fmt.Printf("Insert error: %v", err)
		return nil, err
	}

	return rows, nil
}

func (db *DB) Select(table string, columns []string, where map[string]any) (pgx.Rows, error) {
	// If columns is empty, use "*" as default
	columnsStr := "*"
	if len(columns) > 0 {
		columnsStr = strings.Join(columns, ", ")
	}

	args := []any{}

	query := fmt.Sprintf("SELECT %s FROM %s", columnsStr, table)

	if len(where) > 0 {
		whereClause := ""

		i := 0
		for k, v := range where {
			whereClause += fmt.Sprintf("%s = $%d AND ", k, i+1)
			args = append(args, v)
			i++
		}

		whereClause = whereClause[0 : len(whereClause)-5]

		query += fmt.Sprintf(" WHERE %s", whereClause)
	}

	rows, err := db.conn.Query(context.Background(), query, args...)

	if err != nil {
		fmt.Printf("Select error: %v", err)
		return nil, err
	}

	return rows, nil
}

func (db *DB) Update(table string, columns []string, values []any, where map[string]any) (int64, error) {
	if len(columns) != len(values) {
		return 0, fmt.Errorf("columns and values length mismatch")
	}

	if len(where) == 0 {
		return 0, fmt.Errorf("where is empty")
	}

	query := fmt.Sprintf("UPDATE %s SET ", table)

	args := []any{}

	i := 0
	for _, c := range columns {
		query += fmt.Sprintf("%s = $%d, ", c, i+1)
		i++
	}

	query = query[0 : len(query)-2]

	args = append(args, values...)

	whereClause := ""
	for k, v := range where {
		whereClause += fmt.Sprintf("%s = $%d AND ", k, i+1)
		args = append(args, v)
		i++
	}

	whereClause = whereClause[0 : len(whereClause)-5]

	query += fmt.Sprintf(" WHERE %s", whereClause)

	stmt, err := db.conn.Exec(context.Background(), query, args...)

	if err != nil {
		fmt.Printf("Update error: %v", err)
		return 0, err
	}

	return stmt.RowsAffected(), nil
}

func (db *DB) Delete(table string, where map[string]any) (int64, error) {
	if len(where) == 0 {
		return 0, fmt.Errorf("where is empty")
	}

	query := fmt.Sprintf("DELETE FROM %s ", table)

	args := []any{}

	i := 0
	whereClause := ""
	for k, v := range where {
		whereClause += fmt.Sprintf("%s = $%d AND ", k, i+1)
		args = append(args, v)
		i++
	}

	whereClause = whereClause[0 : len(whereClause)-5]

	query += fmt.Sprintf(" WHERE %s", whereClause)

	stmt, err := db.conn.Exec(context.Background(), query, args...)

	if err != nil {
		fmt.Printf("Delete error: %v", err)
		return 0, err
	}

	return stmt.RowsAffected(), nil
}
