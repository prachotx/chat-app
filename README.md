# Chat App

A full-stack real-time chat application built with Go (Fiber) and Angular, featuring WebSocket-based messaging, JWT authentication, and room-based conversations.

## Tech Stack

**Backend**
- Go with [Fiber v3](https://github.com/gofiber/fiber)
- GORM + PostgreSQL
- WebSocket (via `gofiber/contrib/websocket`)
- JWT authentication (HttpOnly cookie)
- Rate limiting middleware

**Frontend**
- Angular 21 (standalone components, signals)
- Tailwind CSS v4 + DaisyUI
- RxJS

**Infrastructure**
- Docker + Docker Compose (dev & prod)
- Nginx (production reverse proxy)
- pgAdmin (dev database management)

## Features

- Register / Login with JWT stored in HttpOnly cookie
- Create and join chat rooms
- Real-time messaging via WebSocket
- Online presence — see who's currently in a room
- User join/leave notifications broadcast to all room members
- Rate limiting: 100 req/min globally, 10 req/min on auth endpoints

## Project Structure

```
chat-app/
├── api/                  # Go REST API + WebSocket server
│   ├── cmd/server/       # Entrypoint
│   ├── config/           # Env config loader
│   ├── internal/
│   │   ├── dto/          # Request/response shapes
│   │   ├── handler/      # HTTP handlers
│   │   ├── middleware/   # Auth & rate limit middleware
│   │   ├── model/        # GORM models
│   │   ├── repository/   # DB layer
│   │   ├── service/      # Business logic
│   │   └── ws/           # WebSocket hub & client
│   └── pkg/              # Shared utilities (jwt, crypto, response)
├── web/                  # Angular frontend
│   └── src/app/
│       ├── core/         # Services & guards
│       └── pages/        # Login, Register, Rooms, Room
├── docker-compose.dev.yml
├── docker-compose.prod.yml
└── .env.example
```

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- (Optional) Go 1.21+ and Node.js 20+ for local development without Docker

### 1. Configure environment variables

```bash
cp .env.example .env
cp api/.env.example api/.env
```

Edit `api/.env` with your values:

| Variable      | Description                    |
|---------------|--------------------------------|
| `PORT`        | API server port (default 3000) |
| `DB_HOST`     | Postgres host                  |
| `DB_PORT`     | Postgres port                  |
| `DB_USER`     | Postgres user                  |
| `DB_PASSWORD` | Postgres password              |
| `DB_NAME`     | Postgres database name         |
| `SECRET_KEY`  | JWT signing secret             |
| `APP_ENV`     | `development` or `production`  |

Root `.env` controls Docker container names and exposed ports:

| Variable             | Default              |
|----------------------|----------------------|
| `API_PORT`           | `3000`               |
| `WEB_PORT`           | `4200`               |
| `PGADMIN_PORT`       | `5050`               |
| `POSTGRES_DATABASE`  | `real_time_chat_db`  |
| `POSTGRES_USER`      | `user`               |
| `POSTGRES_PASSWORD`  | `password`           |

### 2. Run with Docker

**Development** (hot reload for both API and frontend):

```bash
docker compose -f docker-compose.dev.yml up --build
```

| Service  | URL                     |
|----------|-------------------------|
| Frontend | http://localhost:4200   |
| API      | http://localhost:3000   |
| pgAdmin  | http://localhost:5050   |

**Production:**

```bash
docker compose -f docker-compose.prod.yml up --build
```

The frontend is served via Nginx on `WEB_PORT` and proxies API/WebSocket traffic to the backend on the internal Docker network.

### 3. Run locally (without Docker)

**API:**

```bash
cd api
cp .env.example .env   # fill in DB credentials pointing to a local Postgres
air                    # hot reload — requires github.com/air-verse/air
# or: go run ./cmd/server/main.go
```

**Frontend:**

```bash
cd web
npm install
npm start              # serves at http://localhost:4200
```

## API Overview

All responses follow the format:

```json
{ "message": "...", "data": { ... } }
```

### Auth

| Method | Path               | Description            |
|--------|--------------------|------------------------|
| POST   | `/api/auth/register` | Create a new account |
| POST   | `/api/auth/login`    | Login, sets `access_token` cookie |
| POST   | `/api/auth/logout`   | Clear the auth cookie  |
| GET    | `/api/auth/me`       | Get current user (protected) |

### Rooms

| Method | Path              | Description         |
|--------|-------------------|---------------------|
| GET    | `/api/rooms`      | List all rooms      |
| POST   | `/api/rooms`      | Create a room       |
| GET    | `/api/rooms/:id`  | Get room details    |

### Messages

| Method | Path                       | Description              |
|--------|----------------------------|--------------------------|
| GET    | `/api/rooms/:id/messages`  | Get message history      |

### WebSocket

```
WS /api/ws/:roomId?token=<jwt>
```

Events sent from server to client:

| Event          | Payload                        |
|----------------|--------------------------------|
| `message`      | New chat message               |
| `online_users` | Full list of online user IDs   |
| `user_joined`  | A user connected to the room   |
| `user_left`    | A user disconnected            |

Events sent from client to server:

| Event     | Payload           |
|-----------|-------------------|
| `message` | `{ "content": "..." }` |
