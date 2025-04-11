# -----------------------------
# Этап сборки Go-приложения
# -----------------------------
FROM golang:1.23.4 AS build

WORKDIR /app

# Копируем зависимости и загружаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект и собираем бинарники
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/service/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator/main.go

# -----------------------------
# Финальный образ (runtime)
# -----------------------------
FROM python:3.11-slim AS runtime

# Устанавливаем рабочий каталог
WORKDIR /app

# Устанавливаем системные зависимости, включая ICU (libicu-dev)
RUN apt-get update && apt-get install -y \
    build-essential \
    libgl1-mesa-glx \
    libicu-dev \
    && rm -rf /var/lib/apt/lists/*

# Создаем виртуальное окружение Python и добавляем его в PATH
RUN python3 -m venv /venv
ENV PATH="/venv/bin:$PATH"

# Включаем режим инвариантной глобализации для .NET
ENV DOTNET_SYSTEM_GLOBALIZATION_INVARIANT=1

# Копируем Go-бинарники из этапа сборки
COPY --from=build /app/app /app/app
COPY --from=build /app/migrator /app/migrator

# Копируем остальные необходимые файлы и каталоги
COPY configs/ /app/configs/
COPY scripts/ /app/scripts/
COPY requirements.txt /app/requirements.txt

# Устанавливаем Python-зависимости (например, aspose-cad)
RUN pip install --no-cache-dir -r requirements.txt

EXPOSE 8080 50051

# Запускаем Go-приложение
CMD ["/app/app"]