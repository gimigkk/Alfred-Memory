const container = document.getElementById('mynetwork');

// We use vis.DataSet for dynamic updates
const nodes = new vis.DataSet([]);
const edges = new vis.DataSet([]);

const data = {
    nodes: nodes,
    edges: edges
};

const options = {
    nodes: {
        shape: 'dot',
        font: {
            size: 14,
            color: '#000000',
            face: 'Inter'
        },
        borderWidth: 2,
        chosen: {
            node: function(values, id, selected, hovering) {
                if (hovering) {
                    values.borderWidth = 2;
                    values.borderColor = '#000000';
                }
            },
            label: false
        },
        shadow: {
            enabled: true,
            color: 'rgba(0,0,0,0.15)',
            size: 6,
            x: 0,
            y: 3
        }
    },
    edges: {
        width: 2,
        color: {
            color: '#94a3b8',
            highlight: '#94a3b8'
        },
        chosen: false,
        smooth: {
            type: 'continuous'
        },
        font: {
            size: 12,
            color: '#334155',
            face: 'Inter',
            strokeWidth: 3,
            strokeColor: '#ffffff',
            align: 'middle'
        }
    },
    groups: {
        Task: {
            color: {
                background: '#ef4444',
                border: '#ef4444',
                hover: { background: '#ef4444', border: '#ef4444' },
                highlight: { background: '#ef4444', border: '#ef4444' }
            }
        },
        Event: {
            color: {
                background: '#3b82f6',
                border: '#3b82f6',
                hover: { background: '#3b82f6', border: '#3b82f6' },
                highlight: { background: '#3b82f6', border: '#3b82f6' }
            }
        },
        Person: {
            color: {
                background: '#10b981',
                border: '#10b981',
                hover: { background: '#10b981', border: '#10b981' },
                highlight: { background: '#10b981', border: '#10b981' }
            }
        },
        Insight: {
            color: {
                background: '#f59e0b',
                border: '#f59e0b',
                hover: { background: '#f59e0b', border: '#f59e0b' },
                highlight: { background: '#f59e0b', border: '#f59e0b' }
            }
        },
        Circle: {
            color: {
                background: '#af52de',
                border: '#af52de',
                hover: { background: '#af52de', border: '#af52de' },
                highlight: { background: '#af52de', border: '#af52de' }
            }
        }
    },
    physics: {
        solver: 'forceAtlas2Based',
        forceAtlas2Based: {
            gravitationalConstant: -50,
            centralGravity: 0.01,
            springLength: 100,
            springConstant: 0.08,
            damping: 0.6,
            avoidOverlap: 1
        },
        maxVelocity: 50,
        timestep: 0.15,
        stabilization: { iterations: 150 }
    },
    interaction: {
        hover: true,
        tooltipDelay: 200,
        selectable: false
    }
};

const network = new vis.Network(container, data, options);

// Prevent select actions completely
network.on('selectNode', () => network.unselectAll());
network.on('selectEdge', () => network.unselectAll());

// Draw world-space background dot grid (Figma-like)
network.on('beforeDrawing', (ctx) => {
    const scale = network.getScale();
    if (scale > 0.2) {
        ctx.save();
        // Fade out grid dots slightly as you zoom out
        const opacity = Math.min(1, (scale - 0.2) / 0.3);
        ctx.fillStyle = `rgba(203, 213, 225, ${opacity})`;

        const gridSpacing = 40;
        const topLeft = network.DOMtoCanvas({x: 0, y: 0});
        const bottomRight = network.DOMtoCanvas({x: ctx.canvas.width, y: ctx.canvas.height});

        const startX = Math.floor(topLeft.x / gridSpacing) * gridSpacing;
        const startY = Math.floor(topLeft.y / gridSpacing) * gridSpacing;
        const endX = Math.ceil(bottomRight.x / gridSpacing) * gridSpacing;
        const endY = Math.ceil(bottomRight.y / gridSpacing) * gridSpacing;

        const dotRadius = Math.max(0.5, 1.2 / scale);

        for (let x = startX; x <= endX; x += gridSpacing) {
            for (let y = startY; y <= endY; y += gridSpacing) {
                ctx.beginPath();
                ctx.arc(x, y, dotRadius, 0, 2 * Math.PI);
                ctx.fill();
            }
        }
        ctx.restore();
    }
});

