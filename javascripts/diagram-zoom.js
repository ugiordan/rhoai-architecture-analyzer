// Add fullscreen toggle to mermaid diagrams
// mkdocs-material renders mermaid inside <pre class="mermaid"> elements
document.addEventListener("DOMContentLoaded", function () {
  function addZoomButtons() {
    // Target both .mermaid and pre.mermaid (mkdocs-material format)
    document.querySelectorAll(".mermaid, pre.mermaid").forEach(function (el) {
      if (el.querySelector(".diagram-zoom")) return;
      // Only add button if diagram has rendered (contains SVG)
      if (!el.querySelector("svg") && !el.innerHTML.includes("<svg")) return;

      var btn = document.createElement("button");
      btn.className = "diagram-zoom";
      btn.textContent = "\u26F6";
      btn.title = "Toggle fullscreen";
      btn.addEventListener("click", function (e) {
        e.stopPropagation();
        el.classList.toggle("fullscreen");
        btn.textContent = el.classList.contains("fullscreen") ? "\u2715" : "\u26F6";
      });
      el.style.position = "relative";
      el.appendChild(btn);
    });
  }

  // Run after initial render
  setTimeout(addZoomButtons, 1000);

  // Also watch for dynamic rendering
  var observer = new MutationObserver(function () {
    addZoomButtons();
  });
  observer.observe(document.body, { childList: true, subtree: true });
});

// ESC to exit fullscreen
document.addEventListener("keydown", function (e) {
  if (e.key === "Escape") {
    document.querySelectorAll(".fullscreen").forEach(function (el) {
      el.classList.remove("fullscreen");
      var btn = el.querySelector(".diagram-zoom");
      if (btn) btn.textContent = "\u26F6";
    });
  }
});
