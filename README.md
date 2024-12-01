# Upfluence Coding Challenge

This project is an HTTP API designed to process, aggregate, and provide data sent by the Upfluence public API endpoint streaming social media posts.

## Architectural Design

The project is organized to follow the principles of the "Screaming" architecture. It aims to scream to the developer what they are working on. It abstracts business logic into features and provides low coupling between business logic and incoming/outgoing traffic.

### Folder Organization

* `internal/app/`: Contains code to initialize and launch the server
* `internal/config/`: Simple module to read server configuration from a JSON file
* `internal/features/`: Folder containing all features
    * `internal/features/{feature_name}/`: Actual implementation of the feature containing:
        * Feature interface definition
        * A controller to perform business logic
        * A repository to access data
        * A model
* `internal/interfaces/`: Contains modules to handle incoming traffic and access external services
    * `interfaces/http/`: Module to handle incoming requests. It provides a router using the Gin framework
    * `interfaces/sse/`: Implementation of the SSE client that allows connecting to the streaming server and broadcasting data.
* `internal/logs/`: Module that provides a naive JSON logger


### Architecture Principles

* **Dependency Injection**: Every dependencies is injected. This allows low coupling and eases the testing process.
* **Interfaces <-> Business Logic Decoupling**: By decoupling the interfaces, we isolate the business logic, which makes the project easier to evolve.
* **Simplicity**: Everything has to be as simple as possible.
* **Separation of Concerns**: Each feature is isolated and only provides the necessary components. They don't expose internal implementation.



## What whould I change

## YAML for config

I would use Viper as a config manager to parse the configuration from YAML and having environment variables overriding.

## Logging

I would use a well know library for logging. Indeed my proposal is a naive implementation.

Currentlty there is only the possibilty to print log formatted in JSON in the current terminal, but we can image to provide a way to 
write logs in files with various format.


## What can be added ?

### CORS

If the API will be used by browsers, I would add a CORS middleware to enforce the security. In the context of server to server communication (which I assumed), CORS middleware is not necessary.

### Configuration validation

We can add a mechanism to validates configuration data such as checking ig the `sse_client_config.server_url` is actualy a URL. Open source modules already exists to do that.