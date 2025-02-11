## Code Conventions / Best Practices for Backend project.

1. Naming Conventions
    - Use CamelCase for exported (public) names
    - Use camelCase for unexported (private) names
    - Use short, descriptive names
    - Avoid abbreviations

2. Code Formatting
    - Always use gofmt or goimports to format code
    - Line length should be around 80-120 characters
    - Use 4 spaces for indentation

3. Error Handling
    - Always handle errors
    - Wrap errors with context using fmt.Errorf() or pkg/errors
    
4. Dependency Management
    - Use go mod for dependency management
    - Pin specific versions of dependencies
    - Keep go.mod and go.sum in version control

5. Logging
    - Use structured logging
    - Avoid fmt.Println() in production

6. Testing
    - Write tests for all packages
    - Use table-driven tests
    - Aim for high test coverage

7. Documentation
    - Write self-documented code
    - Use comments for important parts of the code
    - Use godoc for documentation
    - Keep Comments concis and meaningful

8. Code Quality
    - Use linting tools
    - Use static analysis tools

9. Structure

project-root/
│
├── cmd/                # Main application entry points
│   └── server/
│       └── main.go
│
├── pkg/                # Reusable packages
│   ├── models/         # Data structures
│   ├── handlers/       # HTTP request handlers
│   ├── services/       # Business logic
│   └── repository/     # Database interactions
│
├── internal/           # Private packages (not importable by external projects)
│   ├── config/         # Configuration management
│   └── middleware/     # Request middleware
│
├── db/                 # Database-related files
│   ├── migrations/     # Database schema migrations
│   └── seeds/          # Initial data
│
├── scripts/            # Utility scripts
│
├── tests/              # Integration and end-to-end tests
│
├── go.mod              # Dependency management
├── go.sum              # Dependency checksums
└── README.md           # Project documentation