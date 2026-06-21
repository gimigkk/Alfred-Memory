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
      "id": "msg_1781438340_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "Ingin punya pengalaman membuat *\"AI Agent\"*...",
      "timestamp": 1781438340
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781438340_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "ew",
      "timestamp": 1781438340
    }
  }' > /dev/null

echo "Sending Message 3..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781438400_3",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "holy slop @M3-117_Gilang Muhamad W",
      "timestamp": 1781438400
    }
  }' > /dev/null

echo "Sending Message 4..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781438460_4",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Aowkowwkoakwo",
      "timestamp": 1781438460
    }
  }' > /dev/null

echo "Sending Message 5..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439360_5",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Minta medpart apa gmn deh",
      "timestamp": 1781439360
    }
  }' > /dev/null

echo "Sending Message 6..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439420_6",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "gatau woi",
      "timestamp": 1781439420
    }
  }' > /dev/null

echo "Sending Message 7..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439420_7",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "hrus dibales apa gue kacangin aj",
      "timestamp": 1781439420
    }
  }' > /dev/null

echo "Sending Message 8..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439420_8",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "@Nadine IEEE26 betternya gmn din",
      "timestamp": 1781439420
    }
  }' > /dev/null

echo "Sending Message 9..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439480_9",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "tpi itu bukan ajak kerja sama deh?",
      "timestamp": 1781439480
    }
  }' > /dev/null

echo "Sending Message 10..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439480_10",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "dia mintol isi kuesioner buat skripsi",
      "timestamp": 1781439480
    }
  }' > /dev/null

echo "Sending Message 11..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439540_11",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Iya sihh",
      "timestamp": 1781439540
    }
  }' > /dev/null

echo "Sending Message 12..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439540_12",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Dapet apa kita jir klo ngeshare",
      "timestamp": 1781439540
    }
  }' > /dev/null

echo "Sending Message 13..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439540_13",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Ew",
      "timestamp": 1781439540
    }
  }' > /dev/null

echo "Sending Message 14..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439540_14",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "nadine_ieee26@c.us",
      "fromMe": false,
      "body": "ini pribadi sih...",
      "timestamp": 1781439540
    }
  }' > /dev/null

echo "Sending Message 15..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439540_15",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "nadine_ieee26@c.us",
      "fromMe": false,
      "body": "diemin aja ga sih",
      "timestamp": 1781439540
    }
  }' > /dev/null

echo "Sending Message 16..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439540_16",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "mauny sih gitu",
      "timestamp": 1781439540
    }
  }' > /dev/null

echo "Sending Message 17..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439540_17",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "tpi kek tba\" bngt sebut ieeeny mksd gue",
      "timestamp": 1781439540
    }
  }' > /dev/null

echo "Sending Message 18..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439540_18",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "kirain ad yg kenal",
      "timestamp": 1781439540
    }
  }' > /dev/null

echo "Sending Message 19..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439600_19",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "nadine_ieee26@c.us",
      "fromMe": false,
      "body": "iyahh gapapa sih lagian lu jg ga kenal dan itu buat kepentingan dia pribadi juga...",
      "timestamp": 1781439600
    }
  }' > /dev/null

echo "Sending Message 20..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439660_20",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Aowkowwkoakwo iya jir",
      "timestamp": 1781439660
    }
  }' > /dev/null

echo "Sending Message 21..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439660_21",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Ywdah didiemin aja",
      "timestamp": 1781439660
    }
  }' > /dev/null

echo "Sending Message 22..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781439780_22",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "owkheyy",
      "timestamp": 1781439780
    }
  }' > /dev/null

echo "Final block sent! Waiting for debounce queue to commit..."
