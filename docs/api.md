# API Documentation

## Overview
The RiffLog API is a RESTful backend service for tracking guitar practice sessions. It allows users to register an account, authenticate using JSON Web Tokens (JWT), record practice sessions, view practice history, and retrieve summary statistics about their practice habits.

The API is built with Go using the Gin web framework and follows a layered architecture consisting of HTTP handlers, services, and repositories backed by a PostgreSQL database. All endpoints communicate using JSON, and authenticated endpoints require a valid JWT Bearer token.

The API currently supports the following features:

- User registration and authentication
- JWT-protected endpoints
- Browse available practice skills
- Create, update, and delete practice sessions
- List practice sessions with optional filtering
- View aggregate practice statistics

This document describes each available endpoint, the required request format, expected responses, and common error conditions to assist developers integrating with the API.

## Base URL

By default, the development server listens on:

```text
http://localhost:8080
```

## Authentication
Protected endpoints require a JSON Web Token (JWT) obtained from the `POST /login` endpoint.

Include the token in the `Authorization` header using the Bearer authentication scheme:

```text
Authorization: Bearer <JWT>
```

The authenticated user is determined from the JWT. Client requests should not include a user ID, as ownership is enforced by the server.
## Endpoints

### Authentication

#### POST /register

##### Summary
Creates a new Rifflog user.

##### Authentication
None

##### Request Headers
```http
Content-Type: application/json
```

##### Request Body
|Field              |Type               |Description                        |
|:---               |:---               |:---                               |
|email              |string             |User's email addresss              |
|password           |string             |User's password                    |

Example:
```json
{
    "email": "jimmy_user@usermail.com",
    "password": "jimmyRules3"
}
```

##### Success Response
**Status** `201 Created`
```json
{
    "id": 283,
    "email": "jimmy_user@usermail.com",
    "created_at": "2026-07-09T18:00:00Z"
}
```

##### Error Responses
|Status                 |Meaning                                    |
|:---                   |:---                                       |
|400                    |Invalid user data                          |
|500                    |Internal server error                      |

#### POST /login

##### Summary
Log in a user and issue an authentication token.

##### Authentication
None

##### Request Headers
```http
Content-Type: application/json
```

##### Request Body
|Field              |Type               |Description                        |
|:---               |:---               |:---                               |
|email              |string             |User's email addresss              |
|password           |string             |User's password                    |

Example:
```json
{
    "email": "jimmy_user@usermail.com",
    "password": "jimmyRules3"
}
```

