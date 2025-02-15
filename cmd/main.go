package main

import (
	"log"

	"github.com/codepnw/react_go_ecom/config"
	"github.com/codepnw/react_go_ecom/pkg/database"
)

const envFile = "dev.env"

func main() {
	cfg := config.LoadConfig(envFile)

	db, err := database.InitDB(cfg.DBConfig)
	if err != nil {
		log.Fatalf("cant connect to database: %v", err)
	}
	defer db.Close()

	// API Routes
	r := apiRoutes(db, *cfg)

	port := cfg.AppConfig.AppPort
	log.Println("server is running at port", port)
	r.Run(":" + port)
}
