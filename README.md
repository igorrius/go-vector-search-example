# Go Vector Search Example

## Description

This project is a demonstration of how to build a vector search application using Go, following the principles of Clean Architecture. It leverages Google's AI Platform for generating embeddings and providing summarization, with Typesense for efficient vector storage and search.

This example project is designed to showcase a robust, maintainable, and scalable architecture for building AI-powered search applications.

## Project Definition

The core purpose of this project is to provide a working example of a vector search engine. It demonstrates how to:
-   Index text documents by converting them into vector embeddings.
-   Perform semantic searches to find relevant documents based on meaning rather than keywords.
-   Use an AI service to summarize the search results.
-   Structure the application using Clean Architecture principles.

## Responsibilities

The project is divided into several components, each with distinct responsibilities:

-   **API Layer**: Exposes RESTful endpoints for indexing documents and performing searches. It handles HTTP requests and responses.
-   **Application Layer**: Contains the business logic for the application's use cases, such as indexing a document or searching for documents. It orchestrates the domain and infrastructure layers.
-   **Domain Layer**: Represents the core business entities and logic. In this project, it defines the `Document` and the services that operate on it.
-   **Infrastructure Layer**: Implements the external dependencies, such as the database, AI services, and the web server. This layer is responsible for the concrete implementation of the interfaces defined in the application layer.
-   **AI Service**: Responsible for generating vector embeddings for documents and summarizing search results using Google's Generative AI.
-   **Persistence**: Uses Typesense to store and index documents and their vector embeddings for efficient similarity search.

## Architecture

The project follows the **Clean Architecture** pattern, which promotes a separation of concerns and makes the codebase more maintainable and testable. The architecture is divided into three main layers:

-   **Domain**: Contains the core business logic and entities of the application. This is the innermost layer and has no dependencies on other layers.
-   **Application**: Orchestrates the use cases of the application, acting as a bridge between the domain and the infrastructure layers. It defines interfaces that are implemented by the infrastructure layer.
-   **Infrastructure**: Implements the external concerns, such as the database, AI services, and the web server. This layer depends on the application and domain layers.

This architecture ensures that the core business logic is independent of the implementation details of external services, making it easier to test, maintain, and evolve the application.

## Features

-   **Document Indexing**: Upload documents via a RESTful API to be indexed and made searchable.
-   **Vector Search**: Perform semantic searches on the indexed documents.
-   **Content Summarization**: Get AI-powered summaries of search results.
-   **Clean Architecture**: A well-structured codebase with a clear separation of concerns.
-   **Containerized**: The entire application and its dependencies can be run using Docker Compose.

## Tech Stack

-   **Programming Language:** Go
-   **Vector Search:** Typesense
-   **AI/ML:** Google Generative AI
-   **Containerization:** Docker

## Getting Started

### Prerequisites

-   Go 1.21 or higher
-   Docker and Docker Compose
-   Access to Google Cloud Platform with the AI Platform APIs enabled

### Installation

1.  Clone the repository:
    ```sh
    git clone https://github.com/igor-ri/go-vector-search-example.git
    ```
2.  Install the dependencies:
    ```sh
    go mod tidy
    ```
3.  Set up the environment variables:
    ```sh
    cp .env.example .env
    ```
    Then, fill in the `.env` file with your Google Cloud Platform credentials.

## Usage

1.  Start the services:
    ```sh
    docker-compose up -d
    ```
2.  Run the application:
    ```sh
    go run ./cmd/server
    ```

### API Endpoints

#### Index a Document

-   **Endpoint**: `POST /api/v1/documents`
-   **Description**: Indexes a document. The request can be either a JSON payload or multipart form data.

**JSON Payload**

```sh
curl -X POST -H "Content-Type: application/json" -d '{"id": "doc1", "content": "This is a test document."}' http://localhost:8080/api/v1/documents
```

**Multipart Form Data**

```sh
curl -X POST -F "file=@/path/to/your/file.txt" -F "id=doc1" http://localhost:8080/api/v1/documents
```

#### Search Documents

-   **Endpoint**: `GET /api/v1/search`
-   **Description**: Searches for documents based on a query.

**Request**

```sh
curl -X GET "http://localhost:8080/api/v1/search?q=your%20search%20query"
```

**Response**

The response will be a JSON object containing the search results and a summary.

```json
{
  "summary": "This is a summary of the search results.",
  "documents": [
    {
      "id": "doc1",
      "content": "This is a test document.",
      "score": 0.9
    }
  ]
}
```

## Project Conventions

### Code Style

Code is formatted using `gofmt` and linted with `golangci-lint` using the default settings.

### Testing Strategy

Unit tests are required for all business logic in the `domain` and `app` layers. The `testify` library is used for assertions.

### Git Workflow

This project follows the GitHub Flow. New features or fixes should be developed in branches created from `main`. When ready, open a pull request for review and merging.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
