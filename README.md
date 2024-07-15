# Todo App API

This is a simple Todo application API built with Golang. It leverages Couchbase as the database and uses various libraries to enhance functionality, including JWT for authentication, Fiber for the web framework, and Swagger for API documentation.

## Features

- Create, read, update, and delete (CRUD) operations for todos
- User authentication with JWT
- Input validation
- API documentation with Swagger
- Logging with Zap
- Configuration management with Viper

## Technologies Used

- **Couchbase**: Database
- **Fiber**: Web framework
- **JWT**: Authentication
- **Swagger**: API documentation
- **Viper**: Configuration
- **Zap**: Logging
- **Testify**: Testing


## Getting Started

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/harundurmus/go-todo-app.git
    cd todo-app-api
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Configure the application by modifying the `config.yaml` file as needed.

### Running the Application

1. Start the application using Docker Compose:
    ```sh
    docker-compose up
    ```

### API Documentation

Swagger documentation is available at `/swagger/index.html` when the server is running.

### Running Tests

To run the tests, use the following commands:

#### Unit Tests

To run unit tests:
```sh
make unit-test
```

To run unit test coverage:
```sh
make unit-test-coverage
```

To generate a test coverage report:
```sh
make unit-test-coverage-file
```
