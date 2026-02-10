# Backend Template (Golang)

This is a clean backend template built with Golang, featuring a robust authentication and user management system.

## Features

- **Authentication**: Login, Register, Email Verification, Password Reset.
- **User Management**: Profile updates, Role-based access control (Admin/User).
- **Core Structure**:
    - `controllers`: Handlers for API endpoints.
    - `services`: Business logic layer.
    - `models`: Database structs.
    - `dto`: Data Transfer Objects for API requests/responses.
    - `middleware`: Auth and CORS middleware.
    - `utils`: Helper functions (Email, JWT, etc.).
    - `config`: Database and environment configuration.

## Setup

1.  **Environment Variables**:
    Copy `.env.example` (if available) or create `.env` based on `config/load_env.go`.

2.  **Dependencies**:
    ```bash
    go mod tidy
    ```

3.  **Run Application**:
    ```bash
    go run main.go
    ```

## API Endpoints

-   `POST /register`
-   `POST /login`
-   `POST /forgot-password`
-   `POST /reset-password`
-   `GET  /verify`
-   `GET  /profile` (Protected)
-   `PUT  /profile/update` (Protected)
