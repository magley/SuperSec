<script setup>
import * as d3 from 'd3';
import { onMounted, ref } from 'vue';

const hmm = {
    "name": "basic",
    "relations": {
        "owner": {},
        "reviewer": {},
        "editor": {
            "union": [
                {
                    "this": {}
                },
                {
                    "computed_userset": {
                        "relation": "owner"
                    }
                },
                {
                    "computed_userset": {
                        "relation": "reviewer"
                    }
                }
            ]
        },
        "viewer": {
            "union": [
                {
                    "this": {}
                },
                {
                    "computed_userset": {
                        "relation": "editor"
                    }
                }
            ]
        },
        "commenter": {
            "union": [
                {
                    "this": {}
                },
                {
                    "computed_userset": {
                        "relation": "reviewer"
                    }
                }
            ]
        }
    }
}


const namespaceToGraph = (data) => {
    const namespaceName = data.name;

    let idCounter = 1;
    let idMap = {};

    let G = {
        nodes: [],
        links: []
    }

    G.nodes.push({
        "id": idCounter,
        "label": namespaceName,
        "type": "namespace"
    });
    idMap[namespaceName] = idCounter;
    idCounter++;

    for (const [key, value] of Object.entries(data.relations)) {
        G.nodes.push({
            "id": idCounter,
            "label": key,
            "type": "role"
        });
        idMap[key] = idCounter;
        idCounter++;
    }

    for (const [key, value] of Object.entries(data.relations)) {
        if (value.union != undefined) {
            for (let elem of value.union) {
                if (elem.computed_userset != undefined) {
                    const parent = elem.computed_userset.relation;

                    G.links.push({
                        "source": idMap[key],
                        "target": idMap[parent],
                        "dashed": true
                    });
                }
            }
        }

        G.links.push({
            "source": idMap[key],
            "target": idMap[namespaceName]
        });
    }

    return G;
}

