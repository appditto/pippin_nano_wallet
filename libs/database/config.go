package database

import (
	"fmt"
	"path"

	"entgo.io/ent/dialect"
	"github.com/appditto/pippin_nano_wallet/libs/log"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
)

type SqlDBConn interface {
	DSN() string
	Dialect() string
	Driver() string
}

type PostgresConn struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
}

func (c *PostgresConn) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.DBName)
}

func (c *PostgresConn) Dialect() string {
	return dialect.Postgres
}

func (c *PostgresConn) Driver() string {
	return "pgx"
}

type MysqlConn struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
}

func (c *MysqlConn) DSN() string {
	return fmt.Sprintf("mysql://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.DBName)
}

func (c *MysqlConn) Dialect() string {
	return dialect.MySQL
}

func (c *MysqlConn) Driver() string { return c.Dialect() }

type SqliteConn struct {
	FileName string
	Mode     string
}

func (c *SqliteConn) DSN() string {
	// https://github.com/ent/ent/discussions/1667#discussioncomment-4106910
	return fmt.Sprintf("file:%s?cache=shared&mode=%s&_fk=1&_pragma=foreign_keys(1)", c.FileName, c.Mode)
}

func (c *SqliteConn) Dialect() string {
	return dialect.SQLite
}

func (c *SqliteConn) Driver() string { return "sqlite" }

// Gets the DB connection information based on environment variables
func GetSqlDbConn(mock bool) (SqlDBConn, error) {
	if mock {
		return &SqliteConn{FileName: "testing", Mode: "memory"}, nil
	}
	// First see if postgres is confiugred
	postgresDb := utils.GetEnv("POSTGRES_DB", "")
	postgresUser := utils.GetEnv("POSTGRES_USER", "")
	postgresPassword := utils.GetEnv("POSTGRES_PASSWORD", "")
	postgresHost := utils.GetEnv("POSTGRES_HOST", "127.0.0.1")
	postgresPort := utils.GetEnv("POSTGRES_PORT", "5432")

	if postgresDb != "" && postgresUser != "" && postgresPassword != "" {
		log.Infof("Using PostgreSQL database %s@%s:%s", postgresUser, postgresHost, postgresPort)
		return &PostgresConn{
			Host:     postgresHost,
			Port:     postgresPort,
			Password: postgresPassword,
			User:     postgresUser,
			DBName:   postgresDb,
		}, nil
	}

	// See if MySQL is configured
	mysqlDb := utils.GetEnv("MYSQL_DB", "")
	mysqlUser := utils.GetEnv("MYSQL_USER", "")
	mysqlPassword := utils.GetEnv("MYSQL_PASSWORD", "")
	mysqlHost := utils.GetEnv("MYSQL_HOST", "127.0.0.1")
	mysqlPort := utils.GetEnv("MYSQL_PORT", "3306")

	if mysqlDb != "" && mysqlUser != "" && mysqlPassword != "" {
		log.Infof("Using MySQL database %s@%s:%s", mysqlUser, mysqlHost, mysqlPort)
		return &MysqlConn{
			Host:     mysqlHost,
			Port:     mysqlPort,
			Password: mysqlPassword,
			User:     mysqlUser,
			DBName:   mysqlDb,
		}, nil
	}

	// Default to SQLite
	pippinPath, err := utils.GetPippinConfigurationRoot()
	if err != nil {
		return nil, err
	}
	sqliteDb := path.Join(pippinPath, "pippingo.db")
	log.Infof("Using SQLite database at %s", sqliteDb)
	return &SqliteConn{
		FileName: sqliteDb,
		Mode:     "rwc",
	}, nil
}
