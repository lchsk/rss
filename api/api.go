package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lchsk/rss/cache"
	"github.com/lchsk/rss/db"
	"github.com/lchsk/rss/libs/api"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// URLs
	registerUserUrl   = "/api/users"
	authenticationUrl = "/api/authentication"
	refreshTokenUrl   = "/api/authentication/refresh"
	logoutUrl         = "/api/logout"
	fetchUserUrl      = "/api/users/{user_id}"

	// API Errors
	errInvalidUsernameLen = "invalid_username_len"
	errInvalidPasswordLen = "invalid_password_len"
	errInvalidEmailLen    = "invalid_email_len"
	errInvalidEmail       = "invalid_email"
	errInvalidInputFormat = "invalid_input_format"
	// Generic database error, such as violated index
	errDbError = "db_error"
)

var DBA *db.DbAccess
var Cache *cache.CacheAccess

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(registerUserUrl, handlerRegisterUser).Methods(http.MethodPost)
	router.HandleFunc(authenticationUrl, handlerAuthentication).Methods(http.MethodPost)
	router.HandleFunc(logoutUrl, checkValidToken(handlerLogout)).Methods(http.MethodPost)
	router.HandleFunc(refreshTokenUrl, Refresh).Methods(http.MethodPost)
	router.HandleFunc(fetchUserUrl, checkValidToken(handlerFetchUser)).Methods(http.MethodGet)

	return router
}

func setupLogging() {
	lumberjackLog := &lumberjack.Logger{
		Filename:   "./logs/api.log",
		MaxSize:    50, // megabytes
		MaxBackups: 10,
		MaxAge:     90,
		Compress:   true,
		LocalTime:  false,
	}
	log.SetOutput(io.MultiWriter(lumberjackLog, os.Stderr))
	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime | log.Ldate)
}

func setupDB() error {
	dbUser := "rss"
	dbPassword := "rss"
	dbName := "rss_db"
	dbPort := "5432"

	conn, err := db.GetDBConn(dbUser, dbPassword, dbName, dbPort)

	if err != nil {
		return fmt.Errorf("error connecting to DB: %s", err)
	}

	dba, err := db.InitDbAccess(conn)

	if err != nil {
		return fmt.Errorf("error initiating db access: %s", err)
	}

	DBA = dba

	log.Print("Connected to DB")

	return nil
}

func setupCache() error {
	redis, err := cache.GetRedisConn()

	if err != nil {
		return fmt.Errorf("error getting redis conn: %s", err)
	}

	cacheAccess, err := cache.InitCacheAccess(redis)

	if err != nil {
		return fmt.Errorf("error initiating cache access: %s", err)
	}

	Cache = cacheAccess

	log.Print("Connected to Cache")

	return nil
}

func runAPI() {
	router := getRouter()

	log.Print(fmt.Sprintf("Running API on port %d", 8000))
	log.Fatal(http.ListenAndServe(":8000", api.CommonMiddleware(router)))
}

func init() {
	setupLogging()

	err := setupDB()

	if err != nil {
		log.Print(err)
		return
	}

	err = setupCache()

	if err != nil {
		log.Print(err)
		return
	}
}

func main() {
	runAPI()
}
