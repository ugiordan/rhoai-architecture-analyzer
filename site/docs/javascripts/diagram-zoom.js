// Diagram toolbar: zoom in/out, fullscreen, and reset controls
// Wraps mermaid containers in a div so the toolbar survives mermaid's innerHTML replacement
(function () {
  var ZOOM_STEP = 0.25;
  var MIN_SCALE = 0.5;
  var MAX_SCALE = 3;

  function wrapAndAddToolbar(pre) {
    // Already wrapped
    if (pre.parentElement && pre.parentElement.classList.contains("diagram-wrapper")) return;

    var wrapper = document.createElement("div");
    wrapper.className = "diagram-wrapper";
    pre.parentNode.insertBefore(wrapper, pre);
    wrapper.appendChild(pre);

    var toolbar = document.createElement("div");
    toolbar.className = "diagram-toolbar";
    toolbar.innerHTML =
      '<button data-action="zoom-in" title="Zoom in">+</button>' +
      '<button data-action="zoom-out" title="Zoom out">\u2212</button>' +
      '<button data-action="reset" title="Reset zoom">1:1</button>' +
      '<button data-action="fullscreen" title="Fullscreen">\u26F6</button>';
    wrapper.appendChild(toolbar);
  }

  function initAll() {
    document.querySelectorAll("pre.mermaid").forEach(function (pre) {
      wrapAndAddToolbar(pre);
    });
  }

  // Event delegation for all toolbar clicks
  document.addEventListener("click", function (e) {
    var btn = e.target.closest(".diagram-toolbar button");
    if (!btn) return;
    e.stopPropagation();
    e.preventDefault();

    var wrapper = btn.closest(".diagram-wrapper");
    if (!wrapper) return;
    var svg = wrapper.querySelector("svg");
    if (!svg) return;

    if (!wrapper._zoomState) wrapper._zoomState = { scale: 1, fs: false };
    var st = wrapper._zoomState;
    var action = btn.getAttribute("data-action");

    if (action === "zoom-in") {
      st.scale = Math.min(st.scale + ZOOM_STEP, MAX_SCALE);
    } else if (action === "zoom-out") {
      st.scale = Math.max(st.scale - ZOOM_STEP, MIN_SCALE);
    } else if (action === "reset") {
      st.scale = 1;
    } else if (action === "fullscreen") {
      st.fs = !st.fs;
      wrapper.classList.toggle("fullscreen", st.fs);
      btn.textContent = st.fs ? "\u2715" : "\u26F6";
      btn.title = st.fs ? "Exit fullscreen" : "Fullscreen";
      if (!st.fs) st.scale = 1;
    }

    if (st.scale === 1) {
      svg.style.transform = "";
      svg.style.transformOrigin = "";
    } else {
      svg.style.transform = "scale(" + st.scale + ")";
      svg.style.transformOrigin = "top left";
    }
  });

  // ESC exits fullscreen
  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape") {
      document.querySelectorAll(".diagram-wrapper.fullscreen").forEach(function (w) {
        w.classList.remove("fullscreen");
        var btn = w.querySelector('[data-action="fullscreen"]');
        if (btn) { btn.textContent = "\u26F6"; btn.title = "Fullscreen"; }
        var svg = w.querySelector("svg");
        if (svg) { svg.style.transform = ""; svg.style.transformOrigin = ""; }
        if (w._zoomState) { w._zoomState.scale = 1; w._zoomState.fs = false; }
      });
    }
  });

  // Run immediately (pre.mermaid exists in static HTML before mermaid JS runs)
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initAll);
  } else {
    initAll();
  }

  // Re-init on mkdocs-material instant navigation
  if (typeof document$ !== "undefined") {
    document$.subscribe(function () { setTimeout(initAll, 100); });
  }

  // Fallback: watch for new pre.mermaid elements added dynamically
  new MutationObserver(function () { initAll(); })
    .observe(document.body, { childList: true, subtree: true });
})();
