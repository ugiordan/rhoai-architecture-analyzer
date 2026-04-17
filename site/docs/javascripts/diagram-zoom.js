// Diagram toolbar: zoom in/out, fullscreen, and reset controls
// Uses event delegation and polling to work with mkdocs-material's lazy mermaid rendering
(function () {
  var ZOOM_STEP = 0.25;
  var MIN_SCALE = 0.5;
  var MAX_SCALE = 3;

  function getState(container) {
    if (!container._diagramState) {
      container._diagramState = { scale: 1, fullscreen: false };
    }
    return container._diagramState;
  }

  function createToolbar(container) {
    if (container.querySelector(".diagram-toolbar")) return;

    var svg = container.querySelector("svg");
    if (!svg) return;

    container.style.position = "relative";

    var toolbar = document.createElement("div");
    toolbar.className = "diagram-toolbar";
    toolbar.innerHTML =
      '<button data-action="zoom-in" title="Zoom in">+</button>' +
      '<button data-action="zoom-out" title="Zoom out">\u2212</button>' +
      '<button data-action="reset" title="Reset zoom">1:1</button>' +
      '<button data-action="fullscreen" title="Fullscreen">\u26F6</button>';

    container.appendChild(toolbar);
  }

  // Event delegation: handle all toolbar clicks from document
  document.addEventListener("click", function (e) {
    var btn = e.target.closest(".diagram-toolbar button");
    if (!btn) return;

    e.stopPropagation();
    e.preventDefault();

    var container = btn.closest("pre.mermaid, .mermaid");
    if (!container) return;

    var svg = container.querySelector("svg");
    if (!svg) return;

    var state = getState(container);
    var action = btn.getAttribute("data-action");

    if (action === "zoom-in") {
      state.scale = Math.min(state.scale + ZOOM_STEP, MAX_SCALE);
      svg.style.transform = "scale(" + state.scale + ")";
      svg.style.transformOrigin = "top left";
    } else if (action === "zoom-out") {
      state.scale = Math.max(state.scale - ZOOM_STEP, MIN_SCALE);
      svg.style.transform = "scale(" + state.scale + ")";
      svg.style.transformOrigin = "top left";
    } else if (action === "reset") {
      state.scale = 1;
      svg.style.transform = "";
      svg.style.transformOrigin = "";
    } else if (action === "fullscreen") {
      state.fullscreen = !state.fullscreen;
      container.classList.toggle("fullscreen", state.fullscreen);
      btn.textContent = state.fullscreen ? "\u2715" : "\u26F6";
      btn.title = state.fullscreen ? "Exit fullscreen" : "Fullscreen";
      if (!state.fullscreen) {
        state.scale = 1;
        svg.style.transform = "";
        svg.style.transformOrigin = "";
      }
    }
  });

  // ESC exits fullscreen
  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape") {
      document.querySelectorAll(".fullscreen").forEach(function (el) {
        el.classList.remove("fullscreen");
        var btn = el.querySelector('[data-action="fullscreen"]');
        if (btn) {
          btn.textContent = "\u26F6";
          btn.title = "Fullscreen";
        }
        var svg = el.querySelector("svg");
        if (svg) {
          svg.style.transform = "";
          svg.style.transformOrigin = "";
        }
        var state = getState(el);
        state.scale = 1;
        state.fullscreen = false;
      });
    }
  });

  // Poll for rendered mermaid SVGs (robust against timing issues)
  function initToolbars() {
    document.querySelectorAll("pre.mermaid, .mermaid").forEach(function (el) {
      if (el.querySelector("svg") && !el.querySelector(".diagram-toolbar")) {
        createToolbar(el);
      }
    });
  }

  // Poll every 500ms for 10 seconds after page load, then stop
  var attempts = 0;
  var maxAttempts = 20;
  var pollId = setInterval(function () {
    initToolbars();
    attempts++;
    if (attempts >= maxAttempts) clearInterval(pollId);
  }, 500);

  // Also re-init on mkdocs-material instant navigation
  if (typeof document$ !== "undefined") {
    document$.subscribe(function () {
      attempts = 0;
      pollId = setInterval(function () {
        initToolbars();
        attempts++;
        if (attempts >= maxAttempts) clearInterval(pollId);
      }, 500);
    });
  }
})();
