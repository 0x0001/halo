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
	s.Assertions.NoError(conn.Ping())
}
