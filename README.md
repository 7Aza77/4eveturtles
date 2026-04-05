# 📅 GoEvent – Платформа Управления Студенческими Мероприятиями

**GoEvent** — это мощный backend-сервис на языке Go, предоставляющий REST API для централизованного управления студенческими событиями. Проект реализован с использованием **Clean Architecture**, поддерживает JWT-аутентификацию, работу с PostgreSQL, кэширование через Redis и контейнеризацию через Docker.

---

## ✨ Основные возможности (Production-Ready)

-   **🔐 Аутентификация и Безопасность**:
    -   JWT-аутентификация (регистрация, вход).
    -   Хеширование паролей (bcrypt).
    -   **RBAC (Role-Based Access Control)**: Разделение прав для Студентов, Модераторов и Админов.
-   **📅 Управление Мероприятиями**:
    -   Полный CRUD событий (Создание, Просмотр, Редактирование, Удаление).
    -   **Advanced Filtering**: Поиск по названию, локации и диапазону дат.
    -   **Pagination & Sorting**: Поддержка `limit`, `offset` и сортировки по полям.
-   **📝 Регистрации**:
    -   Умная система регистрации участников с проверкой вместимости (Capacity).
    -   Транзакционная логика для надежности данных.
-   **🚀 Производительность и Надежность**:
    -   **Redis Caching**: Кэширование списков и деталей событий для мгновенного отклика.
    -   **Rate Limiting**: Защита от перегрузок (100 запросов в минуту).
    -   **Structured Logging**: Использование `log/slog` (JSON в PROD, Text в LOCAL).
    -   **Graceful Shutdown**: Безопасная остановка сервера без потери данных.
-   **🛠 Инструментарий**:
    -   **Database Migrations**: Автоматическое управление схемой через `golang-migrate`.
    -   **Swagger Documentation**: Интерактивная документация API.

---

## 🛠 Технологический стек

-   **Язык**: Go 1.25+
-   **Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
-   **Database**: [PostgreSQL](https://www.postgresql.org/) + [sqlx](https://github.com/jmoiron/sqlx)
-   **Cache**: [Redis](https://redis.io/)
-   **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
-   **Validation**: [validator/v10](https://github.com/go-playground/validator)
-   **Logging**: `log/slog`
-   **Config**: [cleanenv](https://github.com/ilyakaznacheev/cleanenv)
-   **Documentation**: [Swagger](https://github.com/swaggo/swag)

---

## 📁 Структура проекта (Clean Architecture)

```markdown
├── cmd/app/main.go          # Точка входа в приложение
├── internal/
│   ├── entity/              # Бизнес-модели (User, Event)
│   ├── usecase/             # Бизнес-логика (Register, CreateEvent)
│   ├── repository/          # Слой данных (Postgres, Redis)
│   ├── handler/             # Слой доставки (REST API Handlers & Middleware)
│   └── config/              # Конфигурация приложения
├── migrations/              # SQL-файлы миграций
├── pkg/                     # Переиспользуемые библиотеки
├── config/                  # Файлы настроек (config.yaml)
└── docs/                    # Сгенерированная документация Swagger
```

---

## 🚀 Как запустить

### 1. Быстрый запуск через Docker (Рекомендуется)
Убедитесь, что у вас установлен Docker и Docker Compose.

```bash
docker-compose up -d --build
```
*Миграции применятся автоматически, все сервисы (app, db, redis) будут запущены и связаны.*

### 2. Локальный запуск (Разработка)
1. Установите зависимости:
   ```bash
   go mod tidy
   ```
2. Настройте `config/config.yaml` или используйте ENV-переменные.
3. Запустите приложение:
   ```bash
   go run cmd/app/main.go
   ```

---

## 📖 API Документация (Swagger)

После запуска сервера документация доступна по адресу:
👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

Здесь можно протестировать все эндпоинты в интерактивном режиме.

---

## 🧪 Тестирование

Запуск юнит-тестов для проверки бизнес-логики:

```bash
go test -v ./internal/usecase/...
```

---

## 🔒 Контакты и поддержка
Разработано в рамках проекта **GoEvent Backend Platform**.
 Если у вас есть вопросы или предложения — создавайте Issue или Pull Request!
