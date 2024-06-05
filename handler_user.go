package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Hultyrhik/rssAggregator/internal/database"
	"github.com/google/uuid"
)

// we want to pass additional argument, so we make a method
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	// create empty struct to decode to
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		// name is taken from parameters struct
		Name: params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could't create user: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(user))
}
