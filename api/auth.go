package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/lchsk/rss/libs/user"
)

type AccessDetails struct {
	AccessUuid string
	UserId     string
}

const (
	AccessCookieDuration  = time.Minute * 15
	RefreshCookieDuration = time.Minute * 1500

	AccessTokenHeader = "X-AccessToken"
)

func getCookie(name string, value string, duration time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(duration),
	}

}

func checkValidToken(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := TokenValid(r)

		if err != nil {
			refresh, err := r.Cookie("refresh")

			if err != nil || refresh == nil {
				w.WriteHeader(401)
				return
			}

			if err == nil && refresh != nil {
				tokenData, httpStatus := refreshToken(refresh)

				if httpStatus == 200 {
					r.Header.Set(AccessTokenHeader, tokenData.AccessToken)
					http.SetCookie(w, getCookie("token", tokenData.AccessToken, AccessCookieDuration))
					http.SetCookie(w, getCookie("refresh", tokenData.RefreshToken, RefreshCookieDuration))
				} else {
					w.WriteHeader(httpStatus)
					return
				}
			}
		}

		handler(w, r)
	}
}

func ExtractToken(r *http.Request) string {
	headerToken := r.Header.Get(AccessTokenHeader)

	if headerToken != "" {
		return headerToken
	}

	token, err := r.Cookie("token")

	if err != nil || token == nil {
		log.Print(fmt.Sprintf("error extracting token %s", err))
		return ""
	}

	return token.Value
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_ACCESS_SECRET")), nil
		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func DeleteAuth(uuid string) (int64, error) {
	deleted, err := Cache.Redis.Del(uuid).Result()

	if err != nil {
		return 0, err
	}

	return deleted, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId := claims["user_id"].(string)

		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}

	return nil, err
}

func FetchAuth(auth *AccessDetails) (string, error) {
	userId, err := Cache.Redis.Get(auth.AccessUuid).Result()

	if err != nil {
		return "", err
	}

	return userId, nil
}

func CreateAuth(userId string, t *user.TokenData) error {
	at := time.Unix(t.AccessExpiresAt, 0)
	rt := time.Unix(t.RefreshExpiresAt, 0)
	now := time.Now()

	errAccess := Cache.Redis.Set(t.AccessUuid, userId, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := Cache.Redis.Set(t.RefreshUuid, userId, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func CreateToken(userId string) (*user.TokenData, error) {
	accessUuid := uuid.New().String()
	refreshUuid := uuid.New().String()

	t := &user.TokenData{
		AccessExpiresAt:  time.Now().Add(time.Minute * 15).Unix(),
		AccessUuid:       accessUuid,
		RefreshExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		RefreshUuid:      refreshUuid,
	}

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = userId
	atClaims["access_uuid"] = accessUuid
	atClaims["exp"] = t.AccessExpiresAt
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	t.AccessToken, err = at.SignedString([]byte(os.Getenv("API_ACCESS_SECRET")))

	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}

	rtClaims["user_id"] = userId
	rtClaims["refresh_uuid"] = refreshUuid
	rtClaims["exp"] = t.RefreshExpiresAt

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	t.RefreshToken, err = rt.SignedString([]byte(os.Getenv("API_REFRESH_SECRET")))

	if err != nil {
		return nil, err
	}

	return t, nil
}