// Polling function - full bidirectional sync with vault
async function fetchGraph() {
    try {
        const response = await fetch('/api/vault');
        if (!response.ok) return;
        
        const serverData = await response.json();
        
        const serverNodes = serverData.nodes || [];
        const serverEdges = serverData.edges || [];

        // Build lookup sets of incoming IDs
        const incomingNodeIds = new Set(serverNodes.map(n => n.id));
        const incomingEdgeIds = new Set(serverEdges.map(e => e.id));

        // Remove nodes/edges that no longer exist in the vault
        nodes.getIds().forEach(id => {
            if (!incomingNodeIds.has(id)) nodes.remove(id);
        });
        edges.getIds().forEach(id => {
            if (!incomingEdgeIds.has(id)) edges.remove(id);
        });

        // Calculate degree (number of connected edges) for each node
        const edgeCounts = {};
        serverNodes.forEach(n => {
            edgeCounts[n.id] = 0;
        });
        serverEdges.forEach(e => {
            if (edgeCounts[e.from] !== undefined) {
                edgeCounts[e.from]++;
            }
            if (edgeCounts[e.to] !== undefined) {
                edgeCounts[e.to]++;
            }
        });

// Helper to escape HTML characters
function escapeHtml(unsafe) {
    return String(unsafe)
         .replace(/&/g, "&amp;")
         .replace(/</g, "&lt;")
         .replace(/>/g, "&gt;")
         .replace(/"/g, "&quot;")
         .replace(/'/g, "&#039;");
}

// Generate rich, structured HTML tooltip depending on node type
function generateTooltipHTML(node) {
    const props = node.properties || {};
    const group = node.group;
    
    let html = `<div class="tooltip-content">`;
    html += `<div class="tooltip-header"><span class="node-id">${escapeHtml(node.id)}</span><span class="node-type badge-${group}">${escapeHtml(group)}</span></div>`;
    html += `<div class="tooltip-body">`;
    
    const keys = Object.keys(props).sort();
    let orderedKeys = [];
    if (group === 'Person') {
        orderedKeys = ['name', 'phone_number', 'aliases', 'created_at', 'needs_clarification'];
    } else if (group === 'Task') {
        orderedKeys = ['content', 'status', 'due_date', 'priority', 'aliases', 'verbatim', 'history', 'created_at', 'needs_clarification', 'clarification_basis'];
    } else if (group === 'Event') {
        orderedKeys = ['content', 'status', 'start_date', 'aliases', 'history', 'created_at', 'needs_clarification', 'clarification_basis'];
    } else if (group === 'Insight') {
        orderedKeys = ['content', 'aliases', 'verbatim', 'history', 'created_at', 'needs_clarification', 'clarification_basis'];
    } else if (group === 'Circle') {
        orderedKeys = ['name', 'aliases', 'content', 'history', 'created_at', 'needs_clarification'];
    }

    const displayedKeys = new Set();
    const addKeyRow = (key) => {
        if (displayedKeys.has(key)) return;
        displayedKeys.add(key);
        
        let val = props[key];
        if (val === undefined || val === null) return;
        if (key === 'embedding' || key === 'verbatim_vector') return;
        
        let valStr = '';
        if (Array.isArray(val)) {
            if (val.length === 0) return;
            valStr = `<ul class="tooltip-list">${val.map(item => `<li>${escapeHtml(String(item))}</li>`).join('')}</ul>`;
        } else if (typeof val === 'object') {
            valStr = `<pre class="tooltip-json">${escapeHtml(JSON.stringify(val, null, 2))}</pre>`;
        } else if (typeof val === 'boolean') {
            const badgeClass = val ? 'badge-danger' : 'badge-success';
            valStr = `<span class="badge ${badgeClass}">${val ? 'Yes' : 'No'}</span>`;
        } else {
            const strVal = String(val);
            if (strVal.includes('T') && strVal.includes('Z')) {
                try {
                    valStr = escapeHtml(new Date(strVal).toLocaleString());
                } catch (e) {
                    valStr = escapeHtml(strVal);
                }
            } else {
                valStr = escapeHtml(strVal);
            }
        }
        
        const label = key.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase());
        html += `<div class="tooltip-row">`;
        html += `<span class="tooltip-label">${escapeHtml(label)}</span>`;
        html += `<div class="tooltip-value">${valStr}</div>`;
        html += `</div>`;
    };

    orderedKeys.forEach(addKeyRow);
    keys.forEach(addKeyRow);
    
    html += `</div></div>`;
    return html;
}

// Add or update nodes with dynamic size and mass
        serverNodes.forEach(n => {
            const degree = edgeCounts[n.id] || 0;
            // node size floor size 1 edge, not 0
            const effectiveDegree = Math.max(degree, 1);
            // node size gets bigger with more edges using a log equation
            n.size = 15 + 12 * Math.log(effectiveDegree + 1);
            // make hub nodes have more repulsion (mass)
            n.mass = 1 + degree;

            // Generate HTML tooltip element
            const tooltipEl = document.createElement('div');
            tooltipEl.innerHTML = generateTooltipHTML(n);
            n.title = tooltipEl;

            if (nodes.get(n.id)) {
                nodes.update(n);
            } else {
                nodes.add(n);
            }
        });

        // Add or update edges
        serverEdges.forEach(e => {
            if (edges.get(e.id)) {
                edges.update(e);
            } else {
                edges.add(e);
            }
        });
        
    } catch (e) {
        console.error("Error fetching graph:", e);
    }
}

// Poll every 2 seconds
setInterval(fetchGraph, 2000);
fetchGraph(); // Initial fetch
