const fs = require('fs');
const appJsCode = fs.readFileSync('./public/app.js', 'utf8');

// A very basic mock for generateTooltipHTML evaluation
global.escapeHtml = function(unsafe) {
    return String(unsafe).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
};

// Extract the generateTooltipHTML function
const functionMatch = appJsCode.match(/function generateTooltipHTML[\s\S]*?return html;\n}/);
if (functionMatch) {
    eval(functionMatch[0]);
    
    const mockNode = {
        id: "Task-123",
        group: "Task",
        properties: {
            content: "Test task",
            group_mentions: [
                {
                    speaker: "jeslyn_ieee",
                    phrase: "divisi acara",
                    quote: "rundown dari divisi acara udah fix",
                    note: "ambiguous parent org"
                }
            ]
        }
    };
    
    console.log(generateTooltipHTML(mockNode));
}
