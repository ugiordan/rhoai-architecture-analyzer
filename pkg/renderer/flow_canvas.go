package renderer

// flowCanvasTemplate is the self-contained HTML + Canvas JS rendering engine.
// It reads a DIAGRAM JSON object embedded via Go template and renders an
// interactive animated flow diagram with step-by-step playback.
//
// Modules (all vanilla JS, no external dependencies):
//   Engine:     Canvas setup, coordinate transforms, render loop
//   Theme:      Dark/light color palettes
//   Renderer:   Draw nodes by type (rounded rects, glow, badges, sublabels)
//   Edges:      Edge routing, arrow drawing, permanent connectors
//   Dot:        Animated packet traveling between nodes
//   Playback:   Step-by-step execution with timing and speed control
//   StepsPanel: Sidebar list of flow steps
//   Overlay:    Click-on-node tooltip card
//   Inspector:  Architecture context panel with per-step mutations
//   Snapshot:   State save/restore for instant jump-to-step
const flowCanvasTemplate = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>{{.Title}}</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{background:#0d1117;color:#c9d1d9;font-family:-apple-system,'Segoe UI',sans-serif;overflow:hidden;display:flex;flex-direction:column;height:100vh}
.dark{--bg:#0d1117;--fg:#c9d1d9;--dim:#8b949e;--brt:#e6edf3;--panel:#161b22;--border:#30363d;--btn:#21262d;--btnHover:#30363d;--boxBg:#161b22;--boxBdr:#30363d}
.light{--bg:#ffffff;--fg:#24292f;--dim:#57606a;--brt:#1f2328;--panel:#f6f8fa;--border:#d0d7de;--btn:#f6f8fa;--btnHover:#eaeef2;--boxBg:#f6f8fa;--boxBdr:#d0d7de}
#topbar{display:flex;align-items:center;gap:8px;padding:6px 12px;background:var(--panel);border-bottom:1px solid var(--border);height:42px;flex-shrink:0}
#topbar .title{font-weight:700;font-size:13px;white-space:nowrap}
#topbar button,#topbar select{background:var(--btn);color:var(--fg);border:1px solid var(--border);border-radius:6px;padding:3px 10px;cursor:pointer;font-size:12px}
#topbar button:hover{background:var(--btnHover)}
#topbar button:disabled{opacity:.4;cursor:default}
#topbar select{padding:3px 8px}
#topbar .sep{flex:1}
#topbar .legend{font-size:11px;color:var(--dim);display:flex;gap:8px;align-items:center}
#topbar .legend span{display:flex;align-items:center;gap:3px}
#main{display:flex;flex:1;overflow:hidden}
#canvas-wrap{flex:1;position:relative}
canvas{display:block;width:100%;height:100%}
#right-panel{width:320px;flex-shrink:0;display:flex;flex-direction:column;background:var(--panel);border-left:1px solid var(--border);overflow:hidden}
#steps-panel{flex:1;overflow-y:auto;padding:8px 12px;font-size:12px}
#steps-panel .step{padding:6px 8px;border-radius:4px;margin:2px 0;cursor:pointer;color:var(--dim)}
#steps-panel .step:hover{background:var(--btnHover)}
#steps-panel .step.active{background:var(--btn);color:var(--brt);font-weight:600}
#steps-panel .step.done{color:var(--fg)}
#steps-panel .step .num{display:inline-block;width:20px;font-weight:700;color:var(--dim)}
#inspector-panel{border-top:1px solid var(--border);padding:10px 12px;max-height:40%;overflow-y:auto;font-size:12px}
#inspector-panel .phase{font-size:11px;text-transform:uppercase;letter-spacing:.5px;color:var(--dim);margin-bottom:6px}
#inspector-panel .label{font-weight:600;margin-bottom:8px;font-size:13px}
#inspector-panel .line{padding:2px 4px;margin:1px 0;border-radius:3px;font-family:monospace;font-size:11px;word-break:break-all}
#inspector-panel .line.keep{color:var(--fg)}
#inspector-panel .line.add{color:#3fb950;background:#3fb95015}
#inspector-panel .line.highlight{color:#58a6ff;background:#58a6ff15;font-weight:600}
#inspector-panel .line.err{color:#f85149;background:#f8514915}
#overlay{position:absolute;display:none;background:var(--panel);border:1px solid var(--border);border-radius:8px;padding:14px;max-width:280px;box-shadow:0 8px 24px #00000040;z-index:100;font-size:12px}
#overlay .ov-title{font-weight:700;font-size:14px;margin-bottom:4px}
#overlay .ov-desc{color:var(--dim);margin-bottom:8px}
#overlay .ov-detail{display:flex;justify-content:space-between;padding:2px 0;border-bottom:1px solid var(--border)}
#overlay .ov-detail:last-child{border-bottom:none}
#overlay .ov-key{color:var(--dim)}
#overlay .ov-val{font-weight:600}
#overlay .ov-close{position:absolute;top:8px;right:10px;cursor:pointer;color:var(--dim);font-size:16px;background:none;border:none}
</style>
</head>
<body class="dark">
<div id="topbar">
  <span class="title">{{.Title}}</span>
  <button id="btn-play">&#9654; Play</button>
  <button id="btn-pause" disabled>&#9646;&#9646;</button>
  <button id="btn-reset">&#8634; Reset</button>
  <select id="flow-select"></select>
  <button id="btn-speed">1x</button>
  <button id="btn-loop">Loop</button>
  <span class="sep"></span>
  <span id="legend" class="legend"></span>
  <button id="btn-theme">&#9728;</button>
</div>
<div id="main">
  <div id="canvas-wrap">
    <canvas id="cv"></canvas>
    <div id="overlay"></div>
  </div>
  <div id="right-panel">
    <div id="steps-panel"></div>
    <div id="inspector-panel"></div>
  </div>
</div>
<script>
var D = {{.DiagramJSON}};
(function(){
// ========== THEME ==========
var themes = {
  dark:  {bg:'#0d1117',fg:'#c9d1d9',dim:'#8b949e',brt:'#e6edf3',box:'#161b22',bdr:'#30363d',txt:'#c9d1d9'},
  light: {bg:'#ffffff',fg:'#24292f',dim:'#57606a',brt:'#1f2328',box:'#f6f8fa',bdr:'#d0d7de',txt:'#24292f'}
};
var isDark = true;
function colors(){ return isDark ? themes.dark : themes.light; }
function toggleTheme(){
  isDark = !isDark;
  document.body.className = isDark ? 'dark' : 'light';
}

// ========== ENGINE ==========
var cv = document.getElementById('cv');
var ctx = cv.getContext('2d');
var W, H, sc, ox, oy;
var logW = D.canvas.width || 1200, logH = D.canvas.height || 800;
var panelW = 320;

function resize(){
  var dpr = window.devicePixelRatio || 1;
  W = cv.parentElement.clientWidth;
  H = cv.parentElement.clientHeight;
  cv.width = W * dpr;
  cv.height = H * dpr;
  cv.style.width = W + 'px';
  cv.style.height = H + 'px';
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  sc = Math.min(W / logW, H / logH) * 0.92;
  ox = (W - logW * sc) / 2;
  oy = (H - logH * sc) * 0.1 + 10;
}
function tx(x){ return ox + x * sc; }
function ty(y){ return oy + y * sc; }
function ts(s){ return s * sc; }

// ========== STATE ==========
var activeNodes = {};
var badges = {};
var lines = [];
var glowing = {};
var fading = {};
var dots = [];
var running = false, paused = false, stepIndex = 0;
var playbackSpeed = 1, loopMode = false;
var activeFlowKey = D.defaultFlow || '';
var stepTimer = null;
var snapshots = [];

function resetState(){
  activeNodes = {};
  badges = {};
  lines = [];
  glowing = {};
  fading = {};
  dots = [];
  running = false;
  paused = false;
  stepIndex = 0;
  stepTimer = null;
  snapshots = [];
}

// ========== RENDERER ==========
function roundRect(x, y, w, h, r){
  ctx.beginPath();
  ctx.moveTo(x + r, y);
  ctx.lineTo(x + w - r, y);
  ctx.arcTo(x + w, y, x + w, y + r, r);
  ctx.lineTo(x + w, y + h - r);
  ctx.arcTo(x + w, y + h, x + w - r, y + h, r);
  ctx.lineTo(x + r, y + h);
  ctx.arcTo(x, y + h, x, y + h - r, r);
  ctx.lineTo(x, y + r);
  ctx.arcTo(x, y, x + r, y, r);
  ctx.closePath();
}

function drawNode(key, node){
  var c = colors();
  var x = tx(node.x), y = ty(node.y), w = ts(node.w), h = ts(node.h);
  var isActive = activeNodes[key];
  var hasBadge = badges[key] != null;
  var fadeAlpha = 1;
  if (fading[key]) {
    fadeAlpha = Math.min((Date.now() - fading[key]) / 400, 1);
    if (fadeAlpha >= 1) delete fading[key];
  }

  ctx.save();
  if (isActive || hasBadge) {
    ctx.shadowColor = node.color || c.dim;
    ctx.shadowBlur = ts(16) * fadeAlpha;
    ctx.globalAlpha = 0.4 + 0.6 * fadeAlpha;
  }

  // Fill
  var fillColor = node.color || c.dim;
  var r = ts(8);
  roundRect(x, y, w, h, r);
  ctx.fillStyle = fillColor + '60';
  ctx.fill();
  ctx.strokeStyle = fillColor;
  ctx.lineWidth = isActive ? 2.5 : 1.5;
  ctx.stroke();

  ctx.shadowBlur = 0;
  ctx.shadowColor = 'transparent';
  ctx.globalAlpha = 1;

  // Label
  ctx.fillStyle = c.brt;
  ctx.font = 'bold ' + ts(12) + 'px -apple-system, sans-serif';
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  var label = node.label || '';
  if (label.length > 20) label = label.slice(0, 18) + '...';
  var ly = node.sublabel ? y + h * 0.38 : y + h / 2;
  ctx.fillText(label, x + w / 2, ly);

  // Sublabel
  if (node.sublabel) {
    ctx.fillStyle = c.dim;
    ctx.font = ts(10) + 'px -apple-system, sans-serif';
    ctx.fillText(node.sublabel, x + w / 2, y + h * 0.68);
  }

  // Badge
  if (hasBadge) {
    var bv = badges[key];
    var bArr = Array.isArray(bv) ? bv : [bv];
    for (var bi = 0; bi < bArr.length; bi++) {
      var br = ts(9);
      var bx = x + w - ts(4) - bi * ts(20);
      var by = y - ts(4);
      ctx.beginPath();
      ctx.arc(bx, by, br, 0, Math.PI * 2);
      ctx.fillStyle = node.color || c.dim;
      ctx.fill();
      ctx.fillStyle = '#fff';
      ctx.font = 'bold ' + ts(9) + 'px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText(String(bArr[bi]), bx, by);
    }
  }
  ctx.restore();
}

function renderAll(){
  var c = colors();
  ctx.clearRect(0, 0, W, H);
  ctx.fillStyle = c.bg;
  ctx.fillRect(0, 0, W, H);

  // Draw permanent lines
  for (var li = 0; li < lines.length; li++) {
    drawConnector(lines[li]);
  }

  // Draw nodes
  var keys = Object.keys(D.nodes);
  for (var i = 0; i < keys.length; i++) {
    drawNode(keys[i], D.nodes[keys[i]]);
  }

  // Draw dots
  for (var di = dots.length - 1; di >= 0; di--) {
    var dot = dots[di];
    dot.update();
    if (dot.done) {
      dots.splice(di, 1);
      if (dot.onDone) dot.onDone();
    } else {
      ctx.beginPath();
      ctx.arc(dot.cx, dot.cy, ts(5), 0, Math.PI * 2);
      ctx.fillStyle = dot.color;
      ctx.fill();
      ctx.shadowColor = dot.color;
      ctx.shadowBlur = ts(8);
      ctx.fill();
      ctx.shadowBlur = 0;
      ctx.shadowColor = 'transparent';
    }
  }
}

// ========== EDGES ==========
function edgePt(n, targetX, targetY){
  var cx = n.x + n.w / 2, cy = n.y + n.h / 2;
  var dx = targetX - cx, dy = targetY - cy;
  if (Math.abs(dx) < 1 && Math.abs(dy) < 1) return {x: cx, y: cy};
  var hw = n.w / 2, hh = n.h / 2;
  var ex, ey;
  if (Math.abs(dx) * hh > Math.abs(dy) * hw) {
    ex = cx + (dx > 0 ? hw : -hw);
    ey = cy + dy * (hw / Math.abs(dx));
  } else {
    ey = cy + (dy > 0 ? hh : -hh);
    ex = cx + dx * (hh / Math.abs(dy));
  }
  return {x: ex, y: ey};
}

function drawConnector(line){
  var c = colors();
  var fn = D.nodes[line.from], tn = D.nodes[line.to];
  if (!fn || !tn) return;
  var fp = edgePt(fn, tn.x + tn.w / 2, tn.y + tn.h / 2);
  var tp = edgePt(tn, fn.x + fn.w / 2, fn.y + fn.h / 2);
  var x1 = tx(fp.x), y1 = ty(fp.y), x2 = tx(tp.x), y2 = ty(tp.y);

  ctx.beginPath();
  ctx.moveTo(x1, y1);
  var mx = (x1 + x2) / 2;
  ctx.bezierCurveTo(mx, y1, mx, y2, x2, y2);
  ctx.strokeStyle = line.color || c.dim;
  ctx.lineWidth = 1.5;
  ctx.stroke();

  // Arrowhead
  var angle = Math.atan2(y2 - (y1 + y2) / 2, x2 - (x1 + x2) / 2);
  var sz = ts(6);
  ctx.beginPath();
  ctx.moveTo(x2, y2);
  ctx.lineTo(x2 - sz * Math.cos(angle - 0.4), y2 - sz * Math.sin(angle - 0.4));
  ctx.lineTo(x2 - sz * Math.cos(angle + 0.4), y2 - sz * Math.sin(angle + 0.4));
  ctx.closePath();
  ctx.fillStyle = line.color || c.dim;
  ctx.fill();

  // Number badge on line
  if (line.num) {
    var lx = (x1 + x2) / 2, ly2 = (y1 + y2) / 2;
    var br = ts(8);
    ctx.beginPath();
    ctx.arc(lx, ly2, br, 0, Math.PI * 2);
    ctx.fillStyle = line.color || c.dim;
    ctx.fill();
    ctx.fillStyle = '#fff';
    ctx.font = 'bold ' + ts(8) + 'px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(String(line.num), lx, ly2);
  }
}

// ========== DOT ==========
function Dot(fromKey, toKey, color, speed, onDone){
  var fn = D.nodes[fromKey], tn = D.nodes[toKey];
  if (!fn || !tn) { this.done = true; return; }
  var fp = edgePt(fn, tn.x + tn.w / 2, tn.y + tn.h / 2);
  var tp = edgePt(tn, fn.x + fn.w / 2, fn.y + fn.h / 2);
  this.x1 = tx(fp.x); this.y1 = ty(fp.y);
  this.x2 = tx(tp.x); this.y2 = ty(tp.y);
  this.mx = (this.x1 + this.x2) / 2;
  this.color = color;
  this.t = 0;
  this.speed = speed;
  this.done = false;
  this.onDone = onDone;
  this.cx = this.x1;
  this.cy = this.y1;
}
Dot.prototype.update = function(){
  this.t += this.speed;
  if (this.t >= 1) { this.t = 1; this.done = true; }
  var t = this.t, u = 1 - t;
  this.cx = u*u*this.x1 + 2*u*t*this.mx + t*t*this.x2;
  this.cy = u*u*this.y1 + 2*u*t*((this.y1+this.y2)/2) + t*t*this.y2;
};

// ========== SNAPSHOT ==========
function saveSnapshot(idx){
  snapshots[idx] = {
    activeNodes: JSON.parse(JSON.stringify(activeNodes)),
    badges: JSON.parse(JSON.stringify(badges)),
    lines: JSON.parse(JSON.stringify(lines)),
    glowing: JSON.parse(JSON.stringify(glowing))
  };
}
function restoreSnapshot(idx){
  if (!snapshots[idx]) return false;
  var s = snapshots[idx];
  activeNodes = JSON.parse(JSON.stringify(s.activeNodes));
  badges = JSON.parse(JSON.stringify(s.badges));
  lines = JSON.parse(JSON.stringify(s.lines));
  glowing = JSON.parse(JSON.stringify(s.glowing));
  dots = [];
  return true;
}

// ========== PLAYBACK ==========
function getFlow(){
  return D.flows[activeFlowKey];
}

function execStep(){
  var flow = getFlow();
  if (!flow || stepIndex >= flow.steps.length) {
    running = false;
    updateButtons();
    if (loopMode) {
      var keys = D.flowOrder || Object.keys(D.flows);
      var ci = keys.indexOf(activeFlowKey);
      var next = keys[(ci + 1) % keys.length];
      setTimeout(function(){ setFlow(next); run(); }, 2000 / playbackSpeed);
    }
    return;
  }

  var step = flow.steps[stepIndex];
  updateStepsPanel();
  applyInspectorMutation(stepIndex);

  if (step.mode === 'arrow') {
    var dotSpeed = 0.015 * playbackSpeed;
    var d = new Dot(step.from, step.to, step.color || '#58a6ff', dotSpeed, function(){
      lines.push({from: step.from, to: step.to, color: step.color, num: step.num});
      activeNodes[step.to] = true;
      fading[step.to] = Date.now();
      saveSnapshot(stepIndex);
      stepIndex++;
      stepTimer = setTimeout(execStep, 250 / playbackSpeed);
    });
    dots.push(d);
  } else if (step.mode === 'lightup') {
    activeNodes[step.target] = true;
    fading[step.target] = Date.now();
    if (step.badge) {
      if (!badges[step.target]) badges[step.target] = [];
      if (!Array.isArray(badges[step.target])) badges[step.target] = [badges[step.target]];
      badges[step.target].push(step.badge);
    }
    saveSnapshot(stepIndex);
    stepIndex++;
    stepTimer = setTimeout(execStep, 600 / playbackSpeed);
  }
}

function run(){
  if (running && !paused) { pause(); return; }
  if (paused) { resume(); return; }
  resetState();
  running = true;
  updateButtons();
  updateStepsPanel();
  initInspector();
  stepTimer = setTimeout(execStep, 200);
}

function pause(){
  paused = true;
  if (stepTimer) clearTimeout(stepTimer);
  stepTimer = null;
  updateButtons();
}

function resume(){
  paused = false;
  updateButtons();
  stepTimer = setTimeout(execStep, 100);
}

function reset(){
  if (stepTimer) clearTimeout(stepTimer);
  resetState();
  updateButtons();
  updateStepsPanel();
  initInspector();
}

function jumpToStep(idx){
  if (stepTimer) clearTimeout(stepTimer);
  dots = [];
  if (restoreSnapshot(idx)) {
    stepIndex = idx + 1;
    updateStepsPanel();
    applyInspectorMutation(idx);
  }
}

function setFlow(key){
  activeFlowKey = key;
  document.getElementById('flow-select').value = key;
  reset();
}

function cycleSpeed(){
  if (playbackSpeed < 1) playbackSpeed = 1;
  else if (playbackSpeed < 2) playbackSpeed = 2;
  else playbackSpeed = 0.5;
  document.getElementById('btn-speed').textContent = playbackSpeed + 'x';
}

// ========== STEPS PANEL ==========
function updateStepsPanel(){
  var panel = document.getElementById('steps-panel');
  var flow = getFlow();
  if (!flow) { panel.innerHTML = '<p style="color:var(--dim);padding:12px">Select a flow and press Play</p>'; return; }
  var html = '';
  for (var i = 0; i < flow.steps.length; i++) {
    var s = flow.steps[i];
    var cls = 'step';
    if (i < stepIndex) cls += ' done';
    if (i === stepIndex) cls += ' active';
    html += '<div class="' + cls + '" data-idx="' + i + '">';
    html += '<span class="num">' + (s.num || '') + '</span>';
    html += s.text;
    html += '</div>';
  }
  panel.innerHTML = html;
  var items = panel.querySelectorAll('.step');
  for (var j = 0; j < items.length; j++) {
    items[j].addEventListener('click', (function(idx){ return function(){ jumpToStep(idx); }; })(j));
  }
}

// ========== INSPECTOR ==========
function initInspector(){
  var panel = document.getElementById('inspector-panel');
  if (!D.inspector) { panel.innerHTML = ''; return; }
  var st = D.inspector.initialState;
  renderInspector(st.phase || 'architecture', '', st.headers || [], st.body || []);
}

function applyInspectorMutation(idx){
  if (!D.inspector || !D.inspector.mutations) return;
  var muts = D.inspector.mutations[activeFlowKey];
  if (!muts) return;
  var mut = null;
  for (var i = 0; i < muts.length; i++) {
    if (muts[i].step === idx + 1) { mut = muts[i]; break; }
  }
  if (!mut) return;
  var phase = mut.phase || 'architecture';
  var headers = mut.replaceHeaders || D.inspector.initialState.headers || [];
  var body = mut.replaceBody || D.inspector.initialState.body || [];
  renderInspector(phase, mut.label || '', headers, body);
}

function renderInspector(phase, label, headers, body){
  var panel = document.getElementById('inspector-panel');
  var html = '<div class="phase">' + esc(phase) + '</div>';
  if (label) html += '<div class="label">' + esc(label) + '</div>';
  for (var i = 0; i < headers.length; i++) {
    html += '<div class="line ' + (headers[i].style || 'keep') + '">' + esc(headers[i].value) + '</div>';
  }
  if (body.length > 0) {
    html += '<div style="margin-top:8px;border-top:1px solid var(--border);padding-top:6px">';
    for (var j = 0; j < body.length; j++) {
      html += '<div class="line ' + (body[j].style || 'keep') + '">' + esc(body[j].value) + '</div>';
    }
    html += '</div>';
  }
  panel.innerHTML = html;
}

// ========== OVERLAY ==========
var overlayVisible = false;
function showOverlay(key, screenX, screenY){
  var tt = D.tooltips ? D.tooltips[key] : null;
  if (!tt) return;
  var ov = document.getElementById('overlay');
  var html = '<button class="ov-close" id="ov-close-btn">&#10005;</button>';
  html += '<div class="ov-title">' + esc(tt.title || key) + '</div>';
  if (tt.description) html += '<div class="ov-desc">' + esc(tt.description) + '</div>';
  if (tt.details) {
    for (var i = 0; i < tt.details.length; i++) {
      html += '<div class="ov-detail"><span class="ov-key">' + esc(tt.details[i][0]) + '</span><span class="ov-val">' + esc(tt.details[i][1]) + '</span></div>';
    }
  }
  ov.innerHTML = html;
  ov.style.display = 'block';
  var maxX = W - 290, maxY = H - ov.offsetHeight - 10;
  ov.style.left = Math.min(screenX + 10, maxX) + 'px';
  ov.style.top = Math.min(screenY + 10, maxY) + 'px';
  overlayVisible = true;
  document.getElementById('ov-close-btn').addEventListener('click', hideOverlay);
}
function hideOverlay(){
  document.getElementById('overlay').style.display = 'none';
  overlayVisible = false;
}

// ========== HIT TEST ==========
function hitTest(mx, my){
  var keys = Object.keys(D.nodes);
  for (var i = keys.length - 1; i >= 0; i--) {
    var n = D.nodes[keys[i]];
    var x = tx(n.x), y = ty(n.y), w = ts(n.w), h = ts(n.h);
    if (mx >= x && mx <= x + w && my >= y && my <= y + h) return keys[i];
  }
  return null;
}

// ========== UTILITIES ==========
function esc(s){ return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;').replace(/'/g,'&#39;'); }

function updateButtons(){
  document.getElementById('btn-play').disabled = running && !paused;
  document.getElementById('btn-pause').disabled = !running || paused;
  document.getElementById('btn-play').textContent = paused ? '▶ Resume' : '▶ Play';
}

// ========== INIT ==========
function init(){
  resize();
  window.addEventListener('resize', function(){ resize(); });

  // Populate flow selector
  var sel = document.getElementById('flow-select');
  var order = D.flowOrder || Object.keys(D.flows || {});
  if (order.length === 0) {
    var opt = document.createElement('option');
    opt.textContent = '(no flows)';
    sel.appendChild(opt);
    document.getElementById('btn-play').disabled = true;
  } else {
    for (var i = 0; i < order.length; i++) {
      var f = D.flows[order[i]];
      if (!f) continue;
      var opt = document.createElement('option');
      opt.value = order[i];
      opt.textContent = f.label || order[i];
      sel.appendChild(opt);
    }
    if (!activeFlowKey) activeFlowKey = order[0];
    sel.value = activeFlowKey;
  }

  // Legend
  var legendEl = document.getElementById('legend');
  var lhtml = '';
  if (D.legend) {
    for (var li = 0; li < D.legend.length; li++) {
      lhtml += '<span><svg width="10" height="10"><rect width="10" height="10" rx="2" fill="' + D.legend[li].color + '"/></svg> ' + D.legend[li].label + '</span>';
    }
  }
  legendEl.innerHTML = lhtml;

  // Controls
  document.getElementById('btn-play').addEventListener('click', run);
  document.getElementById('btn-pause').addEventListener('click', pause);
  document.getElementById('btn-reset').addEventListener('click', reset);
  document.getElementById('btn-speed').addEventListener('click', cycleSpeed);
  document.getElementById('btn-loop').addEventListener('click', function(){
    loopMode = !loopMode;
    this.style.opacity = loopMode ? '1' : '0.5';
  });
  document.getElementById('btn-theme').addEventListener('click', toggleTheme);
  document.getElementById('flow-select').addEventListener('change', function(){
    setFlow(this.value);
  });

  // Canvas click
  cv.addEventListener('click', function(e){
    var rect = cv.getBoundingClientRect();
    var mx = e.clientX - rect.left, my = e.clientY - rect.top;
    var key = hitTest(mx, my);
    if (key) { showOverlay(key, e.clientX - cv.parentElement.getBoundingClientRect().left, e.clientY - cv.parentElement.getBoundingClientRect().top); }
    else if (overlayVisible) { hideOverlay(); }
  });

  // Keyboard
  document.addEventListener('keydown', function(e){
    if (e.key === ' ') { e.preventDefault(); run(); }
    if (e.key === 'r') reset();
    if (e.key === 't') toggleTheme();
  });

  // Init panels
  updateStepsPanel();
  initInspector();
  document.getElementById('btn-loop').style.opacity = '0.5';

  // Render loop
  function frame(){ renderAll(); requestAnimationFrame(frame); }
  requestAnimationFrame(frame);
}

if (document.readyState === 'complete') { init(); }
else { window.addEventListener('load', init); }
})();
</script>
</body>
</html>`
