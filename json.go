package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// formats error message into consistent JSON object
func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error:", msg)
	}

	// `json:"error"` is a JSON tag that means that key we Marshal is error
	// in go in struct typically json reflect tags are used
	// to specify how json.Marshal or json.Unmarshal converts struct
	// to json object or reversed
	/*
		So the struct below wil lMarshal to JSON like this:
		{
			"error": "something went wrong"
		}
	*/
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

// code is status code of the response
// payload is interface that marshal to JSON struct
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Attempt to marshal from given JSON string and returns it as bytes
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}
	// adds header to response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	// write a respose body
	w.Write(data)
}
