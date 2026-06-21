#!/bin/bash

echo "==========================================="
echo " BLOCK 1: Creating the Planned Task"
echo "==========================================="
echo "Sending Message 1..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "Guys, gw butuh desain logo buat SoTQ ntar ya",
      "timestamp": 1718960400
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rendi_ramadana_ieee25@c.us",
      "fromMe": false,
      "body": "Oke aman, gw yg handle. Ntar gw kerjain abis makul.",
      "timestamp": 1718960460
    }
  }' > /dev/null

echo "Block 1 sent! Waiting 30 seconds for the agent to commit it..."
sleep 30

echo "==========================================="
echo " BLOCK 2: Updating the Task to Completed"
echo "==========================================="
echo "Sending Message 3..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_2_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rendi_ramadana_ieee25@c.us",
      "fromMe": false,
      "body": "@Apta IEEE25 logo sotq udah beres ya, udah gw taro di drive",
      "timestamp": 1718974200
    }
  }' > /dev/null

echo "Sending Message 4..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_2_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "Mantap tengkyu ren",
      "timestamp": 1718974320
    }
  }' > /dev/null

echo "Block 2 sent! Waiting for debounce queue to commit..."
