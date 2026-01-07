package main

import (
	"context"
	"fmt"
	"team-flow/internal/routes"
	"team-flow/pkg/database/postgres"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("env.env")

	ctx := context.Background()

	_, err := postgres.ConnectDB(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connect to PostgresDB")

	r := gin.Default()
	routes.RegisterHealthRoutes(r)
	r.Run(":5050")

}
