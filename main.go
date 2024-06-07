package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Hultyrhik/rssAggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	// need to import DB driver to our program,
	// _ is placed means we don't need to call anything from it
	_ "github.com/lib/pq"
)

// holds a connection to DB
type apiConfig struct {
	DB *database.Queries
}

func main() {
	// loads .env file
	godotenv.Load(".env")

	// gets PORT value var
	portString := os.Getenv("PORT")
	if portString == "" {
		// exits program immediatly and logs a message
		log.Fatal("Port is not found in the enviroment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		// exits program immediatly and logs a message
		log.Fatal("dbURL is not found in the enviroment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	// convert conn to needed type
	queries := database.New(conn)

	// api Config now can be passed to handler to use DB after this struct
	apiCfg := apiConfig{
		DB: queries,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// new router to mount so that full path /v1/healthz
	// so that there is 2 version if one is broken after some change
	// for REST API
	// It should respond if the server is alive and running
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server staring on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
