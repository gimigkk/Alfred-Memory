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
      "id": "msg_1780887600_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Guys, buat urusan desain UI/UX dan konten sosmed, kita kumpulin di satu wadah aja ya biar ga kecampur sama BPH inti.",
      "timestamp": 1780887600
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780887660_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Gua bikin grup baru namanya 'Creative BPH'. Isinya anak anak divisi creative doang ya, biar gampang koordinasinya.",
      "timestamp": 1780887660
    }
  }' > /dev/null

echo "Sending Message 3..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780887720_3",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "Oke sip bang rafid. Berarti gua masuk situ ya.",
      "timestamp": 1780887720
    }
  }' > /dev/null

echo "Sending Message 4..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780887780_4",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "Oiya btw ntar malem kan kumpul bahas maskot. Apta lu jangan telat lagi ya elah.",
      "timestamp": 1780887780
    }
  }' > /dev/null

echo "Sending Message 5..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780887840_5",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "Apta mah emang udah pasti telat mulu tiap kali disuruh kumpul. Kemarin telat sejam, minggu lalu telat 2 jam. Jangan ngarep dia datang on time.",
      "timestamp": 1780887840
    }
  }' > /dev/null

echo "Sending Message 6..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780887900_6",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "Iya sori naufal, emang kebiasaan buruk gua kalau janjian jam 7 pasti nyampe jam 8 wkwkwk. Susah bangun cuy.",
      "timestamp": 1780887900
    }
  }' > /dev/null

echo "Sending Message 7..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780887960_7",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Yaudah pokoknya lo semua anggota divisi creative, jangan malu maluin.",
      "timestamp": 1780887960
    }
  }' > /dev/null

echo "Final block sent! Waiting for debounce queue to commit..."
