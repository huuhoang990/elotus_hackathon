docker-compose down -v

docker-compose build && docker-compose up -d

docker-compose exec api go run /app/cmd/internal/database/migrate.go
