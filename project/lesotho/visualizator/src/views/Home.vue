<script setup>
import * as d3 from 'd3';
import { onMounted, ref } from 'vue';

let Svg = ref();
const inputFile = ref(null);

const Nodes = ref([]);
const Links = ref([]);
const IdCounter = ref(0);

const nodeRadius = 15;
const nodeLabelSize = "10pt";
const nodeIconSize = 1.25; // Relative to node size.
const calcNodeRadius = (d) => d.type == 'role' ? nodeRadius : nodeRadius * 2; 
const arrowH = nodeRadius * 1.25;

const namespaceToGraph = (data) => {
    const namespaceName = data.name;

    let idMap = {};

    let G = {
        nodes: [],
        links: []
    }

    G.nodes.push({
        "id": IdCounter.value,
        "label": namespaceName,
        "type": "namespace"
    });
    idMap[namespaceName] = IdCounter.value;
    IdCounter.value++;

    for (const [key, value] of Object.entries(data.relations)) {
        G.nodes.push({
            "id": IdCounter.value,
            "label": key,
            "type": "role"
        });
        idMap[key] = IdCounter.value;
        IdCounter.value++;
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
            "target": idMap[namespaceName],
            "dashed": false
        });
    }

    return G;
}

const initSvg = () => {
    // Specify the dimensions of the chart.
    const width = 928;
    const height = 520;

    // Create the SVG container.
    var svg = d3.select("#my_dataviz")
        .append("svg")
            .attr("width", "100%")
            .attr("height", "100%")
            .attr("viewBox", [-width / 2, -height / 2, width, height])
            .attr("style", "max-width: 100%; height: 100%;");

    Svg.value = svg;
}

onMounted(() => {
    initSvg();
    createGraphInitial();
});

const createGraphFromNamespaceJson = (namespaceJSON) => {
    let g = JSON.parse(namespaceJSON);
    createGraphFromNamespaceObject(g);
}

const createGraphInitial = () => {
    var svg = Svg.value;
    svg.append("g").attr("id", "graph-links");
    svg.append("g").attr("id", "graph-nodes");
    svg.append("g").attr("id", "graph-node-labels");
    svg.append("g").attr("id", "graph-node-icons");

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
        .attr('d', d3.line()([[-arrowH-5, 0], [-arrowH-5, arrowH], [-5, arrowH / 2]]))
        .attr('stroke', 'rgb(140, 146, 156)')
        .attr('fill', 'rgb(140, 146, 156)');
}

