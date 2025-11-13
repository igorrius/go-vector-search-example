# Go Vector Search Example

## Description

This project is an example of how to build a vector search application using Go. It uses Google's AI Platform for embedding generation and summarization, and Typesense for vector storage and search.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Access to Google Cloud Platform with the AI Platform APIs enabled

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