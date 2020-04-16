package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/lchsk/rss/cache"
	"github.com/lchsk/rss/db"
	"github.com/lchsk/rss/libs/api"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// URLs
	registerUserUrl     = "/api/users"
	authenticationUrl   = "/api/authentication"
	refreshTokenUrl     = "/api/authentication/refresh"
	logoutUrl           = "/api/logout"
	fetchCurrentUserUrl = "/api/users/current"

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

var DEBUG bool

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(registerUserUrl, handlerRegisterUser).Methods(http.MethodPost)
	router.HandleFunc(authenticationUrl, handlerAuthentication).Methods(http.MethodPost)
	router.HandleFunc(logoutUrl, checkValidToken(handlerLogout)).Methods(http.MethodPost)
	// router.HandleFunc(refreshTokenUrl, Refresh).Methods(http.MethodPost)
	router.HandleFunc(fetchCurrentUserUrl, checkValidToken(handlerFetchCurrentUser)).Methods(http.MethodGet)

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
	dbHost := os.Getenv("POSTGRES_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort := os.Getenv("POSTGRES_PORT")

	logDebug(fmt.Sprintf(
		"Attempting to connect to Postgres on host=%s user=%s pass=%s name=%s port=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort))

	conn, err := db.GetDBConn(dbHost, dbUser, dbPassword, dbName, dbPort)

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
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	logDebug(fmt.Sprintf("Attempting to connect to Redis on %s:%s", redisHost, redisPort))

	redis, err := cache.GetRedisConn(redisHost, redisPort)

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

func setupEnv() error {
	return godotenv.Load("../.env")
}

func runAPI() {
	router := getRouter()

	apiPort := os.Getenv("API_PORT")

	logDebug("Running in DEBUG mode")
	log.Print(fmt.Sprintf("Running API on port %s", apiPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", apiPort), api.CommonMiddleware(router)))
}

func logDebug(msg string) {
	if DEBUG {
		log.Println("DEBUG", msg)
	}
}

func init() {
	err := setupEnv()

	if err != nil {
		log.Fatal(err)
	}

	DEBUG, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	setupLogging()

	err = setupDB()

	if err != nil {
		log.Fatal(err)
	}

	err = setupCache()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	runAPI()
}
