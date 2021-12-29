package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/config"
	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/graph"
	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/graph/generated"
	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/repository"
	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/services"
	"github.com/go-chi/chi"

	"github.com/joho/godotenv"
)

func main() {
	// Check .env.dev file is existing
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Fatal("failed to load env vars ", err)
	}

	// Create a database connection
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}
	defer config.CloseDatabase(db)

	//init routers
	r := initRoutes(db)

	// Start server
	log.Printf("connect to http://localhost:8080/ for GraphQL playground")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Server error %v", err)
	}
}

func initRoutes(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	dbRepo := repository.NewDBRepo(db)
	friendService := services.NewFriendService(dbRepo)

	//GraphQL
	graphqlServer := handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Service: friendService,
	}}))

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", graphqlServer)
	return r
}
