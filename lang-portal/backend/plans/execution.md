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
    - [x] Word-Groups table
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
- [ ] Define Word struct in `pkg/models/word.go`
  - [ ] Add fields: ID, Hindi, Scrambled, Hinglish, English
  - [ ] Implement validation methods
    - [ ] Validate word text length
    - [ ] Validate language constraints
  - [ ] Implement JSON marshaling/unmarshaling
  - [ ] Add database tags
  - [ ] Implement custom validation logic

- [ ] Write Word model tests
  - [ ] Test struct creation
  - [ ] Test validation logic
  - [ ] Test JSON serialization

#### 1.2 Word Repository
- [ ] Create Word repository in `pkg/repository/word_repository.go`
  - [ ] Implement Create method
  - [ ] Implement GetByID method
  - [ ] Implement Update method
  - [ ] Implement Delete method
  - [ ] Implement List method with pagination
  - [ ] Implement Search method
  - [ ] Implement GetRandom method
  - [ ] Implement GetByLanguage method
  - [ ] Add error handling for each method

- [ ] Write Word repository tests
  - [ ] Test CRUD operations
  - [ ] Test pagination
  - [ ] Test search functionality
  - [ ] Test random word retrieval
  - [ ] Test language-specific retrieval
  - [ ] Test error scenarios

#### 1.3 Word Service
- [ ] Create Word service in `pkg/services/word_service.go`
  - [ ] Define service interface
  - [ ] Implement service methods
    - [ ] Word creation with validation
    - [ ] Word retrieval methods
    - [ ] Word search and filtering
    - [ ] Random word selection
    - [ ] Language-specific word retrieval
  - [ ] Implement business logic
  - [ ] Add comprehensive error handling

- [ ] Write Word service tests
  - [ ] Test business logic
  - [ ] Test error scenarios
  - [ ] Test service method interactions

#### 1.4 Word Handlers
- [ ] Create Word handlers in `pkg/handlers/word_handler.go`
  - [ ] Implement HTTP handlers for:
    - [ ] GET /api/words (list words)
    - [ ] GET /api/words/:id (get word details)
    - [ ] GET /api/words/random (get random word)
    - [ ] POST /api/words (create word)
    - [ ] PUT /api/words/:id (update word)
    - [ ] DELETE /api/words/:id (delete word)
  - [ ] Add request validation
  - [ ] Implement error responses
  - [ ] Add logging

- [ ] Write handler integration tests
  - [ ] Test each endpoint
  - [ ] Test request validation
  - [ ] Test error handling

### Phase 2: Groups Module Development
#### 2.1 Group Model
- [ ] Define Group struct in `pkg/models/group.go`
  - [ ] Add fields: ID, Group, CreatedAt
  - [ ] Implement validation methods
    - [ ] Validate group name length
  - [ ] Implement JSON marshaling/unmarshaling
  - [ ] Add database tags

- [ ] Write Group model tests
  - [ ] Test struct creation
  - [ ] Test validation logic
  - [ ] Test JSON serialization

#### 2.2 Group Repository
- [ ] Create Group repository in `pkg/repository/group_repository.go`
  - [ ] Implement Create method
  - [ ] Implement GetByID method
  - [ ] Implement Update method
  - [ ] Implement Delete method
  - [ ] Implement List method with pagination
  - [ ] Implement Search method
  - [ ] Add error handling for each method

- [ ] Write Group repository tests
  - [ ] Test CRUD operations
  - [ ] Test pagination
  - [ ] Test search functionality
  - [ ] Test error scenarios

#### 2.3 Group Service
- [ ] Create Group service in `pkg/services/group_service.go`
  - [ ] Define service interface
  - [ ] Implement service methods
    - [ ] Group creation with validation
    - [ ] Group retrieval methods
    - [ ] Group search and filtering
  - [ ] Implement business logic
  - [ ] Add comprehensive error handling

- [ ] Write Group service tests
  - [ ] Test business logic
  - [ ] Test error scenarios
  - [ ] Test service method interactions

#### 2.4 Group Handlers
- [ ] Create Group handlers in `pkg/handlers/group_handler.go`
  - [ ] Implement HTTP handlers for:
    - [ ] GET /api/groups (list groups)
    - [ ] GET /api/groups/:id (get group details)
    - [ ] GET /api/groups/:id/words (get words in group)
    - [ ] GET /api/groups/:id/words/random (get random words in group)
    - [ ] POST /api/groups (create group)
    - [ ] PUT /api/groups/:id (update group)
    - [ ] DELETE /api/groups/:id (delete group)
  - [ ] Add request validation
  - [ ] Implement error responses
  - [ ] Add logging

