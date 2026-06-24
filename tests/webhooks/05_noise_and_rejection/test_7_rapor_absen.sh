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
      "id": "msg_1780817700_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "@~Syazana Aqila mau liat rapor penilaian buat eval divisi gua",
      "timestamp": 1780817700
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780817820_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "https://docs.google.com/spreadsheets/d/1wvbqIMv_qiH8NCDQN_FDG6T9Dk72ms43/edit?usp=sharing",
      "timestamp": 1780817820
    }
  }' > /dev/null

echo "Sending Message 3..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780817820_3",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "baru ada nilainya doang, blm kita kasih comment",
      "timestamp": 1780817820
    }
  }' > /dev/null

echo "Sending Message 4..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780817820_4",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "ini rapor penilaian tiap proker",
      "timestamp": 1780817820
    }
  }' > /dev/null

echo "Sending Message 5..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780818000_5",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "okeoke thankk uu yak",
      "timestamp": 1780818000
    }
  }' > /dev/null

echo "Time gap of > 30 minutes detected. Waiting 2 minutes for ingestion agent to process block..."
sleep 120

echo "==========================================="
echo " BLOCK 2"
echo "==========================================="
echo "Sending Message 6..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780819920_6",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "@Rendi Ramadana IEEE²⁵ udh ttd absen blm",
      "timestamp": 1780819920
    }
  }' > /dev/null

echo "Sending Message 7..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780821180_7",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rendi_ramadana_ieee__@c.us",
      "fromMe": false,
      "body": "belum njirr lupa🤦🏻♂️",
      "timestamp": 1780821180
    }
  }' > /dev/null

echo "Sending Message 8..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780821180_8",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rendi_ramadana_ieee__@c.us",
      "fromMe": false,
      "body": "tipsen",
      "timestamp": 1780821180
    }
  }' > /dev/null

echo "Final block sent! Waiting for debounce queue to commit..."
