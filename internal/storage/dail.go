package storage

import (
	"context"
	"database/sql"
	"github.com/LasseJacobs/go-starter-kit/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"strings"
)

type Database struct {
	*sqlx.DB
}

type Transaction struct {
	*sqlx.Tx
}

// Connection is the interface a storage provider must implement.
// wrapper around sqlx to allow Databases and transaction to be used interchangably
type Connection interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	DriverName() string
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	Transaction(fn func(tx Connection) error) error
	Close() error
}

func Dial(config *config.Database) (Connection, error) {
	db, err := sqlx.Open(config.Driver, strings.Replace(config.URL, "$password", config.Password, 1))
	if err != nil {
		return nil, err
	}
	//err = bootstrap(db)
	return &Database{db}, err
}

func (db *Database) Transaction(fn func(tx Connection) error) error {
	var dberr error
	cn, err := db.Beginx()
	if err != nil {
		return err
	}
	err = fn(Transaction{cn})
	if err != nil {
		dberr = cn.Rollback()
	} else {
		dberr = cn.Commit()
	}

	if dberr != nil {
		return errors.Wrap(dberr, "error committing or rolling back transaction")
	}

	return err
}

func (tx Transaction) Transaction(fn func(tx Connection) error) error {
	return fn(tx)
}
func (tx Transaction) Close() error {
	//in a transaction do nothing
	return nil
}

/*
func bootstrap(database *sqlx.DB) error {
	//todo: have to be able to bootstrap for tests
	var bootstrapFile = "storage/schema.sql"
	if _, err := os.Stat(bootstrapFile); errors.Is(err, os.ErrNotExist) {
		// for testing the relative location is different
		bootstrapFile = "../storage/schema.sql"
	}
	dat, err := os.ReadFile(bootstrapFile)
	if err != nil {
		return err
	}
	_, err = database.Exec(string(dat))
	return err
}
*/
