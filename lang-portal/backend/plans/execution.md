# Language Portal Backend - Execution Plan

## General Instructions
- Only follow the database schema/structure defined in "lang-portal/README.md"
- At each phase follow the conventions defined in "code-conventions.md"

## Project Development Phases

### Phase 0: Project Setup and Initial Configuration
- [x] Create project directory structure
  - [x] Create `cmd/` directory for main application
    - [x] Create `server/main.go` entry point
  - [x] Create `pkg/` directories
    - [x] `models/`
    - [x] `handlers/`
    - [x] `services/`
    - [x] `repository/`
  - [x] Create `internal/` directories
    - [x] `config/`
    - [x] `middleware/`
  - [x] Create `db/` directories
    - [x] `migrations/`
    - [x] `seeds/`
  - [x] Create `scripts/` directory
  - [x] Create `tests/` directory

- [x] Initialize Go Module and Dependencies
  - [x] Run `go mod init`
  - [x] Add core dependencies
    - [x] Echo framework
    - [x] SQLite driver
    - [x] Migration tool
    - [x] Logging library
  - [x] Configure `go.mod` and `go.sum`
  - [x] Run `go mod tidy`

- [x] Database Configuration
  - [x] Design initial database schema
  - [x] Create migration scripts for:
    - [x] Words table
    - [x] Groups table
    - [x] Study Activities table
    - [x] Sessions table
    - [x] Session Activities table
  - [x] Prepare seed data CSV files
  - [x] Create database initialization script

- [x] Server Configuration
  - [x] Create main application entry point
  - [x] Configure Echo framework
  - [x] Set up server to run on port 3000
  - [x] Implement graceful shutdown
  - [x] Add basic logging middleware
  - [x] Set up error handling middleware
  - [x] Create a /api endpoint to show the server is running
  - [x] The root /api endpoint should return a "healthy" status

- [x] Development Tooling
  - [x] Set up linting (golangci-lint)
  - [x] Configure CI/CD pipeline
  - [x] Add Makefile for common tasks

- [x] Implement graceful shutdown
- [x] Start and run the Server on port 3000

### Phase 1: Words Module Development
#### 1.1 Word Model
- [x] Define Word struct in `pkg/models/word.go`
  - [x] Add fields: ID, Hindi, Scrambled, Hinglish, English
  - [x] Implement validation methods
    - [x] Validate word text length
    - [x] Validate language constraints
  - [x] Implement JSON marshaling/unmarshaling
  - [x] Add database tags
  - [x] Implement custom validation logic

- [x] Write Word model tests
  - [x] Test struct creation
  - [x] Test validation logic
  - [x] Test JSON serialization

#### 1.2 Word Repository
- [x] Create Word repository in `pkg/repository/word_repository.go`
  - [x] Implement Create method
  - [x] Implement GetByID method
  - [x] Implement Update method
  - [x] Implement Delete method
  - [x] Implement List method with pagination
  - [x] Implement Search method
  - [x] Implement GetRandom method
  - [x] Implement GetByLanguage method
  - [x] Add error handling for each method

- [x] Write Word repository tests
  - [x] Test CRUD operations
  - [x] Test pagination
  - [x] Test search functionality
  - [x] Test random word retrieval
  - [x] Test language-specific retrieval
  - [x] Test error scenarios

#### 1.3 Word Service
- [x] Create Word service in `pkg/services/word_service.go`
  - [x] Define service interface
  - [x] Implement service methods
    - [x] Word creation with validation
    - [x] Word retrieval methods
    - [x] Word search and filtering
    - [x] Random word selection
    - [x] Language-specific word retrieval
  - [x] Implement business logic
  - [x] Add comprehensive error handling

- [x] Write Word service tests
  - [x] Test business logic
  - [x] Test error scenarios
  - [x] Test service method interactions

