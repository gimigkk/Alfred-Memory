#!/bin/bash

# Ensure we run from project root
cd "$(dirname "$0")/.."
# Kill any existing server
kill -9 $(cat server.pid 2>/dev/null) 2>/dev/null
pkill -f "bin/alfred" 2>/dev/null

# Clean DB and logs
rm -rf .lbug/
rm -f server.log

# Start new server
./bin/alfred > server.log 2>&1 &
echo $! > server.pid

echo "Waiting for server to start..."
sleep 2

TESTS=(
  "03_advanced_hubbing/test_12_implicit_circle.sh"
  "01_core_extraction/test_1_sambutan.sh"
)

for test_file in "${TESTS[@]}"; do
  echo "=== $test_file START ===" >> server.log
  echo "Running $test_file..."
  ./tests/webhooks/$test_file
  echo "Waiting 180s for agent to finish processing..."
  sleep 180
done

echo "=== ALL TESTS DONE ===" >> server.log
echo "All tests finished!"
