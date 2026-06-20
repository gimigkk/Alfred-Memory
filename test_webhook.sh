#!/bin/bash

echo "Sending Message 1..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1",
      "from": "bahlil@c.us",
      "to": "me@c.us",
      "participant": "",
      "fromMe": false,
      "body": "Bro",
      "timestamp": 1690000000
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_2",
      "from": "bahlil@c.us",
      "to": "me@c.us",
      "participant": "",
      "fromMe": false,
      "body": "Bisa tolong siapkan slide design DPP?",
      "timestamp": 1690000001
    }
  }' > /dev/null

echo "Sending Message 3..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_3",
      "from": "bahlil@c.us",
      "to": "me@c.us",
      "participant": "",
      "fromMe": false,
      "body": "Buat event hari Jumat besok ya, penting nih.",
      "timestamp": 1690000002
    }
  }' > /dev/null

echo "Webhooks sent! Waiting for debounce queue to commit..."
