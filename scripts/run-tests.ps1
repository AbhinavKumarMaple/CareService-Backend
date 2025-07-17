# PowerShell script to run tests and display coverage

# Set the working directory to the project root
$projectRoot = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $projectRoot

# Create output directory for coverage reports if it doesn't exist
$coverageDir = Join-Path $projectRoot "coverage"
if (-not (Test-Path $coverageDir)) {
    New-Item -ItemType Directory -Path $coverageDir | Out-Null
}

# Define coverage output files
$coverageOut = Join-Path $coverageDir "coverage.out"
$coverageHtml = Join-Path $coverageDir "coverage.html"

# Display Go version
Write-Host "Go version:" -ForegroundColor Cyan
go version
Write-Host ""

# Run tests with coverage
Write-Host "Running tests with coverage..." -ForegroundColor Cyan
go test -v -coverprofile=$coverageOut ./src/...

# Check if tests passed
if ($LASTEXITCODE -eq 0) {
    Write-Host "All tests passed!" -ForegroundColor Green
    
    # Generate HTML coverage report
    Write-Host "Generating HTML coverage report..." -ForegroundColor Cyan
    go tool cover -html=$coverageOut -o $coverageHtml
    
    # Display coverage statistics
    Write-Host "Coverage statistics:" -ForegroundColor Cyan
    go tool cover -func=$coverageOut
    
    # Open the coverage report in the default browser
    Write-Host "Opening coverage report in browser..." -ForegroundColor Cyan
    Start-Process $coverageHtml
} else {
    Write-Host "Tests failed!" -ForegroundColor Red
}