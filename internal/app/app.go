package app

import (
	"github.com/upassed/upassed-account-service/internal/caching"
	"github.com/upassed/upassed-account-service/internal/messanging"
	studentRabbit "github.com/upassed/upassed-account-service/internal/messanging/student"
	teacherRabbit "github.com/upassed/upassed-account-service/internal/messanging/teacher"
	"github.com/upassed/upassed-account-service/internal/repository"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
	"reflect"
	"runtime"

	"github.com/upassed/upassed-account-service/internal/config"
	groupRepo "github.com/upassed/upassed-account-service/internal/repository/group"
	studentRepo "github.com/upassed/upassed-account-service/internal/repository/student"
	teacherRepo "github.com/upassed/upassed-account-service/internal/repository/teacher"
	"github.com/upassed/upassed-account-service/internal/server"
	"github.com/upassed/upassed-account-service/internal/service/group"
	"github.com/upassed/upassed-account-service/internal/service/student"
	"github.com/upassed/upassed-account-service/internal/service/teacher"
)

type App struct {
	Server     *server.AppServer
	RabbitConn *rabbitmq.Conn
}

func New(config *config.Config, log *slog.Logger) (*App, error) {
	op := runtime.FuncForPC(reflect.ValueOf(New).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
	)

	db, err := repository.OpenGormDbConnection(config, log)
	if err != nil {
		return nil, err
	}

	redis, err := caching.OpenRedisConnection(config, log)
	if err != nil {
		return nil, err
	}

	teacherRepository := teacherRepo.New(db, redis, config, log)
	studentRepository := studentRepo.New(db, redis, config, log)
	groupRepository := groupRepo.New(db, redis, config, log)

	rabbit, err := messanging.OpenRabbitConnection(config, log)
	if err != nil {
		return nil, err
	}

	studentRabbit.Initialize(rabbit, config, log)
	teacherRabbit.Initialize(rabbit, config, log)

	appServer := server.New(server.AppServerCreateParams{
		Config:         config,
		Log:            log,
		TeacherService: teacher.New(config, log, teacherRepository),
		StudentService: student.New(config, log, studentRepository, groupRepository),
		GroupService:   group.New(config, log, groupRepository),
	})

	log.Info("app successfully created")
	return &App{
		Server:     appServer,
		RabbitConn: rabbit,
	}, nil
}
