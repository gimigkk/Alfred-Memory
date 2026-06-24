#!/bin/bash

echo "Sending Message 1..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886323_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "S1 Ilmu Komputer mengundang Anda untuk bergabung ke rapat Zoom yang terjadwal.",
      "timestamp": 1781969887
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886323_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "disini jam 07.30 pm, 19.30 malem gitu?",
      "timestamp": 1781969888
    }
  }' > /dev/null

echo "Sending Message 3..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886323_3",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Lah bejir",
      "timestamp": 1781969889
    }
  }' > /dev/null

echo "Sending Message 4..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886323_4",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Kynya sih",
      "timestamp": 1781969890
    }
  }' > /dev/null

echo "Sending Message 5..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886323_5",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "sapa belom punya kelompok metkuan",
      "timestamp": 1781969891
    }
  }' > /dev/null

echo "Sending Message 6..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886323_6",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Gwej",
      "timestamp": 1781969892
    }
  }' > /dev/null

echo "Sending Message 7..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886323_7",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Gwej",
      "timestamp": 1781969893
    }
  }' > /dev/null

echo "Sending Message 8..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886323_8",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "sapa belom punya kelompok metkuan",
      "timestamp": 1781969894
    }
  }' > /dev/null

echo "Webhooks sent! Waiting for debounce queue to commit..."
