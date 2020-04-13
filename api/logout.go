package main

import "net/http"

func handlerLogout(w http.ResponseWriter, req *http.Request) {
	tokenAuth, errToken := ExtractTokenMetadata(req)
	if errToken != nil {
		w.WriteHeader(401)
		return
	}

	deleted, delErr := DeleteAuth(tokenAuth.AccessUuid)

	if delErr != nil || deleted == 0 {
		w.WriteHeader(422)
		return
	}

	http.SetCookie(w, getCookie("token", "", -1))
	http.SetCookie(w, getCookie("refresh", "", -1))

	w.WriteHeader(200)
}
