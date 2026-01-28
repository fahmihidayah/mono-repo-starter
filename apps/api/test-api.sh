#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080"

echo -e "${BLUE}=== Testing Blog API ===${NC}\n"

# Test 1: Health Check
echo -e "${BLUE}1. Testing Health Check...${NC}"
HEALTH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/health)
if [ "$HEALTH_RESPONSE" -eq 200 ]; then
    echo -e "${GREEN}✓ Health check passed${NC}\n"
else
    echo -e "${RED}✗ Health check failed (HTTP $HEALTH_RESPONSE)${NC}\n"
    exit 1
fi

# Test 2: Register User
echo -e "${BLUE}2. Testing User Registration...${NC}"
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test'$(date +%s)'@example.com",
    "password": "password123"
  }')

echo "Response: $REGISTER_RESPONSE"
if echo "$REGISTER_RESPONSE" | grep -q "User registered successfully"; then
    echo -e "${GREEN}✓ Registration successful${NC}\n"
    TEST_EMAIL=$(echo "$REGISTER_RESPONSE" | grep -o '"email":"[^"]*"' | cut -d'"' -f4)
else
    echo -e "${RED}✗ Registration failed${NC}\n"
fi

# Test 3: Login with created user
if [ ! -z "$TEST_EMAIL" ]; then
    echo -e "${BLUE}3. Testing User Login...${NC}"
    LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/api/users/login \
      -H "Content-Type: application/json" \
      -d "{
        \"email\": \"$TEST_EMAIL\",
        \"password\": \"password123\"
      }")

    echo "Response: $LOGIN_RESPONSE"
    if echo "$LOGIN_RESPONSE" | grep -q "Login successful"; then
        echo -e "${GREEN}✓ Login successful${NC}\n"
    else
        echo -e "${RED}✗ Login failed${NC}\n"
    fi
fi

# Test 4: Login with wrong password
echo -e "${BLUE}4. Testing Login with wrong password...${NC}"
WRONG_LOGIN=$(curl -s -X POST $BASE_URL/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "wrongpassword"
  }')

if echo "$WRONG_LOGIN" | grep -q "invalid email or password"; then
    echo -e "${GREEN}✓ Correctly rejected wrong password${NC}\n"
else
    echo -e "${RED}✗ Wrong password test failed${NC}\n"
fi

echo -e "${BLUE}=== All Tests Completed ===${NC}"
