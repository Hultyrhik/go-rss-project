package main

import (
	"fmt"
	"net/http"

	"github.com/Hultyrhik/rssAggregator/internal/auth"
	"github.com/Hultyrhik/rssAggregator/internal/database"
)

// looks almost like http handler, but has User
// But this signature doesnt match http handler func because of User
type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// Works on apiConfig so it has access to DB
// This func takes authedHandler and returns HadlerFunc that can be used with Chi router
func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	// return a closure (anonymous function)
	// with the same func signature as http handler
	// but we will have access to apiConfig
	// so we can query the DB
	return func(w http.ResponseWriter, r *http.Request) {
		// get api key from the request
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		// grab the user with that API key
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		// now we should run the handler we are given
		handler(w, r, user)
	}
}
