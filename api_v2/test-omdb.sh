#!/bin/bash

# OMDb API Test Script
# Run this after starting the server with ./start-server.sh

echo "🧪 Testing OMDb Integration"
echo "================================"
echo ""

# Test 1: Connection Test
echo "1️⃣  Testing OMDb Connection..."
curl -s "http://localhost:8080/api/v1/omdb/test" | jq '.'
echo ""
echo "---"
echo ""

# Test 2: Get Movie by IMDb ID (The Matrix)
echo "2️⃣  Getting The Matrix by IMDb ID (tt0133093)..."
curl -s "http://localhost:8080/api/v1/omdb/tt0133093" | jq '.'
echo ""
echo "---"
echo ""

# Test 3: Get Movie by Title
echo "3️⃣  Getting Inception by Title..."
curl -s "http://localhost:8080/api/v1/omdb/title?title=Inception&year=2010" | jq '.'
echo ""
echo "---"
echo ""

# Test 4: Search Movies
echo "4️⃣  Searching for 'Batman'..."
curl -s "http://localhost:8080/api/v1/omdb/search?q=Batman&page=1" | jq '.'
echo ""
echo "---"
echo ""

# Test 5: Search Movies by Type
echo "5️⃣  Searching for 'Star Wars' movies only..."
curl -s "http://localhost:8080/api/v1/omdb/search-by-type?q=Star%20Wars&type=movie&page=1" | jq '.'
echo ""
echo "---"
echo ""

# Test 6: Get Interstellar
echo "6️⃣  Getting Interstellar (tt0816692)..."
curl -s "http://localhost:8080/api/v1/omdb/tt0816692" | jq '{title: .title, year: .year, imdb_rating: .imdb_rating, plot: .plot, provider: .provider}'
echo ""
echo "---"
echo ""

echo "✅ All tests completed!"
