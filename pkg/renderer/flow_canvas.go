package renderer

// flowCanvasTemplate is the self-contained HTML + Canvas JS rendering engine.
// Built from scratch, no external dependencies.
// Reads a DIAGRAM JSON object (flowstory-compatible schema) and renders an
// interactive animated flow diagram with step-by-step playback.
//
// Architecture: 11 JS modules in a single IIFE (~2000 lines):
//   1. Theme      — dark/light palettes
//   2. Engine     — Canvas setup, DPR, coordinate transforms, rAF loop
//   3. Renderer   — 6-layer node drawing (boundary, container, section, stack, box, icon)
//   4. Edges      — edge routing, arrow drawing, numbered badges, waypoints
//   5. Dot        — animated pentagon packet with trail and glow
//   6. Playback   — step-by-step execution (arrow + lightup), timing, speed
//   7. Snapshot   — state save/restore for jump-to-step
//   8. StepsPanel — sidebar DOM list with ◯/●/✓ marks
//   9. Inspector  — mutation engine for request/response display
//  10. Overlay    — click-on-node tooltip card
//  11. Init       — DOM wiring, keyboard, legend, title
const flowCanvasTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.Title}}</title>
<style>
:root{
--bg:#0d1117;--surf:#161b22;--surfA:#1a2332;--bdr:#30363d;
--txt:#c9d1d9;--dim:#8b949e;--muted:#484f58;--brt:#fff;
--accent:#58a6ff;--accentHov:#58a6ff22;
--ok:#3fb950;--warn:#f0883e;--err:#f85149;--purple:#d2a8ff;
--addBg:#1a2e1a;--addC:#3fb950;--delBg:#2d1517;--delC:#f85149;
--hlBg:#1a2332;--hlC:#58a6ff;--plugBg:#d2a8ff22;--plugC:#d2a8ff;
--stepHovBg:#58a6ff22;--stepHovC:#58a6ff;--stepHovBdr:#58a6ff;
}
body.light{
--bg:#fff;--surf:#f6f8fa;--surfA:#ddf4ff;--bdr:#d0d7de;
--txt:#24292f;--dim:#656d76;--muted:#b0b8c1;--brt:#000;
--accent:#0969da;--accentHov:#ddf4ff88;
--ok:#1a7f37;--warn:#bf8700;--err:#82071e;--purple:#8250df;
--addBg:#dafbe1;--addC:#116329;--delBg:#ffebe9;--delC:#82071e;
--hlBg:#ddf4ff;--hlC:#0969da;--plugBg:#8250df15;--plugC:#8250df;
--stepHovBg:#ddf4ff88;--stepHovC:#0969da;--stepHovBdr:#0969da;
}
*{box-sizing:border-box;margin:0;padding:0}
body{background:var(--bg);color:var(--txt);font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;overflow:hidden}
canvas{display:block;position:fixed;top:0;left:0}
.panel{position:fixed;top:0;right:0;width:380px;height:100vh;display:flex;flex-direction:column;background:var(--surf);border-left:1px solid var(--bdr);z-index:100}
.flow-bar{padding:12px 16px;border-bottom:1px solid var(--bdr);display:flex;gap:8px;align-items:center;flex-wrap:wrap}
.flow-bar select{flex:1;padding:7px 10px;border:1px solid var(--bdr);border-radius:6px;background:var(--bg);color:var(--txt);font-size:13px}
.flow-bar button{padding:7px 14px;border:1px solid var(--bdr);border-radius:6px;background:var(--surf);color:var(--txt);font-size:13px;cursor:pointer;transition:all .15s}
.flow-bar button:hover{background:var(--accentHov)}
.btn-start{border-color:var(--ok)!important;color:var(--ok)!important}
.btn-reset{border-color:var(--warn)!important;color:var(--warn)!important}
.steps{padding:10px 16px;border-bottom:1px solid var(--bdr);flex:2;overflow-y:auto}
.steps-title{font-size:12px;font-weight:700;margin-bottom:8px;color:var(--dim)}
.step{margin:2px 0;color:var(--muted);display:flex;align-items:flex-start;gap:6px;font-size:12px;line-height:1.4;padding:3px 4px;border-radius:4px;cursor:pointer;transition:all .2s;border-left:2px solid transparent}
.step:hover{background:var(--stepHovBg);color:var(--stepHovC);transform:translateX(4px);border-left-color:var(--stepHovBdr)}
.step.active{color:var(--accent);font-weight:600}
.step.done{color:var(--txt)}
.step-mark{width:16px;text-align:center;flex-shrink:0}
.insp{flex:3;overflow-y:auto;padding:10px 16px;font-family:'SF Mono',Menlo,Consolas,monospace;font-size:11px}
.insp-title{font-weight:700;font-size:14px;margin-bottom:6px;padding-bottom:6px;border-bottom:1px solid var(--bdr)}
.insp-label{font-weight:600;margin:6px 0 4px;font-size:12px;color:var(--txt)}
.insp-section{margin:6px 0 2px;font-weight:600;color:var(--accent);font-size:12px}
.insp-line{margin:1px 0;padding:2px 6px;border-radius:3px;white-space:pre;line-height:1.6;transition:all .4s}
.insp-line.keep{color:var(--txt)}
.insp-line.add{background:var(--addBg);color:var(--addC);border-left:3px solid var(--addC)}
.insp-line.del{background:var(--delBg);color:var(--delC);text-decoration:line-through;opacity:.7;border-left:3px solid var(--delC)}
.insp-line.highlight{background:var(--hlBg);color:var(--hlC);animation:pulse .8s}
.insp-line.plugin{background:var(--plugBg);color:var(--plugC);font-weight:600}
@keyframes pulse{0%{opacity:.3}50%{opacity:1}100%{opacity:1}}
#overlay{display:none;position:fixed;top:0;left:0;right:380px;bottom:0;z-index:300;background:rgba(0,0,0,.6)}
#ov-card{position:absolute;background:var(--surf);border:1px solid var(--bdr);border-radius:14px;max-width:420px;padding:22px 26px;box-shadow:0 12px 40px rgba(0,0,0,.5)}
#ov-card h2{margin:0 0 6px;font-size:18px;color:var(--txt)}
#ov-card p{margin:0 0 14px;font-size:13px;color:var(--dim);line-height:1.5}
#ov-card .detail{display:flex;justify-content:space-between;padding:3px 0;border-bottom:1px solid var(--bdr);font-size:12px;font-family:'SF Mono',Menlo,monospace}
#ov-card .detail:last-child{border-bottom:none}
#ov-card .dk{color:var(--dim)}.dv{font-weight:600}
#ov-close{position:absolute;top:8px;right:14px;background:none;border:none;color:var(--dim);font-size:18px;cursor:pointer;padding:4px 8px;border-radius:6px}
#ov-close:hover{color:var(--txt)}
#ov-resume{margin-top:14px;padding:7px 20px;border:1px solid var(--ok);border-radius:8px;background:transparent;color:var(--ok);font-size:13px;cursor:pointer}
#ov-resume:hover{background:var(--ok);color:var(--bg)}
#ov-accent{width:4px;height:24px;border-radius:2px;position:absolute;left:0;top:22px}
.theme-btn{position:fixed;top:12px;right:396px;z-index:200;padding:6px 12px;border:1px solid var(--bdr);border-radius:8px;background:var(--surf);color:var(--dim);font-size:16px;cursor:pointer}
.title{position:fixed;top:14px;left:0;right:380px;text-align:center;font-size:22px;font-weight:700;z-index:200;background:linear-gradient(135deg,#58a6ff,#d2a8ff,#f0883e);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text;pointer-events:none}
.legend{position:fixed;bottom:16px;left:16px;padding:10px 14px;background:var(--surf);border:1px solid var(--bdr);border-radius:10px;font-size:13px;z-index:100;display:flex;flex-direction:column;gap:4px}
.legend-item{display:flex;align-items:center;gap:8px}
.legend-dot{width:11px;height:11px;border-radius:50%;flex-shrink:0}
::-webkit-scrollbar{width:6px}
::-webkit-scrollbar-track{background:transparent}
::-webkit-scrollbar-thumb{background:var(--bdr);border-radius:3px}
</style>
</head>
<body>
<button class="theme-btn" id="theme-btn">&#9790;</button>
<div class="title" id="diagram-title"></div>
<div class="panel">
  <div class="flow-bar">
    <select id="flow-sel"></select>
    <button class="btn-start" id="btn-play">&#9654; Start</button>
    <button id="btn-speed">1x</button>
    <button id="btn-loop" style="border-color:var(--accent);color:var(--accent)">&#8635; Loop</button>
    <button class="btn-reset" id="btn-reset">Reset</button>
  </div>
  <div class="steps">
    <div class="steps-title">Flow Steps <span style="font-weight:400;font-size:11px;color:var(--dim)">&#183; click to jump</span></div>
    <div id="steps-container"></div>
  </div>
  <div class="insp">
    <div class="insp-title" id="insp-title" style="color:var(--accent)">Request</div>
    <div id="insp-content"></div>
  </div>
</div>
<div id="overlay">
  <div id="ov-card">
    <button id="ov-close">&#10005;</button>
    <div id="ov-accent"></div>
    <h2 id="ov-title"></h2>
    <p id="ov-desc"></p>
    <div id="ov-details"></div>
    <button id="ov-resume">&#9654; Resume Flow</button>
  </div>
</div>
<div class="legend" id="legend"></div>
<canvas id="cv"></canvas>

<script>
var D = {{.DiagramJSON}};
(function(){
'use strict';

// ========== 1. THEME ==========
var DARK={bg:'#0d1117',box:'#161b22',boxA:'#1a2332',bdr:'#30363d',txt:'#c9d1d9',dim:'#8b949e',brt:'#fff'};
var LIGHT={bg:'#fff',box:'#f6f8fa',boxA:'#ddf4ff',bdr:'#d0d7de',txt:'#24292f',dim:'#656d76',brt:'#000'};
var isDark=true;
function C(){return isDark?DARK:LIGHT}
function toggleTheme(){isDark=!isDark;document.body.classList.toggle('light',!isDark)}

// ========== 2. ENGINE ==========
var cv=document.getElementById('cv'),ctx=cv.getContext('2d');
var W,H,sc,ox,oy;
var logW=D.canvas?D.canvas.width:1250,logH=D.canvas?D.canvas.height:1050;
var PW=380;

function resize(){
  var dpr=devicePixelRatio||1;
  W=innerWidth-PW;H=innerHeight;
  cv.width=W*dpr;cv.height=H*dpr;
  cv.style.width=W+'px';cv.style.height=H+'px';
  ctx.setTransform(dpr,0,0,dpr,0,0);
  sc=Math.min(W/logW,H/logH)*0.92;
  ox=(W-logW*sc)/2;
  oy=(H-logH*sc)*0.04+10;
}
function tx(x){return ox+x*sc}
function ty(y){return oy+y*sc}
function ts(s){return s*sc}

// ========== STATE ==========
var activeNodes={},badges={},lines=[],glowing={},fading={},dots=[];
var running=false,paused=false,stepIdx=0,speed=1,loopMode=false;
var activeFlow=D.defaultFlow||'';
var stepTimer=null,snapshots=[];

function resetState(){
  activeNodes={};badges={};lines=[];glowing={};fading={};dots=[];
  running=false;paused=false;stepIdx=0;stepTimer=null;snapshots=[];
}

// ========== 3. RENDERER ==========
function rr(x,y,w,h,r){ctx.beginPath();ctx.roundRect(x,y,w,h,r);ctx.closePath()}

function drawBoundary(k,n,c){
  ctx.save();
  ctx.setLineDash([ts(10),ts(8)]);
  ctx.lineWidth=ts(2.5);
  ctx.strokeStyle=c.bdr;
  rr(tx(n.x),ty(n.y),ts(n.w),ts(n.h),ts(14));
  ctx.stroke();
  ctx.setLineDash([]);
  if(n.label){
    ctx.font=ts(13)+'px -apple-system,sans-serif';
    ctx.fillStyle=c.dim;
    ctx.textAlign=n.labelAlign==='left'?'left':'right';
    ctx.textBaseline='top';
    var lx=n.labelAlign==='left'?tx(n.x)+ts(12):tx(n.x+n.w)-ts(12);
    ctx.fillText(n.label,lx,ty(n.y)+ts(8));
  }
  ctx.restore();
}

function drawContainer(k,n,c){
  ctx.save();
  var isGlow=glowing[k];
  ctx.setLineDash([ts(7),ts(5)]);
  if(isGlow){
    ctx.shadowColor=n.color||c.dim;
    ctx.shadowBlur=ts(22);
    ctx.lineWidth=ts(3);
    ctx.strokeStyle=(n.color||c.dim)+'cc';
  }else{
    ctx.lineWidth=ts(2.5);
    ctx.strokeStyle=(n.color||c.dim)+'55';
  }
  rr(tx(n.x),ty(n.y),ts(n.w),ts(n.h),ts(10));
  ctx.stroke();
  ctx.shadowBlur=0;ctx.shadowColor='transparent';
  ctx.setLineDash([]);
  if(n.label){
    ctx.font='bold '+ts(12)+'px -apple-system,sans-serif';
    ctx.fillStyle=(n.color||c.dim)+'cc';
    ctx.textAlign='left';ctx.textBaseline='top';
    ctx.fillText(n.label,tx(n.x)+ts(12),ty(n.y)+ts(8));
  }
  ctx.restore();
}

function drawSections(n,c){
  if(!n.sections)return;
  for(var i=0;i<n.sections.length;i++){
    var s=n.sections[i];
    ctx.save();
    rr(tx(s.x||n.x+8),ty(s.y||n.y+30),ts(s.w||n.w-16),ts(s.height||40),ts(8));
    ctx.fillStyle=(s.color||c.dim)+(isDark?'0a':'08');
    ctx.fill();
    ctx.setLineDash([ts(4),ts(3)]);
    ctx.strokeStyle=(s.color||c.dim)+(isDark?'30':'25');
    ctx.lineWidth=ts(1);
    ctx.stroke();
    ctx.setLineDash([]);
    if(s.label){
      ctx.font='bold '+ts(10)+'px -apple-system,sans-serif';
      ctx.fillStyle=(s.color||c.dim)+'88';
      ctx.textAlign='left';ctx.textBaseline='top';
      ctx.fillText(s.label,tx(s.labelX||s.x||n.x+16),ty(s.labelY||s.y||n.y+34));
    }
    ctx.restore();
  }
}

function drawStack(n,c){
  if(!n.stackCount)return;
  var dx=n.stackOffset?n.stackOffset.dx:-8;
  var dy=n.stackOffset?n.stackOffset.dy:-5;
  for(var j=n.stackCount;j>=1;j--){
    ctx.save();
    rr(tx(n.x+dx*j),ty(n.y+dy*j),ts(n.w),ts(n.h),ts(6));
    ctx.fillStyle=(n.color||c.dim)+(isDark?'18':'12');
    ctx.fill();
    ctx.strokeStyle=(n.color||c.dim)+(isDark?'44':'33');
    ctx.lineWidth=ts(1);
    ctx.stroke();
    ctx.restore();
  }
}

function drawBox(k,n,c){
  if(!n||n.type==='boundary'||n.type==='container')return;
  var x=tx(n.x),y=ty(n.y),w=ts(n.w),h=ts(n.h);
  var isAct=activeNodes[k],hasBdg=badges[k]!=null;
  var fa=1;
  if(fading[k]){fa=Math.min((Date.now()-fading[k])/500,1);if(fa>=1)delete fading[k]}

  ctx.save();

  // Icon type
  if(n.type==='icon'){
    ctx.font=ts(32)+'px serif';
    ctx.textAlign='center';ctx.textBaseline='middle';
    ctx.fillText('\u{1F464}',x+w/2,y+ts(22));
    ctx.font='bold '+ts(14)+'px -apple-system,sans-serif';
    ctx.fillStyle=c.txt;
    ctx.fillText(n.label||'',x+w/2,y+ts(44));
    ctx.restore();return;
  }

  // Glow
  if(isAct||hasBdg){
    ctx.shadowColor=n.color||c.dim;
    ctx.shadowBlur=ts(18)*fa;
    ctx.globalAlpha=0.3+0.7*fa;
  }

  // Box shape
  var isDash=n.type==='plugin';
  if(isDash)ctx.setLineDash([ts(5),ts(4)]);
  rr(x,y,w,h,ts(8));
  ctx.fillStyle=(isAct||hasBdg)?c.boxA:c.box;
  ctx.fill();
  ctx.lineWidth=ts((isAct||hasBdg)?2.5:1.5);
  ctx.strokeStyle=(isAct||hasBdg)?(n.color||c.bdr):c.bdr;
  ctx.stroke();
  if(isDash)ctx.setLineDash([]);

  ctx.shadowBlur=0;ctx.shadowColor='transparent';ctx.globalAlpha=1;

  // Text
  var fs=n.fontSize||16;
  ctx.textAlign='center';ctx.textBaseline='middle';
  if(n.sublabel){
    ctx.font='bold '+ts(fs)+'px -apple-system,sans-serif';
    ctx.fillStyle=c.brt;
    ctx.fillText(n.label||'',x+w/2,y+h/2-ts(8));
    ctx.font=ts(Math.max(fs-3,11))+'px -apple-system,sans-serif';
    ctx.fillStyle=c.dim;
    ctx.fillText(n.sublabel,x+w/2,y+h/2+ts(8));
  }else{
    ctx.font='bold '+ts(fs)+'px -apple-system,sans-serif';
    ctx.fillStyle=c.brt;
    ctx.fillText(n.label||'',x+w/2,y+h/2);
  }

  // Badges
  if(hasBdg){
    var bv=badges[k],ba=Array.isArray(bv)?bv:[bv];
    for(var bi=0;bi<ba.length;bi++){
      var br=ts(9),bx=x+ts(8)+bi*ts(22),by=y+ts(5);
      ctx.beginPath();ctx.arc(bx+br,by+br,br,0,Math.PI*2);
      ctx.fillStyle=n.color||c.dim;ctx.fill();
      ctx.font='bold '+ts(10)+'px -apple-system,sans-serif';
      ctx.fillStyle=isDark?'#0d1117':'#fff';
      ctx.textAlign='center';ctx.textBaseline='middle';
      ctx.fillText(String(ba[bi]),bx+br,by+br);
    }
  }
  ctx.restore();
}

function renderAll(){
  var c=C();
  ctx.clearRect(0,0,W,H);
  ctx.fillStyle=c.bg;ctx.fillRect(0,0,W,H);

  var keys=Object.keys(D.nodes||{});
  // Layer 1: boundaries
  for(var i=0;i<keys.length;i++){var n=D.nodes[keys[i]];if(n.type==='boundary')drawBoundary(keys[i],n,c)}
  // Layer 2: containers
  for(var i=0;i<keys.length;i++){var n=D.nodes[keys[i]];if(n.type==='container'){drawContainer(keys[i],n,c);drawSections(n,c)}}
  // Layer 5: stacks
  for(var i=0;i<keys.length;i++){var n=D.nodes[keys[i]];if(n.stackCount)drawStack(n,c)}
  // Lines (persisted arrows)
  for(var i=0;i<lines.length;i++)drawConnector(lines[i],c);
  // Layer 6: boxes
  for(var i=0;i<keys.length;i++)drawBox(keys[i],D.nodes[keys[i]],c);
  // Dots
  for(var i=dots.length-1;i>=0;i--){
    dots[i].update();
    if(dots[i].dead){if(dots[i].cb)dots[i].cb();dots.splice(i,1)}
    else dots[i].draw(ctx,c);
  }
}

// ========== 4. EDGES ==========
function edgePt(n,tx2,ty2){
  var cx=n.x+n.w/2,cy=n.y+n.h/2;
  var dx=tx2-cx,dy=ty2-cy;
  if(Math.abs(dx)<1&&Math.abs(dy)<1)return{x:cx,y:cy};
  var hw=n.w/2,hh=n.h/2,ex,ey;
  if(Math.abs(dx)*hh>Math.abs(dy)*hw){ex=cx+(dx>0?hw:-hw);ey=cy+dy*(hw/Math.abs(dx))}
  else{ey=cy+(dy>0?hh:-hh);ex=cx+dx*(hh/Math.abs(dy))}
  return{x:ex,y:ey}
}

function resolveEdge(n,l,isFrom){
  var yo=l.yOff||0,xo=l.xOff||0;
  if(isFrom){
    if(l.fromLeft)return{x:n.x,y:n.y+n.h/2+yo};
    if(l.fromRight)return{x:n.x+n.w,y:n.y+n.h/2+yo};
    if(l.fromTop)return{x:n.x+n.w/2+(l.fromXOff||0),y:n.y+yo};
    if(l.fromBottom)return{x:n.x+n.w/2+(l.fromXOff||0),y:n.y+n.h+yo};
  }else{
    if(l.toLeft)return{x:n.x,y:n.y+n.h/2+yo};
    if(l.toRight)return{x:n.x+n.w,y:n.y+n.h/2+yo};
    if(l.toTop)return{x:n.x+n.w/2+(l.toXOff||0),y:n.y+yo};
    if(l.toBottom)return{x:n.x+n.w/2+(l.toXOff||0),y:n.y+n.h+yo};
  }
  return null;
}

function getEndpoints(l){
  var fn=D.nodes[l.from],tn=D.nodes[l.to];
  if(!fn||!tn)return null;
  var fp=resolveEdge(fn,l,true)||edgePt(fn,tn.x+tn.w/2,tn.y+tn.h/2);
  var tp=resolveEdge(tn,l,false)||edgePt(tn,fn.x+fn.w/2,fn.y+fn.h/2);
  return{fp:fp,tp:tp};
}

function drawConnector(l,c){
  var ep=getEndpoints(l);if(!ep)return;
  var fp=ep.fp,tp=ep.tp;
  var wps=l.waypoints;

  ctx.save();
  ctx.beginPath();
  if(wps&&wps.length){
    var pts=[fp].concat(wps).concat([tp]);
    ctx.moveTo(tx(pts[0].x),ty(pts[0].y));
    for(var i=1;i<pts.length;i++)ctx.lineTo(tx(pts[i].x),ty(pts[i].y));
  }else{
    ctx.moveTo(tx(fp.x),ty(fp.y));
    ctx.lineTo(tx(tp.x),ty(tp.y));
  }
  ctx.strokeStyle=l.color||c.dim;
  ctx.lineWidth=ts(2.5);
  ctx.stroke();

  // Arrowhead
  var ax,ay,px,py;
  if(wps&&wps.length){ax=tp.x;ay=tp.y;px=wps[wps.length-1].x;py=wps[wps.length-1].y}
  else{ax=tp.x;ay=tp.y;px=fp.x;py=fp.y}
  var ang=Math.atan2(ay-py,ax-px),sz=ts(10);
  ctx.beginPath();
  ctx.moveTo(tx(ax),ty(ay));
  ctx.lineTo(tx(ax)-sz*Math.cos(ang-0.4),ty(ay)-sz*Math.sin(ang-0.4));
  ctx.lineTo(tx(ax)-sz*Math.cos(ang+0.4),ty(ay)-sz*Math.sin(ang+0.4));
  ctx.closePath();
  ctx.fillStyle=l.color||c.dim;
  ctx.fill();

  // Badge on line
  if(l.num!=null){
    var mx,my;
    if(wps&&wps.length){mx=(fp.x+tp.x)/2;my=(fp.y+tp.y)/2}
    else{mx=(fp.x+tp.x)/2;my=(fp.y+tp.y)/2}
    var r=ts(12);
    ctx.beginPath();ctx.arc(tx(mx),ty(my),r,0,Math.PI*2);
    ctx.fillStyle=l.color||c.dim;ctx.fill();
    ctx.font='bold '+ts(10)+'px -apple-system,sans-serif';
    ctx.fillStyle=isDark?'#0d1117':'#fff';
    ctx.textAlign='center';ctx.textBaseline='middle';
    ctx.fillText(String(l.num),tx(mx),ty(my));
  }
  ctx.restore();
}

// ========== 5. DOT ==========
function Dot(fromKey,toKey,color,spd,cb,opts){
  var fn=D.nodes[fromKey],tn=D.nodes[toKey];
  if(!fn||!tn){this.dead=true;return}
  opts=opts||{};
  var l={from:fromKey,to:toKey,fromLeft:opts.fromLeft,fromRight:opts.fromRight,fromTop:opts.fromTop,fromBottom:opts.fromBottom,toLeft:opts.toLeft,toRight:opts.toRight,toTop:opts.toTop,toBottom:opts.toBottom,yOff:opts.yOff,xOff:opts.xOff,fromXOff:opts.fromXOff,toXOff:opts.toXOff,waypoints:opts.waypoints};
  var ep=getEndpoints(l);if(!ep){this.dead=true;return}
  var wps=opts.waypoints||[];
  this.pts=[ep.fp].concat(wps).concat([ep.tp]);
  this.color=color||'#58a6ff';
  this.cb=cb;this.t=0;this.dead=false;
  // Segment lengths
  this.segLens=[];this.totalLen=0;
  for(var i=0;i<this.pts.length-1;i++){
    var dx=this.pts[i+1].x-this.pts[i].x,dy=this.pts[i+1].y-this.pts[i].y;
    var len=Math.sqrt(dx*dx+dy*dy);
    this.segLens.push(len);this.totalLen+=len;
  }
  this.speed=(4/(this.totalLen||1))*(spd/0.02);
  this.x=this.pts[0].x;this.y=this.pts[0].y;this.segIdx=0;
}
Dot.prototype.posAt=function(t){
  if(t<=0)return{x:this.pts[0].x,y:this.pts[0].y};
  if(t>=1)return{x:this.pts[this.pts.length-1].x,y:this.pts[this.pts.length-1].y};
  var d=t*this.totalLen,acc=0;
  for(var i=0;i<this.segLens.length;i++){
    if(acc+this.segLens[i]>=d){
      var st=(d-acc)/this.segLens[i];
      return{x:this.pts[i].x+(this.pts[i+1].x-this.pts[i].x)*st,y:this.pts[i].y+(this.pts[i+1].y-this.pts[i].y)*st};
    }
    acc+=this.segLens[i];
  }
  return{x:this.pts[this.pts.length-1].x,y:this.pts[this.pts.length-1].y};
};
Dot.prototype.update=function(){
  this.t+=this.speed;if(this.t>=1){this.t=1;this.dead=true}
  var p=this.posAt(this.t);this.x=p.x;this.y=p.y;
  // Find segment index for orientation
  var d=this.t*this.totalLen,acc=0;
  for(var i=0;i<this.segLens.length;i++){if(acc+this.segLens[i]>=d){this.segIdx=i;break}acc+=this.segLens[i]}
};
Dot.prototype.draw=function(ctx,c){
  ctx.save();
  // Trail line
  ctx.beginPath();
  ctx.moveTo(tx(this.pts[0].x),ty(this.pts[0].y));
  for(var i=0;i<this.segIdx;i++)ctx.lineTo(tx(this.pts[i+1].x),ty(this.pts[i+1].y));
  ctx.lineTo(tx(this.x),ty(this.y));
  ctx.strokeStyle=this.color;ctx.lineWidth=ts(2);ctx.globalAlpha=0.3;
  ctx.stroke();
  ctx.globalAlpha=1;

  // Trailing dots
  for(var j=4;j>=1;j--){
    var bt=this.t-j*0.06;if(bt<=0)continue;
    var p=this.posAt(bt);
    ctx.beginPath();ctx.arc(tx(p.x),ty(p.y),ts(2.5),0,Math.PI*2);
    ctx.globalAlpha=0.15+(1-j/4)*0.4;
    ctx.fillStyle=this.color;ctx.fill();
  }
  ctx.globalAlpha=1;

  // Pentagon packet
  var si=this.segIdx;
  var ang=Math.atan2(this.pts[Math.min(si+1,this.pts.length-1)].y-this.pts[si].y,this.pts[Math.min(si+1,this.pts.length-1)].x-this.pts[si].x);
  ctx.translate(tx(this.x),ty(this.y));
  ctx.rotate(ang);
  var pw=ts(13),ph=ts(7);
  ctx.beginPath();
  ctx.moveTo(pw/2,0);ctx.lineTo(pw/6,-ph/2);ctx.lineTo(-pw/2,-ph/2);ctx.lineTo(-pw/2,ph/2);ctx.lineTo(pw/6,ph/2);
  ctx.closePath();
  ctx.fillStyle=this.color;
  ctx.shadowColor=this.color;ctx.shadowBlur=ts(16);
  ctx.fill();
  ctx.restore();
};

// ========== 6. PLAYBACK ==========
function getFlow(){return D.flows?D.flows[activeFlow]:null}

function stepOpts(s){return{fromLeft:s.fromLeft,toLeft:s.toLeft,fromRight:s.fromRight,toRight:s.toRight,fromBottom:s.fromBottom,fromTop:s.fromTop,toTop:s.toTop,toBottom:s.toBottom,yOff:s.yOff,xOff:s.xOff,fromXOff:s.fromXOff,toXOff:s.toXOff,waypoints:s.waypoints}}

function execStep(){
  var flow=getFlow();
  if(!flow||stepIdx>=flow.steps.length){
    running=false;updateBtns();
    if(loopMode){
      var ks=D.flowOrder||Object.keys(D.flows||{});
      var ci=ks.indexOf(activeFlow);
      var nxt=ks[(ci+1)%ks.length];
      stepTimer=setTimeout(function(){setFlow(nxt);run()},2000/speed);
    }
    return;
  }
  var s=flow.steps[stepIdx];
  updateSteps();applyMutation(stepIdx);

  if(s.mode==='arrow'){
    var dotSpd=0.012*speed;
    var d=new Dot(s.from,s.to,s.color||'#58a6ff',dotSpd,function(){
      lines.push({from:s.from,to:s.to,color:s.color,num:s.num,fromLeft:s.fromLeft,toLeft:s.toLeft,fromRight:s.fromRight,toRight:s.toRight,fromBottom:s.fromBottom,fromTop:s.fromTop,toTop:s.toTop,toBottom:s.toBottom,yOff:s.yOff,xOff:s.xOff,fromXOff:s.fromXOff,toXOff:s.toXOff,waypoints:s.waypoints});
      activeNodes[s.to]=true;fading[s.to]=Date.now();
      if(s.glow)glowing[s.glow]=true;
      saveSnap(stepIdx);stepIdx++;
      stepTimer=setTimeout(execStep,250/speed);
    },stepOpts(s));
    dots.push(d);
  }else if(s.mode==='lightup'){
    activeNodes[s.target]=true;fading[s.target]=Date.now();
    if(s.badge!=null){
      if(!badges[s.target])badges[s.target]=[];
      if(!Array.isArray(badges[s.target]))badges[s.target]=[badges[s.target]];
      badges[s.target].push(s.badge);
    }
    if(s.errColor&&D.nodes[s.target])D.nodes[s.target]._origColor=D.nodes[s.target].color,D.nodes[s.target].color=s.errColor;
    saveSnap(stepIdx);stepIdx++;
    stepTimer=setTimeout(execStep,600/speed);
  }
}

function run(){
  if(running&&!paused){pause();return}
  if(paused){resume();return}
  resetState();running=true;updateBtns();updateSteps();initInsp();
  stepTimer=setTimeout(execStep,200);
}
function pause(){paused=true;if(stepTimer)clearTimeout(stepTimer);stepTimer=null;updateBtns()}
function resume(){paused=false;updateBtns();stepTimer=setTimeout(execStep,100)}
function reset(){if(stepTimer)clearTimeout(stepTimer);resetState();updateBtns();updateSteps();initInsp()}
function setFlow(k){activeFlow=k;document.getElementById('flow-sel').value=k;reset()}
function cycleSpeed(){if(speed<1)speed=1;else if(speed<2)speed=2;else speed=0.5;document.getElementById('btn-speed').textContent=speed+'x'}

// ========== 7. SNAPSHOT ==========
function saveSnap(idx){snapshots[idx]={an:JSON.parse(JSON.stringify(activeNodes)),bd:JSON.parse(JSON.stringify(badges)),ln:JSON.parse(JSON.stringify(lines)),gl:JSON.parse(JSON.stringify(glowing))}}
function restoreSnap(idx){if(!snapshots[idx])return false;var s=snapshots[idx];activeNodes=JSON.parse(JSON.stringify(s.an));badges=JSON.parse(JSON.stringify(s.bd));lines=JSON.parse(JSON.stringify(s.ln));glowing=JSON.parse(JSON.stringify(s.gl));dots=[];return true}
function jumpTo(idx){if(stepTimer)clearTimeout(stepTimer);dots=[];if(restoreSnap(idx)){stepIdx=idx+1;updateSteps();applyMutation(idx)}}

// ========== 8. STEPS PANEL ==========
function updateSteps(){
  var el=document.getElementById('steps-container');
  var flow=getFlow();
  if(!flow){el.innerHTML='<p style="color:var(--dim);padding:8px">Select a flow and press Start</p>';return}
  var h='';
  for(var i=0;i<flow.steps.length;i++){
    var s=flow.steps[i];
    var cls='step';
    if(i<stepIdx)cls+=' done';
    if(i===stepIdx&&running)cls+=' active';
    var mark=i<stepIdx?'✓':i===stepIdx&&running?'●':'○';
    h+='<div class="'+cls+'" data-i="'+i+'"><span class="step-mark">'+mark+'</span><span>'+esc(s.text||'')+'</span></div>';
  }
  el.innerHTML=h;
  var items=el.querySelectorAll('.step');
  for(var j=0;j<items.length;j++)items[j].addEventListener('click',(function(idx){return function(){jumpTo(idx)}})(j));
}

// ========== 9. INSPECTOR ==========
function initInsp(){
  if(!D.inspector){document.getElementById('insp-content').innerHTML='';return}
  var st=D.inspector.initialState||{};
  renderInsp(st.phase||'request','',st.headers||[],st.body||[]);
}

function applyMutation(idx){
  if(!D.inspector||!D.inspector.mutations)return;
  var muts=D.inspector.mutations[activeFlow];if(!muts)return;
  var mut=null;
  for(var i=0;i<muts.length;i++){if(muts[i].step===idx+1){mut=muts[i];break}}
  if(!mut)return;

  // Process actions
  if(mut.actions){
    for(var a=0;a<mut.actions.length;a++){
      var act=mut.actions[a];
      if(act.action==='add'){
        var line={value:act.value||'',style:act.style||'add',id:act.id||''};
        var target=act.target||'body';
        addInspLine(target,line);continue;
      }
      if(act.action==='remove'){removeInspLine(act.id);continue}
      if(act.id){styleInspLine(act.id,act.style||'highlight')}
    }
  }

  var ph=mut.phase||null;
  var headers=mut.replaceHeaders||null;
  var body=mut.replaceBody||null;
  if(ph||headers||body){
    var el=document.getElementById('insp-content');
    if(ph)document.getElementById('insp-title').textContent=ph.charAt(0).toUpperCase()+ph.slice(1);
    if(ph)document.getElementById('insp-title').style.color=ph==='response'?'var(--ok)':ph==='error'?'var(--err)':'var(--accent)';
    if(headers)replaceSection('headers',headers);
    if(body)replaceSection('body',body);
  }
  if(mut.label){
    var lbl=document.getElementById('insp-label');
    if(!lbl){lbl=document.createElement('div');lbl.id='insp-label';lbl.className='insp-label';
    var content=document.getElementById('insp-content');
    content.insertBefore(lbl,content.firstChild)}
    lbl.textContent=mut.label;
  }
}

function renderInsp(phase,label,headers,body){
  var el=document.getElementById('insp-title');
  el.textContent=phase.charAt(0).toUpperCase()+phase.slice(1);
  el.style.color=phase==='response'?'var(--ok)':phase==='error'?'var(--err)':'var(--accent)';
  var h='';
  if(label)h+='<div class="insp-label">'+esc(label)+'</div>';
  h+='<div class="insp-section">Headers</div><div id="insp-headers">';
  for(var i=0;i<headers.length;i++)h+='<div class="insp-line '+(headers[i].style||'keep')+'" id="il-'+(headers[i].id||i)+'">'+esc(headers[i].value)+'</div>';
  h+='</div>';
  if(body.length){
    h+='<div class="insp-section">Body</div><div id="insp-body">';
    for(var j=0;j<body.length;j++)h+='<div class="insp-line '+(body[j].style||'keep')+'" id="il-'+(body[j].id||'b'+j)+'">'+esc(body[j].value)+'</div>';
    h+='</div>';
  }
  document.getElementById('insp-content').innerHTML=h;
}

function replaceSection(name,lines2){
  var el=document.getElementById('insp-'+name);
  if(!el)return;
  var h='';
  for(var i=0;i<lines2.length;i++)h+='<div class="insp-line '+(lines2[i].style||'keep')+'" id="il-'+(lines2[i].id||name+i)+'">'+esc(lines2[i].value)+'</div>';
  el.innerHTML=h;
}

function addInspLine(target,line){
  var el=document.getElementById('insp-'+target);
  if(!el)return;
  var div=document.createElement('div');
  div.className='insp-line '+(line.style||'add');
  div.id='il-'+(line.id||'');
  div.textContent=line.value;
  el.appendChild(div);
}

function removeInspLine(id){var el=document.getElementById('il-'+id);if(el)el.remove()}
function styleInspLine(id,style){var el=document.getElementById('il-'+id);if(el){el.className='insp-line '+style}}

// ========== 10. OVERLAY ==========
function hitTest(mx,my){
  var keys=Object.keys(D.nodes||{});
  for(var i=keys.length-1;i>=0;i--){
    var n=D.nodes[keys[i]];
    if(n.type==='boundary'||n.type==='container')continue;
    var x=tx(n.x),y=ty(n.y),w=ts(n.w),h=ts(n.h);
    if(mx>=x&&mx<=x+w&&my>=y&&my<=y+h)return keys[i];
  }
  return null;
}

function showOverlay(key){
  var tt=D.tooltips?D.tooltips[key]:null;
  if(!tt)return;
  var n=D.nodes[key];
  document.getElementById('ov-title').textContent=tt.title||key;
  document.getElementById('ov-desc').textContent=tt.description||'';
  document.getElementById('ov-accent').style.background=n?n.color||'var(--accent)':'var(--accent)';
  var dh='';
  if(tt.details){for(var i=0;i<tt.details.length;i++){
    dh+='<div class="detail"><span class="dk">'+esc(tt.details[i][0])+'</span><span class="dv">'+esc(tt.details[i][1])+'</span></div>';
  }}
  document.getElementById('ov-details').innerHTML=dh;
  document.getElementById('overlay').style.display='block';
  // Position card near node
  if(n){
    var card=document.getElementById('ov-card');
    card.style.left=Math.min(tx(n.x+n.w)+20,W-440)+'px';
    card.style.top=Math.max(ty(n.y)-20,10)+'px';
  }
  if(running&&!paused)pause();
}

function hideOverlay(){document.getElementById('overlay').style.display='none'}

// ========== 11. INIT ==========
function esc(s){return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;').replace(/'/g,'&#39;')}

function updateBtns(){
  var pb=document.getElementById('btn-play');
  pb.textContent=paused?'▶ Resume':running?'⏸ Pause':'▶ Start';
  pb.className='btn-start';
}

function init(){
  resize();
  addEventListener('resize',resize);

  // Title
  if(D.meta&&D.meta.title)document.getElementById('diagram-title').textContent=D.meta.title;

  // Flow selector
  var sel=document.getElementById('flow-sel');
  var order=D.flowOrder||Object.keys(D.flows||{});
  for(var i=0;i<order.length;i++){
    var f=D.flows[order[i]];if(!f)continue;
    var opt=document.createElement('option');
    opt.value=order[i];opt.textContent=f.label||order[i];
    sel.appendChild(opt);
  }
  if(!activeFlow&&order.length)activeFlow=order[0];
  sel.value=activeFlow;

  // Legend
  var leg=document.getElementById('legend');
  if(D.legend){var lh='';for(var i=0;i<D.legend.length;i++){
    lh+='<div class="legend-item"><div class="legend-dot" style="background:'+D.legend[i].color+'"></div>'+esc(D.legend[i].label)+'</div>';
  }leg.innerHTML=lh}

  // Controls
  document.getElementById('btn-play').addEventListener('click',run);
  document.getElementById('btn-reset').addEventListener('click',reset);
  document.getElementById('btn-speed').addEventListener('click',cycleSpeed);
  document.getElementById('btn-loop').addEventListener('click',function(){
    loopMode=!loopMode;
    this.textContent=loopMode?'↻ Loop ON':'↻ Loop';
    this.style.borderColor=loopMode?'var(--ok)':'var(--accent)';
    this.style.color=loopMode?'var(--ok)':'var(--accent)';
  });
  document.getElementById('flow-sel').addEventListener('change',function(){setFlow(this.value)});
  document.getElementById('theme-btn').addEventListener('click',toggleTheme);
  document.getElementById('ov-close').addEventListener('click',hideOverlay);
  document.getElementById('ov-resume').addEventListener('click',function(){hideOverlay();if(paused)resume()});

  // Canvas click
  cv.addEventListener('click',function(e){
    var r=cv.getBoundingClientRect();
    var key=hitTest(e.clientX-r.left,e.clientY-r.top);
    if(key)showOverlay(key);
  });

  // Keyboard
  addEventListener('keydown',function(e){
    if(e.key===' '){e.preventDefault();run()}
    if(e.key==='r')reset();
    if(e.key==='t')toggleTheme();
  });

  updateSteps();initInsp();

  // Render loop
  function frame(){renderAll();requestAnimationFrame(frame)}
  requestAnimationFrame(frame);
}

if(document.readyState==='complete')init();
else addEventListener('load',init);
})();
</script>
</body>
</html>`
