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
        size: 20,
        font: {
            size: 14,
            color: '#ffffff',
            face: 'Inter'
        },
        borderWidth: 2,
        shadow: {
            enabled: true,
            color: 'rgba(0,0,0,0.5)',
            size: 10,
            x: 0,
            y: 5
        }
    },
    edges: {
        width: 2,
        color: {
            color: 'rgba(255, 255, 255, 0.3)',
            highlight: '#bb86fc'
        },
        smooth: {
            type: 'continuous'
        },
        font: {
            size: 12,
            color: '#aaaaaa',
            face: 'Inter',
            strokeWidth: 0,
            background: 'none',
            align: 'middle'
        }
    },
    groups: {
        Task: {
            color: { background: '#cf6679', border: '#ff7597' }
        },
        Event: {
            color: { background: '#03dac6', border: '#04fce5' }
        },
        Person: {
            color: { background: '#bb86fc', border: '#d5b2ff' }
        },
        Insight: {
            color: { background: '#ffb74d', border: '#ffcc80' }
        }
    },
    physics: {
        forceAtlas2Based: {
            gravitationalConstant: -50,
            centralGravity: 0.01,
            springLength: 100,
            springConstant: 0.08
        },
        maxVelocity: 50,
        solver: 'forceAtlas2Based',
        timestep: 0.35,
        stabilization: { iterations: 150 }
    },
    interaction: {
        hover: true,
        tooltipDelay: 200
    }
};

const network = new vis.Network(container, data, options);

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

        // Add or update nodes
        serverNodes.forEach(n => {
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
