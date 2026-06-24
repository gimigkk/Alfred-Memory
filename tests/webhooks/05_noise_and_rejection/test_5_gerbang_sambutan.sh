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
      "id": "msg_1780808220_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "udahh harusnya",
      "timestamp": 1780808220
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780808220_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "tapi td anak gua lewat bara katanya gerbangnya ditutup",
      "timestamp": 1780808220
    }
  }' > /dev/null

echo "Sending Message 3..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780808520_3",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "gerbang belakang jg buka katanya",
      "timestamp": 1780808520
    }
  }' > /dev/null

echo "Sending Message 4..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780808880_4",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "eh anjir gebang depan tutup....",
      "timestamp": 1780808880
    }
  }' > /dev/null

echo "Sending Message 5..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780808880_5",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "yang buka bni",
      "timestamp": 1780808880
    }
  }' > /dev/null

echo "Sending Message 6..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780808880_6",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "bni tu dimana",
      "timestamp": 1780808880
    }
  }' > /dev/null

echo "Sending Message 7..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780808940_7",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "eh buka deng depann",
      "timestamp": 1780808940
    }
  }' > /dev/null

echo "Sending Message 8..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780808940_8",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "duh sorry miskom",
      "timestamp": 1780808940
    }
  }' > /dev/null

echo "Sending Message 9..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780810620_9",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "@Rafid Harsyah dmn jir lu kan ada sambutan",
      "timestamp": 1780810620
    }
  }' > /dev/null

echo "Sending Message 10..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780810980_10",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "eh rapit mana anjir",
      "timestamp": 1780810980
    }
  }' > /dev/null

echo "Sending Message 11..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780810980_11",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "woi @Rafid Harsyah akwowkowkwok",
      "timestamp": 1780810980
    }
  }' > /dev/null

echo "Sending Message 12..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780810980_12",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "tolong back up pls vp",
      "timestamp": 1780810980
    }
  }' > /dev/null

echo "Sending Message 13..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780810980_13",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "call aja",
      "timestamp": 1780810980
    }
  }' > /dev/null

echo "Sending Message 14..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780810980_14",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "kalo dia blm dtg",
      "timestamp": 1780810980
    }
  }' > /dev/null

echo "Sending Message 15..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780810980_15",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "emng idiot kesel gw",
      "timestamp": 1780810980
    }
  }' > /dev/null

echo "Sending Message 16..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780810980_16",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "tolong maju jon",
      "timestamp": 1780810980
    }
  }' > /dev/null

echo "Sending Message 17..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780810980_17",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "sambutan bntr aja",
      "timestamp": 1780810980
    }
  }' > /dev/null

echo "Sending Message 18..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780810980_18",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "panggil aja",
      "timestamp": 1780810980
    }
  }' > /dev/null

echo "Sending Message 19..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780811040_19",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "tadi udah otw dah orangnya",
      "timestamp": 1780811040
    }
  }' > /dev/null

echo "Final block sent! Waiting for debounce queue to commit..."
