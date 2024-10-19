package server

import (
	"errors"
	"fmt"
	groupSvc "github.com/upassed/upassed-account-service/internal/service/group"
	studentSvc "github.com/upassed/upassed-account-service/internal/service/student"
	teacherSvc "github.com/upassed/upassed-account-service/internal/service/teacher"
	"log/slog"
	"net"
	"reflect"
	"runtime"

	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"github.com/upassed/upassed-account-service/internal/server/group"
	"github.com/upassed/upassed-account-service/internal/server/student"
	"github.com/upassed/upassed-account-service/internal/server/teacher"
	"google.golang.org/grpc"
)

var (
	errStartingTcpConnection = errors.New("unable to start tcp connection")
	errStartingServer        = errors.New("unable to start gRPC server")
)

type AppServer struct {
	config *config.Config
	log    *slog.Logger
	server *grpc.Server
}

type AppServerCreateParams struct {
	Config         *config.Config
	Log            *slog.Logger
	TeacherService teacherSvc.Service
	StudentService studentSvc.Service
	GroupService   groupSvc.Service
}

func New(params AppServerCreateParams) *AppServer {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDMiddlewareInterceptor(),
			middleware.PanicRecoveryMiddlewareInterceptor(params.Log),
			middleware.LoggingMiddlewareInterceptor(params.Log),
		),
	)

	teacher.Register(server, params.TeacherService)
	student.Register(server, params.StudentService)
	group.Register(server, params.GroupService)

	return &AppServer{
		config: params.Config,
		log:    params.Log,
		server: server,
	}
}

func (server *AppServer) Run() error {
	op := runtime.FuncForPC(reflect.ValueOf(server.Run).Pointer()).Name()

	log := server.log.With(
		slog.String("op", op),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.GrpcServer.Port))
	if err != nil {
		return fmt.Errorf("%s -> %w; %w", op, errStartingTcpConnection, err)
	}

	log.Info("gRPC server is running", slog.String("address", listener.Addr().String()))
	if err := server.server.Serve(listener); err != nil {
		return fmt.Errorf("%s -> %w; %w", op, errStartingServer, err)
	}

	return nil
}

func (server *AppServer) GracefulStop() {
	op := runtime.FuncForPC(reflect.ValueOf(server.GracefulStop).Pointer()).Name()

	log := server.log.With(
		slog.String("op", op),
	)

	log.Info("gracefully stopping gRPC server...")
	server.server.GracefulStop()
}
