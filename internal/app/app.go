package app

import (
	"github.com/upassed/upassed-account-service/internal/caching"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/messanging"
	studentRabbit "github.com/upassed/upassed-account-service/internal/messanging/student"
	teacherRabbit "github.com/upassed/upassed-account-service/internal/messanging/teacher"
	"github.com/upassed/upassed-account-service/internal/repository"
	groupRepo "github.com/upassed/upassed-account-service/internal/repository/group"
	studentRepo "github.com/upassed/upassed-account-service/internal/repository/student"
	teacherRepo "github.com/upassed/upassed-account-service/internal/repository/teacher"
	"github.com/upassed/upassed-account-service/internal/server"
	"github.com/upassed/upassed-account-service/internal/service/group"
	"github.com/upassed/upassed-account-service/internal/service/student"
	"github.com/upassed/upassed-account-service/internal/service/teacher"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
)

type App struct {
	Server     *server.AppServer
	RabbitConn *rabbitmq.Conn
}

func New(config *config.Config, log *slog.Logger) (*App, error) {
	log = logging.Wrap(log, logging.WithOp(New))

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

	studentService := student.New(config, log, studentRepository, groupRepository)
	teacherService := teacher.New(config, log, teacherRepository)

	studentRabbit.Initialize(studentService, rabbit, config, log)
	teacherRabbit.Initialize(teacherService, rabbit, config, log)

	appServer := server.New(server.AppServerCreateParams{
		Config:         config,
		Log:            log,
		TeacherService: teacherService,
		StudentService: studentService,
		GroupService:   group.New(config, log, groupRepository),
	})

	log.Info("app successfully created")
	return &App{
		Server:     appServer,
		RabbitConn: rabbit,
	}, nil
}
