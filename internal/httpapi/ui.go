package httpapi

import "net/http"

func serveUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(uiHTML))
}

const uiHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>FlashDeal API — ELK Observability Demo</title>
  <style>
    @import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap');

    :root {
      --bg:     #0d1117;
      --bg2:    #161b22;
      --bg3:    #21262d;
      --border: #30363d;
      --text:   #e6edf3;
      --muted:  #7d8590;
      --accent: #2f81f7;
      --green:  #3fb950;
      --yellow: #d29922;
      --red:    #f85149;
      --purple: #a371f7;
      --orange: #f0883e;
      --teal:   #39d353;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body { background: var(--bg); color: var(--text); font-family: Inter, sans-serif; line-height: 1.6; }

    nav {
      background: var(--bg2);
      border-bottom: 1px solid var(--border);
      padding: 14px 32px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      position: sticky; top: 0; z-index: 10;
    }
    .nav-brand { font-weight: 600; font-size: 15px; display: flex; align-items: center; gap: 10px; }
    .live-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--red); box-shadow: 0 0 8px var(--red); animation: pulse 1.5s infinite; }
    @keyframes pulse { 0%,100%{opacity:1}50%{opacity:0.3} }
    .nav-links { display: flex; gap: 20px; }
    .nav-links a { color: var(--muted); text-decoration: none; font-size: 13px; transition: color 0.2s; }
    .nav-links a:hover { color: var(--text); }

    .hero { text-align: center; padding: 70px 24px 50px; border-bottom: 1px solid var(--border); }
    .hero-badge {
      display: inline-flex; align-items: center; gap: 6px;
      background: rgba(248,81,73,0.1); border: 1px solid rgba(248,81,73,0.3);
      color: var(--red); font-size: 11px; padding: 4px 12px; border-radius: 20px;
      letter-spacing: 1px; text-transform: uppercase; margin-bottom: 20px;
      font-family: 'JetBrains Mono', monospace;
    }
    h1 { font-size: clamp(26px, 5vw, 42px); font-weight: 700; margin-bottom: 12px; letter-spacing: -0.5px; }
    h1 span { color: var(--orange); }
    .hero p { color: var(--muted); max-width: 580px; margin: 0 auto 28px; font-size: 15px; }
    .badges { display: flex; flex-wrap: wrap; gap: 8px; justify-content: center; }
    .badge { font-family: 'JetBrains Mono', monospace; font-size: 11px; padding: 4px 12px; border-radius: 4px; border: 1px solid; }
    .badge-blue   { color: var(--accent);  border-color: rgba(47,129,247,0.4);  background: rgba(47,129,247,0.08); }
    .badge-green  { color: var(--green);   border-color: rgba(63,185,80,0.4);   background: rgba(63,185,80,0.08); }
    .badge-purple { color: var(--purple);  border-color: rgba(163,113,247,0.4); background: rgba(163,113,247,0.08); }
    .badge-orange { color: var(--orange);  border-color: rgba(240,136,62,0.4);  background: rgba(240,136,62,0.08); }
    .badge-red    { color: var(--red);     border-color: rgba(248,81,73,0.4);   background: rgba(248,81,73,0.08); }
    .badge-yellow { color: var(--yellow);  border-color: rgba(210,153,34,0.4);  background: rgba(210,153,34,0.08); }

    .page { max-width: 1160px; margin: 0 auto; padding: 0 24px 60px; }
    .section { margin-top: 48px; }
    .section-label { font-size: 11px; letter-spacing: 2px; text-transform: uppercase; color: var(--muted); margin-bottom: 16px; }

    /* ARCH */
    .arch-flow {
      background: var(--bg2); border: 1px solid var(--border); border-radius: 12px;
      padding: 28px; display: flex; align-items: center; justify-content: center; flex-wrap: wrap; gap: 0;
    }
    .arch-node { text-align: center; padding: 14px 16px; background: var(--bg3); border: 1px solid var(--border); border-radius: 8px; min-width: 100px; }
    .arch-node .icon { font-size: 20px; margin-bottom: 4px; }
    .arch-node .name { font-size: 11px; font-weight: 600; }
    .arch-node .sub  { font-size: 10px; color: var(--muted); font-family: 'JetBrains Mono', monospace; }
    .arch-arrow { color: var(--muted); font-size: 18px; padding: 0 8px; flex-shrink: 0; }

    /* MAIN GRID */
    .main-grid { display: grid; grid-template-columns: 1fr 380px; gap: 20px; }
    @media(max-width:800px){ .main-grid { grid-template-columns: 1fr; } }

    .card { background: var(--bg2); border: 1px solid var(--border); border-radius: 12px; overflow: hidden; }
    .card-header {
      padding: 14px 18px; border-bottom: 1px solid var(--border);
      font-size: 13px; font-weight: 600; display: flex; align-items: center; gap: 8px;
      background: var(--bg3);
    }
    .pill { font-size: 10px; padding: 2px 8px; border-radius: 3px; font-family: 'JetBrains Mono', monospace; }
    .pill-green  { background: rgba(63,185,80,0.15);  color: var(--green); }
    .pill-blue   { background: rgba(47,129,247,0.15); color: var(--accent); }
    .pill-orange { background: rgba(240,136,62,0.15); color: var(--orange); }
    .pill-red    { background: rgba(248,81,73,0.15);  color: var(--red); }
    .card-body { padding: 18px; }

    input, select {
      width: 100%; background: var(--bg3); border: 1px solid var(--border); color: var(--text);
      padding: 10px 12px; border-radius: 6px; font-family: inherit; font-size: 13px;
      margin-bottom: 10px; outline: none; transition: border-color 0.2s;
    }
    input:focus, select:focus { border-color: var(--accent); }
    select option { background: var(--bg3); }

    .btn { display: inline-flex; align-items: center; gap: 6px; padding: 9px 18px; border-radius: 6px; font-size: 13px; font-weight: 500; cursor: pointer; border: none; transition: all 0.15s; font-family: inherit; width: 100%; justify-content: center; }
    .btn-primary { background: var(--accent); color: #fff; }
    .btn-primary:hover { background: #388bfd; box-shadow: 0 0 12px rgba(47,129,247,0.35); }
    .btn-orange { background: var(--orange); color: #000; margin-top: 4px; }
    .btn-orange:hover { opacity: 0.85; }

    /* DEALS */
    #deals-list { display: flex; flex-direction: column; gap: 10px; }
    .deal-card {
      background: var(--bg3); border: 1px solid var(--border); border-radius: 8px;
      padding: 14px 16px; display: flex; align-items: center; justify-content: space-between; gap: 12px;
    }
    .deal-card:hover { border-color: var(--orange); }
    .deal-id   { font-family: 'JetBrains Mono', monospace; font-size: 10px; color: var(--muted); }
    .deal-name { font-size: 14px; font-weight: 600; }
    .deal-price { font-size: 18px; font-weight: 700; color: var(--orange); }
    .deal-stock { font-size: 11px; color: var(--muted); }
    .deal-badge { background: rgba(63,185,80,0.15); color: var(--green); font-size: 10px; padding: 2px 8px; border-radius: 10px; border: 1px solid rgba(63,185,80,0.3); }

    /* SSE FEED */
    #sse-feed {
      background: #010409; border: 1px solid var(--border); border-radius: 8px;
      padding: 14px; height: 300px; overflow-y: auto;
      font-family: 'JetBrains Mono', monospace; font-size: 11px;
      display: flex; flex-direction: column; gap: 6px;
    }
    .sse-entry { display: flex; flex-direction: column; gap: 2px; border-left: 2px solid; padding-left: 8px; }
    .sse-entry.deal  { border-color: var(--orange); }
    .sse-entry.order { border-color: var(--accent); }
    .sse-event { font-weight: 700; }
    .sse-event.deal  { color: var(--orange); }
    .sse-event.order { color: var(--accent); }
    .sse-data { color: var(--muted); white-space: pre-wrap; word-break: break-all; font-size: 10px; }
    .sse-ts { color: #3d444d; font-size: 10px; }
    .sse-status { color: var(--muted); font-style: italic; }

    /* SKILLS */
    .skills-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 14px; }
    .skill-card { background: var(--bg2); border: 1px solid var(--border); border-radius: 10px; padding: 18px; transition: border-color 0.2s; }
    .skill-card:hover { border-color: var(--accent); }
    .skill-icon { font-size: 22px; margin-bottom: 8px; }
    .skill-name { font-size: 13px; font-weight: 600; margin-bottom: 4px; }
    .skill-desc { font-size: 12px; color: var(--muted); line-height: 1.5; }

    footer { border-top: 1px solid var(--border); padding: 24px; text-align: center; font-size: 12px; color: var(--muted); }
    footer a { color: var(--accent); text-decoration: none; }

    .empty-state { color: var(--muted); font-size: 13px; text-align: center; padding: 24px; }
    .label { font-size: 12px; color: var(--muted); margin-bottom: 4px; }
  </style>
</head>
<body>

<nav>
  <div class="nav-brand">
    <div class="live-dot"></div>
    FlashDeal API
  </div>
  <div class="nav-links">
    <a href="#demo">Live Demo</a>
    <a href="#architecture">Architecture</a>
    <a href="#skills">Skills</a>
    <a href="https://github.com/Gbolahan-Aziz/FlashDeal-API-ELK" target="_blank">GitHub</a>
  </div>
</nav>

<div class="hero">
  <div class="hero-badge">&#9679; Live</div>
  <h1>FlashDeal API<br>with <span>ELK Observability</span></h1>
  <p>A real-time flash sales API in Go with Server-Sent Events. Full log observability via Elasticsearch, Logstash, and Kibana running in Docker Compose.</p>
  <div class="badges">
    <span class="badge badge-blue">Go 1.22</span>
    <span class="badge badge-orange">Elasticsearch</span>
    <span class="badge badge-purple">Logstash</span>
    <span class="badge badge-red">Kibana</span>
    <span class="badge badge-green">Filebeat</span>
    <span class="badge badge-blue">Docker Compose</span>
    <span class="badge badge-yellow">Server-Sent Events</span>
    <span class="badge badge-orange">SQLite</span>
  </div>
</div>

<div class="page">

  <!-- ARCHITECTURE -->
  <div class="section" id="architecture">
    <div class="section-label">Observability Pipeline</div>
    <div class="arch-flow">
      <div class="arch-node">
        <div class="icon">⚡</div>
        <div class="name">Go API</div>
        <div class="sub">Zerolog JSON</div>
      </div>
      <div class="arch-arrow">→</div>
      <div class="arch-node">
        <div class="icon">📦</div>
        <div class="name">Filebeat</div>
        <div class="sub">Docker autodiscover</div>
      </div>
      <div class="arch-arrow">→</div>
      <div class="arch-node">
        <div class="icon">🔧</div>
        <div class="name">Logstash</div>
        <div class="sub">Parse + enrich</div>
      </div>
      <div class="arch-arrow">→</div>
      <div class="arch-node">
        <div class="icon">🔍</div>
        <div class="name">Elasticsearch</div>
        <div class="sub">Index + search</div>
      </div>
      <div class="arch-arrow">→</div>
      <div class="arch-node">
        <div class="icon">📊</div>
        <div class="name">Kibana</div>
        <div class="sub">Dashboard</div>
      </div>
    </div>
    <p style="text-align:center;color:var(--muted);font-size:12px;margin-top:10px;">
      Every API request produces structured JSON logs. Filebeat collects them via Docker autodiscover, Logstash parses and enriches them, and Kibana visualises latency, status codes, and request volume in real time.
    </p>
  </div>

  <!-- LIVE DEMO -->
  <div class="section" id="demo">
    <div class="section-label">Live Demo — Try it now</div>
    <div class="main-grid">

      <!-- LEFT COLUMN -->
      <div style="display:flex;flex-direction:column;gap:18px;">

        <!-- CREATE DEAL -->
        <div class="card">
          <div class="card-header">
            Create Flash Deal
            <span class="pill pill-orange">POST /deals</span>
          </div>
          <div class="card-body">
            <div class="label">Title</div>
            <input id="deal-title" type="text" placeholder="e.g. MacBook Air M3 — 40% off">
            <div style="display:grid;grid-template-columns:1fr 1fr;gap:10px;">
              <div>
                <div class="label">Price ($)</div>
                <input id="deal-price" type="number" placeholder="299.99" step="0.01">
              </div>
              <div>
                <div class="label">Stock</div>
                <input id="deal-stock" type="number" placeholder="50" min="1">
              </div>
            </div>
            <button class="btn btn-orange" onclick="createDeal()">Launch Deal</button>
          </div>
        </div>

        <!-- DEALS BOARD -->
        <div class="card">
          <div class="card-header">
            Active Deals
            <span class="pill pill-blue">GET /deals</span>
            <button onclick="loadDeals()" style="margin-left:auto;background:var(--bg);border:1px solid var(--border);color:var(--muted);padding:4px 10px;border-radius:4px;cursor:pointer;font-size:11px;">Refresh</button>
          </div>
          <div class="card-body">
            <div id="deals-list"><div class="empty-state">No deals yet. Launch one above.</div></div>
          </div>
        </div>

        <!-- CREATE ORDER -->
        <div class="card">
          <div class="card-header">
            Place Order
            <span class="pill pill-green">POST /orders</span>
          </div>
          <div class="card-body">
            <div class="label">Deal</div>
            <select id="order-deal">
              <option value="">Select a deal...</option>
            </select>
            <div class="label">Quantity</div>
            <input id="order-qty" type="number" placeholder="1" min="1" value="1">
            <button class="btn btn-primary" onclick="createOrder()">Place Order</button>
          </div>
        </div>

      </div>

      <!-- RIGHT COLUMN: SSE FEED -->
      <div style="display:flex;flex-direction:column;gap:18px;">
        <div class="card" style="flex:1;">
          <div class="card-header">
            <div class="live-dot"></div>
            Real-time Event Feed
            <span class="pill pill-red">GET /events (SSE)</span>
          </div>
          <div class="card-body" style="padding:0">
            <div id="sse-feed">
              <span class="sse-status">Connecting to event stream...</span>
            </div>
          </div>
        </div>
        <div class="card">
          <div class="card-header" style="font-size:12px;">How SSE works here</div>
          <div class="card-body" style="font-size:12px;color:var(--muted);line-height:1.7;">
            When you create a deal or place an order, the Go API broadcasts a JSON event to all connected clients via <code style="color:var(--orange);background:var(--bg3);padding:1px 5px;border-radius:3px;">GET /events</code>.
            The feed on the left updates in real time — no polling, no WebSocket, just a plain HTTP stream.
            In the full Docker Compose stack, every event also flows through Filebeat into Elasticsearch.
          </div>
        </div>
      </div>

    </div>
  </div>

  <!-- SKILLS -->
  <div class="section" id="skills">
    <div class="section-label">DevOps Skills Demonstrated</div>
    <div class="skills-grid">
      <div class="skill-card">
        <div class="skill-icon">📊</div>
        <div class="skill-name">ELK Stack Observability</div>
        <div class="skill-desc">Full log pipeline: structured JSON from Zerolog, collected by Filebeat with Docker autodiscover, parsed by Logstash, indexed in Elasticsearch, visualised in Kibana.</div>
      </div>
      <div class="skill-card">
        <div class="skill-icon">⚡</div>
        <div class="skill-name">High-Performance Go</div>
        <div class="skill-desc">Concurrent API built with Chi router, interface-driven design, graceful shutdown, panic recovery, and dual-mode storage (in-memory or SQLite).</div>
      </div>
      <div class="skill-card">
        <div class="skill-icon">📡</div>
        <div class="skill-name">Real-time SSE</div>
        <div class="skill-desc">Thread-safe SSE hub with RWMutex, non-blocking broadcast, 15-second heartbeat keepalive, and client disconnection cleanup via context cancellation.</div>
      </div>
      <div class="skill-card">
        <div class="skill-icon">🐳</div>
        <div class="skill-name">Docker Compose Orchestration</div>
        <div class="skill-desc">Five-service stack: Go API, Elasticsearch, Kibana, Logstash, Filebeat. Custom bridge network, persistent volumes, health checks, and startup ordering.</div>
      </div>
      <div class="skill-card">
        <div class="skill-icon">🔐</div>
        <div class="skill-name">Container Security</div>
        <div class="skill-desc">Multi-stage Docker build from golang:1.22 to debian:bookworm-slim. Non-root appuser (UID 1001), minimal runtime image, CA certificate inclusion.</div>
      </div>
      <div class="skill-card">
        <div class="skill-icon">🔧</div>
        <div class="skill-name">Logstash Pipeline Design</div>
        <div class="skill-desc">Logstash filter extracts latency_ms and HTTP status fields, applies type conversion, parses ISO8601 timestamps, and writes to date-partitioned Elasticsearch indices.</div>
      </div>
      <div class="skill-card">
        <div class="skill-icon">📋</div>
        <div class="skill-name">Structured Logging</div>
        <div class="skill-desc">Zerolog produces consistent JSON logs with service name, environment, request ID, latency, and status fields on every request. Makes log parsing deterministic.</div>
      </div>
      <div class="skill-card">
        <div class="skill-icon">🏗️</div>
        <div class="skill-name">Clean Architecture</div>
        <div class="skill-desc">Repository pattern (Store interface), dependency injection (Deps struct), publisher-subscriber (events), and hub pattern for SSE. Kafka-ready event bus.</div>
      </div>
    </div>
  </div>

</div>

<footer>
  Built by <a href="https://github.com/Gbolahan-Aziz" target="_blank">Azeez Razaq</a>
  &nbsp;·&nbsp;
  <a href="https://github.com/Gbolahan-Aziz/FlashDeal-API-ELK" target="_blank">GitHub Repo</a>
  &nbsp;·&nbsp;
  <a href="/healthz">Health Check</a>
</footer>

<script>
  let deals = [];

  // ── SSE ──────────────────────────────────────────────────────
  function connectSSE() {
    const feed = document.getElementById('sse-feed');
    feed.innerHTML = '';
    const es = new EventSource('/events');

    es.onopen = () => {
      const el = document.createElement('div');
      el.className = 'sse-status';
      el.textContent = 'Connected. Waiting for events...';
      feed.appendChild(el);
    };

    es.onmessage = (e) => {
      try {
        const d = JSON.parse(e.data);
        if (d.type === 'ping') return;
        const isOrder = d.event === 'order.created';
        const div = document.createElement('div');
        div.className = 'sse-entry ' + (isOrder ? 'order' : 'deal');
        div.innerHTML =
          '<span class="sse-ts">' + new Date().toLocaleTimeString() + '</span>' +
          '<span class="sse-event ' + (isOrder ? 'order' : 'deal') + '">' + (d.event || 'event') + '</span>' +
          '<span class="sse-data">' + JSON.stringify(d.payload || d, null, 2) + '</span>';
        feed.appendChild(div);
        feed.scrollTop = feed.scrollHeight;
        if (!isOrder) loadDeals();
      } catch (_) {}
    };

    es.onerror = () => {
      setTimeout(connectSSE, 3000);
    };
  }

  // ── API CALLS ─────────────────────────────────────────────────
  async function loadDeals() {
    const res = await fetch('/deals');
    deals = await res.json() || [];
    renderDeals();
    renderDealSelect();
  }

  function renderDeals() {
    const el = document.getElementById('deals-list');
    if (!deals.length) {
      el.innerHTML = '<div class="empty-state">No deals yet. Launch one above.</div>';
      return;
    }
    el.innerHTML = deals.map(d => ` + "`" + `
      <div class="deal-card">
        <div>
          <div class="deal-id">${d.id ? d.id.substring(0,8) : ''}...</div>
          <div class="deal-name">${esc(d.title)}</div>
          <div class="deal-stock">${d.stock} units remaining</div>
        </div>
        <div style="text-align:right">
          <div class="deal-price">$${Number(d.price).toFixed(2)}</div>
          <div class="deal-badge">ACTIVE</div>
        </div>
      </div>
    ` + "`").join('');
  }

  function renderDealSelect() {
    const sel = document.getElementById('order-deal');
    sel.innerHTML = '<option value="">Select a deal...</option>' +
      deals.map(d => '<option value="' + d.id + '">' + esc(d.title) + ' — $' + Number(d.price).toFixed(2) + '</option>').join('');
  }

  async function createDeal() {
    const title = document.getElementById('deal-title').value.trim();
    const price = parseFloat(document.getElementById('deal-price').value);
    const stock = parseInt(document.getElementById('deal-stock').value);
    if (!title || isNaN(price) || isNaN(stock)) { alert('Fill in all fields'); return; }
    const res = await fetch('/deals', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title, price, stock })
    });
    if (res.ok) {
      document.getElementById('deal-title').value = '';
      document.getElementById('deal-price').value = '';
      document.getElementById('deal-stock').value = '';
      loadDeals();
    } else { alert('Failed: ' + res.status); }
  }

  async function createOrder() {
    const deal_id = document.getElementById('order-deal').value;
    const qty = parseInt(document.getElementById('order-qty').value);
    if (!deal_id) { alert('Select a deal'); return; }
    const res = await fetch('/orders', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ deal_id, qty: qty || 1 })
    });
    if (res.ok) { loadDeals(); }
    else if (res.status === 409) { alert('Insufficient stock'); }
    else if (res.status === 404) { alert('Deal not found'); }
    else { alert('Error: ' + res.status); }
  }

  function esc(s) { return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); }

  connectSSE();
  loadDeals();
</script>
</body>
</html>`
