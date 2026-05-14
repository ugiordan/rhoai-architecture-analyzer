// Kiali-style interactive network topology using Cytoscape.js
// Components shown as single nodes (not expanded subgraphs).
// Services visible on click as tooltip detail.
(function () {
  function initGraph(wrapper) {
    var dataEl = wrapper.querySelector("script[type='application/json']");
    if (!dataEl) return;

    var data;
    try {
      data = JSON.parse(dataEl.textContent);
    } catch (e) {
      return;
    }

    // Build service lookup for tooltips
    var servicesByComp = {};
    (data.services || []).forEach(function (svc) {
      if (!servicesByComp[svc.parent]) servicesByComp[svc.parent] = [];
      servicesByComp[svc.parent].push(svc);
    });

    // Create clean render container
    var container = document.createElement("div");
    container.style.cssText = "width:100%;height:100%;";
    wrapper.innerHTML = "";
    wrapper.appendChild(container);

    var elements = [];
    var knownNodes = {};

    // Components as single nodes (NOT compound parents)
    data.components.forEach(function (comp) {
      knownNodes[comp.id] = true;
      var svcCount = comp.serviceCount || 0;
      var label = comp.name;
      if (svcCount > 0) label += "\n(" + svcCount + " svc)";
      elements.push({
        data: {
          id: comp.id,
          label: label,
          shortLabel: comp.name,
          nodeType: "component",
          serviceCount: svcCount,
          netpolCount: comp.netpolCount || 0,
          hasIngress: comp.hasIngress || false,
        },
      });
    });

    // External nodes
    (data.externals || []).forEach(function (ext) {
      knownNodes[ext.id] = true;
      elements.push({
        data: {
          id: ext.id,
          label: ext.name,
          nodeType: "external",
          extType: ext.type,
        },
      });
    });

    // Edges (skip dangling)
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

    var isDark =
      document.body.getAttribute("data-md-color-scheme") === "slate";
    var bgColor = isDark ? "#1e1e1e" : "#fafafa";

    var cy = cytoscape({
      container: container,
      elements: elements,
      minZoom: 0.3,
      maxZoom: 3,
      wheelSensitivity: 0.3,

      style: [
        // Component nodes: filled circles like Kiali
        {
          selector: 'node[nodeType="component"]',
          style: {
            shape: "ellipse",
            width: 70,
            height: 70,
            "background-color": "#3498db",
            "border-width": 3,
            "border-color": "#2980b9",
            label: "data(label)",
            "text-valign": "bottom",
            "text-halign": "center",
            "font-size": "11px",
            "font-weight": "600",
            color: isDark ? "#ccc" : "#333",
            "text-margin-y": 8,
            "text-wrap": "wrap",
            "text-max-width": "140px",
            "text-background-color": bgColor,
            "text-background-opacity": 0.7,
            "text-background-padding": "2px",
            "text-background-shape": "round-rectangle",
          },
        },
        // Components with network policies: orange ring
        {
          selector: "node[netpolCount > 0]",
          style: {
            "border-color": "#f39c12",
            "border-width": 4,
          },
        },
        // Components with ingress: green fill
        {
          selector: "node[hasIngress]",
          style: {
            "background-color": "#27ae60",
            "border-color": "#1e8449",
          },
        },
        // Components with both netpol + ingress: green fill, orange ring
        {
          selector: "node[hasIngress][netpolCount > 0]",
          style: {
            "background-color": "#27ae60",
            "border-color": "#f39c12",
            "border-width": 4,
          },
        },
        // External nodes: smaller diamonds
        {
          selector: 'node[nodeType="external"]',
          style: {
            shape: "diamond",
            width: 45,
            height: 45,
            "background-color": "#e74c3c",
            "border-width": 2,
            "border-color": "#c0392b",
            label: "data(label)",
            "text-valign": "bottom",
            "text-halign": "center",
            "font-size": "10px",
            color: isDark ? "#ccc" : "#555",
            "text-margin-y": 6,
            "text-wrap": "wrap",
            "text-max-width": "100px",
            "text-background-color": bgColor,
            "text-background-opacity": 0.7,
            "text-background-padding": "2px",
            "text-background-shape": "round-rectangle",
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
            "font-size": "9px",
            color: isDark ? "#999" : "#777",
            "text-rotation": "autorotate",
            "text-margin-y": -10,
            "text-background-color": bgColor,
            "text-background-opacity": 0.85,
            "text-background-padding": "3px",
            "text-background-shape": "round-rectangle",
          },
        },
        {
          selector: 'edge[edgeType="watches"]',
          style: {
            width: 3,
            "line-color": "#e74c3c",
            "target-arrow-color": "#e74c3c",
          },
        },
        {
          selector: 'edge[edgeType="sidecar"]',
          style: {
            width: 2.5,
            "line-color": "#9b59b6",
            "target-arrow-color": "#9b59b6",
          },
        },
        {
          selector: 'edge[edgeType="module"]',
          style: {
            width: 1.5,
            "line-color": "#95a5a6",
            "target-arrow-color": "#95a5a6",
            "line-style": "dashed",
          },
        },
        {
          selector: 'edge[edgeType="external"]',
          style: {
            width: 1.5,
            "line-color": "#e67e22",
            "target-arrow-color": "#e67e22",
            "line-style": "dotted",
          },
        },
        // Hover/select
        {
          selector: "node:selected",
          style: {
            "border-width": 5,
            "border-color": "#f1c40f",
            "overlay-opacity": 0,
          },
        },
        // Dimmed state for non-hovered
        {
          selector: ".dimmed",
          style: { opacity: 0.15 },
        },
        {
          selector: ".highlighted",
          style: { opacity: 1 },
        },
      ],

      layout: {
        name: "cose",
        animate: true,
        animationDuration: 1000,
        animationEasing: "ease-out",
        nodeRepulsion: function () {
          return 12000;
        },
        idealEdgeLength: function () {
          return 180;
        },
        edgeElasticity: function () {
          return 100;
        },
        gravity: 0.4,
        numIter: 1500,
        randomize: false,
        componentSpacing: 100,
        nodeDimensionsIncludeLabels: true,
        padding: 40,
      },
    });

    // Tooltip
    var tooltip = document.createElement("div");
    tooltip.className = "topology-tooltip";
    tooltip.style.cssText =
      "position:absolute;display:none;background:" +
      (isDark ? "#2d2d2d" : "#fff") +
      ";color:" +
      (isDark ? "#eee" : "#333") +
      ";border:1px solid " +
      (isDark ? "#555" : "#ccc") +
      ";border-radius:8px;padding:10px 14px;font-size:12px;line-height:1.5;" +
      "pointer-events:none;z-index:100;max-width:350px;" +
      "box-shadow:0 4px 12px rgba(0,0,0,0.2);";
    container.style.position = "relative";
    container.appendChild(tooltip);

    cy.on("tap", "node", function (evt) {
      var node = evt.target;
      var d = node.data();
      var html = '<strong style="font-size:13px;">' + d.shortLabel + "</strong>";

      if (d.nodeType === "component") {
        html += "<br><br>";
        var badges = [];
        if (d.serviceCount > 0) badges.push(d.serviceCount + " services");
        if (d.netpolCount > 0)
          badges.push(d.netpolCount + " network policies");
        if (d.hasIngress) badges.push("ingress routing");
        html += badges.join(" &middot; ");

        // List services
        var svcs = servicesByComp[d.id];
        if (svcs && svcs.length > 0) {
          html +=
            '<div style="margin-top:6px;padding-top:6px;border-top:1px solid ' +
            (isDark ? "#444" : "#eee") +
            ';font-size:11px;">';
          svcs.forEach(function (svc) {
            html +=
              '<div style="padding:1px 0;"><span style="color:#2ecc71;font-weight:600;">&#9679;</span> ' +
              svc.name;
            if (svc.ports) html += ' <span style="opacity:0.6;">:' + svc.ports + "</span>";
            html += "</div>";
          });
          html += "</div>";
        }

        // List connections
        var conns = cy.edges().filter(function (e) {
          return e.data("source") === d.id || e.data("target") === d.id;
        });
        if (conns.length > 0) {
          html +=
            '<div style="margin-top:6px;padding-top:6px;border-top:1px solid ' +
            (isDark ? "#444" : "#eee") +
            ';font-size:11px;">';
          conns.forEach(function (e) {
            var other =
              e.data("source") === d.id ? e.data("target") : e.data("source");
            var dir = e.data("source") === d.id ? "&rarr;" : "&larr;";
            var typeColors = {
              watches: "#e74c3c",
              sidecar: "#9b59b6",
              module: "#95a5a6",
              external: "#e67e22",
            };
            var c = typeColors[e.data("edgeType")] || "#999";
            html +=
              '<div style="padding:1px 0;"><span style="color:' +
              c +
              ';">' +
              dir +
              "</span> " +
              other.replace(/_/g, "-") +
              ' <span style="opacity:0.5;">(' +
              e.data("edgeType") +
              ")</span></div>";
          });
          html += "</div>";
        }
      } else if (d.nodeType === "external") {
        html += "<br>Type: " + d.extType;
      }

      tooltip.innerHTML = html;
      tooltip.style.display = "block";
      var pos = evt.renderedPosition;
      var tLeft = pos.x + 15;
      var tTop = pos.y + 15;
      // Keep tooltip in bounds
      if (tLeft + 350 > container.clientWidth) tLeft = pos.x - 360;
      if (tTop + 200 > container.clientHeight) tTop = pos.y - 200;
      tooltip.style.left = Math.max(0, tLeft) + "px";
      tooltip.style.top = Math.max(0, tTop) + "px";
    });

    cy.on("tap", function (evt) {
      if (evt.target === cy) tooltip.style.display = "none";
    });

    // Hover: highlight connected, dim rest
    cy.on("mouseover", "node", function (evt) {
      var node = evt.target;
      var neighborhood = node.closedNeighborhood();
      cy.elements().addClass("dimmed");
      neighborhood.removeClass("dimmed").addClass("highlighted");
    });

    cy.on("mouseout", "node", function () {
      cy.elements().removeClass("dimmed").removeClass("highlighted");
      tooltip.style.display = "none";
    });

    // Fit on double-click background
    cy.on("dbltap", function (evt) {
      if (evt.target === cy) {
        cy.animate({ fit: { padding: 40 }, duration: 400 });
      }
    });

    wrapper._cy = cy;
  }

  function initAll() {
    document.querySelectorAll(".cytoscape-topology").forEach(function (el) {
      if (el._cy) return;
      initGraph(el);
    });
  }

  // Toolbar handler
  document.addEventListener(
    "click",
    function (e) {
      var btn = e.target.closest(".topology-toolbar button");
      if (!btn) return;
      e.preventDefault();
      var toolbar = btn.closest(".topology-toolbar");
      var graphEl = toolbar.nextElementSibling;
      if (!graphEl) return;
      var action = btn.dataset.action;

      // Fullscreen toggle
      if (action === "fullscreen") {
        if (graphEl.classList.contains("topology-fullscreen")) {
          graphEl.classList.remove("topology-fullscreen");
          toolbar.classList.remove("topology-toolbar-fullscreen");
          var legend = graphEl.nextElementSibling;
          if (legend && legend.classList.contains("topology-legend"))
            legend.classList.remove("topology-legend-fullscreen");
          document.body.style.overflow = "";
          btn.textContent = "Fullscreen";
        } else {
          graphEl.classList.add("topology-fullscreen");
          toolbar.classList.add("topology-toolbar-fullscreen");
          var legend = graphEl.nextElementSibling;
          if (legend && legend.classList.contains("topology-legend"))
            legend.classList.add("topology-legend-fullscreen");
          document.body.style.overflow = "hidden";
          btn.textContent = "Exit";
        }
        if (graphEl._cy)
          setTimeout(function () {
            graphEl._cy.resize();
            graphEl._cy.fit(null, 40);
          }, 100);
        return;
      }

      if (!graphEl._cy) return;
      var cy = graphEl._cy;
      if (action === "fit") cy.animate({ fit: { padding: 40 }, duration: 400 });
      if (action === "zoom-in")
        cy.animate({ zoom: cy.zoom() * 1.4, duration: 200 });
      if (action === "zoom-out")
        cy.animate({ zoom: cy.zoom() / 1.4, duration: 200 });
      if (action === "relayout") {
        cy.layout({
          name: "cose",
          animate: true,
          animationDuration: 800,
          nodeRepulsion: function () {
            return 12000;
          },
          idealEdgeLength: function () {
            return 180;
          },
          gravity: 0.4,
          numIter: 1500,
          randomize: true,
          componentSpacing: 100,
          nodeDimensionsIncludeLabels: true,
        }).run();
      }
    },
    true
  );

  // ESC to exit fullscreen
  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape") {
      var fs = document.querySelector(".topology-fullscreen");
      if (fs) {
        fs.classList.remove("topology-fullscreen");
        var toolbar = fs.previousElementSibling;
        if (toolbar) toolbar.classList.remove("topology-toolbar-fullscreen");
        var legend = fs.nextElementSibling;
        if (legend) legend.classList.remove("topology-legend-fullscreen");
        document.body.style.overflow = "";
        var btn = toolbar && toolbar.querySelector('[data-action="fullscreen"]');
        if (btn) btn.textContent = "Fullscreen";
        if (fs._cy) {
          fs._cy.resize();
          fs._cy.fit(null, 40);
        }
      }
    }
  });

  function waitAndInit() {
    if (typeof cytoscape !== "undefined") {
      initAll();
    } else {
      setTimeout(waitAndInit, 200);
    }
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", waitAndInit);
  } else {
    waitAndInit();
  }

  if (typeof document$ !== "undefined") {
    document$.subscribe(function () {
      setTimeout(waitAndInit, 300);
    });
  }
})();
