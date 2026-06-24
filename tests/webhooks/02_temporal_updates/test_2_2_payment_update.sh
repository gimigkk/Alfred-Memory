#!/bin/bash

echo "Sending Message 13..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969910_13",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "nadine_ieee26@c.us",
      "fromMe": false,
      "body": "guyssss ini yg belom bayar siapaa?? baru aqila, rapid, dan rapip, btw @Rendi Ramadana IEEE²⁵ lu 27rb ke gopay gua yaahh",
      "timestamp": 1781969910
    }
  }' > /dev/null

echo "Sending Message 14..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886236_14",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rendi_ramadana@c.us",
      "fromMe": false,
      "body": "Ok nad, gw tf skrg ya",
      "timestamp": 1781969915
    }
  }' > /dev/null

echo "Sending Message 15..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886236_15",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "nadine_ieee26@c.us",
      "fromMe": false,
      "body": "Sip, oh ya sekalian ini kan reimburse buat rapat bph gacoan kemarin. Talangannya digabung sama punya lu kan apta?",
      "timestamp": 1781969920
    }
  }' > /dev/null

echo "Sending Message 16..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886236_16",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "Yoi, tf ke gopay gw aja",
      "timestamp": 1781969925
    }
  }' > /dev/null

echo "Sending Message 17..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886236_17",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Lah trus si jeslyn ngapain tadi ngirim rekening, kan dia ga nalangin apa2",
      "timestamp": 1781969930
    }
  }' > /dev/null

echo "Sending Message 18..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886236_18",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "wkwkwk becanda doang itu meme",
      "timestamp": 1781969935
    }
  }' > /dev/null

echo "Sending Message 19..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1781969886236_19",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "gajelas wkwk",
      "timestamp": 1781969940
    }
  }' > /dev/null

echo "Webhooks sent! Waiting for debounce queue to commit..."
