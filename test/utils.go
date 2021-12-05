// +acceptance

package test

import (
	"log"

	"google.golang.org/grpc"

	rkt "github.com/wfen/go-grpc-services-course/proto/rocket/v1"
)

func GetClient() rkt.RocketServiceClient {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}

	rocketClient := rkt.NewRocketServiceClient(conn)
	return rocketClient
}
