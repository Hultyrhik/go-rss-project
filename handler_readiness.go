package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// respond with empty struct to marshal empty json object
	respondWithJSON(w, 200, struct{}{})
}
