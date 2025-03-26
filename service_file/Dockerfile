FROM golang:1.23.4 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

FROM debian:latest

WORKDIR /root/

COPY --from=build /app/app .

COPY configs/ configs

COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

EXPOSE 8080 50051

CMD ["wait-for-it.sh", "postgres:5432", "--", "./app"]