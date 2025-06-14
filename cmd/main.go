package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/st-ember/ecommerceprocessor/internal/db"
	"github.com/st-ember/ecommerceprocessor/internal/processor/redis"
	"github.com/st-ember/ecommerceprocessor/internal/processor/worker"
)

func main() {
	fmt.Println("Connecting to Postgres")

	r := chi.NewRouter()

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	db.InitDB(os.Getenv("CONN_STR"))
	defer db.CloseDB()

	fmt.Println("Connecting to Redis")
	redis.Connect()

	http.ListenAndServe(os.Getenv("API_ROOT"), r)

	worker.StartRequestWorker()
}
