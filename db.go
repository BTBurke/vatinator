package vatinator

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/BTBurke/vatinator/bundled"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type DB struct {
	mu   sync.RWMutex
	conn *sqlx.DB
}

func NewDB(path string) (*DB, error) {
	conn, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}
	return &DB{
		conn: conn,
	}, nil
}

func (d *DB) Get(dest interface{}, query string, args ...interface{}) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.conn.Get(dest, query, args...)
}

func (d *DB) Select(dest interface{}, query string, args ...interface{}) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.conn.Select(dest, query, args...)
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.conn.Exec(query, args...)
}

func (d *DB) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.conn.Close()
}

func Migrate(db *DB, migrations ...string) error {
	if db == nil {
		return fmt.Errorf("no database provided for migration")
	}
	for _, migration := range migrations {
		stmt, err := bundled.Asset("assets/" + migration)
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(stmt)); err != nil {
			return errors.Wrapf(err, "migration failed for %s", migration)
		}
	}
	return nil
}
