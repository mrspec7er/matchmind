package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/mrspec7er/matchmind/server/internal/server"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {

	config := &server.Config{}

	server := server.NewInstance(config)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

}
