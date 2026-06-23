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
      "id": "msg_1780836780_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "@M Naufal IEEE²⁵ @Rafid Harsyah",
      "timestamp": 1780836780
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836840_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "tf mn",
      "timestamp": 1780836840
    }
  }' > /dev/null

echo "Sending Message 3..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836840_3",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "es gobak sodor apaan dh",
      "timestamp": 1780836840
    }
  }' > /dev/null

echo "Sending Message 4..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836840_4",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "0954959895 - BCA a.n Apta Adi Nur Fiansah ",
      "timestamp": 1780836840
    }
  }' > /dev/null

echo "Sending Message 5..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836900_5",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "makanya ikut",
      "timestamp": 1780836900
    }
  }' > /dev/null

echo "Sending Message 6..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836900_6",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "901158617350 - Jeslyn Angelica Widjaja - Seabank",
      "timestamp": 1780836900
    }
  }' > /dev/null

echo "Sending Message 7..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836900_7",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "pikm",
      "timestamp": 1780836900
    }
  }' > /dev/null

echo "Sending Message 8..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836900_8",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Ongkeh",
      "timestamp": 1780836900
    }
  }' > /dev/null

echo "Sending Message 9..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780837020_9",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "2",
      "timestamp": 1780837020
    }
  }' > /dev/null

echo "Sending Message 10..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780837080_10",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "AI ini",
      "timestamp": 1780837080
    }
  }' > /dev/null

echo "Sending Message 11..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780837320_11",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "asli ini ini",
      "timestamp": 1780837320
    }
  }' > /dev/null

echo "Sending Message 12..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780837380_12",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "lu tf gopay ye?",
      "timestamp": 1780837380
    }
  }' > /dev/null

echo "Final block sent! Waiting for debounce queue to commit..."
