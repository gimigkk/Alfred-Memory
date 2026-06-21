const fs = require('fs');
const path = require('path');

const inputFile = process.argv[2];
const outputFile = process.argv[3];

if (!inputFile || !outputFile) {
    console.error("Usage: node generate_test.js <input_txt> <output_sh>");
    process.exit(1);
}

const rawText = fs.readFileSync(inputFile, 'utf-8');
const lines = rawText.split('\n').filter(line => line.trim().length > 0);

let bashScript = `#!/bin/bash\n\n`;
let msgCounter = 1;
let lastTimestamp = null;
let blockCounter = 1;

bashScript += `echo "==========================================="\n`;
bashScript += `echo " BLOCK ${blockCounter}"\n`;
bashScript += `echo "==========================================="\n`;

for (const line of lines) {
    // regex to match: [6/7, 11:57] Sender Name: message body
    const match = line.match(/^\[([^\]]+)\]\s([^:]+):\s(.*)$/);
    if (!match) continue;
    
    let [_, timeStr, sender, body] = match;
    
    // Parse timeStr e.g. "6/7, 11:57" or "6/11, 19:19"
    // Assuming year is 2026 based on Alfred.md
    const [datePart, timePart] = timeStr.split(', ');
    const [month, day] = datePart.split('/');
    const [hour, minute] = timePart.split(':');
    
    // Create Date object in UTC or Local
    const msgDate = new Date(2026, parseInt(month) - 1, parseInt(day), parseInt(hour), parseInt(minute));
    const currentTimestamp = Math.floor(msgDate.getTime() / 1000);

    // If gap is > 30 minutes (1800 seconds), insert a sleep
    if (lastTimestamp !== null && (currentTimestamp - lastTimestamp) > 1800) {
        bashScript += `echo "Time gap of > 30 minutes detected. Waiting 2 minutes for ingestion agent to process block..."\n`;
        bashScript += `sleep 120\n\n`;
        blockCounter++;
        bashScript += `echo "==========================================="\n`;
        bashScript += `echo " BLOCK ${blockCounter}"\n`;
        bashScript += `echo "==========================================="\n`;
    }
    
    lastTimestamp = currentTimestamp;
    
    // basic phone number formatting or name to ID
    let senderId = sender.trim().replace(/[^a-zA-Z0-9]/g, '_').toLowerCase() + "@c.us";
    
    // Escape quotes for bash and JSON
    let safeBody = body.replace(/\\/g, '\\\\').replace(/"/g, '\\"');

    bashScript += `echo "Sending Message ${msgCounter}..."\n`;
    bashScript += `curl -s -X POST http://localhost:8080/api/webhook \\\n`;
    bashScript += `  -H "Content-Type: application/json" \\\n`;
    bashScript += `  -d '{\n`;
    bashScript += `    "event": "message",\n`;
    bashScript += `    "payload": {\n`;
    bashScript += `      "id": "msg_${currentTimestamp}_${msgCounter}",\n`;
    bashScript += `      "from": "1234567890@g.us",\n`;
    bashScript += `      "to": "me@c.us",\n`;
    bashScript += `      "participant": "${senderId}",\n`;
    bashScript += `      "fromMe": false,\n`;
    bashScript += `      "body": "${safeBody}",\n`;
    bashScript += `      "timestamp": ${currentTimestamp}\n`;
    bashScript += `    }\n`;
    bashScript += `  }' > /dev/null\n\n`;
    
    msgCounter++;
}

bashScript += `echo "Final block sent! Waiting for debounce queue to commit..."\n`;

fs.writeFileSync(outputFile, bashScript, { mode: 0o755 });
console.log(`Generated ${outputFile} with ${msgCounter - 1} messages across ${blockCounter} blocks.`);
