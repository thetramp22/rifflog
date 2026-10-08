# RiffLog API Documentation

## Overview

The RiffLog API is a RESTful backend service for tracking guitar practice sessions. It allows users to register an account, authenticate using JSON Web Tokens (JWT), record practice sessions, view practice history, and retrieve summary statistics about their practice habits.

The API is built with Go using the Gin web framework and uses PostgreSQL for persistent storage.

### Available Features

- User registration and authentication
- JWT-protected endpoints
- Browse available practice skills
- Create, update, and delete practice sessions
- List practice sessions with optional filtering
- View aggregate practice statistics
- Health and service information endpoints

## Base URLs

### Production

```text
https://api.rifflog.scottstarks.dev
```

### Local Development

```text
http://localhost:8080
```

## Authentication

Protected endpoints require a JSON Web Token (JWT) obtained from `POST /login`.

Include the token in the `Authorization` header using the Bearer authentication scheme:

```http
Authorization: Bearer <JWT>
```

The authenticated user's ID is derived from the JWT. Clients should not send a user ID when creating or modifying practice sessions.

## Common Error Format

Unless otherwise noted, failed requests return JSON in the following format:

```json
{
  "error": "description of the error"
}
```

## Service Endpoints

### GET /

Returns basic information about the API service.

**Authentication:** None

#### Success Response

`200 OK`

```json
{
  "service": "RiffLog API",
  "version": "1.0.0",
  "documentation": "https://github.com/thetramp22/rifflog/blob/main/docs/api.md"
}
```

### GET /health

Returns the health status of the API process.

**Authentication:** None

#### Success Response

`200 OK`

```json
{
  "status": "ok"
}
```

## Authentication Endpoints

### POST /register

Creates a new RiffLog user account.

**Authentication:** None

#### Request Headers

```http
Content-Type: application/json
```

#### Request Body

| Field | Type | Description |
| --- | --- | --- |
| `email` | string | User's email address |
| `password` | string | User's password |

#### Example Request

```json
{
  "email": "jimmy_user@usermail.com",
  "password": "jimmyRules3"
}
```

#### Success Response

`201 Created`

```json
{
  "id": 283,
  "email": "jimmy_user@usermail.com",
  "created_at": "2026-07-09T18:00:00Z"
}
```

#### Error Responses

| Status | Meaning |
| --- | --- |
| `400` | Invalid user data |
| `500` | Internal server error |

### POST /login

Authenticates a user and returns a JWT.

**Authentication:** None

#### Request Headers

```http
Content-Type: application/json
```

#### Request Body

| Field | Type | Description |
| --- | --- | --- |
| `email` | string | User's email address |
| `password` | string | User's password |

#### Example Request

```json
{
  "email": "jimmy_user@usermail.com",
  "password": "jimmyRules3"
}
```

#### Success Response

`200 OK`

```json
{
  "token": "<JWT>",
  "user": {
    "id": 283,
    "email": "jimmy_user@usermail.com",
    "created_at": "2026-07-09T18:00:00Z"
  }
}
```

#### Error Responses

| Status | Meaning |
| --- | --- |
| `400` | Invalid email or request data |
| `401` | Invalid password |
| `404` | User not found |
| `500` | Internal server error |

## Skills

### GET /skills

Returns the list of practice skills available to users.

**Authentication:** None

#### Success Response

`200 OK`

```json
[
  {
    "id": 1,
    "name": "Ear Training",
    "description": "Try playing to identify chords and melodies by ear.",
    "created_at": "2026-07-09T18:00:00Z"
  },
  {
    "id": 2,
    "name": "Scales",
    "description": "Memorize note locations and scale patterns.",
    "created_at": "2026-07-09T18:00:00Z"
  }
]
```

#### Error Responses

| Status | Meaning |
| --- | --- |
| `500` | Internal server error |

# Practice Sessions

All practice-session endpoints are protected by JWT authentication and are prefixed with `/api`.

## POST /api/practice-sessions

Creates a practice session for the authenticated user.

**Authentication:** Required

#### Request Headers

```http
Authorization: Bearer <JWT>
Content-Type: application/json
```

#### Request Body

| Field | Type | Description |
| --- | --- | --- |
| `skill_id` | integer | ID of the practiced skill |
| `duration_minutes` | integer | Length of the session in minutes |
| `practiced_at` | string (RFC3339) | When the session occurred |
| `notes` | string | Optional practice notes |

#### Example Request

```json
{
  "skill_id": 2,
  "duration_minutes": 30,
  "practiced_at": "2026-07-09T18:00:00Z",
  "notes": "Worked on scales"
}
```

#### Success Response

`201 Created`

```json
{
  "id": 2839,
  "skill_id": 2,
  "duration_minutes": 30,
  "notes": "Worked on scales",
  "practiced_at": "2026-07-09T18:00:00Z",
  "created_at": "2026-07-09T19:00:00Z",
  "user_id": 22
}
```

#### Error Responses

| Status | Meaning |
| --- | --- |
| `400` | Invalid request or session data |
| `401` | Missing or invalid authentication |
| `404` | Referenced skill not found |
| `500` | Internal server error |

## GET /api/practice-sessions