- [ ] Write handler integration tests
  - [ ] Test each endpoint
  - [ ] Test request validation
  - [ ] Test error handling

### Phase 3: Word-Groups Module Development
#### 3.1 Word-Groups Model
- [ ] Define WordGroup struct in `pkg/models/word_group.go`
  - [ ] Add fields: ID, GroupID, WordID
  - [ ] Implement validation methods
  - [ ] Implement JSON marshaling/unmarshaling
  - [ ] Add database tags

- [ ] Write Word-Groups model tests

#### 3.2 Word-Groups Repository
- [ ] Create Word-Groups repository
  - [ ] Implement methods to associate words with groups
  - [ ] Implement methods to retrieve words by group
  - [ ] Implement methods to retrieve groups for a word

- [ ] Write Word-Groups repository tests

#### 3.3 Word-Groups Service
- [ ] Create Word-Groups service
  - [ ] Implement business logic for word-group associations
  - [ ] Add validation and error handling

- [ ] Write Word-Groups service tests

#### 3.4 Word-Groups Handlers
- [ ] Create Word-Groups handlers
  - [ ] Implement endpoints for managing word-group relationships

- [ ] Write Word-Groups handler tests

### Phase 4: Study Activities Module Development
#### 4.1 Study Activities Model
- [ ] Define StudyActivity struct in `pkg/models/study_activity.go`
  - [ ] Add fields: ID, Name, Description, Image
  - [ ] Implement validation methods
  - [ ] Implement JSON marshaling/unmarshaling
  - [ ] Add database tags

- [ ] Write Study Activities model tests

#### 4.2 Study Activities Repository
- [ ] Create Study Activities repository
  - [ ] Implement CRUD methods
  - [ ] Implement list and search methods

- [ ] Write Study Activities repository tests

#### 4.3 Study Activities Service
- [ ] Create Study Activities service
  - [ ] Implement business logic
  - [ ] Add validation and error handling

- [ ] Write Study Activities service tests

#### 4.4 Study Activities Handlers
- [ ] Create Study Activities handlers
  - [ ] Implement GET /api/study-activities endpoint
  - [ ] Add request validation
  - [ ] Implement error responses

- [ ] Write Study Activities handler tests

### Phase 5: Sessions Module Development
#### 5.1 Sessions Model
- [ ] Define Session struct in `pkg/models/session.go`
  - [ ] Add fields: ID, ActivityID, GroupID, StartTime, EndTime, Score
  - [ ] Implement validation methods
  - [ ] Implement JSON marshaling/unmarshaling
  - [ ] Add database tags

- [ ] Write Sessions model tests

#### 5.2 Sessions Repository
- [ ] Create Sessions repository
  - [ ] Implement Create method
  - [ ] Implement Update method
  - [ ] Implement List method with pagination
  - [ ] Implement GetByID method
  - [ ] Implement methods to start/end session

- [ ] Write Sessions repository tests

#### 5.3 Sessions Service
- [ ] Create Sessions service
  - [ ] Implement session creation logic
  - [ ] Implement session management methods
  - [ ] Add validation and error handling

- [ ] Write Sessions service tests

#### 5.4 Sessions Handlers
- [ ] Create Sessions handlers
  - [ ] Implement GET /api/sessions endpoints
  - [ ] Implement POST /api/sessions endpoint
  - [ ] Implement PUT /api/sessions endpoint
  - [ ] Add request validation
  - [ ] Implement error responses

- [ ] Write Sessions handler tests

### Phase 6: Session Activities Module Development
#### 6.1 Session Activities Model
- [ ] Define SessionActivity struct in `pkg/models/session_activity.go`
  - [ ] Add fields: ID, SessionID, ActivityID, Question, Answer, Result, Score
  - [ ] Implement validation methods
  - [ ] Implement JSON marshaling/unmarshaling
  - [ ] Add database tags

- [ ] Write Session Activities model tests

#### 6.2 Session Activities Repository
- [ ] Create Session Activities repository
  - [ ] Implement Create method
  - [ ] Implement List method
  - [ ] Implement methods to track activity progress

- [ ] Write Session Activities repository tests

#### 6.3 Session Activities Service
- [ ] Create Session Activities service
  - [ ] Implement logic for adding activities to sessions
  - [ ] Implement scoring and result tracking
  - [ ] Add validation and error handling

- [ ] Write Session Activities service tests

#### 6.4 Session Activities Handlers
- [ ] Create Session Activities handlers
  - [ ] Implement POST /api/sessions/:id/activity endpoint
  - [ ] Implement GET /api/sessions/:id/activity endpoint
  - [ ] Add request validation
  - [ ] Implement error responses

- [ ] Write Session Activities handler tests

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