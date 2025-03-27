package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"service-file/internal/domain"
	"service-file/pkg/config"
	"service-file/pkg/logger"

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
		err = db.AutoMigrate(
			&domain.Directory{},
			&domain.File{},
			&domain.UserDirectory{},
			&domain.UserFile{},
		)

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
	tables := []string{
		"user_directories", // Связующие таблицы
		"user_files",
		"directories",
		"files",
	}

	// Отключаем проверку внешних ключей
	db.Exec("SET CONSTRAINTS ALL DEFERRED")

	// Удаляем таблицы в обратном порядке (сначала дочерние)
	for _, table := range tables {
		db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
	}

	log.Println("All tables dropped successfully")
}

func seedData(db *gorm.DB) {
	// 4. Дерево директорий
	rootDir := domain.Directory{
		Name:       "ROOT",
		WorkflowID: 1,
		Status:     "draft",
	}
	db.Where("name = ?", rootDir.Name).FirstOrCreate(&rootDir)

	folder1 := domain.Directory{
		Name:         "Folder1",
		ParentPathID: &rootDir.ID,
		WorkflowID:   3,
		Status:       "draft",
	}
	db.Where("name = ?", folder1.Name).FirstOrCreate(&folder1)

	folder2 := domain.Directory{
		Name:         "Folder2",
		ParentPathID: &rootDir.ID,
		WorkflowID:   2,
		Status:       "draft",
	}
	db.Where("name = ?", folder2.Name).FirstOrCreate(&folder2)

	folder3 := domain.Directory{
		Name:         "Folder3",
		ParentPathID: &folder1.ID,
		WorkflowID:   3,
		Status:       "draft",
	}
	db.Where("name = ?", folder3.Name).FirstOrCreate(&folder3)

	// 5. Файлы
	files := []domain.File{
		{Name: "File1", DirectoryID: rootDir.ID, Status: "draft"},
		{Name: "File2", DirectoryID: rootDir.ID, Status: "draft"},
		{Name: "File3", DirectoryID: folder1.ID, Status: "draft"},
		{Name: "File4", DirectoryID: folder2.ID, Status: "draft"},
		{Name: "File5", DirectoryID: folder3.ID, Status: "draft"},
		{Name: "File6", DirectoryID: folder3.ID, Status: "draft"},
	}
	if err := db.Create(&files).Error; err != nil {
		log.Fatalf("Failed to create files: %v", err)
	}

	// 6. Назначение прав доступа
	// User1: Folder3, File3, File5, File6
	db.Create(&domain.UserDirectory{UserID: 1, DirectoryID: folder3.ID}) // Добавляем директорию
	db.Create(&domain.UserFile{UserID: 1, FileID: files[2].ID})          // File3
	db.Create(&domain.UserFile{UserID: 1, FileID: files[4].ID})          // File5
	db.Create(&domain.UserFile{UserID: 1, FileID: files[5].ID})          // File6

	// User2: Folder2, File4
	db.Create(&domain.UserDirectory{UserID: 2, DirectoryID: folder2.ID}) // Добавляем директорию
	db.Create(&domain.UserFile{UserID: 2, FileID: files[3].ID})          // File4

	// User3: все директории и файлы
	directories := []domain.Directory{rootDir, folder1, folder2, folder3}
	for _, dir := range directories {
		db.Create(&domain.UserDirectory{UserID: 3, DirectoryID: dir.ID})
	}

	for _, file := range files {
		db.Create(&domain.UserFile{UserID: 3, FileID: file.ID})
	}
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
