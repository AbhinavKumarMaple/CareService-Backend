# Testing Guide

This document provides instructions on how to run tests for the microservices-go project.

## Running Tests

### Using the Provided Scripts

#### Windows (PowerShell)

```powershell
# Run all tests with coverage
.\scripts\run-tests.ps1
```

#### Unix-like Systems (Linux, macOS, WSL)

```bash
# Make the script executable (first time only)
chmod +x scripts/run-tests.sh

# Run all tests with coverage
./scripts/run-tests.sh
```

### Manual Test Commands

If you prefer to run tests manually, you can use the following commands:

```bash
# Run all tests
go test ./src/...

# Run tests with verbose output
go test -v ./src/...

# Run tests with coverage
go test -cover ./src/...

# Run tests with coverage and generate HTML report
go test -coverprofile=coverage.out ./src/...
go tool cover -html=coverage.out -o coverage.html
```

### Running Specific Tests

To run tests for a specific package:

```bash
# Run tests for the schedule usecase
go test -v ./src/application/usecases/schedule

# Run tests for the schedule controller
go test -v ./src/infrastructure/rest/controllers/schedule
```

To run a specific test:

```bash
# Run a specific test function
go test -v ./src/application/usecases/schedule -run TestGetSchedules
```

## Test Structure

The tests are organized according to the project's clean architecture:

1. **Domain Layer**: Tests for domain models and business rules
2. **Application Layer**: Tests for usecases that implement business logic
3. **Infrastructure Layer**: Tests for controllers, repositories, and other infrastructure components

### Mock Objects

Mock objects are used to isolate the unit being tested from its dependencies. For example:

- `mockScheduleRepository`: Mocks the `IScheduleRepository` interface
- `mockUserRepository`: Mocks the `IUserRepository` interface
- `mockScheduleUseCase`: Mocks the `IScheduleUseCase` interface

### Test Helpers

Several helper functions are provided to simplify test setup:

- `setupLogger`: Creates a logger instance for testing
- `setupTestScheduleUseCase`: Creates a Schedule usecase with mock repositories
- `setupTestController`: Creates a Schedule controller with mock usecase
- `createTestSchedule`: Creates a test schedule with predefined values
- `createTestUser`: Creates a test user with predefined values

## Test Coverage

The test coverage report shows which parts of the code are covered by tests. After running the tests with coverage, you can view the HTML report to see detailed coverage information.

The coverage report highlights:

- Lines that are covered by tests (in green)
- Lines that are not covered by tests (in red)
- The percentage of code coverage for each file and package

## Writing New Tests

When writing new tests, follow these guidelines:

1. **Use Table-Driven Tests**: For testing multiple scenarios with the same logic
2. **Test Both Success and Error Cases**: Ensure that error handling is properly tested
3. **Mock Dependencies**: Use mock implementations to isolate the unit being tested
4. **Use Descriptive Test Names**: Make it clear what is being tested
5. **Follow the AAA Pattern**: Arrange, Act, Assert
   - Arrange: Set up the test data and conditions
   - Act: Perform the operation being tested
   - Assert: Verify the results

## Continuous Integration

Tests are automatically run as part of the CI/CD pipeline. Make sure all tests pass before submitting a pull request.
