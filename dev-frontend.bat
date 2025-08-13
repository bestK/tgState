@echo off
echo Starting frontend development server...
cd frontend
start cmd /k "npm run dev"
cd ..
echo Starting Go backend...
go run main.go