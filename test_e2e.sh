#!/bin/bash

# E2E test script for rgt file filtering
# Tests that golang and python file type filters work correctly

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================="
echo "Starting E2E tests for rgt"
echo "========================================="

# Clean up function
cleanup() {
    echo -e "\n${YELLOW}Cleaning up...${NC}"
    if [ ! -z "$RGT_PID" ]; then
        kill $RGT_PID 2>/dev/null || true
    fi
    rm -f test_e2e_temp.go test_e2e_temp.py
}

trap cleanup EXIT

# Build the binary
echo -e "\n${YELLOW}Step 1: Building binary...${NC}"
make build

# Test 1: Golang file filtering
echo -e "\n${YELLOW}Step 2: Testing Golang file filtering...${NC}"
echo "  - Starting rgt with --test-type golang"
echo "  - Expecting: .go files trigger tests, .py files don't"

# Create temporary files for testing
echo "package main" > test_e2e_temp.go
echo "def test(): pass" > test_e2e_temp.py

# Start rgt in background with golang type
./bin/rgt start --test-type golang > /tmp/rgt_golang_output.txt 2>&1 &
RGT_PID=$!
echo "  - rgt started (PID: $RGT_PID)"

# Wait for rgt to initialize
sleep 2

# Touch a Python file - should NOT trigger tests
echo "  - Touching .py file (should be ignored)..."
echo "# modified" >> test_e2e_temp.py
sleep 2

# Touch a Go file - should trigger tests
echo "  - Touching .go file (should trigger tests)..."
echo "// modified" >> test_e2e_temp.go
sleep 3

# Kill the process
kill $RGT_PID 2>/dev/null || true
wait $RGT_PID 2>/dev/null || true
RGT_PID=""

# Check if .go file was detected (note: path might have ./ prefix)
GO_FILE_DETECTED=$(grep -c "File changed:.*test_e2e_temp.go" /tmp/rgt_golang_output.txt 2>/dev/null || true)
GO_FILE_DETECTED=${GO_FILE_DETECTED:-0}
# Check if .py file was detected (should be 0)
PY_FILE_DETECTED=$(grep -c "File changed:.*test_e2e_temp.py" /tmp/rgt_golang_output.txt 2>/dev/null || true)
PY_FILE_DETECTED=${PY_FILE_DETECTED:-0}

echo "  - .go file detections: $GO_FILE_DETECTED"
echo "  - .py file detections: $PY_FILE_DETECTED"

if [ "$GO_FILE_DETECTED" -gt "0" ] && [ "$PY_FILE_DETECTED" -eq "0" ]; then
    echo -e "${GREEN}✓ Golang filtering test PASSED${NC}"
    echo "  .go files triggered tests, .py files were ignored"
else
    echo -e "${RED}✗ Golang filtering test FAILED${NC}"
    if [ "$GO_FILE_DETECTED" -eq "0" ]; then
        echo "  .go file was NOT detected (expected to be detected)"
    fi
    if [ "$PY_FILE_DETECTED" -gt "0" ]; then
        echo "  .py file was detected (expected to be ignored)"
    fi
    echo "  Output:"
    cat /tmp/rgt_golang_output.txt
    exit 1
fi

# Test 2: Python file filtering
echo -e "\n${YELLOW}Step 3: Testing Python file filtering...${NC}"
echo "  - Starting rgt with --test-type python"
echo "  - Expecting: .py files trigger tests, .go files don't"

# Reset files
echo "package main" > test_e2e_temp.go
echo "def test(): pass" > test_e2e_temp.py

# Start rgt in background with python type
./bin/rgt start --test-type python > /tmp/rgt_python_output.txt 2>&1 &
RGT_PID=$!
echo "  - rgt started (PID: $RGT_PID)"

# Wait for rgt to initialize
sleep 2

# Touch a Go file - should NOT trigger tests
echo "  - Touching .go file (should be ignored)..."
echo "// modified again" >> test_e2e_temp.go
sleep 2

# Touch a Python file - should trigger tests
echo "  - Touching .py file (should trigger tests)..."
echo "# modified again" >> test_e2e_temp.py
sleep 3

# Kill the process
kill $RGT_PID 2>/dev/null || true
wait $RGT_PID 2>/dev/null || true
RGT_PID=""

# Check if .py file was detected (note: path might have ./ prefix)
PY_FILE_DETECTED=$(grep -c "File changed:.*test_e2e_temp.py" /tmp/rgt_python_output.txt 2>/dev/null || true)
PY_FILE_DETECTED=${PY_FILE_DETECTED:-0}
# Check if .go file was detected (should be 0)
GO_FILE_DETECTED=$(grep -c "File changed:.*test_e2e_temp.go" /tmp/rgt_python_output.txt 2>/dev/null || true)
GO_FILE_DETECTED=${GO_FILE_DETECTED:-0}

echo "  - .py file detections: $PY_FILE_DETECTED"
echo "  - .go file detections: $GO_FILE_DETECTED"

if [ "$PY_FILE_DETECTED" -gt "0" ] && [ "$GO_FILE_DETECTED" -eq "0" ]; then
    echo -e "${GREEN}✓ Python filtering test PASSED${NC}"
    echo "  .py files triggered tests, .go files were ignored"
else
    echo -e "${RED}✗ Python filtering test FAILED${NC}"
    if [ "$PY_FILE_DETECTED" -eq "0" ]; then
        echo "  .py file was NOT detected (expected to be detected)"
    fi
    if [ "$GO_FILE_DETECTED" -gt "0" ]; then
        echo "  .go file was detected (expected to be ignored)"
    fi
    echo "  Output:"
    cat /tmp/rgt_python_output.txt
    exit 1
fi

echo -e "\n${GREEN}=========================================${NC}"
echo -e "${GREEN}All E2E tests PASSED!${NC}"
echo -e "${GREEN}=========================================${NC}"
