// Diagram controls: zoom in/out, fullscreen overlay
// Strategy: don't wrap or modify mermaid elements before rendering.
// Wait for SVGs to appear, then insert a sibling action bar AFTER each diagram.
// Fullscreen uses a modal overlay (no CSS position:fixed on the pre).
(function () {
  var ZOOM_STEP = 0.25;
  var MIN_SCALE = 0.5;
  var MAX_SCALE = 5;

  function addControls(pre) {
    if (pre.nextElementSibling && pre.nextElementSibling.classList.contains("diagram-actions")) return;
    if (!pre.querySelector("svg")) return;

    var bar = document.createElement("div");
    bar.className = "diagram-actions";
    bar.innerHTML =
      '<button data-action="zoom-in" title="Zoom in">&#43;</button>' +
      '<button data-action="zoom-out" title="Zoom out">&#8722;</button>' +
      '<button data-action="reset" title="Reset zoom">1:1</button>' +
      '<button data-action="fullscreen" title="View full screen">&#9974;</button>';
    pre.after(bar);
    pre._zoomScale = 1;
  }

  function scanAll() {
    document.querySelectorAll("pre.mermaid").forEach(addControls);
  }

  function openFullscreen(svg) {
    var overlay = document.createElement("div");
    overlay.className = "diagram-overlay";

    var toolbar = document.createElement("div");
    toolbar.className = "diagram-overlay-toolbar";
    toolbar.innerHTML =
      '<button data-action="ol-zoom-in" title="Zoom in">&#43;</button>' +
      '<button data-action="ol-zoom-out" title="Zoom out">&#8722;</button>' +
      '<button data-action="ol-reset" title="Reset">1:1</button>' +
      '<button data-action="ol-close" title="Close">&#10005;</button>';
    overlay.appendChild(toolbar);

    var container = document.createElement("div");
    container.className = "diagram-overlay-content";
    var svgClone = svg.cloneNode(true);
    svgClone.removeAttribute("style");
    svgClone.style.maxWidth = "none";
    svgClone.style.maxHeight = "none";
    svgClone.style.width = "auto";
    svgClone.style.height = "auto";
    container.appendChild(svgClone);
    overlay.appendChild(container);

    document.body.appendChild(overlay);
    overlay._scale = 1;

    function close() {
      overlay.remove();
      document.removeEventListener("keydown", escHandler);
    }
    function escHandler(e) { if (e.key === "Escape") close(); }
    document.addEventListener("keydown", escHandler);
    overlay.addEventListener("click", function (e) {
      if (e.target === overlay) close();
    });

    toolbar.addEventListener("click", function (e) {
      var btn = e.target.closest("button");
      if (!btn) return;
      var action = btn.dataset.action;
      if (action === "ol-close") { close(); return; }
      if (action === "ol-zoom-in") overlay._scale = Math.min(overlay._scale + ZOOM_STEP, MAX_SCALE);
      if (action === "ol-zoom-out") overlay._scale = Math.max(overlay._scale - ZOOM_STEP, MIN_SCALE);
      if (action === "ol-reset") overlay._scale = 1;
      svgClone.style.transform = overlay._scale === 1 ? "" : "scale(" + overlay._scale + ")";
      svgClone.style.transformOrigin = "top left";
    });
  }

  // Click delegation for inline action bars (capture phase)
  document.addEventListener("click", function (e) {
    var btn = e.target.closest(".diagram-actions button");
    if (!btn) return;
    e.preventDefault();
    e.stopPropagation();
    e.stopImmediatePropagation();

    var bar = btn.closest(".diagram-actions");
    var pre = bar.previousElementSibling;
    if (!pre) return;
    var svg = pre.querySelector("svg");
    var action = btn.dataset.action;

    if (action === "fullscreen") {
      if (svg) openFullscreen(svg);
      return;
    }

    if (!svg) return;
    if (!pre._zoomScale) pre._zoomScale = 1;

    if (action === "zoom-in") pre._zoomScale = Math.min(pre._zoomScale + ZOOM_STEP, MAX_SCALE);
    if (action === "zoom-out") pre._zoomScale = Math.max(pre._zoomScale - ZOOM_STEP, MIN_SCALE);
    if (action === "reset") pre._zoomScale = 1;

    svg.style.transform = pre._zoomScale === 1 ? "" : "scale(" + pre._zoomScale + ")";
    svg.style.transformOrigin = "top left";
  }, true);

  // Poll for rendered SVGs (mermaid renders asynchronously)
  var pollCount = 0;
  var pollId = setInterval(function () {
    scanAll();
    if (++pollCount > 30) clearInterval(pollId);
  }, 1000);

  // Also catch dynamic renders via MutationObserver (throttled)
  var raf = false;
  new MutationObserver(function () {
    if (raf) return;
    raf = true;
    requestAnimationFrame(function () {
      raf = false;
      scanAll();
    });
  }).observe(document.body, { childList: true, subtree: true });

  // mkdocs-material instant navigation
  if (typeof document$ !== "undefined") {
    document$.subscribe(function () {
      pollCount = 0;
      pollId = setInterval(function () {
        scanAll();
        if (++pollCount > 30) clearInterval(pollId);
      }, 1000);
    });
  }
})();
