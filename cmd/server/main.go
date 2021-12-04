package main

import (
	"log"

	"github.com/wfen/go-grpc-services-course/internal/db"
	"github.com/wfen/go-grpc-services-course/internal/rocket"
)

func Run() error {
	// responsible for initializing and starting
	// our gRPC server
	rocketStore, err := db.New()
	if err != nil {
		return err
	}
	err = rocketStore.Migrate()
	if err != nil {
		log.Println("failed to run migrations")
		return err
	}

	_ = rocket.New(rocketStore)
	return nil
}
func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
