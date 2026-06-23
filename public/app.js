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
        tooltipDelay: 0,
        selectable: true,
        selectConnectedEdges: false
    }
};

const network = new vis.Network(container, data, options);

let selectedNodeId = null;
let isUpdatingStyles = false;

const groupColors = {
    Task: { background: '#ef4444', border: '#ef4444' },
    Event: { background: '#3b82f6', border: '#3b82f6' },
    Person: { background: '#10b981', border: '#10b981' },
    Insight: { background: '#f59e0b', border: '#f59e0b' },
    Circle: { background: '#af52de', border: '#af52de' }
};

function hexToRGBA(hex, alpha) {
    if (!hex) return `rgba(0, 0, 0, ${alpha})`;
    const cleanHex = hex.replace('#', '');
    let r, g, b;
    if (cleanHex.length === 3) {
        r = parseInt(cleanHex[0] + cleanHex[0], 16);
        g = parseInt(cleanHex[1] + cleanHex[1], 16);
        b = parseInt(cleanHex[2] + cleanHex[2], 16);
    } else {
        r = parseInt(cleanHex.substring(0, 2), 16);
        g = parseInt(cleanHex.substring(2, 4), 16);
        b = parseInt(cleanHex.substring(4, 6), 16);
    }
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}

let currentTransition = null;

