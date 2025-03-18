package main

import (
	"backend/internal/domain"
	"backend/pkg/config"
	"backend/pkg/logger"
	"backend/pkg/utils"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Настройка флагов
	resetFlag := flag.Bool("reset", false, "Очистить базу данных")
	migrateFlag := flag.Bool("migrate", false, "Применить миграции")
	seedFlag := flag.Bool("seed", false, "Заполнить тестовыми данными")
	flag.Parse()

	cfg := config.MustLoad()
	log := setupLogger()

	// Подключение к БД
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Error("failed to open database", slog.String("error", err.Error()))
		return
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Логика выполнения команд
	if *resetFlag {
		resetDatabase(db)
		log.Info("database cleared successfully")
	}

	if *migrateFlag {
		err = db.AutoMigrate(&domain.Role{}, &domain.User{}, &domain.Directory{}, &domain.File{})
		if err != nil {
			log.Error("failed to auto migrate", slog.String("error", err.Error()))
			return
		}
		log.Info("migrations applied successfully")	
	}

	if *seedFlag {
		seedData(db)
		log.Info("seed data inserted successfully")
	}
}

func resetDatabase(db *gorm.DB) {
	// Отключаем проверку внешних ключей
	db.Exec("SET CONSTRAINTS ALL DEFERRED")

	// Очистка таблиц в правильном порядке
	tables := []interface{}{
		&domain.User{},
		&domain.Role{},
		&domain.File{},
		&domain.Directory{},
		&domain.UserDirectory{},
		&domain.UserFile{},
	}

	for _, table := range tables {
		db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(table)
	}

	// Сброс автоинкрементных счетчиков
	db.Exec("TRUNCATE TABLE user_directories, user_files, directories, files, users, roles RESTART IDENTITY CASCADE")
}

func seedData(db *gorm.DB) {
	passHash, _ := utils.HashPassword("12345678")

	// Создаем роли
	adminRole := domain.Role{RoleName: "admin"}
	db.Where("role_name = ?", "admin").FirstOrCreate(&adminRole)

	constructorRole := domain.Role{RoleName: "constructor"}
	db.Where("role_name = ?", "constructor").FirstOrCreate(&constructorRole)

	// Создаем пользователей
	user1 := domain.User{
		Login:    "snowwy",
		PassHash: passHash,
		RoleID:   adminRole.ID,
	}
	db.Where("login = ?", "snowwy").FirstOrCreate(&user1)

	user2 := domain.User{
		Login:    "nubik_snowwy",
		PassHash: passHash,
		RoleID:   constructorRole.ID,
	}
	db.Where("login = ?", "nubik_snowwy").FirstOrCreate(&user2)

	// Создаем директории
	rootDir := domain.Directory{Name: "ROOT"}
	db.Where("name = ?", "ROOT").Attrs(domain.Directory{Status: "archive"}).FirstOrCreate(&rootDir)

	archivedDir := domain.Directory{Name: "Archived Directory"}
	db.Where("name = ?", "Archived Directory").Attrs(domain.Directory{
		ParentPathID: &rootDir.ID,
		Status:       "archive",
	}).FirstOrCreate(&archivedDir)

	wipDir1 := domain.Directory{Name: "WIP Directory 1"}
	db.Where("name = ?", "WIP Directory").Attrs(domain.Directory{
		ParentPathID: &rootDir.ID,
		Status:       "wip",
	}).FirstOrCreate(&wipDir1)

	wipDir2 := domain.Directory{Name: "WIP Directory 2"}
	db.Where("name = ?", "WIP Directory").Attrs(domain.Directory{
		ParentPathID: &rootDir.ID,
		Status:       "wip",
	}).FirstOrCreate(&wipDir2)

	// Создаем файлы
	files := []domain.File{
		{Name: "Archived1.txt", DirectoryID: rootDir.ID, Status: "archive"},
		{Name: "Archived2.txt", DirectoryID: archivedDir.ID, Status: "archive"},
		{Name: "809.txt", DirectoryID: rootDir.ID, Status: "wip"},
		{Name: "810.txt", DirectoryID: rootDir.ID, Status: "wip"},
		{Name: "811.txt", DirectoryID: wipDir1.ID, Status: "wip"},
		{Name: "812.txt", DirectoryID: wipDir2.ID, Status: "wip"},
	}

	for i := range files {
		db.Where("name = ?", files[i].Name).Attrs(files[i]).FirstOrCreate(&files[i])
	}

	// Связываем пользователей с директориями и файлами
	db.Model(&user1).Association("Directories").Replace([]domain.Directory{rootDir, archivedDir, wipDir1})
	db.Model(&user1).Association("Files").Replace(files[:3])

	db.Model(&user2).Association("Directories").Replace([]domain.Directory{wipDir2})
	db.Model(&user2).Association("Files").Replace([]domain.File{files[3], files[4]})
}

func setupLogger() *slog.Logger {
	opts := logger.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
