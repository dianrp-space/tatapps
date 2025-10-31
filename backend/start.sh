#!/bin/bash

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}🚀 Starting Backend Server...${NC}"

# Check if already running
if lsof -ti:8080 > /dev/null 2>&1; then
    echo -e "${RED}❌ Backend already running on port 8080${NC}"
    echo -e "${YELLOW}Use ./restart.sh to restart${NC}"
    exit 1
fi

# Start backend
cd /var/www/html/tatapps/backend
nohup go run cmd/api/main.go > /tmp/backend.log 2>&1 &

# Wait for server to start
sleep 3

# Check if server is running
if lsof -ti:8080 > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Backend server started!${NC}"
    echo -e "${GREEN}📝 Log: tail -f /tmp/backend.log${NC}"
    echo -e "${GREEN}🌐 Health: http://localhost:8080/health${NC}"
else
    echo -e "${RED}❌ Failed to start. Check: cat /tmp/backend.log${NC}"
    exit 1
fi