Returns practice sessions belonging to the authenticated user.

**Authentication:** Required

### Query Parameters

| Parameter | Type | Description |
| --- | --- | --- |
| `skill` | integer | Return only sessions for the specified skill |
| `from` | date | Return sessions on or after this date |
| `to` | date | Return sessions on or before this date |

The `to` date is treated as inclusive.

#### Example

```text
GET /api/practice-sessions?skill=2&from=2026-07-01&to=2026-07-31
```

#### Success Response

`200 OK`

```json
[
  {
    "id": 3748,
    "skill_id": 2,
    "skill_name": "Scales",
    "skill_description": "Memorize note locations and scale patterns.",
    "duration_minutes": 25,
    "notes": "Scales practice",
    "practiced_at": "2026-07-09T18:00:00Z",
    "created_at": "2026-07-09T19:00:00Z",
    "user_id": 22
  },
  {
    "id": 3921,
    "skill_id": 2,
    "skill_name": "Scales",
    "skill_description": "Memorize note locations and scale patterns.",
    "duration_minutes": 25,
    "notes": "More scales practice",
    "practiced_at": "2026-07-15T18:00:00Z",
    "created_at": "2026-07-15T19:00:00Z",
    "user_id": 22
  }
]
```

#### Error Responses

| Status | Meaning |
| --- | --- |
| `400` | Invalid query parameters |
| `401` | Missing or invalid authentication |
| `500` | Internal server error |

## GET /api/practice-sessions/stats

Returns aggregate practice statistics for the authenticated user.

**Authentication:** Required

### Statistics

- Total minutes practiced
- Total number of sessions
- Most practiced skill by total minutes
- Longest practice session

#### Success Response

`200 OK`

```json
{
  "total_minutes": 384,
  "total_sessions": 12,
  "most_practiced_skill": {
    "name": "Scales",
    "total_minutes": 98
  },
  "longest_session": 43
}
```

If the user has no practice sessions, `most_practiced_skill` is `null`:

```json
{
  "total_minutes": 0,
  "total_sessions": 0,
  "most_practiced_skill": null,
  "longest_session": 0
}
```

#### Error Responses

| Status | Meaning |
| --- | --- |
| `401` | Missing or invalid authentication |
| `500` | Internal server error |

## PUT /api/practice-sessions/{id}

Updates a practice session belonging to the authenticated user.

**Authentication:** Required

### Path Parameters

| Parameter | Type | Description |
| --- | --- | --- |
| `id` | integer | Practice session ID |

### Request Body

| Field | Type | Description |
| --- | --- | --- |
| `skill_id` | integer | ID of the practiced skill |
| `duration_minutes` | integer | Length of the session in minutes |
| `practiced_at` | string (RFC3339) | When the session occurred |
| `notes` | string | Optional practice notes |

#### Example Request

```json
{
  "skill_id": 2,
  "duration_minutes": 30,
  "practiced_at": "2026-07-09T18:00:00Z",
  "notes": "Worked on scales"
}
```

#### Success Response

`200 OK`

```json
{
  "id": 2839,
  "skill_id": 2,
  "duration_minutes": 30,
  "notes": "Worked on scales",
  "practiced_at": "2026-07-09T18:00:00Z",
  "created_at": "2026-07-09T19:00:00Z",
  "user_id": 22
}
```

#### Error Responses

| Status | Meaning |
| --- | --- |
| `400` | Invalid request, session data, or ID |
| `401` | Missing or invalid authentication |
| `404` | Practice session, user, or skill not found |
| `500` | Internal server error |

If the specified practice session belongs to another user, the API also returns `404` rather than exposing the existence of another user's resource.

## DELETE /api/practice-sessions/{id}

Deletes a practice session belonging to the authenticated user.

**Authentication:** Required

### Path Parameters

| Parameter | Type | Description |
| --- | --- | --- |
| `id` | integer | Practice session ID |

#### Success Response

`200 OK`

```json
{
  "message": "practice session 3842 deleted"
}
```

#### Error Responses

| Status | Meaning |
| --- | --- |
| `400` | Invalid or missing ID |
| `401` | Missing or invalid authentication |
| `404` | Practice session not found |
| `500` | Internal server error |

If the specified practice session belongs to another user, the API returns `404`.

## Authentication and Authorization Notes

The API uses JWT authentication to establish the identity of the requesting user.

For protected practice-session endpoints:

1. The JWT middleware validates the token.
2. The authenticated user ID is placed into the request context.
3. The handler retrieves that user ID from the context.
4. Services and repositories use that ID when operating on user-owned resources.

This prevents clients from selecting another user's ID in the request body or URL.

## Example Workflow

A typical client workflow is:

1. `POST /register` to create an account.
2. `POST /login` to authenticate and receive a JWT.
3. `GET /skills` to retrieve available practice skills.
4. Include the JWT as a Bearer token for protected requests.
5. `POST /api/practice-sessions` to record practice.
6. `GET /api/practice-sessions` to retrieve practice history.
7. `GET /api/practice-sessions/stats` to retrieve aggregate statistics.
8. Use `PUT` or `DELETE` on a practice-session resource as needed.