function animateTransition(targetNodes, targetEdges, duration = 250) {
    if (currentTransition) {
        cancelAnimationFrame(currentTransition.frameId);
    }

    const startTime = performance.now();
    const startNodes = {};
    nodes.get().forEach(n => {
        startNodes[n.id] = {
            color: {
                background: n.color?.background || '#97c2fc',
                border: n.color?.border || '#2b7ce9'
            },
            borderWidth: n.borderWidth !== undefined ? n.borderWidth : 2,
            fontColor: n.font?.color || '#000000'
        };
    });

    const startEdges = {};
    edges.get().forEach(e => {
        startEdges[e.id] = {
            color: e.color?.color || '#94a3b8',
            fontColor: e.font?.color || '#334155',
            strokeColor: e.font?.strokeColor || '#ffffff'
        };
    });

    function parseColor(colorStr) {
        if (!colorStr) return { r: 0, g: 0, b: 0, a: 1 };
        if (colorStr.startsWith('rgba')) {
            const parts = colorStr.match(/[\d\.]+/g);
            if (parts) {
                return {
                    r: parseFloat(parts[0]),
                    g: parseFloat(parts[1]),
                    b: parseFloat(parts[2]),
                    a: parts[3] !== undefined ? parseFloat(parts[3]) : 1
                };
            }
        } else if (colorStr.startsWith('#')) {
            const cleanHex = colorStr.replace('#', '');
            let r, g, b;
            if (cleanHex.length === 3) {
                r = parseInt(cleanHex[0] + cleanHex[0], 16);
                g = parseInt(cleanHex[1] + cleanHex[1], 16);
                b = parseInt(cleanHex[2] + cleanHex[2], 16);
            } else {
                r = parseInt(cleanHex.substring(0, 2), 16);
                g = parseInt(cleanHex.substring(2, 4), 16);
                b = parseInt(cleanHex.substring(4, 6), 16);
            }
            return { r, g, b, a: 1 };
        } else if (colorStr.startsWith('rgb')) {
            const parts = colorStr.match(/\d+/g);
            if (parts) {
                return { r: parseInt(parts[0]), g: parseInt(parts[1]), b: parseInt(parts[2]), a: 1 };
            }
        }
        return { r: 148, g: 163, b: 184, a: 1 };
    }

    const parsedTargetsNodes = {};
    targetNodes.forEach(tn => {
        parsedTargetsNodes[tn.id] = {
            color: {
                background: parseColor(tn.color?.background),
                border: parseColor(tn.color?.border)
            },
            borderWidth: tn.borderWidth !== undefined ? tn.borderWidth : 2,
            fontColor: parseColor(tn.font?.color)
        };
    });

    const parsedTargetsEdges = {};
    targetEdges.forEach(te => {
        parsedTargetsEdges[te.id] = {
            color: parseColor(te.color?.color),
            fontColor: parseColor(te.font?.color),
            strokeColor: parseColor(te.font?.strokeColor)
        };
    });

    const parsedStartsNodes = {};
    Object.keys(startNodes).forEach(id => {
        const sn = startNodes[id];
        parsedStartsNodes[id] = {
            color: {
                background: parseColor(sn.color.background),
                border: parseColor(sn.color.border)
            },
            borderWidth: sn.borderWidth,
            fontColor: parseColor(sn.fontColor)
        };
    });

    const parsedStartsEdges = {};
    Object.keys(startEdges).forEach(id => {
        const se = startEdges[id];
        parsedStartsEdges[id] = {
            color: parseColor(se.color),
            fontColor: parseColor(se.fontColor),
            strokeColor: parseColor(se.strokeColor)
        };
    });

    function interpolateColor(start, target, t) {
        const r = Math.round(start.r + (target.r - start.r) * t);
        const g = Math.round(start.g + (target.g - start.g) * t);
        const b = Math.round(start.b + (target.b - start.b) * t);
        const a = start.a + (target.a - start.a) * t;
        return `rgba(${r}, ${g}, ${b}, ${a})`;
    }

    function tick(now) {
        const elapsed = now - startTime;
        const progress = Math.min(1, elapsed / duration);
        // Cubic ease-out curve
        const t = 1 - Math.pow(1 - progress, 3);

        const nodesToUpdate = [];
        Object.keys(parsedStartsNodes).forEach(id => {
            const start = parsedStartsNodes[id];
            const target = parsedTargetsNodes[id] || start;

            const currentBg = interpolateColor(start.color.background, target.color.background, t);
            const currentBorder = interpolateColor(start.color.border, target.color.border, t);
            const currentFontColor = interpolateColor(start.fontColor, target.fontColor, t);
            const currentBorderWidth = start.borderWidth + (target.borderWidth - start.borderWidth) * t;

            nodesToUpdate.push({
                id: id,
                color: {
                    background: currentBg,
                    border: currentBorder,
                    hover: { background: currentBg, border: currentBorder },
                    highlight: { background: currentBg, border: currentBorder }
                },
                borderWidth: currentBorderWidth,
                font: { color: currentFontColor }
            });
        });

        const edgesToUpdate = [];
        Object.keys(parsedStartsEdges).forEach(id => {
            const start = parsedStartsEdges[id];
            const target = parsedTargetsEdges[id] || start;

            const currentColor = interpolateColor(start.color, target.color, t);
            const currentFontColor = interpolateColor(start.fontColor, target.fontColor, t);
            const currentStrokeColor = interpolateColor(start.strokeColor, target.strokeColor, t);

            edgesToUpdate.push({
                id: id,
                color: {
                    color: currentColor,
                    highlight: currentColor
                },
                font: {
                    color: currentFontColor,
                    strokeColor: currentStrokeColor
                }
            });
        });

        isUpdatingStyles = true;
        nodes.update(nodesToUpdate);
        edges.update(edgesToUpdate);
        isUpdatingStyles = false;

        if (progress < 1) {
            currentTransition.frameId = requestAnimationFrame(tick);
        } else {
            currentTransition = null;
        }
    }

    currentTransition = {
        frameId: requestAnimationFrame(tick)
    };
}

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
        orderedKeys = ['content', 'status', 'start_date', 'aliases', 'verbatim', 'history', 'created_at', 'needs_clarification', 'clarification_basis'];
    } else if (group === 'Insight') {
        orderedKeys = ['content', 'aliases', 'verbatim', 'history', 'created_at', 'needs_clarification', 'clarification_basis'];
    } else if (group === 'Circle') {
        orderedKeys = ['name', 'aliases', 'content', 'verbatim', 'history', 'created_at', 'needs_clarification'];
    }

    const displayedKeys = new Set();
    const addKeyRow = (key) => {
        if (displayedKeys.has(key)) return;
        displayedKeys.add(key);
        
        let val = props[key];
        if (val === undefined || val === null) return;
        if (key === 'embedding' || key === 'verbatim_vector') return;
        
        let valStr = '';
        if (key === 'history' && Array.isArray(val)) {
            if (val.length === 0) return;
            valStr = `<div class="history-timeline">` + val.map(item => {
                const itemStr = String(item);
                const parts = itemStr.split(' - ');
                if (parts.length >= 2) {
                    const time = escapeHtml(parts[0]);
                    const content = escapeHtml(parts.slice(1).join(' - '));
                    return `<div class="timeline-item"><div class="timeline-time">${time}</div><div class="timeline-content">${content}</div></div>`;
                }
                return `<div class="timeline-item"><div class="timeline-content">${escapeHtml(itemStr)}</div></div>`;
            }).join('') + `</div>`;
        } else if (Array.isArray(val)) {
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

function updateGraphStyles(bringToFront = false) {
    if (isUpdatingStyles) return;
    isUpdatingStyles = true;

    try {
        const allNodes = nodes.get();
        const allEdges = edges.get();

        if (selectedNodeId === null) {
            // Restore everything to default style
            const nodesToUpdate = allNodes.map(n => {
                const originalColor = groupColors[n.group] || { background: '#97c2fc', border: '#2b7ce9' };
                return {
                    id: n.id,
                    color: {
                        background: originalColor.background,
                        border: originalColor.border,
                        hover: { background: originalColor.background, border: originalColor.border },
                        highlight: { background: originalColor.background, border: originalColor.border }
                    },
                    borderWidth: 2,
                    font: { color: '#000000' }
                };
            });

            const edgesToUpdate = allEdges.map(e => {
                return {
                    id: e.id,
                    color: {
                        color: '#94a3b8',
                        highlight: '#94a3b8'
                    },
                    font: {
                        color: '#334155',
                        strokeColor: '#ffffff'
                    }
                };
            });

            if (bringToFront) {
                isUpdatingStyles = false;
                animateTransition(nodesToUpdate, edgesToUpdate, 250);
            } else {
                nodes.update(nodesToUpdate);
                edges.update(edgesToUpdate);
                isUpdatingStyles = false;
            }
            document.getElementById('side-panel').classList.add('hidden');
            return;
        }

        // Active node selection using BFS for depth-based highlighting
        const nodeDepths = new Map();
        nodeDepths.set(selectedNodeId, 0);

        let currentLevel = [selectedNodeId];
        for (let depth = 1; depth <= 3; depth++) {
            let nextLevel = [];
            for (const nodeId of currentLevel) {
                const neighbors = network.getConnectedNodes(nodeId);
                for (const neighbor of neighbors) {
                    if (!nodeDepths.has(neighbor)) {
                        nodeDepths.set(neighbor, depth);
                        nextLevel.push(neighbor);
                    }
                }
            }
            currentLevel = nextLevel;
        }

        const nodesToUpdate = [];
        const nodesToBringToFront = [];

        allNodes.forEach(n => {
            const originalColor = groupColors[n.group] || { background: '#97c2fc', border: '#2b7ce9' };
            const depth = nodeDepths.has(n.id) ? nodeDepths.get(n.id) : Infinity;

            if (depth === 0) {
                // Selected node: bold dark border, 100% opacity
                nodesToUpdate.push({
                    id: n.id,
                    color: {
                        background: originalColor.background,
                        border: '#0f172a',
                        hover: { background: originalColor.background, border: '#0f172a' },
                        highlight: { background: originalColor.background, border: '#0f172a' }
                    },
                    borderWidth: 4,
                    font: { color: '#000000' }
                });
                nodesToBringToFront.push(n.id);
            } else if (depth === 1) {
                // Direct connection: 100% opacity
                nodesToUpdate.push({
                    id: n.id,
                    color: {
                        background: originalColor.background,
                        border: originalColor.border,
                        hover: { background: originalColor.background, border: originalColor.border },
                        highlight: { background: originalColor.background, border: originalColor.border }
                    },
                    borderWidth: 2,
                    font: { color: '#000000' }
                });
                nodesToBringToFront.push(n.id);
            } else if (depth === 2) {
                // Depth 2: 40% opacity
                const bg = hexToRGBA(originalColor.background, 0.4);
                const border = hexToRGBA(originalColor.border, 0.4);
                nodesToUpdate.push({
                    id: n.id,
                    color: {
                        background: bg,
                        border: border,
                        hover: { background: bg, border: border },
                        highlight: { background: bg, border: border }
                    },
                    borderWidth: 2,
                    font: { color: 'rgba(15, 23, 42, 0.4)' }
                });
                nodesToBringToFront.push(n.id);
            } else if (depth === 3) {
                // Depth 3: 20% opacity
                const bg = hexToRGBA(originalColor.background, 0.2);
                const border = hexToRGBA(originalColor.border, 0.2);
                nodesToUpdate.push({
                    id: n.id,
                    color: {
                        background: bg,
                        border: border,
                        hover: { background: bg, border: border },
                        highlight: { background: bg, border: border }
                    },
                    borderWidth: 1.5,
                    font: { color: 'rgba(15, 23, 42, 0.2)' }
                });
                nodesToBringToFront.push(n.id);
            } else {
                // Depth > 3: heavily faded
                const fadedBg = hexToRGBA(originalColor.background, 0.1);
                const fadedBorder = hexToRGBA(originalColor.border, 0.15);
                nodesToUpdate.push({
                    id: n.id,
                    color: {
                        background: fadedBg,
                        border: fadedBorder,
                        hover: { background: fadedBg, border: fadedBorder },
                        highlight: { background: fadedBg, border: fadedBorder }
                    },
                    borderWidth: 1.5,
                    font: { color: 'rgba(15, 23, 42, 0.15)' }
                });
            }
        });

        // Bring active nodes (the selected one and its neighbors) to front
        if (bringToFront && nodesToBringToFront.length > 0) {
            const positions = network.getPositions(nodesToBringToFront);
            
            // Gather all connected edges first
            const connectedEdgeIdsSet = new Set();
            nodesToBringToFront.forEach(nodeId => {
                const edgeIds = network.getConnectedEdges(nodeId);
                edgeIds.forEach(id => connectedEdgeIdsSet.add(id));
            });
            const connectedEdgeIds = Array.from(connectedEdgeIdsSet);
            const edgeDatas = connectedEdgeIds.map(id => edges.get(id)).filter(Boolean);

            // Remove edges and then nodes
            edges.remove(connectedEdgeIds);
            nodes.remove(nodesToBringToFront);

            // Re-add nodes at their current positions with original colors to start the transition
            const nodesToAdd = nodesToBringToFront.map(nodeId => {
                const originalNode = allNodes.find(n => n.id === nodeId);
                const pos = positions[nodeId];
                if (originalNode && pos) {
                    return {
                        ...originalNode,
                        x: pos.x,
                        y: pos.y,
                        color: originalNode.color,
                        borderWidth: originalNode.borderWidth,
                        font: originalNode.font
                    };
                }
                return originalNode;
            }).filter(Boolean);

            nodes.add(nodesToAdd);
            edges.add(edgeDatas);

            // Re-select the selected node programmatically
            network.selectNodes([selectedNodeId]);
        }

        // Handle edges color and font labels based on depth
        const edgesToUpdate = allEdges.map(e => {
            const depthFrom = nodeDepths.has(e.from) ? nodeDepths.get(e.from) : Infinity;
            const depthTo = nodeDepths.has(e.to) ? nodeDepths.get(e.to) : Infinity;
            const edgeDepth = Math.max(depthFrom, depthTo);

            if (edgeDepth <= 1) {
                return {
                    id: e.id,
                    color: { color: '#64748b', highlight: '#64748b' },
                    font: { color: '#334155', strokeColor: '#ffffff' }
                };
            } else if (edgeDepth === 2) {
                return {
                    id: e.id,
                    color: { color: 'rgba(100, 116, 139, 0.4)', highlight: 'rgba(100, 116, 139, 0.4)' },
                    font: { color: 'rgba(51, 65, 85, 0.4)', strokeColor: 'rgba(255, 255, 255, 0.4)' }
                };
            } else if (edgeDepth === 3) {
                return {
                    id: e.id,
                    color: { color: 'rgba(100, 116, 139, 0.2)', highlight: 'rgba(100, 116, 139, 0.2)' },
                    font: { color: 'rgba(51, 65, 85, 0.2)', strokeColor: 'rgba(255, 255, 255, 0.2)' }
                };
            } else {
                return {
                    id: e.id,
                    color: { color: 'rgba(148, 163, 184, 0.1)', highlight: 'rgba(148, 163, 184, 0.1)' },
                    font: { color: 'rgba(51, 65, 85, 0.1)', strokeColor: 'rgba(255, 255, 255, 0.1)' }
                };
            }
        });

        if (bringToFront) {
            isUpdatingStyles = false;
            animateTransition(nodesToUpdate, edgesToUpdate, 250);
        } else {
            nodes.update(nodesToUpdate);
            edges.update(edgesToUpdate);
            isUpdatingStyles = false;
        }

        const selectedNodeData = allNodes.find(n => n.id === selectedNodeId);
        if (selectedNodeData) {
            document.getElementById('side-panel-content').innerHTML = generateTooltipHTML(selectedNodeData);
            document.getElementById('side-panel').classList.remove('hidden');
        }

    } catch (e) {
        console.error("Error updating graph styles:", e);
        isUpdatingStyles = false;
    }
}

network.on('click', function (params) {
    if (isUpdatingStyles) return;

    const clickedNodes = params.nodes;
    
    if (clickedNodes.length > 0) {
        const clickedNodeId = clickedNodes[0];
        if (clickedNodeId === selectedNodeId) {
            // Clicked the already selected node -> deselect (toggle off)
            selectedNodeId = null;
            isUpdatingStyles = true;
            network.unselectAll();
            isUpdatingStyles = false;
            updateGraphStyles(true);
        } else {
            // Clicked a new/different node -> select
            selectedNodeId = clickedNodeId;
            updateGraphStyles(true);
        }
    } else {
        // Clicked background or edge -> deselect
        selectedNodeId = null;
        isUpdatingStyles = true;
        network.unselectAll();
        isUpdatingStyles = false;
        updateGraphStyles(true);
    }
});

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

        // Reapply selection styles to any new/updated nodes/edges
        updateGraphStyles(false);
        
    } catch (e) {
        console.error("Error fetching graph:", e);
    }
}

// Poll every 2 seconds
setInterval(fetchGraph, 2000);
fetchGraph(); // Initial fetch
