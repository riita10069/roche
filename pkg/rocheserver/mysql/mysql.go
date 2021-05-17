package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
)

type (
	Conn struct {
		User     string `required:"true"`
		Password string `required:"true"`
		Protocol string `required:"true"`
		Address  string `required:"true"`
		Schema   string `required:"true"`
	}
)


func GetMySQL() (*sql.DB, error) {
	var conn Conn

	DBMS := "mysql"
	USER := "root"
	PASS := "password"
	PROTOCOL := "tcp"
	ADDRESS := "mysql"
	DBNAME := "test"

	err := envconfig.Process("db", &conn)

	if err == nil {
		USER = conn.User
		PASS = conn.Password
		PROTOCOL = conn.Protocol
		ADDRESS = conn.Address
		DBNAME = conn.Schema
	}

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "(" + ADDRESS + "" + ")" + "/" + DBNAME

	db, err := sql.Open(DBMS, CONNECT+"?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, xerrors.Errorf("Cannot Establish Connection to MySQL: %w")
	}
	
	return db, nil
}
