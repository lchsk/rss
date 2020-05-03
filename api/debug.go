package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const TestChannelsDir = "./test_channels"

var testChannels = map[string]string{
	"lchsk_1":      "lchsk_1.xml",
	"538_features": "538_features.xml",
	"538_nate":     "538_nate.xml",
	"538_politics": "538_politics.xml",
	"nyt_politics": "nyt_politics.xml",
	"nyt_science":  "nyt_science.xml",
}

func handlerServeTestChannels(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	channel := params["channel"]

	channelPath, ok := testChannels[channel]

	if !ok {
		w.WriteHeader(404)
		return
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", TestChannelsDir, channelPath))

	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	w.WriteHeader(200)
	w.Write(data)
}
