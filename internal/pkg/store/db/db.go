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

func connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, err
	}

	return dbPool, nil
}

func NewDB(ctx context.Context, config config.DB) (*DB, error) {
	db := &DB{
		config: config,
		conn:   nil,
	}

	dsn := db.GetDSN()

	conn, err := connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	db.conn = conn

	return db, nil
}

func (db *DB) Connection() *pgxpool.Pool {
	return db.conn
}

func (db *DB) Close() {
	db.conn.Close()
}

func (db *DB) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db.config.User, db.config.Password, db.config.Host, db.config.Port, db.config.Database)
}

func (db *DB) Insert(table string, columns []string, values []any) (pgx.Rows, error) {
	if len(columns) != len(values) {
		return nil, fmt.Errorf("columns and values length mismatch")
	}

	var q strings.Builder

	q.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES ", table, strings.Join(columns, ", ")))
	args := []any{}

	q.WriteString("(")
	format := "$%d, "
	lastIndex := len(values) - 1
	for i := range values {
		if i == lastIndex {
			format = "$%d"
		}

		q.WriteString(fmt.Sprintf(format, i+1))
	}

	q.WriteString(") RETURNING *")

	query := q.String()

	args = append(args, values...)

	rows, err := db.conn.Query(context.Background(), query, args...)

	if err != nil {
		fmt.Printf("Insert error: %v", err)
		return nil, err
	}

	return rows, nil
}

// TODO: Make where clause more flexible (IN, LIKE, OR, etc)
func (db *DB) Select(table string, columns []string, where map[string]any) (pgx.Rows, error) {
	// If columns is empty, use "*" as default
	columnsStr := "*"
	if len(columns) > 0 {
		columnsStr = strings.Join(columns, ", ")
	}

	var q strings.Builder

	args := []any{}

	q.WriteString(fmt.Sprintf("SELECT %s FROM %s", columnsStr, table))

	if len(where) > 0 {
		var w strings.Builder

		q.WriteString(" WHERE ")

		i := 0
		format := "%s = $%d AND "
		lastIndex := len(where) - 1
		for k, v := range where {
			if i == lastIndex {
				format = "%s = $%d"
			}

			w.WriteString(fmt.Sprintf(format, k, i+1))
			args = append(args, v)
			i++
		}

		q.WriteString(w.String())
	}

	query := q.String()

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

	var q strings.Builder

	q.WriteString(fmt.Sprintf("UPDATE %s SET ", table))

	args := []any{}

	i := 0
	format := "%s = $%d, "
	lastIndex := len(values) - 1
	for _, c := range columns {
		if i == lastIndex {
			format = "%s = $%d"
		}

		q.WriteString(fmt.Sprintf(format, c, i+1))
		i++
	}

	args = append(args, values...)

	if len(where) > 0 {
		var w strings.Builder

		q.WriteString(" WHERE ")

		wi := 0
		format := "%s = $%d AND "
		lastIndex := len(where) - 1
		for k, v := range where {
			if wi == lastIndex {
				format = "%s = $%d"
			}

			w.WriteString(fmt.Sprintf(format, k, i+1))
			args = append(args, v)
			wi++
		}

		q.WriteString(w.String())
	}

	query := q.String()

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

	var q strings.Builder

	q.WriteString(fmt.Sprintf("DELETE FROM %s ", table))

	args := []any{}

	if len(where) > 0 {
		var w strings.Builder

		q.WriteString(" WHERE ")

		i := 0
		format := "%s = $%d AND "
		lastIndex := len(where) - 1
		for k, v := range where {
			if i == lastIndex {
				format = "%s = $%d"
			}

			w.WriteString(fmt.Sprintf(format, k, i+1))
			args = append(args, v)
			i++
		}

		q.WriteString(w.String())
	}

	query := q.String()

	stmt, err := db.conn.Exec(context.Background(), query, args...)

	if err != nil {
		fmt.Printf("Delete error: %v", err)
		return 0, err
	}

	return stmt.RowsAffected(), nil
}