#### 1.4 Word Handlers
- [x] Create Word handlers in `pkg/handlers/word_handler.go`
  - [x] Implement HTTP handlers for:
    - [x] GET /api/words (list words)
    - [x] GET /api/words/:id (get word details)
    - [x] GET /api/words/random (get random word)
    - [x] POST /api/words (create word)
    - [x] PUT /api/words/:id (update word)
    - [x] DELETE /api/words/:id (delete word)
  - [x] Add request validation
  - [x] Implement error responses
  - [x] Add logging

- [x] Write handler integration tests
  - [x] Test each endpoint
  - [x] Test request validation
  - [x] Test error handling

- [x] Start the server and run it on port 3000

### Phase 2: Groups Module Development
#### 2.1 Group Model
- [x] Define Group struct in `pkg/models/group.go`
  - [x] Add fields: ID, Group, CreatedAt, Description
  - [x] Implement validation methods
    - [x] Validate group name length
  - [x] Implement JSON marshaling/unmarshaling
  - [x] Add database tags

- [x] Write Group model tests
  - [x] Test struct creation
  - [x] Test validation logic
  - [x] Test JSON serialization

#### 2.2 Group Repository
- [x] Create Group repository in `pkg/repository/group_repository.go`
  - [x] Implement Create method
  - [x] Implement GetByID method
  - [x] Implement Update method
  - [x] Implement Delete method
  - [x] Implement List method with pagination
  - [x] Implement Search method
  - [x] Add error handling for each method

- [x] Write Group repository tests
  - [x] Test CRUD operations
  - [x] Test pagination
  - [x] Test search functionality
  - [x] Test error scenarios

#### 2.3 Group Service
- [x] Create Group service in `pkg/services/group_service.go`
  - [x] Define service interface
  - [x] Implement service methods
    - [x] Group creation with validation
    - [x] Group retrieval methods
    - [x] Group search and filtering
  - [x] Implement business logic
  - [x] Add comprehensive error handling

- [x] Write Group service tests
  - [x] Test business logic
  - [x] Test error scenarios
  - [x] Test service method interactions

#### 2.4 Group Handlers
- [x] Create Group handlers in `pkg/handlers/group_handler.go`
  - [x] Implement HTTP handlers for:
    - [x] GET /api/groups (list groups)
    - [x] GET /api/groups/:id (get group details)
    - [x] GET /api/groups/:id/words (get words in group)
    - [x] GET /api/groups/:id/words/random (get random words in group)
    - [x] POST /api/groups (create group)
    - [x] PUT /api/groups/:id (update group)
    - [x] DELETE /api/groups/:id (delete group)
  - [x] Add request validation
  - [x] Implement error responses
  - [x] Add logging

- [x] Write handler integration tests
  - [x] Test each endpoint
  - [x] Test request validation
  - [x] Test error handling

### Phase 4: Study Activities Module Development
#### 4.1 Study Activities Model
- [x] Define StudyActivity struct in `pkg/models/study_activity.go`
  - [x] Add fields: ID, Name, Description, Image, Score
  - [x] Implement validation methods
  - [x] Implement JSON marshaling/unmarshaling
  - [x] Add database tags

- [x] Write Study Activities model tests
  - [x] Test struct creation
  - [x] Test validation logic
  - [x] Test JSON serialization

#### 4.2 Study Activities Repository
- [x] Create Study Activities repository
  - [x] Implement GetAll method
  - [x] Add error handling for method

- [x] Write Study Activities repository tests
  - [x] Test retrieval of study activities
  - [x] Test error scenarios

#### 4.3 Study Activities Service
- [x] Create Study Activities service
  - [x] Implement GetStudyActivities method
  - [x] Add validation and error handling

- [x] Write Study Activities service tests
  - [x] Test business logic
  - [x] Test error scenarios
  - [x] Test service method interactions

#### 4.4 Study Activities Handlers
- [x] Create Study Activities handlers
  - [x] Implement GET /api/study-activities endpoint
  - [x] Add request validation
  - [x] Implement error responses

