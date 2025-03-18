package postgresrepo

import (
	"backend/internal/domain"
	"backend/pkg/config"
	"fmt"

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

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := DropAllTables(db); err != nil {
		return nil, fmt.Errorf("%s: failed to drop tables: %w", op, err)
	}

	// Применение миграций через golang-migrate
	if err := applyMigrations(db); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Database{db: db}, nil
}

// applyMigrations выполняет автоматические миграции для всех моделей
func applyMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.Directory{},
		&domain.File{},
		&domain.UserDirectory{},
		&domain.UserFile{},
		&domain.Approval{},
		&domain.Annotation{},
		&domain.Workflow{},
	)
}

// TODO: разобраться с миграциями и инициализацией БД
// DropAllTables очищает базу данных
func DropAllTables(db *gorm.DB) error {
	tables := []string{"users", "roles", "folders", "files", "directories", "approvals", "annotations", "workflows"}
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}
	return nil
}

// GetDB возвращает экземпляр GORM DB
func (d *Database) GetDB() *gorm.DB {
	return d.db
}
