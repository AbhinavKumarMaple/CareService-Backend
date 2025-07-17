#!/bin/bash
# Script to run tests and display coverage

# Set the working directory to the project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# Create output directory for coverage reports if it doesn't exist
COVERAGE_DIR="$PROJECT_ROOT/coverage"
mkdir -p "$COVERAGE_DIR"

# Define coverage output files
COVERAGE_OUT="$COVERAGE_DIR/coverage.out"
COVERAGE_HTML="$COVERAGE_DIR/coverage.html"

# Display Go version
echo -e "\033[36mGo version:\033[0m"
go version
echo ""

# Run tests with coverage
echo -e "\033[36mRunning tests with coverage...\033[0m"
go test -v -coverprofile="$COVERAGE_OUT" ./src/...

# Check if tests passed
if [ $? -eq 0 ]; then
    echo -e "\033[32mAll tests passed!\033[0m"
    
    # Generate HTML coverage report
    echo -e "\033[36mGenerating HTML coverage report...\033[0m"
    go tool cover -html="$COVERAGE_OUT" -o "$COVERAGE_HTML"
    
    # Display coverage statistics
    echo -e "\033[36mCoverage statistics:\033[0m"
    go tool cover -func="$COVERAGE_OUT"
    
    # Try to open the coverage report in the default browser
    echo -e "\033[36mOpening coverage report in browser...\033[0m"
    if command -v xdg-open &> /dev/null; then
        xdg-open "$COVERAGE_HTML"
    elif command -v open &> /dev/null; then
        open "$COVERAGE_HTML"
    else
        echo -e "\033[33mCould not open browser automatically. Please open $COVERAGE_HTML manually.\033[0m"
    fi
else
    echo -e "\033[31mTests failed!\033[0m"
fi