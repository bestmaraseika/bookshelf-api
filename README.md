# BookShelf API

RESTful API для управления личной библиотекой книг, разработанный на Go.

Учебный backend-проект с PostgreSQL, Docker Compose, миграциями базы данных и unit-тестированием.

## Возможности

- Создание книги.
- Получение книги по ID.
- Получение списка книг с пагинацией.
- Обновление книги.
- Удаление книги.
- Получение топ-5 рекомендованных книг.
- Отметка книги как закончившейся.
- Валидация JSON и входных данных.
- Обработка HTTP-ошибок.

## Технологии

| Технология | Назначение |
|---|---|
| Go | Язык разработки backend |
| Chi | HTTP-маршрутизация |
| PostgreSQL 17 | Реляционная база данных |
| database/sql + pgx | Работа с PostgreSQL |
| golang-migrate | SQL-миграции |
| Docker | Контейнеризация |
| Docker Compose | Запуск и управление сервисами |
| testing, httptest | Unit-тестирование |
| go-sqlmock | Тестирование SQL-операций без PostgreSQL |
| Postman | Ручное тестирование API |

## Архитектура

Приложение использует многослойную архитектуру:

```text
HTTP Client
    |
    v
Router (Chi)
    |
    v
Handler
    |
    v
Service
    |
    v
Repository
    |
    v
PostgreSQL
```

**Handler** отвечает за HTTP-запросы, декодирование JSON, валидацию и формирование ответов.

**Service** отвечает за бизнес-логику приложения.

**Repository** выполняет SQL-запросы и преобразует результаты в Go-структуры.

## Структура проекта

```text
bookshelf/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   ├── routes/
│   └── services/
├── migrations/
├── .dockerignore
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Требования для запуска

- Docker Desktop или Docker Engine с Docker Compose.
- Свободные локальные порты `8182` и `5432`.

Go SDK не требуется для запуска через Docker Compose: приложение компилируется внутри контейнера сборки.

## Запуск через Docker Compose

### 1. Подготовить конфигурацию

Создать локальный `.env` на основе `.env.example`.

Для PowerShell:

```powershell
Copy-Item .env.example .env
```

Указать необходимые значения:

```dotenv
POSTGRES_USER=bookshelf_user
POSTGRES_PASSWORD=change_me
POSTGRES_DB=bookshelf

DATABASE_URL=postgres://bookshelf_user:change_me@localhost:5432/bookshelf?sslmode=disable

