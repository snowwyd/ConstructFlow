package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"service-core/internal/domain"
	"service-core/pkg/config"
	"service-core/pkg/logger"
	"service-core/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Настройка флагов
	resetFlag := flag.Bool("reset", false, "Очистить базу данных")
	migrateFlag := flag.Bool("migrate", false, "Применить миграции")
	seedFlag := flag.Bool("seed", false, "Заполнить тестовыми данными")
	flag.Parse()

	cfg := config.MustLoadEnv()
	log := setupLogger()

	// Подключение к БД
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
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
			&domain.Approval{},
			&domain.Workflow{}, // Явное указание связующих таблиц
		)

		db.Exec("CREATE SEQUENCE IF NOT EXISTS workflows_workflow_id_seq")

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
		"approvals",
		"workflows",
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
		{WorkflowID: 1, UserID: users[0].ID, WorkflowOrder: 1, WorkflowName: "Процедура согласования 1"},
		{WorkflowID: 1, UserID: users[1].ID, WorkflowOrder: 2, WorkflowName: "Процедура согласования 1"},
		{WorkflowID: 1, UserID: users[2].ID, WorkflowOrder: 3, WorkflowName: "Процедура согласования 1"},

		// Workflow 2: user2 (order 1), user2 (order 2), user3 (order 3)
		{WorkflowID: 2, UserID: users[1].ID, WorkflowOrder: 1, WorkflowName: "Процедура согласования 2"},
		{WorkflowID: 2, UserID: users[1].ID, WorkflowOrder: 2, WorkflowName: "Процедура согласования 2"},
		{WorkflowID: 2, UserID: users[2].ID, WorkflowOrder: 3, WorkflowName: "Процедура согласования 2"},

		// Workflow 3: user3 (order 1, 2, 3)
		{WorkflowID: 3, UserID: users[2].ID, WorkflowOrder: 1, WorkflowName: "Тестовая процедура согласования"},
		{WorkflowID: 3, UserID: users[2].ID, WorkflowOrder: 2, WorkflowName: "Тестовая процедура согласования"},
		{WorkflowID: 3, UserID: users[2].ID, WorkflowOrder: 3, WorkflowName: "Тестовая процедура согласования"},
	}
	if err := db.Create(&workflows).Error; err != nil {
		log.Fatalf("Failed to create workflows: %v", err)
	}

	db.Exec("SELECT setval('workflows_workflow_id_seq', (SELECT MAX(workflow_id) FROM workflows))")

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