### Phase 5: Sessions Module Development
#### 5.1 Sessions Model
- [x] Define Session struct in `pkg/models/session.go`
  - [x] Add fields: ID, ActivityID, GroupID, StartTime, EndTime, Score
  - [x] Implement validation methods
  - [x] Implement JSON marshaling/unmarshaling
  - [x] Add database tags

- [x] Write Sessions model tests
  - [x] Test struct creation
  - [x] Test validation logic
  - [x] Test JSON serialization
  - [x] Test session duration calculation
  - [x] Test session completion status

#### 5.2 Sessions Repository
- [x] Create Sessions repository
  - [x] Implement Create method
  - [x] Implement Update method
  - [x] Implement List method with pagination
  - [x] Implement GetByID method
  - [x] Implement methods to start/end session
  - [x] Add error handling for each method

- [x] Write Sessions repository tests
  - [x] Test CRUD operations
  - [x] Test pagination
  - [x] Test error scenarios

#### 5.3 Sessions Service
- [x] Create Sessions service
  - [x] Implement session creation logic
  - [x] Implement session management methods
  - [x] Add validation and error handling

- [x] Write Sessions service tests
  - [x] Test business logic
  - [x] Test error scenarios
  - [x] Test service method interactions

#### 5.4 Sessions Handlers
- [x] Create Sessions handlers
  - [x] Implement GET /api/sessions endpoints
  - [x] Implement POST /api/sessions endpoint
  - [x] Implement PUT /api/sessions endpoint
  - [x] Implement DELETE /api/sessions endpoint
  - [x] Add request validation
  - [x] Implement error responses

- [x] Write Sessions handler tests
  - [x] Test each endpoint
  - [x] Test request validation
  - [x] Test error handling

### Phase 6: Session Activities Module Development
#### 6.1 Session Activities Model
- [x] Define SessionActivity struct in `pkg/models/session_activity.go`
  - [x] Add fields: ID, SessionID, ActivityID, Question, Answer, Result, Score
  - [x] Implement validation methods
  - [x] Implement JSON marshaling/unmarshaling
  - [x] Add database tags
  - [x] Implement helper methods for result checking

#### 6.2 Session Activities Repository
- [x] Create Session Activities repository in `pkg/repository/session_activity_repository.go`
  - [x] Implement Create method
  - [x] Implement List method
  - [x] Implement methods to track activity progress
  - [x] Add error handling for each method

- [x] Write Session Activities repository tests
  - [x] Test CRUD operations
  - [x] Test error scenarios

#### 6.3 Session Activities Service
- [x] Create Session Activities service in `pkg/services/session_activity_service.go`
  - [x] Implement logic for adding activities to sessions
  - [x] Implement scoring and result tracking
  - [x] Add validation and error handling

- [x] Write Session Activities service tests
  - [x] Test business logic
  - [x] Test error scenarios
  - [x] Test service method interactions

#### 6.4 Session Activities Handlers
- [x] Create Session Activities handlers in `pkg/handlers/session_activity_handler.go`
  - [x] Implement POST /api/sessions/:id/activity endpoint
  - [x] Implement GET /api/sessions/:id/activity endpoint
  - [x] Add request validation
  - [x] Implement error responses

- [x] Write Session Activities handler tests
  - [x] Test each endpoint
  - [x] Test request validation
  - [x] Test error handling

### Phase 7: Integration and Final Testing
- [ ] Comprehensive integration testing
- [ ] Performance testing
- [ ] Security review
- [ ] Documentation
- [ ] Final code review and refactoring

## Development Principles
- Maintain clean, modular code
- Write comprehensive tests
- Follow SOLID principles
- Minimize external dependencies
- Prioritize code readability

## Micro-Task Guidelines
- Break down each phase into small, manageable tasks
- Use Test-Driven Development (TDD)
- Implement continuous integration
- Perform code reviews
- Maintain clear documentation

## Potential Challenges
- Efficient random word selection
- Pagination performance
- Activity tracking complexity
- Data reset implementation
- Maintaining code modularity