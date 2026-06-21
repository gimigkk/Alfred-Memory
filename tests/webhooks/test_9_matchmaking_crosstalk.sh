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
      "id": "msg_1780834800_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Btw harus matchmakieeeng apa nggak",
      "timestamp": 1780834800
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780834860_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Pake aja",
      "timestamp": 1780834860
    }
  }' > /dev/null

echo "Sending Message 3..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780834860_3",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Biar sekalian ada manfaatnya",
      "timestamp": 1780834860
    }
  }' > /dev/null

echo "Sending Message 4..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780834860_4",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Maks brp org itu",
      "timestamp": 1780834860
    }
  }' > /dev/null

echo "Sending Message 5..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780834920_5",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Infonya dmn eh",
      "timestamp": 1780834920
    }
  }' > /dev/null

echo "Sending Message 6..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780834920_6",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Dh",
      "timestamp": 1780834920
    }
  }' > /dev/null

echo "Sending Message 7..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780834920_7",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Bingung gw webnya gada isinya",
      "timestamp": 1780834920
    }
  }' > /dev/null

echo "Sending Message 8..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780835280_8",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Kata jipiti 4",
      "timestamp": 1780835280
    }
  }' > /dev/null

echo "Sending Message 9..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780835280_9",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Tapi kata jipiti dan blm ada info tentang tahun ini kan",
      "timestamp": 1780835280
    }
  }' > /dev/null

echo "Sending Message 10..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780835280_10",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Trs geratis btw",
      "timestamp": 1780835280
    }
  }' > /dev/null

echo "Sending Message 11..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780835280_11",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Gas",
      "timestamp": 1780835280
    }
  }' > /dev/null

echo "Sending Message 12..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780835280_12",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "IEEE gaperlu nangis² mikirin dana",
      "timestamp": 1780835280
    }
  }' > /dev/null

echo "Sending Message 13..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780835280_13",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Oke guys yang token ainya banyak mohon tunjuk tangan",
      "timestamp": 1780835280
    }
  }' > /dev/null

echo "Sending Message 14..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780835340_14",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Saya saya",
      "timestamp": 1780835340
    }
  }' > /dev/null

echo "Sending Message 15..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780835460_15",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Eh kita blm ada tutor matchmakieeeng ya",
      "timestamp": 1780835460
    }
  }' > /dev/null

echo "Sending Message 16..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780835460_16",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Hooh",
      "timestamp": 1780835460
    }
  }' > /dev/null

echo "Sending Message 17..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780835460_17",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Td kan dijelaisn candra msh bajyak yg blm tau",
      "timestamp": 1780835460
    }
  }' > /dev/null

echo "Sending Message 18..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780835460_18",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Itu matchmakieeeng pake automation apa nggak dah",
      "timestamp": 1780835460
    }
  }' > /dev/null

echo "Sending Message 19..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836300_19",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "konteks automation? maksud lu kalo ada yang daftar otomatis ada bot discord yang announce kah?",
      "timestamp": 1780836300
    }
  }' > /dev/null

echo "Sending Message 20..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836360_20",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Semuanya via bot discord",
      "timestamp": 1780836360
    }
  }' > /dev/null

echo "Sending Message 21..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836360_21",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Ga google form",
      "timestamp": 1780836360
    }
  }' > /dev/null

echo "Sending Message 22..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836360_22",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Atau kenapa ga bikin web dah",
      "timestamp": 1780836360
    }
  }' > /dev/null

echo "Sending Message 23..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836360_23",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Wahai anak ilmu komputer ieee",
      "timestamp": 1780836360
    }
  }' > /dev/null

echo "Sending Message 24..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836360_24",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "tadi pas presentasi ga disampein kah",
      "timestamp": 1780836360
    }
  }' > /dev/null

echo "Sending Message 25..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836360_25",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "divisi gua",
      "timestamp": 1780836360
    }
  }' > /dev/null

echo "Sending Message 26..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836360_26",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Apa",
      "timestamp": 1780836360
    }
  }' > /dev/null

echo "Sending Message 27..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836360_27",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Web buat matchmakieeeng?",
      "timestamp": 1780836360
    }
  }' > /dev/null

echo "Sending Message 28..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836360_28",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "gua dah ngajuin ke tech infor tapi katanya 'nanti dulu' ceunah",
      "timestamp": 1780836360
    }
  }' > /dev/null

echo "Sending Message 29..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836420_29",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "AJAK CREATIVE PLEASEEEE",
      "timestamp": 1780836420
    }
  }' > /dev/null

echo "Sending Message 30..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836420_30",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "JADI DESIGNER UIUX",
      "timestamp": 1780836420
    }
  }' > /dev/null

