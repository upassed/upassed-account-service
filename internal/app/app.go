package app

import (
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
	Server *server.AppServer
}

func New(config *config.Config, log *slog.Logger) (*App, error) {
	op := runtime.FuncForPC(reflect.ValueOf(New).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
	)

	teacherRepository, err := teacherRepo.New(config, log)
	if err != nil {
		return nil, err
	}

	studentRepository, err := studentRepo.New(config, log)
	if err != nil {
		return nil, err
	}

	groupRepository, err := groupRepo.New(config, log)
	if err != nil {
		return nil, err
	}

	appServer := server.New(server.AppServerCreateParams{
		Config:         config,
		Log:            log,
		TeacherService: teacher.New(config, log, teacherRepository),
		StudentService: student.New(config, log, studentRepository, groupRepository),
		GroupService:   group.New(config, log, groupRepository),
	})

	log.Info("app successfully created")
	return &App{
		Server: appServer,
	}, nil
}
