# Task Manager API

A simple RESTful API for managing tasks, built with Go and the Gin framework.

This project provides basic CRUD (Create, Read, Update, Delete) operations for a task management system. It is intended as a simple demonstration of building a web service in Go.

**Note:** This API uses an in-memory data store. All tasks will be reset every time the application is restarted.

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.18 or newer)

## Getting Started

1.  **Clone the repository** (or navigate to the existing project directory).

2.  **Navigate to the API directory:**
    ```sh
    cd task_manager_api
    ```

3.  **Install dependencies:**
    Go will automatically handle dependencies when you build or run the project.

4.  **Run the server:**
    ```sh
    go run .
    ```
    The server will start and listen on `http://localhost:8080`.

## API Endpoints

The following endpoints are available:

### Get All Tasks

-   **URL:** `/tasks`
-   **Method:** `GET`
-   **Description:** Retrieves a list of all tasks.
-   **Example `curl`:**
    ```sh
    curl http://localhost:8080/tasks
    ```

### Get Task by ID

-   **URL:** `/tasks/:id`
-   **Method:** `GET`
-   **Description:** Retrieves a single task by its ID.
-   **Example `curl`:**
    ```sh
    curl http://localhost:8080/tasks/1
    ```

### Add a New Task

-   **URL:** `/tasks`
-   **Method:** `POST`
-   **Description:** Creates a new task.
-   **Headers:** `Content-Type: application/json`
-   **Body (raw JSON):**
    ```json
    {
        "id": "4",
        "title": "New Task from API",
        "description": "A task created via POST request.",
        "duedate": "2025-12-01T15:00:00Z",
        "status": "Pending"
    }
    ```
-   **Example `curl`:**
    ```sh
    curl -X POST http://localhost:8080/tasks \
    -H "Content-Type: application/json" \
    -d '{"id": "4", "title": "New Task", "description": "A new task", "duedate": "2025-12-01T15:00:00Z", "status": "Pending"}'
    ```

### Update a Task

-   **URL:** `/tasks/:id`
-   **Method:** `PUT`
-   **Description:** Updates an existing task's title or description.
-   **Headers:** `Content-Type: application/json`
-   **Body (raw JSON):**
    ```json
    {
        "title": "Updated Task Title",
        "description": "This task has been updated."
    }
    ```
-   **Example `curl`:**
    ```sh
    curl -X PUT http://localhost:8080/tasks/1 \
    -H "Content-Type: application/json" \
    -d '{"title": "Updated Title"}'
    ```

### Delete a Task

-   **URL:** `/tasks/:id`
-   **Method:** `DELETE`
-   **Description:** Deletes a task by its ID.
-   **Example `curl`:**
    ```sh
    curl -X DELETE http://localhost:8080/tasks/1
    ```
