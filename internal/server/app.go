package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net"

	config "github.com/upassed/upassed-account-service/internal/config/app"
	"google.golang.org/grpc"
)

var (
	ErroStartingTcpConnection error = errors.New("unable to start tcp connection")
	ErrStartingServer         error = errors.New("unable to start gRPC server")
)

type AppServer struct {
	config *config.Config
	log    *slog.Logger
	server *grpc.Server
}

type AppServerCreateParams struct {
	Config         *config.Config
	Log            *slog.Logger
	TeacherService teacherService
}

func New(params AppServerCreateParams) *AppServer {
	server := grpc.NewServer()
	registerTeacherServer(server, params.TeacherService)

	return &AppServer{
		config: params.Config,
		log:    params.Log,
		server: server,
	}
}

func (server *AppServer) Run() error {
	const op = "server.Run()"

	log := server.log.With(
		slog.String("op", op),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.GrpcServer.Port))
	if err != nil {
		return fmt.Errorf("%s -> %w; %w", op, ErroStartingTcpConnection, err)
	}

	log.Info("gRPC server is running", slog.String("address", listener.Addr().String()))
	if err := server.server.Serve(listener); err != nil {
		return fmt.Errorf("%s -> %w; %w", op, ErrStartingServer, err)
	}

	return nil
}

func (server *AppServer) GracefulStop() {
	const op = "server.GracefulStop()"

	log := server.log.With(
		slog.String("op", op),
	)

	log.Info("gracefully stopping gRPC server...")
	server.server.GracefulStop()
}
