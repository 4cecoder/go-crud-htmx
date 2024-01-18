@echo off
go mod tidy >nul 2>&1
if errorlevel 1 (
    echo Error: Go is not installed or not reachable.
    exit /b
)
echo Go is installed and reachable.

echo Running the api...
go run main.go