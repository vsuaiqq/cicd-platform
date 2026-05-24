# CI/CD Client

Frontend для CI/CD системы: Vite, React, TypeScript, Emotion, Redux Toolkit. Обращается к **api-gateway**.

## Запуск

```bash
npm install
npm run dev
```

По умолчанию приложение на `http://localhost:5173`. API-запросы уходят на `http://localhost:8080` (см. переменные окружения).

## Переменные окружения

- `VITE_API_URL` — базовый URL api-gateway (без слэша в конце). По умолчанию: `http://localhost:8080`.

Скопируй `.env.example` в `.env` и при необходимости измени значение.

## Сборка

```bash
npm run build
```

Артефакты в `dist/`.

## Роуты и API

- **/** — Dashboard (требуется авторизация).
- **/login** — вход.
- **/register** — регистрация.

Используются эндпоинты api-gateway:

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/validate`
- `POST /api/v1/auth/refresh`

Токены хранятся в `localStorage` и автоматически подставляются в запросы; при истечении access token выполняется refresh.
