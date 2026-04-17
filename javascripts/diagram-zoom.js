// Diagram controls for mkdocs-material mermaid diagrams.
//
// mkdocs-material replaces <pre class="mermaid"> with <div class="mermaid">
// containing a CLOSED shadow DOM with the SVG inside. We cannot access the
// SVG directly, so we operate on the outer <div> element for zoom/fullscreen.
(function () {
  var ZOOM_STEP = 0.25;
  var MIN_SCALE = 0.5;
  var MAX_SCALE = 5;

  function addControls(div) {
    // Skip if already has controls
    if (div.nextElementSibling && div.nextElementSibling.classList.contains("diagram-actions")) return;
    // Only target rendered divs (not unprocessed pre elements)
    if (div.tagName !== "DIV") return;

    var bar = document.createElement("div");
    bar.className = "diagram-actions";
    bar.innerHTML =
      '<button data-action="zoom-in" title="Zoom in">&#43;</button>' +
      '<button data-action="zoom-out" title="Zoom out">&#8722;</button>' +
      '<button data-action="reset" title="Reset zoom">1:1</button>' +
      '<button data-action="fullscreen" title="View full screen">&#9974;</button>';
    div.after(bar);
    div._zoomScale = 1;
  }

  function scanAll() {
    // After mkdocs-material renders, the element is <div class="mermaid">
    document.querySelectorAll("div.mermaid").forEach(addControls);
  }

  function openFullscreen(mermaidDiv) {
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
    // Clone the entire div (the closed shadow DOM renders visually within it)
    var clone = mermaidDiv.cloneNode(true);
    clone.style.cssText = "max-width:none; max-height:none; width:auto; height:auto; overflow:visible;";
    container.appendChild(clone);
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
      clone.style.transform = overlay._scale === 1 ? "" : "scale(" + overlay._scale + ")";
      clone.style.transformOrigin = "top left";
    });
  }

  // Click delegation (capture phase to beat other handlers)
  document.addEventListener("click", function (e) {
    var btn = e.target.closest(".diagram-actions button");
    if (!btn) return;
    e.preventDefault();
    e.stopPropagation();
    e.stopImmediatePropagation();

    var bar = btn.closest(".diagram-actions");
    var mermaidDiv = bar.previousElementSibling;
    if (!mermaidDiv) return;
    var action = btn.dataset.action;

    if (action === "fullscreen") {
      openFullscreen(mermaidDiv);
      return;
    }

    // Zoom: transform the outer div (shadow DOM content scales with it)
    if (!mermaidDiv._zoomScale) mermaidDiv._zoomScale = 1;
    if (action === "zoom-in") mermaidDiv._zoomScale = Math.min(mermaidDiv._zoomScale + ZOOM_STEP, MAX_SCALE);
    if (action === "zoom-out") mermaidDiv._zoomScale = Math.max(mermaidDiv._zoomScale - ZOOM_STEP, MIN_SCALE);
    if (action === "reset") mermaidDiv._zoomScale = 1;

    mermaidDiv.style.transform = mermaidDiv._zoomScale === 1 ? "" : "scale(" + mermaidDiv._zoomScale + ")";
    mermaidDiv.style.transformOrigin = "top left";
  }, true);

  // Poll for rendered diagrams (mkdocs-material renders asynchronously)
  var pollCount = 0;
  var pollId = setInterval(function () {
    scanAll();
    if (++pollCount > 30) clearInterval(pollId);
  }, 1000);

  // Catch dynamic renders via MutationObserver (throttled)
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
