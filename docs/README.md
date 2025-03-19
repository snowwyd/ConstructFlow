# ConstructFlow Backend API Documentation

## Описание проекта
Backend для системы управления файлами и папками с аутентификацией и авторизацией.  
Реализованы следующие методы:
1. Взаимодействие с пользователями
- Авторизация (`/api/v1/auth/login [POST]`)
- Регистрация (`/api/v1/auth/register [POST]`)
- Получение данных текущего пользователя (`/api/v1/auth/me [GET]`)
- Создание роли (`/api/v1/auth/role [POST]`)

2. Взаимодействие с файлами и папками
- Создание директории (`/api/v1/directories/upload [POST]`)
- Удаление директории (`/api/v1/directories [DELETE]`)
- Получение дерева файлов и директорий (`/api/v1/directories [GET]`)

- Создание файла (`/api/v1/files/upload [POST]`)
- Удаление файла (`/api/v1/files [DELETE]`)
- Получение информации о файле (`api/v1/files/:file_id [GET]`)

А также:
- Swagger (`/swagger/index.html`)
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

### 3. Примените миграции для тестирования (опционально)
```bash
# Запуск в директории backend/
# Очищает БД, а затем создает тестовые данные для проверки методов
task seeddb
go run cmd/migrator/main.go -reset -migrate -seed # если отсутствует go-task
```

## API Endpoints

Все запросы и ответы будут совпадать с описанными здесь, если <u>применить миграции</u> для тестирования (пункт 3)

### 1. Авторизация (`POST /api/v1/auth/login`)
**Запрос:**
```json
{
  "login": "snowwy",
  "password": "12345678"
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
  "role_id": 2
}
```

**Ответ (201 Created):**
```json
{
  "user_id": 3
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
  "id": 3,
  "login": "new_user",
  "role": "constructor"
}
```

**Ошибки:**
- `401 Unauthorized`: Токен отсутствует или неверен.
- `404 Not Found`: Пользователь не найден.
- `500 Internal Server Error`: Ошибка сервера.

---

### 4. Создание роли (`POST /api/v1/auth/role`)
**Запрос:**
```json
{
  "role_name": "new_role"
}
```

**Ответ (201 Created):**
```json
{
  "role_id": 3
}
```

**Ошибки:**
- `400 Bad Request`: Некорректный формат запроса.
- `409 Conflict`: Роль с таким названием уже существует.
- `500 Internal Server Error`: Ошибка сервера.

---

### 5. Создание директории (`POST /api/v1/directories/upload`)
**Заголовки:**
```http
Authorization: Bearer <JWT_TOKEN>
```
В заголовок вставьте токен, который был сгененрирован в п.1 по образцу (для snowwy)

**Запрос:**
```json
{
  "parent_path_id": 3,
  "name": "Test project"
}
```

**Ответ (200 OK):**
```json
{
  "id": 5
}
```

**Ошибки:**
- `400 Bad Request`: Некорректный формат запроса.
- `401 Unauthorized`: Токен отсутствует или неверен.
- `403 Forbidden`: У пользователя нет доступа к созданию директории.
- `500 Internal Server Error`: Ошибка сервера.

---

### 6. Удаление директории (`DELETE /api/v1/directories`)
**Заголовки:**
```http
Authorization: Bearer <JWT_TOKEN>
```
В заголовок вставьте токен, который был сгененрирован в п.1 по образцу (для snowwy)

**Запрос:**
```json
{
  "directory_id": 5
}
```

**Ответ (200 OK):**
```json
{
  "message": "Directory deleted successfully"
}
```

**Ошибки:**
- `400 Bad Request`: Некорректный формат запроса.
- `401 Unauthorized`: Токен отсутствует или неверен.
- `403 Forbidden`: У пользователя нет доступа к удалению директории.
- `500 Internal Server Error`: Ошибка сервера.

---