echo "Sending Message 31..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836420_31",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "lah udah bijin desainnya",
      "timestamp": 1780836420
    }
  }' > /dev/null

echo "Sending Message 32..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836420_32",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "bikin",
      "timestamp": 1780836420
    }
  }' > /dev/null

echo "Sending Message 33..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836420_33",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "mau",
      "timestamp": 1780836420
    }
  }' > /dev/null

echo "Sending Message 34..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836420_34",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "AGHHHH",
      "timestamp": 1780836420
    }
  }' > /dev/null

echo "Sending Message 35..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836420_35",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Pen liat",
      "timestamp": 1780836420
    }
  }' > /dev/null

echo "Sending Message 36..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836480_36",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "gabisa dibuka website ormnya @M Naufal IEEE²⁵",
      "timestamp": 1780836480
    }
  }' > /dev/null

echo "Sending Message 37..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836480_37",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "meskipun pake ipb access",
      "timestamp": 1780836480
    }
  }' > /dev/null

echo "Sending Message 38..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836480_38",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "katanya bisa coo??",
      "timestamp": 1780836480
    }
  }' > /dev/null

echo "Sending Message 39..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836480_39",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "gabisa anak anak gua juga nyoba gabisa",
      "timestamp": 1780836480
    }
  }' > /dev/null

echo "Sending Message 40..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836480_40",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "auh pk sony blm ttd dah",
      "timestamp": 1780836480
    }
  }' > /dev/null

echo "Sending Message 41..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836480_41",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "vroo",
      "timestamp": 1780836480
    }
  }' > /dev/null

echo "Sending Message 42..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836720_42",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Gw bisa",
      "timestamp": 1780836720
    }
  }' > /dev/null

echo "Sending Message 43..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836720_43",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Tadi",
      "timestamp": 1780836720
    }
  }' > /dev/null

echo "Sending Message 44..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836720_44",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "linuk only",
      "timestamp": 1780836720
    }
  }' > /dev/null

echo "Sending Message 45..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836720_45",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "mpruy",
      "timestamp": 1780836720
    }
  }' > /dev/null

echo "Sending Message 46..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836780_46",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "@M Naufal IEEE²⁵ @Rafid Harsyah",
      "timestamp": 1780836780
    }
  }' > /dev/null

echo "Sending Message 47..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836840_47",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "tf mn",
      "timestamp": 1780836840
    }
  }' > /dev/null

echo "Sending Message 48..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836840_48",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "es gobak sodor apaan dh",
      "timestamp": 1780836840
    }
  }' > /dev/null

echo "Sending Message 49..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836840_49",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "0954959895 - BCA a.n Apta Adi Nur Fiansah ",
      "timestamp": 1780836840
    }
  }' > /dev/null

echo "Sending Message 50..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836900_50",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "makanya ikut",
      "timestamp": 1780836900
    }
  }' > /dev/null

echo "Sending Message 51..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836900_51",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "901158617350 - Jeslyn Angelica Widjaja - Seabank",
      "timestamp": 1780836900
    }
  }' > /dev/null

echo "Sending Message 52..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836900_52",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "pikm",
      "timestamp": 1780836900
    }
  }' > /dev/null

echo "Sending Message 53..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780836900_53",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafid_harsyah@c.us",
      "fromMe": false,
      "body": "Ongkeh",
      "timestamp": 1780836900
    }
  }' > /dev/null

echo "Sending Message 54..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780837020_54",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "2",
      "timestamp": 1780837020
    }
  }' > /dev/null

echo "Sending Message 55..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780837080_55",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "AI ini",
      "timestamp": 1780837080
    }
  }' > /dev/null

echo "Sending Message 56..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780837320_56",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m_naufal_ieee__@c.us",
      "fromMe": false,
      "body": "asli ini ini",
      "timestamp": 1780837320
    }
  }' > /dev/null

echo "Sending Message 57..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780837380_57",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "apta_ieee25@c.us",
      "fromMe": false,
      "body": "lu tf gopay ye?",
      "timestamp": 1780837380
    }
  }' > /dev/null

echo "Time gap of > 30 minutes detected. Waiting 2 minutes for ingestion agent to process block..."
sleep 120

echo "==========================================="
echo " BLOCK 2"
echo "==========================================="
echo "Sending Message 58..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780840020_58",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rafif_ilmany_ieee25@c.us",
      "fromMe": false,
      "body": "eh ayo gas ini",
      "timestamp": 1780840020
    }
  }' > /dev/null

echo "Sending Message 59..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780840260_59",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "m3_117_gilang_muhamad_w@c.us",
      "fromMe": false,
      "body": "Gweh mau",
      "timestamp": 1780840260
    }
  }' > /dev/null

echo "Final block sent! Waiting for debounce queue to commit..."
