#!/bin/bash
# UnoReverse Docker Cleanup Script
# Use this to manually clean up old containers, images, and volumes

set -e

echo "╔════════════════════════════════════════════╗"
echo "║  UnoReverse Docker System Cleanup          ║"
echo "╚════════════════════════════════════════════╝"
echo ""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Stop all UnoReverse containers
echo -e "${YELLOW}[1/7]${NC} Stopping UnoReverse containers..."
docker-compose down 2>/dev/null || true
sleep 2

# Step 2: Remove old containers by name
echo -e "${YELLOW}[2/7]${NC} Removing old containers..."
docker rm -f postgres-uno backend-uno frontend-uno prometheus-uno grafana-uno 2>/dev/null || true

# Step 3: Remove dangling volumes
echo -e "${YELLOW}[3/7]${NC} Cleaning dangling volumes..."
docker volume prune -f 2>/dev/null || true

# Step 4: Remove dangling images
echo -e "${YELLOW}[4/7]${NC} Removing unused images..."
docker image prune -f 2>/dev/null || true

# Step 5: Show disk space reclaimed
echo -e "${YELLOW}[5/7]${NC} Disk space status..."
docker system df

# Step 6: Verify containers are gone
echo -e "${YELLOW}[6/7]${NC} Verifying cleanup..."
REMAINING=$(docker ps -a --filter "name=postgres-uno\|backend-uno\|frontend-uno\|prometheus-uno\|grafana-uno" -q | wc -l)
if [ $REMAINING -eq 0 ]; then
    echo -e "${GREEN}✓ All old containers removed${NC}"
else
    echo -e "${RED}✗ Warning: $REMAINING containers still exist${NC}"
fi

# Step 7: Ready for fresh deployment
echo -e "${YELLOW}[7/7]${NC} Ready for fresh deployment..."
echo ""
echo -e "${GREEN}╔════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  System cleaned successfully!              ║${NC}"
echo -e "${GREEN}║  You can now run: docker-compose up -d     ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════╝${NC}"
