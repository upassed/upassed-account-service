package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/google/uuid"
	config "github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"github.com/upassed/upassed-account-service/internal/server/group"
	"github.com/upassed/upassed-account-service/internal/server/student"
	"github.com/upassed/upassed-account-service/internal/server/teacher"
	business "github.com/upassed/upassed-account-service/internal/service/model"
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
	StudentService studentService
	GroupService   groupService
}

type teacherService interface {
	Create(ctx context.Context, teacher business.Teacher) (business.TeacherCreateResponse, error)
	FindByID(ctx context.Context, teacherID uuid.UUID) (business.Teacher, error)
}

type studentService interface {
	Create(context.Context, business.Student) (business.StudentCreateResponse, error)
	FindByID(ctx context.Context, studentID uuid.UUID) (business.Student, error)
}

type groupService interface {
	FindStudentsInGroup(context.Context, uuid.UUID) ([]business.Student, error)
	FindByID(context.Context, uuid.UUID) (business.Group, error)
	FindByFilter(context.Context, business.GroupFilter) ([]business.Group, error)
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
