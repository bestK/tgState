@echo off
echo Building frontend...
cd frontend
call npm install
echo Running type check...
call npm run type-check
if %errorlevel% neq 0 (
    echo Type check failed, but continuing with build...
)
call npm run build
cd ..
echo Frontend build complete!
echo Building Go application...
go build -o tgstate.exe
echo Build complete!