const updateGraph = () => {
    const node = d3.select("#graph-nodes")
        .selectAll("circle")
        .data(Nodes.value)
        .join("circle")
        .attr("r", (d) => calcNodeRadius(d))
        .attr("fill", (d) => d.type == 'role' ? 'rgb(229, 113, 208)' : 'rgb(84, 154, 246)' );
    
    const link = d3.select("#graph-links")
        .selectAll("line")
        .data(Links.value)
        .join("line")
        .attr("stroke", "rgb(140, 146, 156)")
        .attr("stroke-opacity", 0.6)
        .attr("stroke-dasharray", d => d.dashed ? "4" : null)
        .attr('marker-start', d => d.dashed ? 'url(#arrow)' : '');

    const node_label = d3.select("#graph-node-labels")
        .selectAll("text")
        .data(Nodes.value)
        .join("text")
        .attr("dy", nodeRadius * 1.35)
        .attr("fill", "rgb(204, 209, 218)")
        .attr("text-anchor", "middle")
        .attr("font-family", "monospace")
        .attr("font-size", nodeLabelSize)
        .text(d => d.label);

    const node_icon = d3.select("#graph-node-icons")
        .selectAll("image")
        .data(Nodes.value)
        .join("image")
        .attr("href", (d) => d.type == 'role' ? "public/ico_user.png" : "public/ico_file.png")
        .attr("width", (d) => calcNodeRadius(d) * nodeIconSize)
        .attr("height", (d) => calcNodeRadius(d) * nodeIconSize)
        .attr("x", (d) => -calcNodeRadius(d) * nodeIconSize / 2)
        .attr("y", (d) => -calcNodeRadius(d) * nodeIconSize / 2);

    //-------------------------------------------------------------------------
    // Interact
    //-------------------------------------------------------------------------

    node.call(d3.drag()
        .on("start", dragstarted)
        .on("drag", dragged)
        .on("end", dragended));
    node_icon.call(d3.drag()
        .on("start", dragstarted)
        .on("drag", dragged)
        .on("end", dragended));

    function dragstarted(event) {
        if (!event.active) {
            simulation.alphaTarget(0.2).restart();
        }
        event.subject.fx = event.subject.x;
        event.subject.fy = event.subject.y;
    }

    function dragged(event) {
        event.subject.fx = event.x;
        event.subject.fy = event.y;
    }

    function dragended(event) {
        if (!event.active) {
            simulation.alphaTarget(0);
        }
        event.subject.fx = null;
        event.subject.fy = null;
    }

    function handleZoom(e) {
        d3.select('#graph-links').attr('transform', e.transform);
        d3.select('#graph-nodes').attr('transform', e.transform);
        d3.select('#graph-node-labels').attr('transform', e.transform);
        d3.select('#graph-node-icons').attr('transform', e.transform);
    }
    let zoom = d3.zoom().on('zoom', (handleZoom));
    d3.select('svg').call(zoom);

    //-------------------------------------------------------------------------
    // Force layout
    //-------------------------------------------------------------------------

    const simulation = d3.forceSimulation(Nodes.value)
        .force("link", d3.forceLink(Links.value).id(d => d.id).distance(100))
        .force("charge", d3.forceManyBody().strength(-1500))
        .force("x", d3.forceX())
        .force("y", d3.forceY());

    console.log(Nodes.value)

    simulation.on("tick", () => {
        link
            .attr("x1", d => d.source.x)
            .attr("y1", d => d.source.y)
            .attr("x2", d => d.target.x)
            .attr("y2", d => d.target.y);

        node
            .attr("cx", d => d.x)
            .attr("cy", d => d.y);

        node_label
            .attr("x", d => d.x)
            .attr("y", d => d.y + nodeRadius * 1.35);

        node_icon
            .attr("x", d => d.x - calcNodeRadius(d) * nodeIconSize / 2)
            .attr("y", d => d.y - calcNodeRadius(d) * nodeIconSize / 2);
    });
}

const createGraphFromNamespaceObject = (namespace) => {
    const G = namespaceToGraph(namespace);

    for (let node of G.nodes) {
        Nodes.value.push(node);
    }
    for (let link of G.links) {
        Links.value.push(link);
    }

    updateGraph();
};

const onUploadFile = (e) => {
    var reader = new FileReader();
    reader.readAsText(e.target.files[0], "UTF-8");
    reader.onload = function (evt) {
        createGraphFromNamespaceJson(evt.target.result);
        inputFile.value = null;
    }
    reader.onerror = function (evt) {
        console.error(evt);
        inputFile.value = null;
    }
}

</script>

<template>
<div id="body">
    <h1 class="font">Lesotho visualizator</h1>
    <h3 class="font2">Upload a namespace file</h3>
    <form method="post" enctype="multipart/form-data" >
        <input type="file" name="file" accept=".json" @change="(e) => onUploadFile(e)" ref="inputFile">
    </form>
    <div id="my_dataviz"></div>
</div>
</template>

<style scoped>

.font {
    font-family: monospace;
    font-size: 4em;
    color:lightgray;
    font-weight: 100;
    margin: 0;
    padding: 0;
}

.font2 {
    font-family: monospace;
    font-size: 2em;
    color:gray;
    font-weight: 100;
    margin: 0;
    padding: 0;
}

#body {
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