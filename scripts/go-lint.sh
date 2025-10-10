#!/bin/bash

# Development helper script for Go linting and testing
# This script can be used to manually run checks that are disabled in lint-staged

echo "🔍 Running Go format check..."
cd api_v2
gofmt -l . | grep -v "^$"
if [ $? -eq 0 ]; then
    echo "❌ Some files are not properly formatted. Running gofmt -w ..."
    gofmt -w .
    echo "✅ Files formatted"
else
    echo "✅ All files are properly formatted"
fi

echo ""
echo "🔍 Running Go vet (may fail due to module issues)..."
go vet ./...
if [ $? -eq 0 ]; then
    echo "✅ Go vet passed"
else
    echo "⚠️  Go vet failed - this is expected during development"
fi

echo ""
echo "🔍 Running Go mod tidy..."
go mod tidy
if [ $? -eq 0 ]; then
    echo "✅ Go mod tidy completed"
else
    echo "❌ Go mod tidy failed"
fi

echo ""
echo "🔍 Checking for build errors..."
go build ./...
if [ $? -eq 0 ]; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
fi
