# Language Portal Backend API Routes

## Words API Endpoints

### Create a Word
- **Method:** POST
- **Endpoint:** `/api/words`
- **Description:** Create a new word in the language portal
- **Request Body:**
  ```json
  {
    "hindi": "string",
    "english": "string",
    "hinglish": "string"
  }
  ```
- **Success Response:** 
  - Status Code: 200 OK
  - Returns the created word object
- **Error Responses:**
  - 400 Bad Request: Invalid word data
  - 500 Internal Server Error: Failed to create word

### List Words
- **Method:** GET
- **Endpoint:** `/api/words`
- **Description:** Retrieve a list of words with optional pagination
- **Query Parameters:**
  - `page`: Page number (default: 1)
  - `pageSize`: Number of words per page (default: 10)
  - `language`: Filter by language (optional)
- **Success Response:**
  - Status Code: 200 OK
  - Returns an array of word objects and total count

### Search Words
- **Method:** GET
- **Endpoint:** `/api/words/search`
- **Description:** Search words based on a query string
- **Query Parameters:**
  - `query`: Search term
  - `language`: Language to search in (optional)
- **Success Response:**
  - Status Code: 200 OK
  - Returns:
    ```json
    {
      "words": [],
      "totalCount": 0
    }
    ```

### Get Random Word
- **Method:** GET
- **Endpoint:** `/api/words/random`
- **Description:** Retrieve a random word from the database
- **Success Response:**
  - Status Code: 200 OK
  - Returns a single word object

### Get Word by ID
- **Method:** GET
- **Endpoint:** `/api/words/:id`
- **Description:** Retrieve a specific word by its unique identifier
- **Success Response:**
  - Status Code: 200 OK
  - Returns the word object
- **Error Responses:**
  - 404 Not Found: Word not found
  - 500 Internal Server Error: Database error

### Get Words by Group
- **Method:** GET
- **Endpoint:** `/api/words/groups/:group-id`
- **Description:** Retrieve all words associated with a specific group
- **Success Response:**
  - Status Code: 200 OK
  - Returns an array of word objects
  - Returns empty array if no words found
- **Error Responses:**
  - 400 Bad Request: Invalid group ID format
  - 500 Internal Server Error: Failed to retrieve words

### Update Word
- **Method:** PUT
- **Endpoint:** `/api/words/:id`
- **Description:** Update an existing word
- **Request Body:**
  ```json
  {
    "hindi": "string",
    "english": "string",
    "hinglish": "string"
  }
  ```
- **Success Response:**
  - Status Code: 200 OK
  - Returns the updated word object
- **Error Responses:**
  - 400 Bad Request: Invalid word data
  - 404 Not Found: Word not found
  - 500 Internal Server Error: Failed to update word

### Delete Word
- **Method:** DELETE
- **Endpoint:** `/api/words/:id`
- **Description:** Delete a word by its unique identifier
- **Success Response:**
  - Status Code: 200 OK
- **Error Responses:**
  - 404 Not Found: Word not found
  - 500 Internal Server Error: Failed to delete word

## Groups API Endpoints

### Create a Group
- **Method:** POST
- **Endpoint:** `/api/groups`
- **Description:** Create a new word group
- **Request Body:**
  ```json
  {
    "name": "string",
    "description": "string"
  }
  ```
- **Success Response:**
  - Status Code: 200 OK
  - Returns the created group object

### List Groups
- **Method:** GET
- **Endpoint:** `/api/groups`
- **Description:** Retrieve a list of word groups
- **Query Parameters:**
  - `page`: Page number (default: 1)
  - `pageSize`: Number of groups per page (default: 10)
- **Success Response:**
  - Status Code: 200 OK
  - Returns an array of group objects and total count

### Get Group by ID
- **Method:** GET
- **Endpoint:** `/api/groups/:id`
- **Description:** Retrieve a specific group by its unique identifier
- **Success Response:**
  - Status Code: 200 OK
  - Returns the group object
- **Error Responses:**
  - 404 Not Found: Group not found
  - 500 Internal Server Error: Database error

### Update Group
- **Method:** PUT
- **Endpoint:** `/api/groups/:id`
- **Description:** Update an existing group
- **Request Body:**
  ```json
  {
    "name": "string",
    "description": "string"
  }
  ```
- **Success Response:**
  - Status Code: 200 OK
  - Returns the updated group object
- **Error Responses:**
  - 400 Bad Request: Invalid group data
  - 404 Not Found: Group not found
  - 500 Internal Server Error: Failed to update group

### Delete Group
- **Method:** DELETE
- **Endpoint:** `/api/groups/:id`
- **Description:** Delete a group by its unique identifier
- **Success Response:**
  - Status Code: 200 OK
- **Error Responses:**
  - 404 Not Found: Group not found
  - 500 Internal Server Error: Failed to delete group

## System Endpoints

### Health Check
- **Method:** GET
- **Endpoint:** `/health`
- **Description:** Check the health and status of the API
- **Success Response:**
  ```json
  {
    "status": "healthy"
  }
  ```

### API Information
- **Method:** GET
- **Endpoint:** `/api`
- **Description:** Get basic information about the Language Portal Backend
- **Success Response:**
  ```json
  {
    "status": "healthy",
    "message": "Language Portal Backend is running successfully!"
  }
  ```

## Authentication and Authorization
*Authentication mechanisms will be added in future updates.*

## Error Handling
- All error responses include a descriptive error message
- Consistent error response format across all endpoints
- Appropriate HTTP status codes for different error scenarios

## Versioning
*API versioning strategy will be implemented in future updates.*
