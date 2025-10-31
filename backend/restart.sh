#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔄 Restarting Backend Server...${NC}"

# Stop existing backend process
echo -e "${YELLOW}Stopping existing backend...${NC}"
pkill -f "go run cmd/api/main.go" 2>/dev/null
lsof -ti:8080 | xargs kill -9 2>/dev/null
sleep 2

# Check if port is free
if lsof -ti:8080 > /dev/null 2>&1; then
    echo -e "${RED}❌ Port 8080 still in use. Please check manually.${NC}"
    exit 1
fi

# Start backend
echo -e "${YELLOW}Starting backend server...${NC}"
cd /var/www/html/tatapps/backend
nohup go run cmd/api/main.go > /tmp/backend.log 2>&1 &

# Wait for server to start
sleep 3

# Check if server is running
if lsof -ti:8080 > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Backend server started successfully!${NC}"
    echo -e "${GREEN}📝 Log: tail -f /tmp/backend.log${NC}"
    echo -e "${GREEN}🌐 API: http://localhost:8080/health${NC}"
    
    # Show last few lines of log
    echo -e "\n${YELLOW}Last 10 lines of log:${NC}"
    tail -10 /tmp/backend.log
else
    echo -e "${RED}❌ Failed to start backend server.${NC}"
    echo -e "${RED}Check log: cat /tmp/backend.log${NC}"
    exit 1
fi
