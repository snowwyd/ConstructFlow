package postgresrepo

import (
	"fmt"
	"service-core/pkg/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func New(cfg *config.Config) (*Database, error) {
	const op = "database.postgres.New"

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Database{db: db}, nil
}
