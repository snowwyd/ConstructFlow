package main

import (
	"backend/internal/domain"
	"backend/pkg/config"
	"backend/pkg/logger"
	"backend/pkg/utils"
	"flag"
	"fmt"
	"log"
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
		err = db.AutoMigrate(
			&domain.Role{},
			&domain.User{},
			&domain.Directory{},
			&domain.File{},
			&domain.Approval{},
			&domain.Workflow{},
			&domain.UserDirectory{}, // Явное указание связующих таблиц
			&domain.UserFile{},      // Явное указание связующих таблиц
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
		"approvals",
		"workflows",
		"files",
		"directories",
		"users",
		"roles",
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
	// TODO: сделать трех пользовтелей: 2 конструктора и 1 админ

	// TODO: создать workflows:
	// 1: user1, user2, user3
	// 2: user2, user2, user3
	// 3: user3, user3, user3

	// TODO: сделать дерево файлов:
	// ROOT (Folder1, Folder2, File1, File2) WorkflowID = 1
	// Folder1 (Folder3, File3) WorkflowID = 3
	// Folder2 (File4) WorkflowID = 2
	// Folder3 (File5, File6) WorkflowID = 3

	// TODO: выдать права:
	// User1 (Folder3, File3, File5, File6)
	// User2 (Folder2, File4)
	// User3 (Folder1-3, File1-6, ROOT)

	passHash, _ := utils.HashPassword("12345678")

	// 1. Создаем роли
	adminRole := domain.Role{RoleName: "admin"}
	db.Where(adminRole).FirstOrCreate(&adminRole)

	constructorRole := domain.Role{RoleName: "constructor"}
	db.Where(constructorRole).FirstOrCreate(&constructorRole)

	// 2. Создаем пользователей
	users := []domain.User{
		{Login: "user1", PassHash: passHash, RoleID: constructorRole.ID},
		{Login: "user2", PassHash: passHash, RoleID: constructorRole.ID},
		{Login: "user3", PassHash: passHash, RoleID: adminRole.ID},
	}
	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("Failed to create users: %v", err)
	}

	// 3. Создаем workflows
	workflows := []domain.Workflow{
		// Workflow 1: user1 (order 1), user2 (order 2), user3 (order 3)
		{WorkflowID: 1, UserID: users[0].ID, WorkflowOrder: 1},
		{WorkflowID: 1, UserID: users[1].ID, WorkflowOrder: 2},
		{WorkflowID: 1, UserID: users[2].ID, WorkflowOrder: 3},

		// Workflow 2: user2 (order 1), user2 (order 2), user3 (order 3)
		{WorkflowID: 2, UserID: users[1].ID, WorkflowOrder: 1},
		{WorkflowID: 2, UserID: users[1].ID, WorkflowOrder: 2},
		{WorkflowID: 2, UserID: users[2].ID, WorkflowOrder: 3},

		// Workflow 3: user3 (order 1, 2, 3)
		{WorkflowID: 3, UserID: users[2].ID, WorkflowOrder: 1},
		{WorkflowID: 3, UserID: users[2].ID, WorkflowOrder: 2},
		{WorkflowID: 3, UserID: users[2].ID, WorkflowOrder: 3},
	}
	if err := db.Create(&workflows).Error; err != nil {
		log.Fatalf("Failed to create workflows: %v", err)
	}

	// 4. Дерево директорий
	rootDir := domain.Directory{
		Name:       "ROOT",
		WorkflowID: 1,
	}
	db.Where("name = ?", rootDir.Name).FirstOrCreate(&rootDir)

	folder1 := domain.Directory{
		Name:         "Folder1",
		ParentPathID: &rootDir.ID,
		WorkflowID:   3,
	}
	db.Where("name = ?", folder1.Name).FirstOrCreate(&folder1)

	folder2 := domain.Directory{
		Name:         "Folder2",
		ParentPathID: &rootDir.ID,
		WorkflowID:   2,
	}
	db.Where("name = ?", folder2.Name).FirstOrCreate(&folder2)

	folder3 := domain.Directory{
		Name:         "Folder3",
		ParentPathID: &folder1.ID,
		WorkflowID:   3,
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
	db.Model(&users[0]).Association("Directories").Append([]domain.Directory{folder3}) // Добавляем директорию
	db.Model(&users[0]).Association("Files").Append([]domain.File{files[2], files[4], files[5]})

	// User2: Folder2, File4
	db.Model(&users[1]).Association("Directories").Append([]domain.Directory{folder2})
	db.Model(&users[1]).Association("Files").Append([]domain.File{files[3]})

	// User3: все директории и файлы
	directories := []domain.Directory{rootDir, folder1, folder2, folder3}
	db.Model(&users[2]).Association("Directories").Append(directories)
	db.Model(&users[2]).Association("Files").Append(files)
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
