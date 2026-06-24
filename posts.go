package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(wr http.ResponseWriter, req *http.Request) {
	type reqBody struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	post := reqBody{}
	err := decoder.Decode(&post)
	if err != nil {
		// TODO: call return error func
		return
	}

	if len(post.Body) > 140 {
		// TODO: call return error func
	}

	// TODO: call return json func
}
