package main

import (
	"catalog-api/gql"
	"catalog-api/internal/models"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultPort = "8080"

func main() {
	dsn := "host=localhost user=user password=pass dbname=catalog port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}

	db.AutoMigrate(&models.Category{}, &models.Course{})

	if os.Getenv("ENV") == "dev" {
		log.Println("Seeding database...")
		models.SeedDB(db)
	}

	srv := handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{Resolvers: &gql.Resolver{DB: db}}))

	http.Handle("/", playground.Handler("GraphQL", "/query"))
	http.Handle("/query", srv)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	fmt.Printf("\nCatalog API running at http://localhost:%s/\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
