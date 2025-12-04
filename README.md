## Warehouse Management System (Golang + SQLite)

Полнофункциональное веб‑приложение для автоматизации складского учёта.  
Бэкенд на Go (чистая архитектура: Controller–Service–Repository–Models–Middleware–Config), хранилище — SQLite.

---

## Структура проекта

- `backend/src`
  - `main.go` — точка входа, настройка роутов, middleware, DI слоёв.
  - `controllers/` — HTTP‑слой (auth, products, categories, suppliers, warehouse, orders).
  - `services/` — бизнес‑логика (валидация, правила домена, хэширование паролей, JWT).
  - `repositories/` — реализация доступа к данным (SQLite).
  - `models/` — доменные сущности (User, Product, Category, Supplier, Order, StockMovement и т.д.).
  - `middleware/` — CORS, логирование, JWT‑аутентификация, проверка ролей.
  - `config/` — конфигурация приложения и инициализация БД (`db.go`, `config.go`).
- `frontend/public`
  - `index.html` — логин + регистрация.
  - `dashboard.html` — дашборд, метрики по основным сущностям.
  - `products.html` — управление товарами.
  - `warehouse.html` — складские операции.
  - `orders.html` — заказы.

---

## Запуск проекта

### 1. Предварительные требования

- Go **1.22+**
- Git

SQLite используется в виде встроенной библиотеки (`modernc.org/sqlite`), отдельный сервер БД не нужен.

### 2. Настройка окружения

Перейдите в каталог `backend` и установите зависимости:

```bash
cd backend
go mod tidy
```

Переменные окружения (можно задать через `.env` или системные переменные):

- `JWT_SECRET` — секрет для подписи JWT (в разработке, при отсутствии, берётся небезопасное значение по умолчанию).
- `PORT` — порт HTTP‑сервера (по умолчанию `8080`).

### 3. Запуск бэкенда

Находясь в каталоге `backend`:

```bash
go run ./src
```

При первом запуске:

- создаётся файл БД `backend/warehouse.db`;
- выполняются миграции (создаются таблицы `users`, `products`, `categories`, `suppliers`, `orders`, `order_items`, `stock_movements`, `order_status_history` и индексы);
- HTTP‑сервер поднимается на `http://localhost:8080`.

Фронтенд‑страницы раздаются тем же сервером:

- `/` → `frontend/public/index.html` (логин/регистрация)
- `/dashboard` → `dashboard.html`
- `/products` → `products.html`
- `/warehouse` → `warehouse.html`
- `/orders` → `orders.html`

---

## Аутентификация и авторизация

- Используется JWT (HMAC, секрет `JWT_SECRET`).
- Пароли хэшируются через `bcrypt` (пакет `golang.org/x/crypto/bcrypt`).
- Роли пользователей: `admin`, `manager`, `storekeeper`.
- Доступ к защищённым эндпоинтам проверяется через middleware:
  - `AuthMiddleware` — проверка JWT, извлечение `userID` и `role` в context.
  - `RoleMiddleware` — проверка роли по списку разрешённых.

Интерактивная регистрация и вход выполняются с фронтенда:

- Страница **логина**: `GET /` → `index.html` (форма входа, при успехе сохранение JWT в `localStorage` и переход на `/dashboard`).
- Страница **регистрации**: `GET /register` → `register.html` (создание пользователя с выбранной ролью, при успехе такой же вход и переход на `/dashboard`).

### Эндпоинты аутентификации

Все тела/ответы — `application/json`.

- **POST `/api/auth/register`** — регистрация пользователя.
  - Тело:
    ```json
    {
      "email": "user@example.com",
      "password": "secret123",
      "role": "admin"
    }
    ```
  - Ответ `201 Created`:
    ```json
    {
      "user": {
        "id": "u-...",
        "email": "user@example.com",
        "role": "admin",
        "created_at": "...",
        "updated_at": "..."
      }
    }
    ```

- **POST `/api/auth/login`** — вход.
  - Тело:
    ```json
    {
      "token": "AUTH_TOKEN", 
      "email": "user@example.com", 
      "password": "secret123"
    }
    ```
  - Ответ `200 OK`: как у `/register` (token + user).

- **GET `/api/auth/me`** — текущий пользователь по JWT.
  - Заголовок: `Authorization: Bearer <token>`.
  - Ответ `200 OK` — объект пользователя.

