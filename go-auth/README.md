# Go Authentication Service

A simple Authentication and Authorization service built with Go and the Gin framework. It demonstrates user registration, login with JWT generation, and protecting routes using middleware.

## 🚀 Tech Stack

*   **Language:** Go
*   **Framework:** [Gin](https://github.com/gin-gonic/gin)
*   **Authentication:** [JWT-Go](https://github.com/dgrijalva/jwt-go)
*   **Security:** [Bcrypt](https://golang.org/x/crypto/bcrypt) (for password hashing)

## 📂 Project Structure

The project follows a clean MVC-like architecture:

```
go-auth/
├── controllers/      # Request handlers (Register, Login, Secure)
├── middleware/       # JWT validation middleware
├── models/           # Data structures (User)
├── routes/           # Route definitions and setup
├── main.go           # Entry point
├── go.mod            # Dependencies
└── README.md         # Documentation
```

## 🛠️ Getting Started

### Prerequisites

*   [Go](https://golang.org/dl/) (version 1.25 or higher) installed.

### Installation

1.  Navigate to the project directory:
    ```bash
    cd go-auth
    ```

2.  Install dependencies:
    ```bash
    go mod tidy
    ```

### Running the Application

Start the server:
```bash
go run main.go
```
The server will start on `http://localhost:8080`.

## 🔌 API Endpoints

### 1. Register User
Creates a new user account.

*   **URL:** `/register`
*   **Method:** `POST`
*   **Body:**
    ```json
    {
        "email": "user@example.com",
        "password": "password123"
    }
    ```
*   **Response:** `200 OK`
    ```json
    {
        "message": "User registered successfully"
    }
    ```

### 2. Login
Authenticates a user and returns a JWT token.

*   **URL:** `/login`
*   **Method:** `POST`
*   **Body:**
    ```json
    {
        "email": "user@example.com",
        "password": "password123"
    }
    ```
*   **Response:** `200 OK`
    ```json
    {
        "message": "User logged in successfully",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```

### 3. Secure Route
A protected endpoint that requires a valid JWT token.

*   **URL:** `/secure`
*   **Method:** `GET`
*   **Headers:**
    *   `Authorization`: `Bearer <your_token_here>`
*   **Response:** `200 OK`
    ```json
    {
        "message": "This is a secure route"
    }
    ```

## 🧪 Testing

You can test the API using the provided PowerShell script (if you are on Windows) or manually using `curl`.

**Using the Test Script:**
```powershell
# From the root a2sv-backend directory
powershell -ExecutionPolicy Bypass -File .\scripts\test_go_auth.ps1
```

**Using Curl:**
```bash
# 1. Register
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{"email":"test@test.com","password":"123"}'

# 2. Login (Copy the token from response)
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"email":"test@test.com","password":"123"}'

# 3. Access Secure Route
curl -X GET http://localhost:8080/secure -H "Authorization: Bearer <TOKEN>"
```
