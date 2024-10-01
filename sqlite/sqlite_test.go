package sqlite

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
)

type sqliteTestSuite struct {
	suite.Suite
}

func TestSqliteTestSuite(t *testing.T) {
	suite.Run(t, new(sqliteTestSuite))
}

func (s *sqliteTestSuite) TestSqlite() {
	conn, err := sql.Open("sqlite3", ":memory:")
	s.Assertions.NoError(err)
	defer conn.Close()
	s.Assertions.NoError(conn.Ping())

	cases := []struct {
		name string
		stmt string
		args []interface{}
	}{
		{
			name: "create table",
			stmt: `CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)`,
			args: nil,
		},
		{
			name: "insert data",
			stmt: `INSERT INTO test (name) VALUES (?)`,
			args: []interface{}{"test1"},
		},
		{
			name: "select data",
			stmt: `SELECT * FROM test`,
			args: nil,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			_, err := conn.Exec(c.stmt, c.args...)
			s.Assertions.NoError(err, "exec failed")
		})
	}
}
