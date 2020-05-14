// +build integration

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lchsk/rss/libs/cache"
	"github.com/lchsk/rss/libs/db"
)

func setupSchema(conn *sql.DB) error {
	log.Println("Setting up schema")
	err := db.InstallSchema(conn, "../sql/schema.sql")

	if err != nil {
		return fmt.Errorf("installschema err %s", err)
	}

	return nil
}

func setupIntegrationTests() error {
	log.Println("Setting up integration tests")
	err := godotenv.Load("../.env")

	if err != nil {
		return fmt.Errorf("cannot find env file: %s", err)
	}

	conn, err := db.GetDBConn(
		os.Getenv("POSTGRES_TEST_HOST"),
		os.Getenv("POSTGRES_TEST_USER"),
		os.Getenv("POSTGRES_TEST_PASSWORD"),
		os.Getenv("POSTGRES_TEST_DB"),
		os.Getenv("POSTGRES_TEST_PORT"))

	if err != nil {
		return fmt.Errorf("getdbconn err: %v", err)
	}

	err = setupSchema(conn)

	if err != nil {
		return fmt.Errorf("unable to set schema up: %s", err)
	}

	dba, err := db.InitDbAccess(conn)

	if err != nil {
		return fmt.Errorf("initdbaccess err: %v", err)
	}

	DBA = dba

	redis, err := cache.GetRedisConn(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	if err != nil {
		return fmt.Errorf("getredisconn err: %v", err)
	}

	cacheAccess, err := cache.InitCacheAccess(redis)

	if err != nil {
		return fmt.Errorf("initcacheaccess err: %v", err)
	}

	Cache = cacheAccess

	return nil
}
