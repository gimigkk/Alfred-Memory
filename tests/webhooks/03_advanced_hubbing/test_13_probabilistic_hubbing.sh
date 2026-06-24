#!/bin/bash

echo "==========================================="
echo " TEST 13: PROBABILISTIC HUBBING (RULE 2B) "
echo "==========================================="
echo "Context: Rendi and Jeslyn discuss bringing equipment and printed materials."
echo "They do NOT mention 'Rapat Finalisasi Venue' or the exact date."
echo "Agent should infer the link via Domain Overlap + Active Obligations."

echo "Sending Message 1..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780891000_1",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "rendi_ramadana@c.us",
      "fromMe": false,
      "body": "eh gw udah nyewa sound system sama ambil proyektor nih, ntar malam gw bawa ya.",
      "timestamp": 1780891000
    }
  }' > /dev/null

echo "Sending Message 2..."
curl -s -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": "message",
    "payload": {
      "id": "msg_1780891060_2",
      "from": "1234567890@g.us",
      "to": "me@c.us",
      "participant": "jeslyn_ieee@c.us",
      "fromMe": false,
      "body": "mantap ren, gw juga udah ngeprint kertasnya buat dibagi-bagiin ntar.",
      "timestamp": 1780891060
    }
  }' > /dev/null

echo "Final block sent! Waiting for debounce queue to commit..."
