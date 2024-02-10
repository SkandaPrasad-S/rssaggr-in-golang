package main

import (
	// "fmt"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/SkandaPrasad-S/rssaggr/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)
type apiConfig struct{
	DB *database.Queries
}

func main() {
	godotenv.Load()

	portString := os.Getenv("Port")
	if portString == "" {
		log.Fatal("PORT is not found in the env")
	}
	dbURL := os.Getenv("DB_URL")
	if portString == "" {
		log.Fatal("DB_URL is not found in the env")
	}

	conn,err := sql.Open("postgres",dbURL)
	if err!=nil{
		log.Fatal("Cannot connect to database")
	}


	apiCfg := apiConfig{
		DB:database.New(conn),
	}
	router := chi.NewRouter()

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Add your allowed origins here
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)

	v1Router.Get("/err",handlerErr)
	v1Router.Post("/users",apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)
	log.Printf("Starting server here on port %v",portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
