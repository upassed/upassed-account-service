package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	config "github.com/upassed/upassed-account-service/internal/config/app"
	"github.com/upassed/upassed-account-service/internal/logger"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
)

var (
	ErrorOpeningDbConnection error = errors.New("failed to open connection to a database")
	ErrorPingingDatabase     error = errors.New("failed to ping database")
)

type TeacherRepositoryImpl struct {
	log *slog.Logger
	db  *sql.DB
}

func NewTeacherRepository(config *config.Config, log *slog.Logger) (*TeacherRepositoryImpl, error) {
	const op = "repository.NewTeacherRepository()"

	log = log.With(
		slog.String("op", op),
	)

	log.Info("started connecting to postgres database")
	postgresInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Storage.Host,
		config.Storage.Port,
		config.Storage.User,
		config.Storage.Password,
		config.Storage.DatabaseName,
	)

	db, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		log.Error("error while opening connection to a database", logger.Error(err))
		return nil, fmt.Errorf("%s - %w", op, ErrorOpeningDbConnection)
	}

	if err := db.Ping(); err != nil {
		log.Error("error while pinging a database")
		return nil, fmt.Errorf("%s - %w", op, ErrorPingingDatabase)
	}

	log.Debug("database connection established successfully")
	return &TeacherRepositoryImpl{
		db:  db,
		log: log,
	}, nil
}

func (repository *TeacherRepositoryImpl) Save(context.Context, domain.Teacher) error {
	panic("not implemented!")
}
