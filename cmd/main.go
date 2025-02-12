package main

import (
	"log"
	"os"

	"github.com/codepnw/react_go_ecom/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const envFile = "dev.env"

func init() {
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("loading file .env failed: %v", err)
	}
}

func main() {
	dbConfig := database.DBConfig{
		DBAddr:       os.Getenv("DB_URL"),
		MaxOpenConns: 10,
		MaxIdleConns: 10,
		MaxIdleTime:  "15m",
	}

	db, err := database.InitDB(dbConfig)
	if err != nil {
		log.Fatalf("cant connect to database: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	apiRoutes(db, r, "v1")

	port := os.Getenv("APP_PORT")
	log.Println("server is running at port", port)
	r.Run(":" + port)
}
