# ConstructFlow Backend API Documentation

## Описание проекта
Backend для системы управления файлами и папками с аутентификацией и авторизацией.  
Реализованы следующие методы:
- Авторизация (`/auth/login`)
- Регистрация (`/auth/register`)
- Получение данных текущего пользователя (`/auth/me`)

---

## Требования
- Docker и docker-compose
- Go 1.20+ (для локальной разработки)
- [go-task](https://taskfile.dev/) для выполнения задач (опционально)

---

## Локальная разработка

### 1. Установка зависимостей
```bash
# Установите go-task (https://taskfile.dev/install/)
# Для macOS:
brew install go-task

# Для Linux:
curl -sL https://taskfile.dev/install.sh | sh

# Для Windows:
# Скачайте .exe файл с https://github.com/go-task/task/releases
```

### 2. Запуск приложения через Docker
```bash
# Создайте .env файл в директории /backend

# Запустите контейнеры (PostgreSQL + приложение)
docker-compose --build up -d # если отсутствует go-task
task build 

# Откатить контейнеры
docker-compose down # если отсутствует go-task
task composedown
```
Пример .env файла:
```
APP_SECRET = your_app_secret
DB_USER = username
DB_PASSWORD = your_password
CONFIG_PATH = configs/local.yaml
```

## API Endpoints

### 1. Авторизация (`POST /api/v1/auth/login`)
**Запрос:**
```json
{
  "login": "user_login",
  "password": "user_password"
}
```

**Ответ (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ошибки:**
- `400 Bad Request`: Некорректный формат запроса.
- `401 Unauthorized`: Неверные логин или пароль.
- `500 Internal Server Error`: Ошибка сервера.

---

### 2. Регистрация (`POST /api/v1/auth/register`)
**Запрос:**
```json
{
  "login": "new_user",
  "password": "secure_password",
  "role": "user"
}
```

**Ответ (201 Created):**
```json
{
  "user_id": 123
}
```

**Ошибки:**
- `400 Bad Request`: Некорректный формат запроса.
- `409 Conflict`: Пользователь с таким логином уже существует.
- `500 Internal Server Error`: Ошибка сервера.

---

### 3. Получение данных текущего пользователя (`GET /api/v1/auth/me`)
**Заголовки:**
```http
Authorization: Bearer <JWT_TOKEN>
```

**Ответ (200 OK):**
```json
{
  "id": 123,
  "login": "user_login",
  "role": "user"
}
```

**Ошибки:**
- `401 Unauthorized`: Токен отсутствует или неверен.
- `404 Not Found`: Пользователь не найден.
- `500 Internal Server Error`: Ошибка сервера.

---

## Примеры запросов (curl)

### Авторизация
```bash
curl -X POST http://localhost:8080/auth/login \
-H "Content-Type: application/json" \
-d '{"login": "admin", "password": "admin123"}'
```

### Регистрация
```bash
curl -X POST http://localhost:8080/auth/register \
-H "Content-Type: application/json" \
-d '{"login": "new_user", "password": "pass123", "role": "user"}'
```

### Получение данных пользователя
```bash
curl http://localhost:8080/auth/me \
-H "Authorization: Bearer <JWT_TOKEN>"
```

---

## Структура ошибок
Все ошибки возвращаются в формате:
```json
{
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid login or password"
  }
}
```

**Коды ошибок:**
- `INVALID_REQUEST`: Некорректный формат запроса.
- `INVALID_CREDENTIALS`: Неверные логин или пароль.
- `USER_NOT_FOUND`: Пользователь не найден.
- `USER_ALREADY_EXISTS`: Пользователь уже существует.
- `INTERNAL_ERROR`: Внутренняя ошибка сервера.

## Структура проекта
```
/backend
├── cmd/                # Точки входа (HTTP-сервер)
├── internal/           # Бизнес-логика
│   ├── app/            # Инициализация приложения
│   ├── controller/     # HTTP-контроллеры
│   ├── domain/         # Сущности и интерфейсы
│   └── infrastructure/ # Взаимодействие с БД и другими дополнительными модулями
│   └── usecase/        # Бизнес-логика
├── data/               # Миграции и данные
└── pkg/                # Вспомогательные пакеты
```

---

## Контакты
По вопросам обращайтесь к **[Дане](https://github.com/snowwyd)** 😏
