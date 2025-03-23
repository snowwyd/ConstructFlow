# ConstructFlow Backend API Documentation

## Описание проекта
Backend для системы управления файлами и папками с аутентификацией и авторизацией.  


Методы доступные в swagger:
- Swagger (`/swagger/index.html`)

## Запуск бека через докер
```bash
docker-compose up --build -d
docker-compose down -v 
```

## API Endpoints
**Ошибки:**
- `400 Bad Request`: Некорректный формат запроса.
- `401 Unauthorized`: Неверные логин или пароль.
- `403 Forbidden`: У пользователя нет доступа к созданию директории.
- `404 Not Found`: Пользователь не найден.
- `409 Conflict`: Роль с таким названием уже существует.
- `500 Internal Server Error`: Ошибка сервера.

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

### 2. Регистрация (`POST /auth/register`)
**Запрос:**
```json
{
  "login": "new_user",
  "password": "secure_password",
  "role_id": 2
}
```

**Ответ (201 Created):**

### 3. Получение данных текущего пользователя (`GET /auth/me`)
**Ответ (200 OK):**
```json
{
  "id": 3,
  "login": "user3",
  "role": "admin"
}
```

### 4. Создание роли (`POST/auth/role`)
**Запрос:**
```json
{
  "role_name": "new_role"
}
```

**Ответ (201 Created):**

### 5. Получение файлов на подписание для конкретного пользователя (`GET /file-approvals`)
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

### 6. Отправка файла на согласование (`PUT /file-approvals/{file_id}/approve`)
**Ответ (201 OK):**

### 7. Подписание файла (`PUT /approval/{approval_id}/sign`)
**Параметры пути:**
- `approval_id`: Идентификатор одобрения (например, 1).

**Ответ (204 OK):**

### 8. Отправка файла на доработку с аннотацией (`PUT /file-approvals/{approval_id}/annotate`)
**Параметры пути:**
- `approval_id`: Идентификатор одобрения (например, 1).

**Тело запроса:**
```json
{
  "message": "Комментарий к доработке"
}
```

**Ответ (200 OK):**
```json
{
  "message": "Approval annotated successfully"
}
```

### 9. Завершение согласования (`PUT /approval/{approval_id}/finalize`)

Работает только, когда ставится последняя подпись

**Параметры пути:**
- `approval_id`: Идентификатор одобрения (например, 2).

**Ответ (200 OK):**
```json
{
  "message": "Approval finalized successfully"
}
```

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

## Контакты
По вопросам обращайтесь к **[Дане](https://github.com/snowwyd)** 😏

<img src ="https://media.giphy.com/media/7dHKAiRnGDvbSAbT54/giphy.gif" />
  