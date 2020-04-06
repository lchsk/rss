package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
}

func Refresh(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data RefreshTokenInput
	err := decoder.Decode(&data)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	refreshToken := data.RefreshToken

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_REFRESH_SECRET")), nil
	})

	if err != nil {
		w.WriteHeader(401)
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		w.WriteHeader(401)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		w.WriteHeader(401)
		return
	}

	refreshUuid, ok := claims["refresh_uuid"].(string)

	if !ok {
		w.WriteHeader(422)
		return
	}

	userId, ok := claims["user_id"].(string)

	if !ok {
		w.WriteHeader(422)
		return
	}

	deleted, delErr := DeleteAuth(refreshUuid)

	if delErr != nil || deleted == 0 {
		w.WriteHeader(401)
		return
	}

	ts, createErr := CreateToken(userId)

	if createErr != nil {
		w.WriteHeader(403)
		return
	}

	saveErr := CreateAuth(userId, ts)

	if saveErr != nil {
		w.WriteHeader(403)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(AuthenticationResponse{
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
	})
}