### 7. Получение дерева файлов и директорий (`GET /api/v1/directories`)
**Заголовки:**
```http
Authorization: Bearer <JWT_TOKEN>
```
В заголовок вставьте токен, который был сгененрирован в п.1 по образцу (для snowwy)

**Параметры запроса:**
```json
{
  "is_archive": true
}
```

**Ответ (200 OK):**
```json
{
    "data": [
        {
            "id": 1,
            "name_folder": "ROOT",
            "status": "archive",
            "files": [
                {
                    "id": 1,
                    "name_file": "Archived1.txt",
                    "status": "archive",
                    "directory_id": 1
                }
            ]
        },
        {
            "id": 2,
            "name_folder": "Archived Directory",
            "status": "archive",
            "parent_path_id": 1,
            "files": [
                {
                    "id": 2,
                    "name_file": "Archived2.txt",
                    "status": "archive",
                    "directory_id": 2
                }
            ]
        }
    ]
}
```

**Ошибки:**
- `400 Bad Request`: Некорректный формат запроса.
- `401 Unauthorized`: Токен отсутствует или неверен.
- `500 Internal Server Error`: Ошибка сервера.

---

### 8. Создание файла (`POST /api/v1/files/upload`)
**Заголовки:**
```http
Authorization: Bearer <JWT_TOKEN>
```
В заголовок вставьте токен, который был сгененрирован в п.1 по образцу (для snowwy)

**Запрос:**
```json
{
  "directory_id": 3,
  "name": "New File.txt"
}
```

**Ответ (200 OK):**
```json
{
  "id": 7
}
```

**Ошибки:**
- `400 Bad Request`: Некорректный формат запроса.
- `401 Unauthorized`: Токен отсутствует или неверен.
- `403 Forbidden`: У пользователя нет доступа к созданию файла.
- `500 Internal Server Error`: Ошибка сервера.

---

### 9. Удаление файла (`DELETE /api/v1/files`)
**Заголовки:**
```http
Authorization: Bearer <JWT_TOKEN>
```
В заголовок вставьте токен, который был сгененрирован в п.1 по образцу (для snowwy)

**Запрос:**
```json
{
  "file_id": 7
}
```

**Ответ (200 OK):**
```json
{
  "message": "File deleted successfully"
}
```

**Ошибки:**
- `400 Bad Request`: Некорректный формат запроса.
- `401 Unauthorized`: Токен отсутствует или неверен.
- `403 Forbidden`: У пользователя нет доступа к удалению файла.
- `500 Internal Server Error`: Ошибка сервера.

---

### 10. Получение информации о файле (`GET /api/v1/files/:file_id`)
**Заголовки:**
```http
Authorization: Bearer <JWT_TOKEN>
```
В заголовок вставьте токен, который был сгененрирован в п.1 по образцу (для snowwy)

**Параметры пути:**
- `file_id`: Идентификатор файла (например, 1).

**Ответ (200 OK):**
```json
{
    "id": 1,
    "name_file": "Archived1.txt",
    "status": "archive",
    "directory_id": 1
}
```

**Ошибки:**
- `400 Bad Request`: Некорректный формат запроса.
- `401 Unauthorized`: Токен отсутствует или неверен.
- `403 Forbidden`: У пользователя нет доступа к файлу.
- `404 Not Found`: Файл не найден.
- `500 Internal Server Error`: Ошибка сервера.

---

## Дополнительные возможности

### Swagger UI
Swagger-документация доступна по адресу:
```
/swagger/index.html
```

---


## Структура ошибок
Все ошибки возвращаются в формате (пример для invalid credentials):
```json
{
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid login or password"
  }
}
```

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

<div class="tenor-gif-embed" data-postid="26206051" data-share-method="host" data-aspect-ratio="1.77778" data-width="100%"><a href="https://tenor.com/view/homelander-based-the-boys-homelander-the-boys-facts-gif-26206051">Homelander Based GIF</a>from <a href="https://tenor.com/search/homelander-gifs">Homelander GIFs</a></div> <script type="text/javascript" async src="https://tenor.com/embed.js"></script>