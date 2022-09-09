package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSqlDbConnPostgres(t *testing.T) {
	// Postgres
	os.Setenv("POSTGRES_DB", "pippin")
	os.Setenv("POSTGRES_USER", "user")
	os.Setenv("POSTGRES_PASSWORD", "password")
	defer os.Unsetenv("POSTGRES_DB")
	defer os.Unsetenv("POSTGRES_USER")
	defer os.Unsetenv("POSTGRES_PASSWORD")

	conn, err := GetSqlDbConn(false)
	assert.Nil(t, err)

	assert.Equal(t, "postgres://user:password@127.0.0.1:5432/pippin", conn.DSN())
	assert.Equal(t, "pgx", conn.Dialect())
}

func TestGetSqlDbConnMysql(t *testing.T) {
	// Postgres
	os.Setenv("MYSQL_DB", "pippin")
	os.Setenv("MYSQL_USER", "user")
	os.Setenv("MYSQL_PASSWORD", "password")
	defer os.Unsetenv("MYSQL_DB")
	defer os.Unsetenv("MYSQL_USER")
	defer os.Unsetenv("MYSQL_PASSWORD")

	conn, err := GetSqlDbConn(false)
	assert.Nil(t, err)

	assert.Equal(t, "mysql://user:password@127.0.0.1:3306/pippin", conn.DSN())
	assert.Equal(t, "mysql", conn.Dialect())
}

func TestGetSqlDbConnSqlite(t *testing.T) {
	// Postgres
	os.Setenv("HOME", "/home/user")
	defer os.Unsetenv("HOME")

	conn, err := GetSqlDbConn(false)
	assert.Nil(t, err)

	assert.Equal(t, "file:/home/user/PippinData/pippin.db?cache=shared&mode=rwc&_fk=1", conn.DSN())
	assert.Equal(t, "sqlite3", conn.Dialect())
}

func TestGetSqlDbConnMock(t *testing.T) {
	conn, err := GetSqlDbConn(true)
	assert.Nil(t, err)

	assert.Equal(t, "file:testing?cache=shared&mode=memory&_fk=1", conn.DSN())
	assert.Equal(t, "sqlite3", conn.Dialect())
}
