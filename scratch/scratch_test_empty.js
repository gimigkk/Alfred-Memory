const fs = require('fs');
const appJsCode = fs.readFileSync('./public/app.js', 'utf8');

global.escapeHtml = function(unsafe) {
    return String(unsafe).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
};

const functionMatch = appJsCode.match(/function generateTooltipHTML[\s\S]*?return html;\n}/);
if (functionMatch) {
    eval(functionMatch[0]);
    
    const mockNode = {
        id: "Task-456",
        group: "Task",
        properties: {
            content: "Test task 2",
            group_mentions: []
        }
    };
    
    console.log(generateTooltipHTML(mockNode));
}
