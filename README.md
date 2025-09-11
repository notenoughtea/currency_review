# Currency Review

## Быстрый старт

```bash
docker compose up -d --build
# запуск всего проекта в контейнерах
```

- Gateway доступен: http://localhost:8081
- Adminer: http://localhost:8111 (Server: `db`, User: `admin`, Password: `admin`, DB: `currency_db`)

## Наполнение тестовыми данными (сидер) опционально...

```bash
# после успешного запуска db и migrator
docker compose run --rm seed
```

Сидер создаёт записи для валют `USD`, `EUR`, `GBP` примерно за 90 случайных дат последних ~120 дней.

## Конфигурация

- Конфиг: `config.yaml` (монтируется внутрь контейнеров как `/app/config.yaml`).
- Важные параметры:
  - `DB.host: db`, `DB.port: 5432` (используйте `localhost:5111` для локального запуска вне Docker)
  - `outer_currency_url`: база внешнего API (поддерживает `.../<base>.json`)
  - `DefaultBaseCurrency`: базовая валюта для клиента

## Эндпоинты

- Регистрация: `POST /api/v1/register`
- Логин: `POST /api/v1/login`
- Курсы за период: `GET /api/v1/rate?currency=USD&date_from=YYYY-MM-DD&date_to=YYYY-MM-DD`
  - Заголовок: `Authorization: Bearer <JWT>`

## Примеры запросов

Регистрация

```bash
curl -X POST http://localhost:8081/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "Username": "test",
    "Password": "test"
  }'
```

Логин

```bash
curl -X POST http://localhost:8081/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "Username": "test",
    "Password": "test"
  }'
```

Курсы за период

```bash
# замените <JWT> на полученный токен
curl -X GET "http://localhost:8081/api/v1/rate?currency=USD&date_from=2024-10-01&date_to=2024-10-21" \
  -H "Authorization: Bearer <JWT>"
```

Запрос с одинаковыми датами

```bash
curl -X GET "http://localhost:8081/api/v1/rate?currency=USD&date_from=2024-09-09&date_to=2024-09-09" \
  -H "Authorization: Bearer <JWT>"
```

## Полезное

- Просмотр логов:

```bash
docker compose logs -f db migrator currency gateway cron seed
```

- Пересборка отдельного сервиса:

```bash
docker compose up -d --build currency
```
