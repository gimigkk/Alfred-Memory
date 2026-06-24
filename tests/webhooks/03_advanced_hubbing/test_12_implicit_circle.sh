#!/bin/bash

echo "==========================================="
echo " TEST 12: IMPLICIT CIRCLE INFERENCE"
echo "==========================================="

echo "Sending Message 1..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780890000_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "nadine_ieee26@c.us",
      "fromMe": false,
      "body": "Halo panitia inti, besok kita rapat finalisasi venue ya jam 7 malam di sekre.",
      "timestamp": 1780890000
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780890060_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rendi_ramadana@c.us",
      "fromMe": false,
      "body": "Siap nad, dari divisi logistik barang-barang udah di list semua, besok gw bawa catatannya.",
      "timestamp": 1780890060
    }
  }' > /dev/null

echo "Sending Message 3..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780890120_3",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "Aman, rundown dari divisi acara juga udah fix. Tinggal nunggu acc dari BPH.",
      "timestamp": 1780890120
    }
  }' > /dev/null

echo "Sending Message 4..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780890180_4",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "Gue gabisa ikut rapat besok, tolong wakilin ya anak-anak humas.",
      "timestamp": 1780890180
    }
  }' > /dev/null

echo "Final block sent! Waiting for debounce queue to commit..."
