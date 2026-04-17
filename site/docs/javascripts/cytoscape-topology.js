// Kiali-style interactive network topology using Cytoscape.js
// Renders force-directed graph from JSON data embedded in the page.
(function () {
  var COLORS = {
    component: { bg: "#3498db", border: "#2980b9", text: "#fff" },
    service: { bg: "#2ecc71", border: "#27ae60", text: "#fff" },
    external: { bg: "#e74c3c", border: "#c0392b", text: "#fff" },
  };

  var EDGE_COLORS = {
    watches: "#e74c3c",
    sidecar: "#9b59b6",
    module: "#95a5a6",
    external: "#e67e22",
  };

  function initGraph(wrapper) {
    var dataEl = wrapper.querySelector("script[type='application/json']");
    if (!dataEl) return;

    var data;
    try {
      data = JSON.parse(dataEl.textContent);
    } catch (e) {
      console.error("Cytoscape topology: invalid JSON", e);
      return;
    }

    // Remove the script tag and create a clean render container
    var container = document.createElement("div");
    container.style.cssText = "width:100%;height:100%;";
    wrapper.innerHTML = "";
    wrapper.appendChild(container);

    var elements = [];

    // Build a set of known node IDs to filter edges
    var knownNodes = {};

    // Add component nodes (parents)
    data.components.forEach(function (comp) {
      knownNodes[comp.id] = true;
      elements.push({
        data: {
          id: comp.id,
          label: comp.name,
          nodeType: "component",
          serviceCount: comp.serviceCount || 0,
          netpolCount: comp.netpolCount || 0,
          hasIngress: comp.hasIngress || false,
        },
      });
    });

    // Add service nodes as children of components
    (data.services || []).forEach(function (svc) {
      knownNodes[svc.id] = true;
      elements.push({
        data: {
          id: svc.id,
          label: svc.name + (svc.ports ? "\n:" + svc.ports : ""),
          nodeType: "service",
          parent: svc.parent,
          ports: svc.ports || "",
        },
      });
    });

    // Add external nodes
    (data.externals || []).forEach(function (ext) {
      knownNodes[ext.id] = true;
      elements.push({
        data: {
          id: ext.id,
          label: ext.name + "\n" + ext.type,
          nodeType: "external",
        },
      });
    });

    // Add edges (only if both endpoints exist)
    (data.edges || []).forEach(function (edge) {
      if (!knownNodes[edge.from] || !knownNodes[edge.to]) return;
      elements.push({
        data: {
          id: edge.from + "-" + edge.to + "-" + edge.type,
          source: edge.from,
          target: edge.to,
          edgeType: edge.type,
          label: edge.type,
        },
      });
    });

    // Detect dark mode
    var isDark =
      document.body.getAttribute("data-md-color-scheme") === "slate";
    var bgColor = isDark ? "#1e1e1e" : "#fafafa";
    var parentBg = isDark ? "rgba(255,255,255,0.05)" : "rgba(0,0,0,0.03)";

    var cy = cytoscape({
      container: container,
      elements: elements,
      minZoom: 0.3,
      maxZoom: 3,
      wheelSensitivity: 0.3,

      style: [
        // Compound (parent) nodes: components that have services
        {
          selector: 'node[nodeType="component"]',
          style: {
            shape: "round-rectangle",
            "background-color": parentBg,
            "border-width": 2,
            "border-color": COLORS.component.bg,
            "border-opacity": 0.7,
            label: "data(label)",
            "text-valign": "top",
            "text-halign": "center",
            "font-size": "12px",
            "font-weight": "bold",
            color: COLORS.component.bg,
            "text-margin-y": -6,
            padding: "16px",
            "min-width": "60px",
            "min-height": "40px",
          },
        },
        // Leaf component nodes (no services): filled circles
        {
          selector: "node[nodeType='component']:childless",
          style: {
            "background-color": COLORS.component.bg,
            "border-color": COLORS.component.border,
            color: COLORS.component.text,
            "text-valign": "center",
            "text-margin-y": 0,
            padding: "10px",
            shape: "ellipse",
            width: "label",
            height: "label",
          },
        },
        // Service nodes
        {
          selector: 'node[nodeType="service"]',
          style: {
            shape: "round-rectangle",
            "background-color": COLORS.service.bg,
            "border-width": 1,
            "border-color": COLORS.service.border,
            label: "data(label)",
            "text-valign": "center",
            "text-halign": "center",
            "font-size": "9px",
            color: COLORS.service.text,
            width: "label",
            height: "label",
            padding: "6px",
            "text-wrap": "wrap",
            "text-max-width": "120px",
          },
        },
        // External nodes
        {
          selector: 'node[nodeType="external"]',
          style: {
            shape: "diamond",
            "background-color": COLORS.external.bg,
            "border-width": 2,
            "border-color": COLORS.external.border,
            label: "data(label)",
            "text-valign": "center",
            "text-halign": "center",
            "font-size": "10px",
            color: COLORS.external.text,
            width: 50,
            height: 50,
            "text-wrap": "wrap",
          },
        },
        // Edges base
        {
          selector: "edge",
          style: {
            width: 2,
            "line-color": "#bdc3c7",
            "target-arrow-color": "#bdc3c7",
            "target-arrow-shape": "triangle",
            "curve-style": "bezier",
            "arrow-scale": 0.8,
            label: "data(label)",
            "font-size": "8px",
            color: isDark ? "#aaa" : "#666",
            "text-rotation": "autorotate",
            "text-margin-y": -8,
            "text-background-color": bgColor,
            "text-background-opacity": 0.8,
            "text-background-padding": "2px",
          },
        },
        {
          selector: 'edge[edgeType="watches"]',
          style: {
            width: 3,
            "line-color": EDGE_COLORS.watches,
            "target-arrow-color": EDGE_COLORS.watches,
          },
        },
        {
          selector: 'edge[edgeType="sidecar"]',
          style: {
            width: 2,
            "line-color": EDGE_COLORS.sidecar,
            "target-arrow-color": EDGE_COLORS.sidecar,
          },
        },
        {
          selector: 'edge[edgeType="module"]',
          style: {
            width: 1.5,
            "line-color": EDGE_COLORS.module,
            "target-arrow-color": EDGE_COLORS.module,
            "line-style": "dashed",
          },
        },
        {
          selector: 'edge[edgeType="external"]',
          style: {
            width: 1.5,
            "line-color": EDGE_COLORS.external,
            "target-arrow-color": EDGE_COLORS.external,
            "line-style": "dotted",
          },
        },
        // Selection highlight
        {
          selector: "node:active, node:selected",
          style: {
            "border-width": 3,
            "border-color": "#f1c40f",
            "overlay-opacity": 0,
          },
        },
        {
          selector: "edge:active, edge:selected",
          style: { width: 4, "overlay-opacity": 0 },
        },
      ],

      layout: {
        name: "cose",
        animate: true,
        animationDuration: 800,
        nodeRepulsion: function () {
          return 8000;
        },
        idealEdgeLength: function () {
          return 120;
        },
        edgeElasticity: function () {
          return 100;
        },
        nestingFactor: 1.2,
        gravity: 0.25,
        numIter: 1000,
        randomize: false,
        componentSpacing: 80,
        nodeDimensionsIncludeLabels: true,
        padding: 30,
      },
    });

    // Tooltip on tap
    var tooltip = document.createElement("div");
    tooltip.style.cssText =
      "position:absolute;display:none;background:" +
      (isDark ? "#333" : "#fff") +
      ";color:" +
      (isDark ? "#eee" : "#333") +
      ";border:1px solid " +
      (isDark ? "#555" : "#ddd") +
      ";border-radius:6px;padding:8px 12px;font-size:12px;" +
      "pointer-events:none;z-index:100;max-width:300px;box-shadow:0 2px 8px rgba(0,0,0,0.15);";
    container.style.position = "relative";
    container.appendChild(tooltip);

    cy.on("tap", "node", function (evt) {
      var node = evt.target;
      var d = node.data();
      var html = "<strong>" + d.label.replace(/\n/g, " ") + "</strong>";
      if (d.nodeType === "component") {
        html += "<br>Services: " + d.serviceCount;
        if (d.netpolCount > 0) html += "<br>NetworkPolicies: " + d.netpolCount;
        if (d.hasIngress) html += "<br>Has ingress routing";
        var edges = cy.edges().filter(function (e) {
          return e.data("source") === d.id || e.data("target") === d.id;
        });
        if (edges.length > 0) html += "<br>Connections: " + edges.length;
      } else if (d.nodeType === "service" && d.ports) {
        html += "<br>Ports: " + d.ports;
      }
      tooltip.innerHTML = html;
      tooltip.style.display = "block";
      var pos = evt.renderedPosition;
      tooltip.style.left = pos.x + 15 + "px";
      tooltip.style.top = pos.y + 15 + "px";
    });

    cy.on("tap", function (evt) {
      if (evt.target === cy) tooltip.style.display = "none";
    });

    // Highlight connected edges on node hover
    cy.on("mouseover", "node", function (evt) {
      var node = evt.target;
      var connected = node.connectedEdges();
      cy.edges().not(connected).style("opacity", 0.15);
      connected.style("opacity", 1);
      var connectedNodes = connected.connectedNodes();
      cy.nodes().not(connectedNodes).not(node).style("opacity", 0.3);
    });

    cy.on("mouseout", "node", function () {
      cy.elements().style("opacity", 1);
      tooltip.style.display = "none";
    });

    // Fit on double-click
    cy.on("dbltap", function (evt) {
      if (evt.target === cy) {
        cy.animate({ fit: { padding: 30 }, duration: 400 });
      }
    });

    // Store reference for toolbar
    wrapper._cy = cy;
  }

  function initAll() {
    document.querySelectorAll(".cytoscape-topology").forEach(function (el) {
      if (el._cy) return;
      initGraph(el);
    });
  }

  // Toolbar click handler
  document.addEventListener(
    "click",
    function (e) {
      var btn = e.target.closest(".topology-toolbar button");
      if (!btn) return;
      e.preventDefault();
      var toolbar = btn.closest(".topology-toolbar");
      var container = toolbar.nextElementSibling;
      if (!container || !container._cy) return;
      var cy = container._cy;
      var action = btn.dataset.action;

      if (action === "fit")
        cy.animate({ fit: { padding: 30 }, duration: 400 });
      if (action === "zoom-in")
        cy.animate({ zoom: cy.zoom() * 1.3, duration: 200 });
      if (action === "zoom-out")
        cy.animate({ zoom: cy.zoom() / 1.3, duration: 200 });
      if (action === "relayout") {
        cy.layout({
          name: "cose",
          animate: true,
          animationDuration: 800,
          nodeRepulsion: function () {
            return 8000;
          },
          idealEdgeLength: function () {
            return 120;
          },
          nestingFactor: 1.2,
          gravity: 0.25,
          numIter: 1000,
          randomize: true,
          componentSpacing: 80,
          nodeDimensionsIncludeLabels: true,
        }).run();
      }
    },
    true
  );

  // Wait for Cytoscape library to load, then init
  function waitAndInit() {
    if (typeof cytoscape !== "undefined") {
      initAll();
    } else {
      setTimeout(waitAndInit, 200);
    }
  }

  // Start initialization
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", waitAndInit);
  } else {
    waitAndInit();
  }

  // mkdocs-material instant navigation
  if (typeof document$ !== "undefined") {
    document$.subscribe(function () {
      setTimeout(waitAndInit, 300);
    });
  }
})();
