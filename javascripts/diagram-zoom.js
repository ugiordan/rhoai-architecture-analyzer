// Diagram toolbar: zoom in/out, fullscreen, and reset controls
(function () {
  var ZOOM_STEP = 0.25;
  var MIN_SCALE = 0.5;
  var MAX_SCALE = 3;

  function wrapAndAddToolbar(pre) {
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

  // Find the SVG anywhere within or near the wrapper
  function findSvg(wrapper) {
    // Direct descendant search (most common)
    var svg = wrapper.querySelector("svg");
    if (svg) return svg;
    // Check if mermaid created a sibling element with the SVG
    var next = wrapper.nextElementSibling;
    if (next && next.querySelector) {
      svg = next.querySelector("svg");
      if (svg) return svg;
    }
    return null;
  }

  var initDone = false;
  function initAll() {
    document.querySelectorAll("pre.mermaid").forEach(function (pre) {
      wrapAndAddToolbar(pre);
    });
    initDone = true;
  }

  // Use CAPTURE phase to get the event before anything else can stop it
  document.addEventListener("click", function (e) {
    var btn = e.target;
    // Walk up to find if we clicked inside a toolbar button
    while (btn && btn !== document) {
      if (btn.tagName === "BUTTON" && btn.parentElement &&
          btn.parentElement.classList.contains("diagram-toolbar")) {
        break;
      }
      btn = btn.parentElement;
    }
    if (!btn || btn === document || btn.tagName !== "BUTTON") return;

    e.stopPropagation();
    e.preventDefault();
    e.stopImmediatePropagation();

    var wrapper = btn.closest(".diagram-wrapper");
    if (!wrapper) {
      console.warn("[diagram-zoom] No .diagram-wrapper found for button");
      return;
    }

    if (!wrapper._zoomState) wrapper._zoomState = { scale: 1, fs: false };
    var st = wrapper._zoomState;
    var action = btn.getAttribute("data-action");

    // Fullscreen works even without SVG
    if (action === "fullscreen") {
      st.fs = !st.fs;
      wrapper.classList.toggle("fullscreen", st.fs);
      btn.textContent = st.fs ? "\u2715" : "\u26F6";
      btn.title = st.fs ? "Exit fullscreen" : "Fullscreen";
      if (!st.fs) st.scale = 1;
    } else if (action === "zoom-in") {
      st.scale = Math.min(st.scale + ZOOM_STEP, MAX_SCALE);
    } else if (action === "zoom-out") {
      st.scale = Math.max(st.scale - ZOOM_STEP, MIN_SCALE);
    } else if (action === "reset") {
      st.scale = 1;
    }

    // Apply zoom to SVG if found
    var svg = findSvg(wrapper);
    if (svg) {
      if (st.scale === 1) {
        svg.style.transform = "";
        svg.style.transformOrigin = "";
      } else {
        svg.style.transform = "scale(" + st.scale + ")";
        svg.style.transformOrigin = "top left";
      }
    }
  }, true); // <-- capture phase

  // ESC exits fullscreen
  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape") {
      document.querySelectorAll(".diagram-wrapper.fullscreen").forEach(function (w) {
        w.classList.remove("fullscreen");
        var btn = w.querySelector('[data-action="fullscreen"]');
        if (btn) { btn.textContent = "\u26F6"; btn.title = "Fullscreen"; }
        var svg = findSvg(w);
        if (svg) { svg.style.transform = ""; svg.style.transformOrigin = ""; }
        if (w._zoomState) { w._zoomState.scale = 1; w._zoomState.fs = false; }
      });
    }
  });

  // Run on DOMContentLoaded (pre.mermaid exists in static HTML)
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initAll);
  } else {
    initAll();
  }

  // Re-init on mkdocs-material instant navigation (document$ is an RxJS observable)
  if (typeof document$ !== "undefined") {
    document$.subscribe(function () { setTimeout(initAll, 100); });
  }

  // Watch for dynamically added pre.mermaid elements (throttled)
  var pending = false;
  new MutationObserver(function () {
    if (pending) return;
    pending = true;
    requestAnimationFrame(function () {
      pending = false;
      if (initDone) initAll();
    });
  }).observe(document.body, { childList: true, subtree: true });
})();