---

## Модуль товаров

Роуты защищены JWT; создание/редактирование/удаление доступны только ролям `admin` и `manager`.

- **GET `/api/products`**
  - Возвращает массив товаров.
  - Заголовок: `Authorization: Bearer <token>`.
  - Ответ `200 OK`:
    ```json
    [{
      "id": "p-1",
      "sku": "SKU-001",
      "name": "Товар",
      "description": "Описание",
      "category_id": "c-1",
      "supplier_id": "s-1",
      "unit": "pcs",
      "created_at": "...",
      "updated_at": "..."
    }]
    ```

- **POST `/api/products`**
  - Роли: `admin`, `manager`.
  - Тело:
    ```json
    {
      "sku": "SKU-001",
      "name": "Товар",
      "description": "Описание",
      "category_id": "c-1",
      "supplier_id": "s-1",
      "unit": "pcs"
    }
    ```
  - Ответ `201 Created` — созданный товар.

- **GET `/api/products/{id}`** — получить товар по ID.
- **PUT `/api/products/{id}`** — обновить товар (тело как при создании).
- **DELETE `/api/products/{id}`** — удалить товар (роль `admin`).

Аналогичные CRUD‑эндпоинты реализованы для:

- `/api/categories` (`GET, POST, GET {id}, PUT {id}, DELETE {id}`)
- `/api/suppliers`  (`GET, POST, GET {id}, PUT {id}, DELETE {id}`)

---

## Складские операции

Маршруты:

- **POST `/api/warehouse/receipt`** — приёмка товара на склад.
  - Роли: `admin`, `manager`, `storekeeper`.
  - Тело:
    ```json
    {
      "product_id": "p-1",
      "supplier_id": "s-1",
      "quantity": 10,
      "price": 50,
      "expiry_date": "2025-12-31T00:00:00Z"
    }
    ```

- **POST `/api/warehouse/write-off`** — списание товара.
  - Роли: `admin`, `manager`, `storekeeper`.
  - Тело:
    ```json
    {
      "product_id": "p-1",
      "quantity": 2
    }
    ```

- **POST `/api/warehouse/reserve`** — резервирование под заказ.
  - Роли: `admin`, `manager`.
  - Тело:
    ```json
    {
      "product_id": "p-1",
      "order_id": "o-1",
      "quantity": 3
    }
    ```

- **GET `/api/warehouse/inventory`** — текущие остатки.
  - Ответ:
    ```json
    [{
      "product_id": "p-1",
      "quantity": 15
    }]
    ```

---

## Заказы

- **GET `/api/orders`** — список заказов.
- **POST `/api/orders`** — создать заказ (с автоматическим резервированием товара).
  - Тело:
    ```json
    {
      "customer": "ООО Ромашка",
      "items": [
        { "product_id": "p-1", "quantity": 2, "price": 100 },
        { "product_id": "p-2", "quantity": 1, "price": 50 }
      ]
    }
    ```
  - Ответ содержит заказ с полями `items`, `status`, `status_history`.

- **GET `/api/orders/{id}`** — получить заказ по ID.
- **PUT `/api/orders/{id}/status`** — обновить статус заказа.
  - Тело:
    ```json
    { "status": "completed" }
    ```

История статусов сохраняется в таблице `order_status_history` и возвращается в поле `status_history` заказа.

---

## Фронтенд и работа с API

Фронтенд — это набор HTML/JS‑страниц, которые обращаются к REST‑API:

- `index.html` — логин и регистрация, сохраняет JWT в `localStorage` (`wms_token`) и данные пользователя (`wms_user`).
- `dashboard.html` — использует `/api/auth/me`, `/api/products`, `/api/categories`, `/api/suppliers`, `/api/orders` для отображения метрик.
- `products.html` — полный CRUD по товарам через `/api/products`.
- `warehouse.html` — приёмка, списание, резервирование и просмотр остатков через `/api/warehouse/*`.
- `orders.html` — список заказов, создание и смена статуса через `/api/orders`.

Все запросы к защищённым эндпоинтам выполняются с заголовком:

```http
Authorization: Bearer <JWT_TOKEN>
```

Ошибки API стандартизированы: при ошибке сервер возвращает объект вида:

```json
{ "error": "сообщение об ошибке" }
```

и соответствующий HTTP‑статус (`400/401/403/404/409/500`).