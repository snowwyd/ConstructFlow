# ConstructFlow Backend API Documentation

## Описание проекта
Backend для системы управления файлами и папками с аутентификацией и авторизацией.  
Реализованы следующие методы:
- Авторизация (`/api/v1/auth/login`)
- Регистрация (`/api/v1/auth/register`)
- Получение данных текущего пользователя (`/api/v1/auth/me`)
- Создание роли (`/api/v1/auth/role`)

---

## Требования
- Docker и docker-compose
- Go 1.20+ (для локальной разработки)
- [go-task](https://taskfile.dev/) для выполнения часто используемых команд (опционально)

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

# Откатить контейнеры и очистить БД
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
- `404 Not Found`: Пользователь не найден.
- `500 Internal Server Error`: Ошибка сервера.

---

### 2. Регистрация (`POST /api/v1/auth/register`)
**Запрос:**
```json
{
  "login": "new_user",
  "password": "secure_password",
  "role_id": 809
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
- `404 Not Found`: Роль не найдена.
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

### 4. Создание роли (`POST /api/v1/auth/role`)
**Запрос:**
```json
{
  "role_name": "new_role",
}
```

**Ответ (201 Created):**
```json
{
  "role_id": 123
}
```

**Ошибки:**
- `400 Bad Request`: Некорректный формат запроса.
- `409 Conflict`: Роль с таким названием уже существует.
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

### Создание роли
```bash
curl -X POST http://localhost:8080/auth/role \
-H "Content-Type: application/json" \
-d '{"role_name": "new_role"}'
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
- `UNAUTHORIZED`: Пользователь не прошел аутентификацию.
- `MISSING FIELDS`: Не указаны необходимые поля запроса.
- `INVALID_REQUEST`: Некорректный формат запроса.
- `INVALID_CREDENTIALS`: Неверные логин или пароль.
- `USER_NOT_FOUND`: Пользователь не найден.
- `ROLE_NOT_FOUND`: Роль не найдена.
- `USER_ALREADY_EXISTS`: Пользователь уже существует.
- `ROLE_ALREADY_EXISTS`: Роль уже существует.
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