APP_PORT=8182
```

Заменить `change_me` собственным паролем и использовать его последовательно в соответствующих параметрах.

Если проект уже запускался с существующим Docker Volume, необходимо использовать действующие учётные данные базы данных.

Файл `.env` не должен добавляться в Git.

### 2. Запустить приложение

Выполнить из корня проекта:

```powershell
docker compose up --build -d
```

Docker Compose:

1. Запустит PostgreSQL.
2. Дождётся готовности базы данных.
3. Применит необходимые SQL-миграции.
4. Запустит Go API.

Проверить состояние:

```powershell
docker compose ps -a
```

Посмотреть журнал приложения:

```powershell
docker compose logs app
```

API доступен по адресу:

```text
http://localhost:8182
```

### 3. Проверить работу API

```http
GET http://localhost:8182/health
```

Ожидаемый ответ:

```json
{
  "status": "ok"
}
```

HTTP-статус: `200 OK`.

### 4. Остановить приложение

```powershell
docker compose down
```

Команда останавливает и удаляет контейнеры проекта, сохраняя именованный Docker Volume с данными PostgreSQL.

## API Endpoints

| Метод | Endpoint | Назначение |
|---|---|---|
| GET | `/health` | Проверка HTTP-сервера |
| POST | `/books` | Создание книги |
| GET | `/books` | Список книг с пагинацией |
| GET | `/books/{id}` | Получение книги по ID |
| PUT | `/books/{id}` | Обновление книги |
| DELETE | `/books/{id}` | Удаление книги |
| GET | `/books/recommend` | Топ-5 рекомендаций |
| POST | `/books/{id}/mark-out-of-stock` | Отметить книгу как закончившуюся |

## Примеры запросов

### POST /books

Создание новой книги.

```http
POST http://localhost:8182/books
Content-Type: application/json
```

Тело запроса:

```json
{
  "title": "Dune",
  "author": "Frank Herbert",
  "year": 1965,
  "isbn": null,
  "out_of_stock": false,
  "rating": 9
}
```

Успешный ответ: `201 Created`.

Поля `id`, `created_at` и `updated_at` генерируются сервером и базой данных.

### GET /books

Получение списка книг с пагинацией.

```http
GET http://localhost:8182/books?page=1&limit=10
```

Параметры:

- `page` — номер страницы, по умолчанию 1.
- `limit` — количество записей, по умолчанию 10; допустимый диапазон 1–100.

Успешный ответ: `200 OK` и JSON-массив книг.

### GET /books/{id}

```http
GET http://localhost:8182/books/1
```

Успешный ответ: `200 OK`.

Если книга отсутствует: `404 Not Found`.

### PUT /books/{id}

Полное обновление редактируемых полей книги.

```http
PUT http://localhost:8182/books/1
Content-Type: application/json
```

```json
{
  "title": "Dune: Updated",
  "author": "Frank Herbert",
  "year": 1965,
  "isbn": null,
  "out_of_stock": false,
  "rating": 10
}
```

Успешный ответ: `200 OK`.

### DELETE /books/{id}

```http
DELETE http://localhost:8182/books/1
```

Успешный ответ: `204 No Content`.

### GET /books/recommend

```http
GET http://localhost:8182/books/recommend
```

Возвращает до пяти книг.

Правила сортировки:

1. Сначала книги с рейтингом по убыванию.
2. Затем книги без рейтинга по году издания, от новых к старым.
3. При совпадении рейтинга и года — по ID.

Успешный ответ: `200 OK`.

### POST /books/{id}/mark-out-of-stock

```http
POST http://localhost:8182/books/1/mark-out-of-stock
```

Устанавливает `out_of_stock = true` и обновляет `updated_at`.

Успешный ответ: `204 No Content`.

Тело запроса не требуется.

## Модель Book

| Поле | Тип | Описание |
|---|---|---|
| id | uint64 | Уникальный идентификатор |
| title | string | Обязательное название |
| author | string | Обязательный автор |
| year | int | Год издания |
| isbn | string/null | Необязательный уникальный ISBN |
| out_of_stock | bool | Статус наличия |
| rating | int/null | Необязательный рейтинг 1–10 |
| created_at | timestamp | Время создания |
| updated_at | timestamp | Время последнего изменения |

## Обработка ошибок

API возвращает ошибки в формате JSON:

```json
{
  "error": "book not found"
}
```

Основные HTTP-статусы:

| Статус | Значение |
|---|---|
| 200 | Успешное получение или обновление |
| 201 | Ресурс создан |
| 204 | Успешная операция без тела ответа |
| 400 | Некорректные входные данные |
| 404 | Книга не найдена |
| 409 | Конфликт уникальности ISBN |
| 413 | Тело запроса превышает лимит |
| 500 | Внутренняя ошибка сервера |

## Миграции

Миграции хранятся в каталоге `migrations`.

Используется `golang-migrate`.

При запуске Docker Compose сервис `migrate` автоматически применяет ещё не применённые миграции.

Текущая миграция создаёт таблицу `books` с первичным ключом, ограничениями и временными метками.

## Тестирование

Для запуска unit-тестов локально необходим Go SDK.

Запуск всех тестов:

```powershell
go test ./...
```

Подробный вывод:

```powershell
go test ./... -v
```

Статический анализ:

```powershell
go vet ./...
```

Измерение покрытия:

```powershell
go test ./... -cover
```

Тесты используют `testing`, `httptest`, тестовые реализации Repository и `go-sqlmock`.

Для большинства unit-тестов запуск PostgreSQL не требуется.

## Примечания

Проект разработан в учебных целях.

Он демонстрирует работу с HTTP, PostgreSQL, SQL, миграциями, многослойной архитектурой, Docker и автоматизированными тестами.