##### Success Response
**Status** `200 OK`
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
    "user": {
        "id": 283,
        "email": "jimmy_user@usermail.com",
        "created_at": "2026-07-09T18:00:00Z"
    }
}
```

##### Error Responses
|Status                 |Meaning                                    |
|:---                   |:---                                       |
|400                    |Invalid email address                      |
|401                    |Invalid password                           |
|404                    |User not found                             |
|500                    |Internal server error                      |

### Skills

#### GET /skills

##### Summary
Retrieves a list of skills to practice during sessions.

##### Authentication
None

##### Request Headers
```http
Content-Type: application/json
```

##### Request Body
None

##### Success Response
**Status** `200 OK`
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

##### Error Responses
|Status                 |Meaning                                    |
|:---                   |:---                                       |
|500                    |Internal server error                      |

### Practice Sessions

#### POST /practice-sessions

##### Summary
Creates a new practice session for the authenticated user.

##### Authentication
Required (Bearer Token)

##### Request Headers
```http
Authorization: Bearer <JWT>
Content-Type: application/json
```

##### Request Body
|Field              |Type               |Description                        |
|:---               |:---               |:---                               |
|skill_id           |integer            |ID of the practiced skill          |
|duration_minutes   |integer            |Length of the session in minutes   |
|practiced_at       |string (RFC3339)   |When the session occured           |
|notes              |string             |Optional practice notes            |

Example:
```json
{
    "skill_id": 2,
    "duration_minutes": 30,
    "practiced_at": "2026-07-09T18:00:00Z",
    "notes": "Worked on scales"
}
```

##### Success Response
**Status:** `201 Created`
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

##### Error Responses
|Status                 |Meaning                                    |
|:---                   |:---                                       |
|400                    |Invalid session data                       |
|401                    |Missing or invalid authentication token    |
|404                    |Referenced skill not found                 |
|500                    |Internal server error                      |

#### GET /practice-sessions

##### Summary
Retrieve practice sessions for authenticated user.  Can be filtered with query parameters.

##### Authentication
Required (Bearer Token)

##### Request Headers
```http
Authorization: Bearer <JWT>
Content-Type: application/json
```

##### Request Query Parameters
|Parameter          |Type               |Description                        |
|:---               |:---               |:---                               |
|skill              |integer            |Only return sessions for this skill|
|from               |date               |Sessions on or after this date     |
|to                 |date               |Sessions on or before this date    |

##### Request Body
None

##### Success Response
**Status** `200 OK`
```json
[
    {
        "id": 3748,
        "skill_id": 2,
        "skill_name": "Scales",
        "skill_description": "Memorize note locations and scale patterns.",
        "duration_minutes": 25,
        "notes": "Scales practice, my fingers hurt",
        "practiced_at": "2026-07-09T18:00:00Z",
        "created_at": "2026-07-09T19:00:00Z",
        "user_id": "22"
    }
    {
        "id": 3921,
        "skill_id": 2,
        "skill_name": "Scales",
        "skill_description": "Memorize note locations and scale patterns.",
        "duration_minutes": 25,
        "notes": "More scales practice, my fingers hurt a little less",
        "practiced_at": "2026-08-09T18:00:00Z",
        "created_at": "2026-08-09T19:00:00Z",
        "user_id": 22
    }
]
```

##### Error Responses
|Status                 |Meaning                                    |
|:---                   |:---                                       |
|400                    |Invalid query parameters                   |
|401                    |Missing or invalid authentication token    |
|404                    |User invalid or not found                  |
|500                    |Internal server error                      |

#### GET /practice-sessions/stats

##### Summary
Retrieve a set of statistics for the authenticate user.
Statistics reported:
- Total minutes practiced
- Total sessions practiced
- Most practiced skill (calculated by minutes)
- Longest session practiced

##### Authentication
Required (Bearer Token)

##### Request Headers
```http
Authorization: Bearer <JWT>
Content-Type: application/json
```

##### Request Body
None

##### Success Response
**Status** `200 OK`
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

##### Error Responses
|Status                 |Meaning                                    |
|:---                   |:---                                       |
|401                    |Missing or invalid authentication token    |
|500                    |Internal server error                      |

#### PUT /practice-sessions/{id}

##### Summary
Update a session previously created by the authenticated user.

##### Authentication
Required (Bearer Token)

##### Request Headers
```http
Authorization: Bearer <JWT>
Content-Type: application/json
```

##### Request Path Parameters
|Parameter                  |Description                                    |
|:---                       |:---                                           |
|id                         |Practice session ID                            |

##### Request Body
|Field              |Type               |Description                        |
|:---               |:---               |:---                               |
|skill_id           |integer            |ID of the practiced skill          |
|duration_minutes   |integer            |Length of the session in minutes   |
|practiced_at       |string (RFC3339)   |When the session occured           |
|notes              |string             |Optional practice notes            |

Example:
```json
{
    "skill_id": 2,
    "duration_minutes": 30,
    "practiced_at": "2026-07-09T18:00:00Z",
    "notes": "Worked on scales"
}
```

##### Success Response
**Status** `200 OK`
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

##### Error Responses
|Status                 |Meaning                                    |
|:---                   |:---                                       |
|400                    |Invalid or missing id parameter            |
|401                    |Missing or invalid authentication token    |
|404                    |Requested data not found                   |
|500                    |Internal server error                      |

#### DELETE /practice-sessions/{id}

##### Summary
Delete a session previously created by the authenticated user.

##### Authentication
Required (Bearer Token)

##### Request Headers
```http
Authorization: Bearer <JWT>
Content-Type: application/json
```

##### Request Path Parameters
|Parameter                  |Description                                    |
|:---                       |:---                                           |
|id                         |Practice session ID                            |

##### Request Body
None

##### Success Response
**Status** `200 OK`
```json
{
    "message": "practice session 3842 deleted"
}
```

##### Error Responses
|Status                 |Meaning                                    |
|:---                   |:---                                       |
|400                    |Invalid or missing id parameter            |
|401                    |Missing or invalid authentication token    |
|404                    |Requested data not found                   |
|500                    |Internal server error                      |

## Error Response Format

Unless otherwise noted, failed requests return a JSON response in the following format:

```json
{
    "error": "description of the error"
}