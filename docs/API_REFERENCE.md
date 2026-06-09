# API Reference

## Base URL
- **Local:** `http://localhost:8080`
- **Production:** `https://api.taskmaster.com`

## Authentication
All protected endpoints require a Bearer token:
```
Authorization: Bearer <jwt_token>
```

## Endpoints

### Health
```
GET /health
```
**Response:**
```json
{
  "status": "ok",
  "time": "2026-05-23T12:00:00Z"
}
```

### Authentication

#### Register
```
POST /api/v1/auth/register
```
**Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "User Name"
}
```
**Response (201):**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "User Name",
    "role": "member",
    "created_at": "2026-05-23T12:00:00Z"
  }
}
```

#### Login
```
POST /api/v1/auth/login
```
**Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "User Name",
    "role": "member"
  }
}
```

### Tasks

#### Create Task
```
POST /api/v1/tasks
```
**Headers:** `Authorization: Bearer <token>`
**Body:**
```json
{
  "title": "Task title",
  "description": "Optional description",
  "priority": "high",
  "assignee_id": "uuid",
  "due_date": "2026-06-01T00:00:00Z"
}
```
**Response (201):**
```json
{
  "task": {
    "id": "uuid",
    "title": "Task title",
    "status": "todo",
    "priority": "high",
    "creator_id": "uuid",
    "team_id": "uuid",
    "created_at": "2026-05-23T12:00:00Z"
  }
}
```

#### List Tasks
```
GET /api/v1/tasks?limit=20&offset=0
```
**Headers:** `Authorization: Bearer <token>`
**Response (200):**
```json
{
  "tasks": [
    {
      "id": "uuid",
      "title": "Task title",
      "status": "todo",
      "priority": "high",
      "created_at": "2026-05-23T12:00:00Z"
    }
  ]
}
```

#### Get Task
```
GET /api/v1/tasks/:id
```
**Response (200):** Single task object

#### Update Task
```
PUT /api/v1/tasks/:id
```
**Body:** Partial task fields

#### Delete Task
```
DELETE /api/v1/tasks/:id
```
**Response (200):**
```json
{"message": "deleted"}
```

#### Assign Task
```
POST /api/v1/tasks/:id/assign
```
**Body:**
```json
{"assignee_id": "uuid"}
```

### WebSocket
```
GET /ws
```
**Protocol:** WebSocket
**Messages:**
```json
{
  "type": "task:created",
  "payload": { /* task object */ }
}
```

## Error Responses
```json
{
  "error": "Error message"
}
```

**Status Codes:**
- 400: Bad Request (validation error)
- 401: Unauthorized (missing/invalid token)
- 403: Forbidden (insufficient permissions)
- 404: Not Found
- 429: Too Many Requests (rate limit)
- 500: Internal Server Error

## Rate Limits
- 100 requests per minute per IP
- Applies to all `/api/v1/*` endpoints
- Excludes `/health` and auth endpoints
