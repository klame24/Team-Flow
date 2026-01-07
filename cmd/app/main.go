package main

import (
	"context"
	"fmt"
	"team-flow/pkg/database/postgres"

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
}
