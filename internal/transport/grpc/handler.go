package grpc

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/wfen/go-grpc-services-course/internal/rocket"
	rkt "github.com/wfen/go-grpc-services-course/proto/rocket/v1"
)

// RocketService - define the interface that the concrete implementation
// has to adhere to
type RocketService interface {
	GetRocketByID(ctx context.Context, id string) (rocket.Rocket, error)
	InsertRocket(ctx context.Context, rkt rocket.Rocket) (rocket.Rocket, error)
	DeleteRocket(ctx context.Context, id string) error
}

// Handler - will handle incoming gRPC requests
type Handler struct {
	RocketService RocketService
}

// New - returns a new gRPC handler
func New(rktService RocketService) Handler {
	return Handler{
		RocketService: rktService,
	}
}

func (h Handler) Run() error {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Print("could not listen on port 50051")
		return err
	}
	m := cmux.New(l)
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	// Create your protocol servers.
	grpcS := grpc.NewServer()
	rkt.RegisterRocketServiceServer(grpcS, h)
	reflection.Register(grpcS)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gwMux := runtime.NewServeMux()
	if err := rkt.RegisterRocketServiceHandlerServer(ctx, gwMux, h); err != nil {
		log.Printf("error registering service handler: %s", err)
	}
	router := chi.NewRouter()
	statusHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		_, _ = io.WriteString(writer, "OK\n")
	})
	router.Get("/status.html", statusHandler)
	router.Get("/rocket.swagger.json", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(writer, rkt.RocketSwaggerJSON)
	}))
	// GRPC->REST Gateway handles any routes/methods not defined below
	router.MethodNotAllowed(gwMux.ServeHTTP)
	router.NotFound(gwMux.ServeHTTP)

	go grpcS.Serve(grpcL)
	go http.Serve(httpL, router)

	if err := m.Serve(); err != nil {
		log.Printf("failed to serve: %s\n", err)
		return err
	}

	return nil
}

// GetRocket - retrieves a rocket by id and returns the response.
func (h Handler) GetRocket(ctx context.Context, req *rkt.GetRocketRequest) (*rkt.GetRocketResponse, error) {
	log.Print("Get Rocket gRPC Endpoint Hit")

	rocket, err := h.RocketService.GetRocketByID(ctx, req.Id)
	if err != nil {
		log.Print("Failed to retrieve rocket by ID")
		return &rkt.GetRocketResponse{}, err
	}
	return &rkt.GetRocketResponse{
		Rocket: &rkt.Rocket{
			Id:   rocket.ID,
			Name: rocket.Name,
			Type: rocket.Type,
		},
	}, nil
}

func (h Handler) AddRocket(ctx context.Context, req *rkt.AddRocketRequest) (*rkt.AddRocketResponse, error) {
	log.Print("Add Rocket gRPC endpoint hit")

	if _, err := uuid.Parse(req.Rocket.Id); err != nil {
		errorStatus := status.Error(codes.InvalidArgument, "uuid is not valid")
		log.Print("given uuid is not valid")
		return &rkt.AddRocketResponse{}, errorStatus
	}
	newRkt, err := h.RocketService.InsertRocket(ctx, rocket.Rocket{
		ID:   req.Rocket.Id,
		Name: req.Rocket.Name,
		Type: req.Rocket.Type,
	})
	if err != nil {
		log.Print("failed to insert rocket into database")
		return &rkt.AddRocketResponse{}, err
	}
	return &rkt.AddRocketResponse{
		Rocket: &rkt.Rocket{
			Id:   newRkt.ID,
			Name: newRkt.Name,
			Type: newRkt.Type,
		},
	}, nil
}

// DeleteRocket - handler for deleting a rocket
func (h Handler) DeleteRocket(ctx context.Context, req *rkt.DeleteRocketRequest) (*rkt.DeleteRocketResponse, error) {
	log.Print("Delete Rocket gRPC endpoint hit")
	err := h.RocketService.DeleteRocket(ctx, req.Rocket.Id)
	if err != nil {
		return &rkt.DeleteRocketResponse{}, err
	}
	return &rkt.DeleteRocketResponse{
		Status: "successfully deleted rocket",
	}, nil
}
