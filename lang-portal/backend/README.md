# Backend for the Lang Portal Project
A hindi language learning school wants to build a sample project that could help the users learn appropriate language.
The portal in its current form, will do the following:
- Act as a launchpad for different study activities.
- Store Language words and sentences, that would aid in learning activities. 
- Act as a record keeping app for tracking progress and performance.

## Technical Requirements
- Backend will be wriiten using Go, with standard libraries.
- The database used will be SQLite3. 
- The API request and response will be in JSON format.
- The API will be stateless, and will not store any persistent data.
- The API will not be versioned. 

## Database Design

table: words
columns: 
   - id: integer
   - hindi: string 
   - scrambled: string
   - hinglish: string
   - english: string
   - created_at: datetime

table: groups
columns: 
   - id: integer
   - group: string
   - created_at: datetime

table: word_groups
columns: 
   - id: integer
   - group_id: integer
   - word_id: integer
   - created_at: datetime

table: study_activities
columns: 
   - id: integer
   - name: string 
   - description: string
   - image: string
   - score: integer
   - created_at: datetime

table: sessions
columns: 
   - id: integer
   - activity_id: integer
   - start_time: datetime
   - end_time: datetime
   - score: integer
   - created_at: datetime

table: session_activities
columns:
   - id: integer
   - session_id: integer
   - activity_id: integer
   - challenge: string
   - answer: string
   - input: string
   - score: integer
   - created_at: datetime


## ER Diagram

```mermaid
---
title: Lang Portal Database Schema
---

erDiagram
    words {
        integer id PK
        string hindi
        string scrambled
        string hinglish
        string english
        datetime created_at
    }

    groups {
        integer id PK
        string group
        datetime created_at
    }

    word_groups ||--o{ words : contains
    word_groups }o--|| groups : has
    word_groups ||--|| sessions : has
    word_groups {
        integer id PK
        integer group_id FK
        integer word_id FK
        datetime created_at
    }

    study_activities ||--o{ sessions : has
    study_activities {
        integer id PK
        string name
        string description
        string image
        integer score
        datetime created_at
    }

    sessions ||--o{ session_activities : have
    sessions {
        integer id PK
        integer activity_id FK
        datetime start_time
        datetime end_time
        integer score
        datetime created_at
    }

    session_activities {
        integer id PK
        integer session_id FK
        integer activity_id FK
        string challenge
        string answer
        string input
        integer score
        datetime created_at
    }
```


## API Design
- [GET] /api/words
    - lists all words

- [GET] /api/words/random
    - this should take a group_id

- [GET] /api/words/search
    - this should take a search term

- [GET] /api/groups
    - lists all groups

- [GET] /api/words/groups/:group-id  
    - lists all words from a group
    - joins words and groups tables based on word_groups table and filters by group_id

- [GET] /api/study-activities 
    - lists all available study activities

- [POST] /api/sessions
  - this should take activity_id
  - this handler should automatically start_time for session

- [POST] /api/session-activity
  - this should take session_id
  - this should take activity_id
  - the challenge should be added to the session_activity table
  - the answer should be added to the session_activity table
  - the input should be added to the session_activity table
  - the score should be added to the session_activity table
  - this should be a single row in the table

- [PUT] /api/sessions
    - this should allow updating the end_time and score of a session
- [GET] /api/sessions 
    - lists details of all sessions
    - pagination is required
    - it should also provide study activity name of each session by joining the study_activities table
- [GET] /api/sessions/:id
    - lists individual session details
    - joins sessions and session_activities tables based on session_id 
    - lists session details including its study activities
    - no pagination is required.

- [DELETE] /api/sessions/
    - this should delete all sessions
    - this should also delete all session_activities associated with the sessions

## Documentation
- Avoid Littering the codebase with comments. 
- Modify the swagger doc with endpoint changes
- The swagge doc should exactly match the API
- No versioning is required, and should not be used in the API
- All APIs are prefix with /api, documentation should be also follow same
- Swagger UI should be updated after api documentation 
- There is always only a single user, no authentication/authorisation is required
- Parameters will be in url/query for GET requests
- Parameters will be in url/query for  requests
- Keep parameters in body for POST/PUT requests