## skillsRock gRPC auth service

### **Настройка проекта**
1. Создайте файл .env в корне проекта и добавьте следующие параметры:
```bash
# Настройки окружения и логгирования
ENV=env
LOG_LEVEL=info

# Настройка сервера gRPC
GRPC_PORT=44044
GRPC_TIMEOUT=10h

# Настройки подключения к базе данных
DB_HOST=localhost
DB_PORT=port
DB_USERNAME=username
DB_PASSWORD=password
DB_DATABASE=database
DB_SSL_MODE=disable
DB_POOL_MAX_CONNS=10
DB_POOL_MAX_CONN_LIFETIME=300s
DB_POOL_MAX_CONN_IDLE_TIME=150s

# Параметры для генерации JWT
JWT_SECRET=secret
JWT_TTL=2m
```

2. Команды для запуска контейнера и подключения к базе данных PostgreSQL
```bash
docker run --name your_db -e POSTGRES_PASSWORD=your_password -p 5433:5432 -d postgres
psql -h localhost -p 5433 -U postgres -d your_db
```

3. Команды для установки и применения миграций с помощью go goose
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose create имя_миграции sql
goose -dir internal/storage/migrations postgres "postgres://user:pass@localhost:5433/dbname?sslmode=disable" up
```

4. Команды для установки и применения тестов для api и service слоев
```bash
go install github.com/vektra/mockery/v2@v2.53.4
mockery --name=AuthRepository --dir=internal/services --output=internal/mocks --log-level=debug
mockery --name=AuthRepository --dir=internal/api --output=internal/mocks --log-level=debug
```

5. Установите зависимости и запустите проект:
```bash
go run cmd/main.go
```

Сервис будет доступен по адресу: http://localhost:44044. Проверка работоспособности проводилась в Postman с использованием proto файлов

## Примеры запросов
### **Регистрация пользователя**
localhost:44044 - AuthService/Register
```bash
{
  "username": "user1",
  "password": "secret"
}
```

### Ответ:
```bash
{
  "message": "success"
}
```

### **Авторизация пользователя**
localhost:44044 - AuthService/Login
```bash
{
  "username": "user1",
  "password": "secret"
}
```

### Ответ:
```bash
{
  "token": "secret_token"
}
```

Тестирование можно проводить удобно через среду разработки GoLand и проверить работоспособность слоев api и services.

Сервис готов к использованию и может служить основой для полноценного приложения.




