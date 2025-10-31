#!/bin/bash

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}🛑 Stopping Backend Server...${NC}"

# Stop backend process
pkill -f "go run cmd/api/main.go" 2>/dev/null
lsof -ti:8080 | xargs kill -9 2>/dev/null

sleep 2

# Check if stopped
if lsof -ti:8080 > /dev/null 2>&1; then
    echo -e "${RED}❌ Failed to stop backend (port still in use)${NC}"
    exit 1
else
    echo -e "${GREEN}✅ Backend server stopped${NC}"
fi
