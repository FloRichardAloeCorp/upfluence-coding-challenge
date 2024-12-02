# Upfluence Coding Challenge

This project is an HTTP API designed to process, aggregate, and provide data sent by the Upfluence public API endpoint streaming social media posts.

## Overview

The project aggregates statistics about social media posts streamed by the Upfluence public API streaming endpoint.

The API establishes a single connection to the server and broadcasts the stream to its internal subscribers. This design prevents stream duplication and ensures the system can handle higher loads. Each request creates a subscriber to the client, which communicates data through channels.

Posts streamed may or may not contain the desired dimensions, but the server analyzes all posts. Equivalent dimensions are rendered consistently; for example, likes are always represented with the likes JSON field. This enables generic parsing and processing of events without needing to know from which platform the posts were published.

## Installation

### Requirements

Go version >= 1.21

### Get sources and run the project

```bash
go run main.go
```

### Run with Docker

```bash
docker compose up
```

## Configuration

The API is configurable via a JSON file located at the `API_CONFIG` location (default: `./config.json`).

You can override the default configuration location by setting the `API_CONFIG` environment variable.

### Configuration Reference

```json
{
    "sse_client_config": {
        // Streaming endpoint.
        "server_url": "https://stream.upfluence.co/stream",

        // Maximum number of attempts to reconnect if there is a problem.
        "max_reconnection_attempts": 10
    },
    "router": {
        // Listening port of the server.
        "port": 8080,

        // Gin log mode, can be debug or release.
        "gin_mode": "debug",

        // Timeout if closing the server is too long.
        "shutdown_timeout": 5,

        "analysis_handler_config": {
            // List of accepted dimension values.
            "authorized_dimensions": [
                "likes",
                "comments",
                "favorites",
                "retweets"
            ]
        }
    },
    "logger": {
        // Logger level, can be either INFO or ERROR.
        "level": "INFO",

        // Optionnal ouput path, leave it empty to write to stderr.
        // You can specify a file where the logs will be written.
        "output": ""
    }
}
```

## Documentation

A `swagger` file is available [here](./swagger.yaml).

## Testing

### Unit testing
To run unit tests run the following command:
```bash
make test
```

**NOTE:** The Upfluence public API is mocked during unit tests, so no external services are involved in testing.

Current testing states:

```bash
go clean -testcache
go test -timeout 1m -cover ./...
        github.com/FloRichardAloeCorp/upfluence-coding-challenge                coverage: 0.0% of statements
        github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/app           coverage: 0.0% of statements
        github.com/FloRichardAloeCorp/upfluence-coding-challenge/test/mockings          coverage: 0.0% of statements
ok      github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/config        2.689s  coverage: 100.0% of statements
ok      github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/features/aggregate    10.972s coverage: 96.0% of statements
ok      github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/http       0.478s  coverage: 100.0% of statements
ok      github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/http/middlewares   2.253s  coverage: 100.0% of statements
ok      github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/sse        20.452s coverage: 91.2% of statements
ok      github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs  1.824s  coverage: 100.0% of statements
```

## Architectural Design

The project follows the principles of the "Screaming" architecture. It emphasizes clarity, ensuring the codebase explicitly reflects the business logic. Features are abstracted and decoupled from incoming and outgoing traffic, enabling low coupling and better maintainability.

### Folder Organization

* `internal/app/`: Initializes and launches the server
* `internal/config/`: Module to read server configuration from a JSON file
* `internal/features/`: Contains all features
    * `internal/features/{feature_name}/`: Implementation of a feature, including:
        * Feature interface definition
        * Controller for business logic
        * Repository for data access
        * Model
* `internal/interfaces/`: Handles incoming traffic and external service interactions
    * `interfaces/http/`: Manages incoming requests using the Gin framework
    * `interfaces/sse/`: Implements the SSE client to connect to the streaming server and broadcast data
* `internal/logs/`: Provides a basic JSON logger

### Architecture Principles

* **Dependency Injection**: All dependencies are injected, ensuring low coupling and facilitating testing.
* **Interfaces <-> Business Logic Decoupling**: Decoupling interfaces isolates the business logic, making the project easier to evolve.
* **Simplicity**: Solutions are designed to be as simple as possible.
* **Separation of Concerns**: Features are isolated and expose only the necessary components, keeping internal implementations private.

## What Would I Change?

### YAML for Configuration

I would use Viper as a configuration manager to parse configuration files in YAML format and allow overriding values through environment variables.

### Logging

I would replace the naive implementation with a well-known logging library.

## Potential Additions

### CORS Middleware

If the API will be consumed by browsers, adding a CORS middleware would enhance security. However, for server-to-server communication (as assumed here), it is not required.

### Configuration Validation

Adding validation for configuration data would help catch errors, such as ensuring `sse_client_config.server_url` is a valid URL. Existing open-source modules could handle this.

### Enhanced SSE Client

Currently, the client only processes events prefixed with "data:". It could be extended to handle additional prefixes, ensuring broader compatibility with the streaming protocol.

### Rate limiting

Depending on the deployment context, a rate-limiting middleware would be a valuable addition to prevent high load spikes.
