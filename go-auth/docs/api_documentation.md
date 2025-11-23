# Go Auth API Documentation

Base URL: `http://localhost:8080`

## Authentication

### Register
Register a new user in the in-memory store.

- **Endpoint:** `/register`
- **Method:** `POST`
- **Content-Type:** `application/json`

**Request Body:**
| Field | Type | Description |
|-------|------|-------------|
| `email` | string | User's email address |
| `password` | string | User's password |

**Example Request:**
```json
{
  "email": "newuser@example.com",
  "password": "securepassword"
}
```

**Success Response (200 OK):**
```json
{
  "message": "User registered successfully"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Invalid request payload"
}
```

---

### Login
Authenticate a user and receive a JWT token.

- **Endpoint:** `/login`
- **Method:** `POST`
- **Content-Type:** `application/json`

**Request Body:**
| Field | Type | Description |
|-------|------|-------------|
| `email` | string | Registered email address |
| `password` | string | Password |

**Example Request:**
```json
{
  "email": "newuser@example.com",
  "password": "securepassword"
}
```

**Success Response (200 OK):**
```json
{
  "message": "User logged in successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": "Invalid email or password"
}
```

---

## Protected Routes

### Secure Endpoint
A sample protected route that requires a valid JWT.

- **Endpoint:** `/secure`
- **Method:** `GET`
- **Headers:** `Authorization: Bearer <token>`

**Success Response (200 OK):**
```json
{
  "message": "This is a secure route"
}
```

**Error Response (401 Unauthorized):**
- Missing Header:
  ```json
  { "error": "Authorization header is required" }
  ```
- Invalid Format:
  ```json
  { "error": "Invalid authorization header" }
  ```
- Invalid/Expired Token:
  ```json
  { "error": "Invalid JWT" }
  ```
