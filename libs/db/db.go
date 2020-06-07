package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	sq "github.com/Masterminds/squirrel"

	"github.com/lchsk/rss/libs/channel"
	"github.com/lchsk/rss/libs/posts"
	"github.com/lchsk/rss/libs/user"
)

type DbAccess struct {
	DB      *sql.DB
	SQ      *sq.StatementBuilderType
	User    *user.UserAccess
	Channel *channel.ChannelAccess
	Posts   *posts.PostsAccess
}

func GetDBConnection() (*DbAccess, error) {
	dbHostVar := "POSTGRES_HOST"
	dbUserVar := "POSTGRES_USER"
	dbPasswordVar := "POSTGRES_PASSWORD"
	dbNameVar := "POSTGRES_DB"
	dbPortVar := "POSTGRES_PORT"

	if os.Getenv("INTEGRATION_TEST") != "" {
		dbHostVar = "POSTGRES_TEST_HOST"
		dbUserVar = "POSTGRES_TEST_USER"
		dbPasswordVar = "POSTGRES_TEST_PASSWORD"
		dbNameVar = "POSTGRES_TEST_DB"
		dbPortVar = "POSTGRES_TEST_PORT"
	}

	dbHost := os.Getenv(dbHostVar)
	dbUser := os.Getenv(dbUserVar)
	dbPassword := os.Getenv(dbPasswordVar)
	dbName := os.Getenv(dbNameVar)
	dbPort := os.Getenv(dbPortVar)

	log.Printf("Attempting to connect to Postgres on host=%s db=%s port=%s\n", dbHost, dbName, dbPort)

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
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	ua, err := user.InitUserAccess(db)

	if err != nil {
		return nil, err
	}

	ca, err := channel.InitChannelAccess(db, &psql)

	if err != nil {
		return nil, err
	}

	pa, err := posts.InitPostsAccess(db, &psql)

	if err != nil {
		return nil, err
	}

	return &DbAccess{DB: db, SQ: &psql, User: ua, Channel: ca, Posts: pa}, nil
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
