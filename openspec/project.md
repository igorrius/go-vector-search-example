# Project Context

## Purpose
This is an example project to demonstrate vector search capabilities with Go and Typesense.

## Tech Stack
- **Programming Language:** Go
- **Vector Search:** Typesense
- **AI/ML:** Google Generative AI
- **Containerization:** Docker

## Project Conventions

### Code Style
Code is formatted using `gofmt` and linted with `golangci-lint` using the default settings.

### Architecture Patterns
The project follows a **Clean Architecture** pattern, with a clear separation of concerns:
- **`domain`:** Core business logic and entities.
- **`app`:** Application-specific logic and interfaces.
- **`infra`:** Concrete implementations of interfaces (e.g., database, external APIs).

### Testing Strategy
Unit tests are required for all business logic in the `domain` and `app` layers. The `testify` library is used for assertions.

### Git Workflow
This project follows the GitHub Flow. New features or fixes should be developed in branches created from `main`. When ready, open a pull request for review and merging.

## Domain Context
In the context of this project, a "document" refers to a small chunk of text that has been split from a larger corpus. Each document is associated with a vector embedding for semantic search.

## Important Constraints
The entire application, including the Typesense database, must be runnable locally using Docker and the provided `docker-compose.yml` file. This ensures a consistent development and testing environment.

## External Dependencies
- **Typesense:** Used for vector search and document storage.
- **Google Generative AI:** Used for generating embeddings and summarizing content.
