## ADDED Requirements
### Requirement: Upload Document API
The system SHALL provide an API endpoint to upload documents.

#### Scenario: Successful document upload
- **WHEN** a user uploads a document via a POST request to `/upload`
- **THEN** the system SHALL respond with a 200 OK status and a JSON object containing a success message.

### Requirement: Search API
The system SHALL provide an API endpoint to search for documents.

#### Scenario: Successful search
- **WHEN** a user sends a GET request to `/search` with a query parameter `q`
- **THEN** the system SHALL respond with a 200 OK status and a JSON array of matching documents.