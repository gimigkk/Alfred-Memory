#!/bin/bash

echo "==========================================="
echo " BLOCK 1"
echo "==========================================="
echo "Sending Message 1..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780815600_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn@c.us",
      "fromMe": false,
      "body": "Guys, besok sore kita meeting buat bahas design ya jam 3 sore. Jangan lupa.",
      "timestamp": 1780815600
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780815720_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn@c.us",
      "fromMe": false,
      "body": "Eh Gilang Muhamad W, tolong siapin draft presentasi UI nya dong sebelum meeting.",
      "timestamp": 1780815720
    }
  }' > /dev/null

echo "Sending Message 3..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780815900_3",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Sip, aman jes.",
      "timestamp": 1780815900
    }
  }' > /dev/null

echo "Final block sent! Waiting for debounce queue to commit..."
