package main

import (
	"flag"
	"fmt"
    "log"
	"os"

    "github.com/joho/godotenv"
	"github.com/mattmac4241/grouper-api/service"
)

func main() {
    err := godotenv.Load()
    if err != nil {
      log.Fatal("Error loading .env file")
    }


    dbname := os.Getenv("DBNAME")
    user := os.Getenv("DBUSER")
    password := os.Getenv("DBPASSWORD")
    host := os.Getenv("DBHOST")
	redisAddress := os.Getenv("REDIS_ADDRESS")

	service.REDIS, _ = service.InitRedisClient(redisAddress, "")
	service.DB = service.InitDatabase(host, user, dbname, password)
	defer service.CloseDatabase()
	defer service.REDIS.Close()

	handleFlags()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	publish()
	server := service.NewServer()
	server.Run(":" + port)
}

func publish() {
	serviceClient := service.CatalogWebClient{
        RootURL: os.Getenv("PUBLISH_URL"),
    }
	serviceClient.PublishService()
}


func handleFlags() {
	createPTR := flag.Bool("create", false, "creates the models")
	migratePTR := flag.Bool("migrate", false, "migrates the models")
	deletePTR := flag.Bool("delete", false, "deletes the models")
	flag.Parse()

	if *deletePTR == true {
		fmt.Println("DELETE MODELS")
		service.DropModels()
	}
	if *createPTR == true {
		fmt.Println("CREATE MODELS")
		service.CreateModels()
	}
	if *migratePTR == true {
		fmt.Println("MIGRATE MODELS")
		service.MigrateModels()
	}
}
