# Task Manager API

A RESTful API for managing tasks, built with Go, the Gin framework, and MongoDB.

This project provides CRUD (Create, Read, Update, Delete) operations for a task management system, secured with JWT authentication and role-based access control.

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.18 or newer)
- MongoDB Atlas Account (or local instance)

## Getting Started

1.  **Clone the repository** (or navigate to the existing project directory).

2.  **Navigate to the API directory:**
    ```sh
    cd task_manager
    ```

3.  **Install dependencies:**
    ```sh
    go mod tidy
    ```

4.  **Run the server:**
    ```sh
    go run main.go
    ```
    The server will start and listen on `http://localhost:8080`.

## Authentication & Authorization

This API uses JSON Web Tokens (JWT) for authentication.
- **Public Endpoints:** Register, Login
- **Protected Endpoints:** All task operations require a valid JWT token in the `Authorization` header (Bearer scheme).
- **Role-Based Access:**
    - **User:** Can view tasks.
    - **Admin:** Can view, create, update, and delete tasks, and promote users.
    - The first registered user is automatically assigned the **Admin** role.

## API Endpoints

### User Management

#### Register a New User
-   **URL:** `/register`
-   **Method:** `POST`
-   **Description:** Creates a new user account.
-   **Body:**
    ```json
    {
        "username": "johndoe",
        "password": "secretpassword"
    }
    ```

#### Login
-   **URL:** `/login`
-   **Method:** `POST`
-   **Description:** Authenticates a user and returns a JWT token.
-   **Body:**
    ```json
    {
        "username": "johndoe",
        "password": "secretpassword"
    }
    ```
-   **Response:**
    ```json
    {
        "token": "eyJhbGciOiJIUzI1NiIsInR..."
    }
    ```

#### Promote User (Admin Only)
-   **URL:** `/promote`
-   **Method:** `POST`
-   **Headers:** `Authorization: Bearer <token>`
-   **Body:**
    ```json
    {
        "user_id": "6554c7f8a1b2c3d4e5f6g7h8"
    }
    ```

### Task Operations

#### Get All Tasks
-   **URL:** `/tasks`
-   **Method:** `GET`
-   **Headers:** `Authorization: Bearer <token>`
-   **Description:** Retrieves a list of all tasks.

#### Get Task by ID
-   **URL:** `/tasks/:id`
-   **Method:** `GET`
-   **Headers:** `Authorization: Bearer <token>`
-   **Description:** Retrieves a single task by its ID.

#### Add a New Task (Admin Only)
-   **URL:** `/tasks`
-   **Method:** `POST`
-   **Headers:** `Authorization: Bearer <token>`
-   **Body:**
    ```json
    {
        "title": "New Task",
        "description": "Description here",
        "due_date": "2025-12-01T15:00:00Z",
        "status": "Pending"
    }
    ```

#### Update a Task (Admin Only)
-   **URL:** `/tasks/:id`
-   **Method:** `PUT`
-   **Headers:** `Authorization: Bearer <token>`
-   **Body:**
    ```json
    {
        "title": "Updated Title"
    }
    ```

#### Delete a Task (Admin Only)
-   **URL:** `/tasks/:id`
-   **Method:** `DELETE`
-   **Headers:** `Authorization: Bearer <token>`
-   **Description:** Deletes a task by its ID.

