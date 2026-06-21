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
      "id": "msg_1780811160_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "@M3-117_Gilang Muhamad W ini ada yg live report kan...",
      "timestamp": 1780811160
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780811460_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "lu kgk ?",
      "timestamp": 1780811460
    }
  }' > /dev/null

echo "Sending Message 3..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780811820_3",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Ada",
      "timestamp": 1780811820
    }
  }' > /dev/null

echo "Sending Message 4..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780812480_4",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "udah sampe",
      "timestamp": 1780812480
    }
  }' > /dev/null

echo "Sending Message 5..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780812720_5",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "ni gd tanya jwb qil? @~Syazana Aqila",
      "timestamp": 1780812720
    }
  }' > /dev/null

echo "Sending Message 6..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780812780_6",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "tergantung manager yg presentasi sbnrnya itu",
      "timestamp": 1780812780
    }
  }' > /dev/null

echo "Sending Message 7..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780812780_7",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "kita kasih 10 mnt buat progress",
      "timestamp": 1780812780
    }
  }' > /dev/null

echo "Sending Message 8..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780812780_8",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "takut ga cukup aja sih sbnrnya",
      "timestamp": 1780812780
    }
  }' > /dev/null

echo "Sending Message 9..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780812780_9",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "ngapa lu share anyg @M Naufal IEEE²⁵",
      "timestamp": 1780812780
    }
  }' > /dev/null

echo "Sending Message 10..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780812840_10",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "kan dia gada manager",
      "timestamp": 1780812840
    }
  }' > /dev/null

echo "Sending Message 11..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780812960_11",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "@~Syazana Aqila ini ga ada absen qil?",
      "timestamp": 1780812960
    }
  }' > /dev/null

echo "Sending Message 12..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813020_12",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "ada",
      "timestamp": 1780813020
    }
  }' > /dev/null

echo "Sending Message 13..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813020_13",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "kalo gak telat mah",
      "timestamp": 1780813020
    }
  }' > /dev/null

echo "Sending Message 14..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813020_14",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "ada tadi di depan",
      "timestamp": 1780813020
    }
  }' > /dev/null

echo "Sending Message 15..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813020_15",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "2x",
      "timestamp": 1780813020
    }
  }' > /dev/null

echo "Sending Message 16..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813020_16",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "kgk ada orang yang standby maksud gw",
      "timestamp": 1780813020
    }
  }' > /dev/null

echo "Sending Message 17..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813020_17",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "tul",
      "timestamp": 1780813020
    }
  }' > /dev/null

echo "Sending Message 18..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813020_18",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "trus lu kapan kesini",
      "timestamp": 1780813020
    }
  }' > /dev/null

echo "Sending Message 19..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813020_19",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "karena lu telat boyy",
      "timestamp": 1780813020
    }
  }' > /dev/null

echo "Sending Message 20..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813020_20",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "ya kan tetep aja harus ada absen",
      "timestamp": 1780813020
    }
  }' > /dev/null

echo "Sending Message 21..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813080_21",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "yang standby",
      "timestamp": 1780813080
    }
  }' > /dev/null

echo "Sending Message 22..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813080_22",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "buat data",
      "timestamp": 1780813080
    }
  }' > /dev/null

echo "Sending Message 23..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813080_23",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "engak lah",
      "timestamp": 1780813080
    }
  }' > /dev/null

echo "Sending Message 24..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813080_24",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "close gate anyg",
      "timestamp": 1780813080
    }
  }' > /dev/null

echo "Sending Message 25..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813080_25",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "nanti ditagih diakhir buat yg blm absen",
      "timestamp": 1780813080
    }
  }' > /dev/null

echo "Sending Message 26..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813080_26",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "todayy",
      "timestamp": 1780813080
    }
  }' > /dev/null

echo "Sending Message 27..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813140_27",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "gausa jawab kalo jawabannya gobl",
      "timestamp": 1780813140
    }
  }' > /dev/null

echo "Sending Message 28..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813140_28",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "😭😭chill gua datengg",
      "timestamp": 1780813140
    }
  }' > /dev/null