onMounted(() => {
    // const data = {
    //     "nodes": [
    //     {
    //         "id": 1,
    //         "name": "A",
    //         "type": "namespace",
    //     },
    //     {
    //         "id": 2,
    //         "name": "B",
    //         "type": "role",
    //     },
    //     {
    //         "id": 3,
    //         "name": "C",
    //         "type": "role",
    //     },
    //     {
    //         "id": 4,
    //         "name": "D",
    //         "type": "role",
    //     },
    //     {
    //         "id": 5,
    //         "name": "E",
    //         "type": "role",
    //     },
    //     {
    //         "id": 6,
    //         "name": "F",
    //         "type": "role",
    //     },
    //     {
    //         "id": 7,
    //         "name": "G",
    //         "type": "role",
    //     },
    //     {
    //         "id": 8,
    //         "name": "H",
    //         "type": "role",
    //     },
    //     {
    //         "id": 9,
    //         "name": "I",
    //         "type": "role",
    //     },
    //     {
    //         "id": 10,
    //         "name": "J",
    //         "type": "role",
    //     }
    //     ],
    //     "links": [

    //     {
    //         "source": 1,
    //         "target": 2
    //     },
    //     {
    //         "source": 1,
    //         "target": 5
    //     },
    //     {
    //         "source": 1,
    //         "target": 6
    //     },

    //     {
    //         "source": 2,
    //         "target": 3
    //     },
    //             {
    //         "source": 2,
    //         "target": 7
    //     }
    //     ,

    //     {
    //         "source": 3,
    //         "target": 4
    //     },
    //         {
    //         "source": 8,
    //         "target": 3
    //     }
    //     ,
    //     {
    //         "source": 4,
    //         "target": 5
    //     }
    //     ,

    //     {
    //         "source": 4,
    //         "target": 9
    //     },
    //     {
    //         "source": 5,
    //         "target": 10
    //     }
    //     ]
    // };

    const data = namespaceToGraph(hmm);

    // Specify the dimensions of the chart.
    const width = 928;
    const height = 500;
    const nodeRadius = 5;

    const arrowH = 8.0;

    // Specify the color scale.
    const color = d3.scaleOrdinal(d3.schemeCategory10);

    // The force simulation mutates links and nodes, so create a copy
    // so that re-evaluating this cell produces the same result.
    const links = data.links.map(d => ({...d}));
    const nodes = data.nodes.map(d => ({...d}));

    // Create a simulation with several forces.
    const simulation = d3.forceSimulation(nodes)
        .force("link", d3.forceLink(links).id(d => d.id).distance(50))
        .force("charge", d3.forceManyBody())
        .force("x", d3.forceX())
        .force("y", d3.forceY());

    // Create the SVG container.
    var svg = d3.select("#my_dataviz")
        .append("svg")
            .attr("width", "100%")
            .attr("height", "100%")
            .attr("viewBox", [-width / 2, -height / 2, width, height])
            .attr("style", "max-width: 100%; height: auto;");

    // Arrowhead marker definition. 
    svg
        .append('defs')
        .append('marker')
        .attr('id', 'arrow')
        .attr('viewBox', [-arrowH, 0, arrowH, arrowH])
        .attr('refX', arrowH / 2)
        .attr('refY', arrowH / 2)
        .attr('markerWidth', arrowH)
        .attr('markerHeight', arrowH)
        .attr('orient', 'auto-start-reverse')
        .append('path')
        .attr('d', d3.line()([[-arrowH, 0], [-arrowH, arrowH], [0, arrowH / 2]]))
        .attr('stroke', 'rgb(140, 146, 156)')
        .attr('fill', 'rgb(140, 146, 156)');

    // Add a line for each link, and a circle for each node.
    const link = svg.append("g")
        .attr("id", "graph-link")
        .attr("stroke", "rgb(140, 146, 156)")
        .attr("stroke-opacity", 0.6)
        .selectAll("line")
        .data(links)
        .join("line")
        .attr("stroke-width", d => Math.sqrt(d.value))
        .attr("stroke-dasharray", d => d.dashed ? "4" : null)
        .attr('marker-start', d => d.dashed ? 'url(#arrow)' : '');

    const nodeLabel = svg.append("g")
        .attr("id", "graph-node-label")
        .attr("class", "pointer-events-none")
        .selectAll("text")
        .data(nodes)
        .join("text")
        .attr("dy", nodeRadius * 2.25)
        .attr("fill", "rgb(204, 209, 218)")
        .attr("text-anchor", "middle")
        .text(d => d.label);

    const node = svg.append("g")
        .attr("id", "graph-node")
        .selectAll("circle")
        .data(nodes)
        .join("circle")
        .attr("r", (d) => d.type == 'role' ? nodeRadius : nodeRadius * 2)
        .attr("fill", (d) => d.type == 'role' ? 'rgb(229, 113, 208)' : 'rgb(84, 154, 246)' );

    // Add a drag behavior.
    node.call(d3.drag()
            .on("start", dragstarted)
            .on("drag", dragged)
            .on("end", dragended));
    
    // Set the position attributes of links and nodes each time the simulation ticks.
    simulation.on("tick", () => {
        link
            .attr("x1", d => d.source.x)
            .attr("y1", d => d.source.y)
            .attr("x2", d => d.target.x)
            .attr("y2", d => d.target.y);

        node
            .attr("cx", d => d.x)
            .attr("cy", d => d.y);
        nodeLabel
            .attr("x", d => d.x)
            .attr("y", d => d.y + nodeRadius * 2.25);
    });

    // Reheat the simulation when drag starts, and fix the subject position.
    function dragstarted(event) {
        if (!event.active) {
            simulation.alphaTarget(0.3).restart();
        }
        event.subject.fx = event.subject.x;
        event.subject.fy = event.subject.y;
    }

    // Update the subject (dragged node) position during drag.
    function dragged(event) {
        event.subject.fx = event.x;
        event.subject.fy = event.y;
    }

    // Restore the target alpha so the simulation cools after dragging ends.
    // Unfix the subject position now that it’s no longer being dragged.
    function dragended(event) {
        if (!event.active) {
            simulation.alphaTarget(0);
        }
        event.subject.fx = null;
        event.subject.fy = null;
    }

    function handleZoom(e) {
        d3.select('#graph-link').attr('transform', e.transform);
        d3.select('#graph-node').attr('transform', e.transform);
        d3.select('#graph-node-label').attr('transform', e.transform);
    }
    let zoom = d3.zoom().on('zoom', (handleZoom));
    d3.select('svg').call(zoom);
});


</script>

<template>
<div id="my_dataviz"></div>
</template>

<style scoped>

#my_dataviz {
    background: rgb(30, 33, 42);
    background-image: url("bg.png");
    background-repeat: repeat;
}

.pointer-events-none text {
    pointer-events: none;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
    cursor: default;
}

</style>