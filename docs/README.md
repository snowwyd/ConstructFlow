# ConstructFlow Backend API Documentation

## Описание проекта

Backend для системы управления файлами и папками с аутентификацией и авторизацией.  

## Обзор ConstructFlow API
В ролке показан основной функционал ConstructFlow API и рассмотрены некоторые пользовательские сценарии при работе сервиса.

<iframe width="560" height="315" src="https://www.youtube.com/embed/q1TcMKiZBGE?si=C2wT-CH1rlt20Wz4" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>


## Требования

- Docker и docker-compose
- Go 1.20+ (для локальной разработки)


## Локальная разработка

### 1. Запуск приложения через Docker

```bash
# Создайте .env файл в директории /backend

# Запустите контейнеры (PostgreSQL + приложение)
docker-compose up --build -d

# Откатить контейнеры и очистить БД
docker-compose down -v
```

Пример .env файла:

```
APP_SECRET = your_app_secret
DB_USER = username
DB_PASSWORD = your_password
CONFIG_PATH = configs/local.yaml

APP_PORT=8080
```

configs/local.yaml обязательно такое значение, остальные на ваше усмотрение

### 2. Примените миграции для тестирования

```bash
# Запуск в директории backend/
# Очищает БД, а затем создает тестовые данные для проверки методов
docker-compose run --rm migrator
```

## API Endpoints

Все запросы и ответы будут совпадать с описанными здесь, если <u>применить миграции</u> для тестирования (пункт 3)

### 1. Авторизация (`POST /auth/login`)

**Запрос:**

```json
{
	"login": "user3",
	"password": "12345678"
}
```

**Ответ (200 OK):**

```json
{
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
---

### 2. Регистрация (`POST /auth/register`)

**Запрос:**

```json
{
	"login": "new_user",
	"password": "secure_password",
	"role_id": 2
}
```

**Ответ (201 Created)**

---

### 3. Получение данных текущего пользователя (`GET /auth/me`)

**Заголовки:**

```http
Authorization: Bearer <JWT_TOKEN>
```

**Ответ (200 OK):**

```json
{
	"id": 3,
	"login": "user3",
	"role": "admin"
}
```
---

### 4. Создание роли (`POST /auth/role`)

**Запрос:**

```json
{
	"role_name": "new_role"
}
```

**Ответ (201 Created)**

---

### 5. Создание директории (`POST /directories/create`)

**Заголовки:**

```http
Authorization: Bearer <JWT_TOKEN>
```

**Запрос:**

```json
{
	"parent_path_id": 1,
	"name": "Test project"
}
```

**Ответ (201 Created)**

---

### 6. Удаление директории (`DELETE /directories`)

**Заголовки:**

```http
Authorization: Bearer <JWT_TOKEN>
```

**Запрос:**

```json
{
	"directory_id": 5
}
```

**Ответ (204 No content)**

---

### 7. Получение дерева файлов и директорий (`POST /directories`)

**Заголовки:**

```http
Authorization: Bearer <JWT_TOKEN>
```

**Параметры запроса:**

```json
{
	"is_archive": false
}
```

**Ответ (200 OK):**

```json
{
	"data": [
		{
			"id": 1,
			"name_folder": "ROOT",
			"status": "",
			"files": [
				{
					"id": 1,
					"name_file": "File1",
					"status": "draft",
					"directory_id": 1
				},
				{
					"id": 2,
					"name_file": "File2",
					"status": "draft",
					"directory_id": 1
				}
			]
		},
		{
			"id": 2,
			"name_folder": "Folder1",
			"status": "",
			"parent_path_id": 1,
			"files": [
				{
					"id": 3,
					"name_file": "File3",
					"status": "draft",
					"directory_id": 2
				}
			]
		},
		{
			"id": 3,
			"name_folder": "Folder2",
			"status": "",
			"parent_path_id": 1,
			"files": [
				{
					"id": 4,
					"name_file": "File4",
					"status": "draft",
					"directory_id": 3
				}
			]
		},
		{
			"id": 4,
			"name_folder": "Folder3",
			"status": "",
			"parent_path_id": 2,
			"files": [
				{
					"id": 5,
					"name_file": "File5",
					"status": "draft",
					"directory_id": 4
				},
				{
					"id": 6,
					"name_file": "File6",
					"status": "draft",
					"directory_id": 4
				}
			]
		}
	]
}
```
---

### 8. Создание файла (`POST /files/upload`)

**Заголовки:**

```http
Authorization: Bearer <JWT_TOKEN>
```

**Запрос:**

```json
{
	"directory_id": 1,
	"name": "New File.txt"
}
```

**Ответ (201 Created)**

---

### 9. Удаление файла (`DELETE /files`)

**Заголовки:**

```http
Authorization: Bearer <JWT_TOKEN>
```

**Запрос:**

```json
{
	"file_id": 7
}
```

**Ответ (204 No Content)**

---

### 10. Получение информации о файле (`GET /files/:file_id`)

**Заголовки:**

```http
Authorization: Bearer <JWT_TOKEN>
```

**Параметры пути:**

- `file_id`: Идентификатор файла (например, 1).

**Ответ (200 OK):**

```json
{
	"id": 1,
	"name_file": "File1",
	"status": "draft",
	"directory_id": 1
}
```
---

### 11. Отправка файла на согласование (`PUT /files/:file_id/approve`)

**Заголовки:**

```http
Authorization: Bearer <JWT_TOKEN>
```

**Параметры пути:**

- `file_id`: Идентификатор файла (например, 3).

**Ответ (201 Created)**

---

### 12. Получение файлов на подписание для конкретного пользователя (`GET /file-approvals`)

**Заголовки:**

```http
Authorization: Bearer <JWT_TOKEN>
```

**Ответ (200 OK):**

```json
[
	{
		"id": 1,
		"file_id": 3,
		"file_name": "File3",
		"status": "on approval",
		"workflow_order": 1
	}
]
```
---

### 13. Подписание файла (`PUT /file-approvals/:approval_id/sign`)

**Заголовки:**

```http
Authorization: Bearer <JWT_TOKEN>
```

**Параметры пути:**

- `approval_id`: Идентификатор процесса процедуры согласования (например, 1). Этот id можно получить из п.12

**Ответ (204 No Content)**

---

### 14. Отправка файла на доработку с аннотацией (`PUT /file-approvals/:approval_id/annotate`)

**Заголовки:**

```http
Authorization: Bearer <JWT_TOKEN>
```

**Параметры пути:**

- `approval_id`: взят =1 для примера

**Тело запроса:**

```json
{
	"message": "some message"
}
```

**Ответ (204 No Content)**

---

### 15. Завершение согласования (`PUT /file-approvals/:approval_id/finalize`)

Работает только, когда ставится последняя подпись

**Заголовки:**

```http
Authorization: Bearer <JWT_TOKEN>
```

**Параметры пути:**

- `approval_id`

**Ответ (204 No Content)**

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

<img src ="https://media.giphy.com/media/7dHKAiRnGDvbSAbT54/giphy.gif" />