echo "Sending Message 29..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_29",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Bener",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 30..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_30",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Sini masuk",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 31..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_31",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rendi_ramadana_ieee__@c.us",
      "fromMe": false,
      "body": "Gwe gapake jaket IEEE",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 32..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_32",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Lahh mpruyy",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 33..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_33",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "gapapa",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 34..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_34",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "yaudah",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 35..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_35",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "lewat blkng",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 36..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_36",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Aowkowwkoakwo",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 37..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_37",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "rennn lewat blkng ren",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 38..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_38",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Lucu lg",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 39..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_39",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rendi_ramadana_ieee__@c.us",
      "fromMe": false,
      "body": "bukain",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 40..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_40",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Ngefoto dr kaca",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 41..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813320_41",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "mama aku salahs eragam",
      "timestamp": 1780813320
    }
  }' > /dev/null

echo "Sending Message 42..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813380_42",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "tai lu",
      "timestamp": 1780813380
    }
  }' > /dev/null

echo "Sending Message 43..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813380_43",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "manja jir",
      "timestamp": 1780813380
    }
  }' > /dev/null

echo "Sending Message 44..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813380_44",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "tinggal dorong",
      "timestamp": 1780813380
    }
  }' > /dev/null

echo "Sending Message 45..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813380_45",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Rendi",
      "timestamp": 1780813380
    }
  }' > /dev/null

echo "Sending Message 46..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813380_46",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "makanya",
      "timestamp": 1780813380
    }
  }' > /dev/null

echo "Sending Message 47..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813380_47",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "pake jaket ssrs pula",
      "timestamp": 1780813380
    }
  }' > /dev/null

echo "Sending Message 48..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813380_48",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "gw bakar lama2",
      "timestamp": 1780813380
    }
  }' > /dev/null

echo "Sending Message 49..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813500_49",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "wkwokwowk",
      "timestamp": 1780813500
    }
  }' > /dev/null

echo "Sending Message 50..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813500_50",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "guys reminder ke manager yaa nanti presentasinya 10 mnt",
      "timestamp": 1780813500
    }
  }' > /dev/null

echo "Sending Message 51..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813500_51",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "gada timer kah qil didpean?",
      "timestamp": 1780813500
    }
  }' > /dev/null

echo "Sending Message 52..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813620_52",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "ga ada lagi...",
      "timestamp": 1780813620
    }
  }' > /dev/null

echo "Sending Message 53..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813740_53",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "LOLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLll",
      "timestamp": 1780813740
    }
  }' > /dev/null

echo "Sending Message 54..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813740_54",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "apa gaada timekeeper",
      "timestamp": 1780813740
    }
  }' > /dev/null

echo "Sending Message 55..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813800_55",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "iya qil seenggaknya ada satu orang jdi timekeeper",
      "timestamp": 1780813800
    }
  }' > /dev/null

echo "Sending Message 56..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813860_56",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "manrisk si, jangan sampe lewat jam pinjem sampe kena usir",
      "timestamp": 1780813860
    }
  }' > /dev/null

echo "Sending Message 57..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813860_57",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "wait wait",
      "timestamp": 1780813860
    }
  }' > /dev/null

echo "Sending Message 58..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813860_58",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "itu petugas ga dikasi makan?",
      "timestamp": 1780813860
    }
  }' > /dev/null

echo "Sending Message 59..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813860_59",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "ini dah gua bawa makanannya jir",
      "timestamp": 1780813860
    }
  }' > /dev/null

echo "Sending Message 60..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813860_60",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "petugasnya dmn dah",
      "timestamp": 1780813860
    }
  }' > /dev/null

echo "Sending Message 61..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813920_61",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rezonaldo_ieee__@c.us",
      "fromMe": false,
      "body": "wa aja",
      "timestamp": 1780813920
    }
  }' > /dev/null

echo "Sending Message 62..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780813980_62",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "iya tar anak gua yg kasih langsung",
      "timestamp": 1780813980
    }
  }' > /dev/null

echo "Sending Message 63..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780814100_63",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "dah dikasi",
      "timestamp": 1780814100
    }
  }' > /dev/null

echo "Sending Message 64..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780814640_64",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "qiluyy mau minta yg stikk oren lgi boleh gaa",
      "timestamp": 1780814640
    }
  }' > /dev/null

echo "Sending Message 65..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780814640_65",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "bolee",
      "timestamp": 1780814640
    }
  }' > /dev/null

echo "Sending Message 66..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780814640_66",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "_62_896_1099_1799@c.us",
      "fromMe": false,
      "body": "waitt",
      "timestamp": 1780814640
    }
  }' > /dev/null

echo "Sending Message 67..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780814640_67",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "maaci cintahh",
      "timestamp": 1780814640
    }
  }' > /dev/null

echo "Final block sent! Waiting for debounce queue to commit..."
