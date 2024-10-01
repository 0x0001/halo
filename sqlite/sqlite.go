package sqlite

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"modernc.org/sqlite"
)

type sqliteDriver struct {
	*sqlite.Driver
}

func (d sqliteDriver) Open(dsn string) (driver.Conn, error) {
	conn, err := d.Driver.Open(dsn)
	if err != nil {
		return conn, err
	}
	c := conn.(interface {
		Exec(stmt string, args []driver.Value) (driver.Result, error)
	})
	if _, err := c.Exec("PRAGMA foreign_keys = on;", nil); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to enable enable foreign keys: %w", err)
	}
	return conn, nil
}

func init() {
	sql.Register("sqlite3", sqliteDriver{&sqlite.Driver{}})
}
