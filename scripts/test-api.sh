#!/usr/bin/env bash

# MiniBB API Test Script

BASE_URL="http://localhost:8080/api"

echo "Testing MiniBB API..."

# Test 1: Get boards
echo "1. Testing GET /api/boards"
response=$(curl -s "${BASE_URL}/boards")
echo "Response: $response"
echo ""

# Test 2: Get a specific board
echo "2. Testing GET /api/board/general"
response=$(curl -s "${BASE_URL}/board/general")
echo "Response: $response"
echo ""

# Test 3: Get topics for board ID 1
echo "3. Testing GET /api/topics/board/1"
response=$(curl -s "${BASE_URL}/topics/board/1")
echo "Response: $response"
echo ""

# Test 4: Create a test topic
echo "4. Testing POST /api/topic"
response=$(curl -s -X POST "${BASE_URL}/topic" \
  -H "Content-Type: application/json" \
  -d '{"board_id": 1, "title": "Test Topic", "author": "TestUser##test123"}')
echo "Response: $response"

# Extract topic ID from response for next test
topic_id=$(echo $response | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo "Created topic ID: $topic_id"
echo ""

if [ -n "$topic_id" ]; then
    # Test 5: Create a test post
    echo "5. Testing POST /api/post"
    response=$(curl -s -X POST "${BASE_URL}/post" \
      -H "Content-Type: application/json" \
      -d "{\"topic_id\": $topic_id, \"author\": \"TestUser##test123\", \"content\": \"This is a **test post** with *markdown*!\"}")
    echo "Response: $response"
    echo ""

    # Test 6: Get posts for the topic
    echo "6. Testing GET /api/posts/topic/$topic_id"
    response=$(curl -s "${BASE_URL}/posts/topic/${topic_id}")
    echo "Response: $response"
fi

echo ""
echo "API testing complete!"
