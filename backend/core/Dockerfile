FROM golang:1.23.4 AS build

WORKDIR /core

COPY go.mod go.sum ./
RUN go mod download

COPY . .


# RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api/main.go
# RUN CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator/main.go

# FROM debian:latest

# WORKDIR /root/

# COPY --from=build /app/app .
# COPY --from=build /app/migrator .

# COPY configs/ configs

EXPOSE 8080
CMD [ "go", "run", "cmd/api/main.go", "-env=prod"]