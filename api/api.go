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
	"github.com/lchsk/rss/libs/api"
	"github.com/lchsk/rss/libs/cache"
	"github.com/lchsk/rss/libs/db"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// URLs
	registerUserUrl               = "/api/users"
	authenticationUrl             = "/api/authentication"
	refreshTokenUrl               = "/api/authentication/refresh"
	logoutUrl                     = "/api/logout"
	fetchCurrentUserUrl           = "/api/users/current"
	fetchCurrentUserChannelsUrl   = "/api/users/current/channels"
	fetchCurrentUserCategoriesUrl = "/api/users/current/categories"
	changeUserPassword            = "/api/users/current/password"
	addNewChannelUrl              = "/api/channels"
	fetchInboxPostsUrl            = "/api/posts-inbox"
	fetchChannelPostsUrl          = "/api/posts/channels/{channel}"
	fetchCategoryPostsUrl         = "/api/posts/categories/{category}"
	postsUrl                      = "/api/posts/{id}"

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
	router.HandleFunc(fetchCurrentUserChannelsUrl, checkValidToken(handlerFetchCurrentUserChannels)).Methods(http.MethodGet)
	router.HandleFunc(fetchCurrentUserCategoriesUrl, checkValidToken(handlerFetchCurrentUserCategories)).Methods(http.MethodGet)
	router.HandleFunc(changeUserPassword, checkValidToken(handlerChangeUserPassword)).Methods(http.MethodPatch)
	router.HandleFunc(addNewChannelUrl, checkValidToken(handlerAddNewChannelUrl)).Methods(http.MethodPost)
	router.HandleFunc(fetchInboxPostsUrl, checkValidToken(handlerFetchInboxPosts)).Methods(http.MethodGet)

	router.HandleFunc(fetchChannelPostsUrl, checkValidToken(handlerFetchChannelPosts)).Methods(http.MethodGet)
	router.HandleFunc(fetchCategoryPostsUrl, checkValidToken(handlerFetchCategoryPosts)).Methods(http.MethodGet)
	router.HandleFunc(postsUrl, checkValidToken(handlerPosts)).Methods(http.MethodGet, http.MethodPatch)

	if DEBUG {
		const serveTestChannels = "/api/debug/channels/{channel}"
		router.HandleFunc(serveTestChannels, handlerServeTestChannels).Methods(http.MethodGet)
	}

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

func setupCache() error {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	log.Print(fmt.Sprintf("Attempting to connect to Redis on %s:%s", redisHost, redisPort))

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

	log.Print(fmt.Sprintf("DEBUG mode: %v\n", DEBUG))

	if DEBUG {
		log.Println("RSS Channels for debugging")
		for channel, _ := range testChannels {
			log.Printf("http://localhost:%s/api/debug/channels/%s", apiPort, channel)
		}
	}

	log.Print(fmt.Sprintf("Running API on port %s", apiPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", apiPort), api.CommonMiddleware(router)))
}

func init() {
	log.Println("Api init")
	err := setupEnv()

	if err != nil {
		log.Fatal(err)
	}

	DEBUG, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	setupLogging()

	dba, err := db.GetDBConnection()

	if err != nil {
		log.Fatal(err)
	}

	DBA = dba

	err = setupCache()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Api init finished")
}

func main() {
	runAPI()
}
