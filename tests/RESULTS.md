# Результаты проверки

- `go test -race ./...` — успешно во всех трёх Go-модулях.
  Unit-тесты находятся в Order Service: `internal/orders/service` и
  `internal/orders/httpclient`. User и Notification проверены интеграционно.
- `go vet ./...` — успешно во всех трёх модулях.
- `docker compose config --quiet` — успешно.
- `docker compose up -d --build` — образы собраны, три приложения и PostgreSQL запущены.
- `python3 tests/integration.py` — успешно, код завершения 0.

Фактический вывод последнего интеграционного прогона:

```text
PASS: CRUD, validation, HTTP order -> user -> notification scenario
PASS: notification outage preserves order, returns 201 and notification_sent=false
PASS: user outage returns 503, creates no order; other services stay available
PASS: application restarts preserve data
PASS: all three DB roles denied connections to another service database
PASS: 54 HTTP status assertions; test records removed
```

Проверены реальные HTTP-запросы и записи PostgreSQL. При отключении уведомлений
создаётся читаемый заказ с 201 и `notification_sent=false` в ответе создания.
При отключении пользователей POST заказа возвращает 503; число заказов не растёт.
При этом чтение существующих заказов и Notification Service продолжают работать.
После рестарта приложений записи доступны. Все три роли получили запрет доступа
к чужой БД. Тестовые записи удалены, сервисы восстановлены.
