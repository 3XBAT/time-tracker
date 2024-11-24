# Time Tracker Service

Сервис для отслеживания времени выполнения задач пользователями.

## Описание

Time Tracker - это REST API сервис, который позволяет:
- Управлять пользователями (создание, обновление, удаление, получение)
- Создавать задачи для пользователей
- Отслеживать время выполнения задач
- Получать статистику по выполненным задачам

## Особенности работы

### Задачи (Tasks)
- Сервис отображает только **завершенные** задачи (где есть и время начала, и время окончания)
- Длительность задачи отображается в формате: "Xd XXh XXm" (дни, часы, минуты)
- При создании задачи без указания времени начала, оно устанавливается автоматически
- Задачу можно завершить только один раз
- Время должно быть в формате RFC3339 (пример: "2024-03-20T15:30:00.000+03:00")

### Пользователи (Users)
- Каждый пользователь имеет уникальный номер паспорта
- Поддерживается пагинация при получении списка пользователей (параметры Limit и Offset)
- Фильтрация пользователей возможна по всем полям

## API Endpoints

### Health Check
GET /health
Проверка работоспособности сервиса

### Users
GET /users - Получение списка пользователей с фильтрацией
GET /users/:id - Получение пользователя по ID
POST /users - Создание пользователя
PUT /users/:id - Обновление данных пользователя
DELETE /users/:id - Удаление пользователя

### Tasks
GET /tasks - Получение списка выполненных задач
POST /tasks - Создание новой задачи
PUT /tasks/:id - Завершение задачи
DELETE /tasks/:id - Удаление задачи

## Примеры запросов

### Создание задачи
POST /tasks
{
    "name": "Новая задача",
    "start_time": "2024-03-20T15:30:00.000+03:00",
    "user_id": 1
}

### Получение задач пользователя
GET /tasks?user_id=1&start_time=2024-03-01T00:00:00Z&end_time=2024-03-31T23:59:59Z

### Создание пользователя
POST /users
{
    "passport_number": "1234 567890",
    "name": "Иван",
    "patronymic": "Иванович",
    "surname": "Иванов",
    "addr": "г. Москва, ул. Примерная, д. 1"
}

## Конфигурация

Сервис настраивается через переменные окружения (файл .env):

## База данных

Сервис использует PostgreSQL. Структура базы данных:

### Таблица users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    passport_number VARCHAR(11) NOT NULL,
    name VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255) NOT NULL,
    addr VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL
);

### Таблица tasks
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE,
    end_time TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

## Запуск

1. Клонируйте репозиторий
git clone <repository-url>

2. Создайте файл .env

3. Установите зависимости
go mod download

4. Запустите PostgreSQL

5. Создайте базу данных
CREATE DATABASE timeDB;

6. Запустите приложение
go run cmd/main.go

## Swagger документация

Swagger UI доступен по адресу: http://localhost:8080/swagger/index.html

Для регенерации документации:
swag init -g cmd/main.go -o ./docs

## Основные зависимости

- Go 1.22+
- PostgreSQL
- gin-gonic/gin (веб-фреймворк)
- swaggo/swag (документация API)
- jmoiron/sqlx (работа с базой данных)
- golang-migrate/migrate (миграции)