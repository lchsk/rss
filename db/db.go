package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/lchsk/rss/user"
)

type DbAccess struct {
	DB   *sql.DB
	User *user.UserAccess
}

var DBA *DbAccess

func InitDbAccess(db *sql.DB) (*DbAccess, error) {
	ua, err := user.InitUserAccess(db)

	if err != nil {
		return nil, err
	}

	return &DbAccess{DB: db, User: ua}, nil
}

func GetDBConn(host, user, password, dbname, port string) (*sql.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func InstallSchema(db *sql.DB, schemaFile string) error {
	f, err := ioutil.ReadFile(schemaFile)

	if err != nil {
		return err
	}

	schema := string(f)

	if err != nil {
		return err
	}

	_, err = db.Exec(schema)

	return err
}
