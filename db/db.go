package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/lchsk/rss/channel"
	"github.com/lchsk/rss/user"
)

type DbAccess struct {
	DB      *sql.DB
	User    *user.UserAccess
	Channel *channel.ChannelAccess
}

func GetDBConnection() (*DbAccess, error) {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort := os.Getenv("POSTGRES_PORT")

	log.Printf("Attempting to connect to Postgres on host=%s port=%s\n", dbHost, dbPort)

	conn, err := GetDBConn(dbHost, dbUser, dbPassword, dbName, dbPort)

	if err != nil {
		return nil, fmt.Errorf("error connecting to DB: %s", err)
	}

	dba, err := InitDbAccess(conn)

	if err != nil {
		return nil, fmt.Errorf("error initiating db access: %s", err)
	}

	log.Printf("Connected to Postgres host=%s port=%s\n", dbHost, dbPort)

	return dba, nil
}

func InitDbAccess(db *sql.DB) (*DbAccess, error) {
	ua, err := user.InitUserAccess(db)

	if err != nil {
		return nil, err
	}

	ca, err := channel.InitChannelAccess(db)

	if err != nil {
		return nil, err
	}

	return &DbAccess{DB: db, User: ua, Channel: ca}, nil
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
