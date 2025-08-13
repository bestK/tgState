#!/bin/bash
echo "Building frontend..."
cd frontend
npm install
echo "Running type check..."
npm run type-check || echo "Type check failed, but continuing with build..."
npm run build
cd ..
echo "Frontend build complete!"
echo "Building Go application..."
go build -o tgstate
echo "Build complete!"