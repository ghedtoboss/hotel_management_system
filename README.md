# Hotel Management System API

This is a RESTful API for managing hotel operations such as room reservations, customer management, revenue tracking, and user roles. The project is built with Go and designed to be modular, scalable, and easy to use.

## Features

- **Authentication & Authorization**: JWT-based authentication and role-based access control for `admin` and `receptionist` roles.
- **User Management**: Endpoints to manage users, their profiles, and passwords.
- **Room Management**: Create, update, delete, and retrieve hotel rooms and their details.
- **Reservation Management**: Manage reservations with status updates and detailed views.
- **Revenue Tracking**: Calculate daily, monthly, and total revenue with occupancy reports.
- **Swagger Integration**: API documentation is auto-generated and served via Swagger.
- **Modular Design**: Organized with separate controllers, middleware, and route handlers for easy maintenance.

---

## Technologies Used

- **Programming Language**: [Go](https://golang.org/)
- **Router**: [Gorilla Mux](https://github.com/gorilla/mux)
- **Authentication**: JWT-based authorization
- **API Documentation**: [Swaggo](https://github.com/swaggo/swag)
- **Static File Hosting**: Swagger docs served with Go's `http.FileServer`

---

## Installation

### Prerequisites
1. Install [Go](https://golang.org/doc/install) (v1.16+).
2. Install Swagger CLI for generating API documentation:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
