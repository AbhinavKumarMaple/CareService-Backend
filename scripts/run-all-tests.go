package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	startTime := time.Now()
	fmt.Println("=== Running All Tests ===")
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Println("========================")

	// Get the project root directory
	projectRoot, err := getProjectRoot()
	if err != nil {
		fmt.Printf("Error finding project root: %v\n", err)
		os.Exit(1)
	}

	// Change to project root directory
	err = os.Chdir(projectRoot)
	if err != nil {
		fmt.Printf("Error changing to project root directory: %v\n", err)
		os.Exit(1)
	}

	// Run tests with verbose output
	fmt.Println("\n=== Running tests with verbose output ===")
	cmd := exec.Command("go", "test", "-v", "./src/...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	// Check if tests passed
	if err != nil {
		fmt.Printf("\n❌ Tests failed: %v\n", err)
		os.Exit(1)
	}

	// Print summary
	duration := time.Since(startTime)
	fmt.Printf("\n✅ All tests completed successfully in %.2f seconds\n", duration.Seconds())
}

// getProjectRoot finds the project root directory
func getProjectRoot() (string, error) {
	// Get the directory of the current script
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get the current file path")
	}

	// The script is in the scripts directory, so go up one level
	return filepath.Dir(filepath.Dir(filename)), nil
}