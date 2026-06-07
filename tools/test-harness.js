const fs = require('fs');
const path = require('path');
const { spawnSync } = require('child_process');

function parseArgs() {
  const args = process.argv.slice(2);
  const options = {
    script: 'game/smoke.yaml',
    out: 'test-logs/shots',
    help: false
  };

  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--script' || args[i] === '-s') {
      options.script = args[++i];
    } else if (args[i] === '--out' || args[i] === '-o') {
      options.out = args[++i];
    } else if (args[i] === '--help' || args[i] === '-h') {
      options.help = true;
    }
  }
  return options;
}

function printHelp() {
  console.log(`Usage: node tools/test-harness.js [options]
Options:
  -s, --script <path>   Path to the smoke YAML script (default: game/smoke.yaml)
  -o, --out <dir>       Output directory for logs and shots (default: test-logs/shots)
  -h, --help            Show this help message
`);
}

function ensureDir(dir) {
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }
}

function cleanDir(dir) {
  if (fs.existsSync(dir)) {
    fs.rmSync(dir, { recursive: true, force: true });
  }
  ensureDir(dir);
}

function runCommand(cmd, args, options = {}) {
  console.log(`[Harness] Running: ${cmd} ${args.join(' ')}`);
  const localBin = '/home/luke/.local/bin';
  const goBin = '/home/luke/.local/go/bin';
  const tinygoBin = '/home/luke/.local/tinygo/bin';
  const customPath = `${localBin}:${goBin}:${tinygoBin}:${process.env.PATH}`;

  const result = spawnSync(cmd, args, {
    stdio: 'inherit',
    ...options,
    env: { 
      ...process.env, 
      PATH: customPath, 
      GOTOOLCHAIN: 'local',
      ...options.env 
    }
  });
  return result.status === 0;
}

function buildWasm() {
  console.log('\n[Harness] Phase 1: Building WASM Artifact...');
  const success = runCommand('tinygo', [
    'build',
    '-opt=1',
    '-no-debug',
    '-gc=leaking',
    '-o',
    'game.wasm',
    '-target',
    'wasip1',
    '-buildmode=c-shared',
    '.'
  ], { cwd: path.resolve('game') });
  
  if (success) {
    console.log('[Harness] WASM build successful: game/game.wasm');
  } else {
    console.error('[Harness] WASM build failed.');
  }
  return success;
}

function runSmokeTests(scriptPath, outDir) {
  console.log('\n[Harness] Phase 2: Running Smoke Tests...');
  const absScriptPath = path.resolve(scriptPath);
  const absOutDir = path.resolve(outDir);

  if (!fs.existsSync(absScriptPath)) {
    console.error(`[Harness] Smoke script not found at ${absScriptPath}`);
    return false;
  }

  // Ensure output directory is clean
  cleanDir(absOutDir);

  const success = runCommand('go', [
    'run',
    '.',
    '-smoke',
    absScriptPath,
    '-smoke-out',
    absOutDir
  ], { cwd: path.resolve('game') });

  if (success) {
    console.log(`[Harness] Smoke tests completed successfully. Shots written to ${outDir}`);
  } else {
    console.error('[Harness] Smoke test execution failed.');
  }
  return success;
}

function escapeHtml(unsafe) {
  return unsafe
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}

function generatePlaybackHtml(outDir, scriptPath) {
  console.log('\n[Harness] Phase 3: Generating Playback Dashboard...');
  
  const files = fs.readdirSync(outDir);
  const ansiFiles = files.filter(f => f.endsWith('.ansi')).sort();
  
  if (ansiFiles.length === 0) {
    console.error('[Harness] No ANSI shot files found. Cannot generate playback.');
    return;
  }

  // Parse shot names, ordinals, and seat indices
  // Format: NN-name.ansi or NN-name.seatK.ansi
  const shotsMap = new Map();
  const seatsSet = new Set([0]); // Always show at least seat 0

  for (const file of ansiFiles) {
    const filePath = path.join(outDir, file);
    const textPath = filePath.replace('.ansi', '.txt');
    const ansiContent = fs.readFileSync(filePath, 'utf8');
    const textContent = fs.existsSync(textPath) ? fs.readFileSync(textPath, 'utf8') : '';

    const matchSeat = file.match(/^(\d+)-([^\.]+)\.seat(\d+)\.ansi$/);
    const matchCollapsed = file.match(/^(\d+)-([^\.]+)\.ansi$/);

    let ordinal, name, seat, collapsed;

    if (matchSeat) {
      ordinal = parseInt(matchSeat[1], 10);
      name = matchSeat[2];
      seat = parseInt(matchSeat[3], 10);
      collapsed = false;
      seatsSet.add(seat);
    } else if (matchCollapsed) {
      ordinal = parseInt(matchCollapsed[1], 10);
      name = matchCollapsed[2];
      seat = 0; // Collapsed frame is assigned to seat 0
      collapsed = true;
    } else {
      continue;
    }

    if (!shotsMap.has(ordinal)) {
      shotsMap.set(ordinal, {
        ordinal,
        name,
        collapsed,
        seats: {}
      });
    }

    const shot = shotsMap.get(ordinal);
    shot.seats[seat] = {
      ansi: ansiContent,
      txt: textContent
    };
  }

  const shots = Array.from(shotsMap.values()).sort((a, b) => a.ordinal - b.ordinal);
  const seats = Array.from(seatsSet).sort((a, b) => a - b);
  const scriptContent = fs.existsSync(scriptPath) ? fs.readFileSync(scriptPath, 'utf8') : '';

  const htmlContent = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Shellcade Game Test Playback</title>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;600;800&family=Outfit:wght@400;700&family=Fira+Code:wght@400;500&display=swap" rel="stylesheet">
  <style>
    :root {
      --bg-dark: #0b0f19;
      --bg-card: #151d30;
      --border-color: #243557;
      --text-main: #f3f4f6;
      --text-muted: #9ca3af;
      --accent: #3b82f6;
      --accent-glow: rgba(59, 130, 246, 0.4);
      --terminal-bg: #090d16;
    }

    * {
      box-sizing: border-box;
      margin: 0;
      padding: 0;
    }

    body {
      background-color: var(--bg-dark);
      color: var(--text-main);
      font-family: 'Inter', sans-serif;
      padding: 2rem;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
    }

    header {
      width: 100%;
      max-width: 1200px;
      margin-bottom: 2rem;
      border-bottom: 1px solid var(--border-color);
      padding-bottom: 1.5rem;
      display: flex;
      justify-content: space-between;
      align-items: flex-end;
    }

    h1 {
      font-family: 'Outfit', sans-serif;
      font-size: 2.5rem;
      font-weight: 700;
      background: linear-gradient(135deg, #60a5fa, #3b82f6);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      margin-bottom: 0.5rem;
    }

    .meta-text {
      color: var(--text-muted);
      font-size: 0.95rem;
    }

    .container {
      display: grid;
      grid-template-columns: 1.6fr 1fr;
      gap: 2rem;
      width: 100%;
      max-width: 1200px;
    }

    .panel {
      background-color: var(--bg-card);
      border: 1px solid var(--border-color);
      border-radius: 12px;
      padding: 1.5rem;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
    }

    .terminal-container {
      display: flex;
      flex-direction: column;
      align-items: center;
    }

    .terminal-window {
      width: 100%;
      background-color: var(--terminal-bg);
      border: 1px solid var(--border-color);
      border-radius: 8px;
      overflow: hidden;
      box-shadow: 0 0 15px var(--accent-glow);
    }

    .terminal-header {
      background-color: #111827;
      padding: 0.5rem 1rem;
      display: flex;
      align-items: center;
      gap: 6px;
      border-bottom: 1px solid var(--border-color);
    }

    .dot {
      width: 12px;
      height: 12px;
      border-radius: 50%;
    }
    .dot-red { background-color: #ef4444; }
    .dot-yellow { background-color: #f59e0b; }
    .dot-green { background-color: #10b981; }

    .terminal-title {
      margin-left: 1rem;
      color: var(--text-muted);
      font-size: 0.85rem;
      font-family: 'Fira Code', monospace;
    }

    .terminal-screen {
      padding: 1rem;
      font-family: 'Fira Code', monospace;
      font-size: 14px;
      line-height: 1.25;
      overflow-x: auto;
      white-space: pre;
    }

    .controls {
      margin-top: 1.5rem;
      width: 100%;
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .slider-row {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .timeline {
      flex-grow: 1;
      accent-color: var(--accent);
      cursor: pointer;
    }

    .btn-row {
      display: flex;
      justify-content: center;
      gap: 0.5rem;
    }

    button {
      background-color: #1e293b;
      border: 1px solid var(--border-color);
      color: var(--text-main);
      padding: 0.6rem 1.2rem;
      border-radius: 6px;
      font-weight: 600;
      cursor: pointer;
      transition: all 0.2s ease;
    }

    button:hover {
      background-color: var(--accent);
      border-color: var(--accent);
      box-shadow: 0 0 8px var(--accent-glow);
    }

    button.active {
      background-color: var(--accent);
      border-color: var(--accent);
    }

    .info-panel {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }

    .info-title {
      font-family: 'Outfit', sans-serif;
      font-size: 1.25rem;
      font-weight: 600;
      border-bottom: 1px solid var(--border-color);
      padding-bottom: 0.5rem;
    }

    .script-view {
      background-color: var(--terminal-bg);
      border: 1px solid var(--border-color);
      border-radius: 8px;
      padding: 1rem;
      font-family: 'Fira Code', monospace;
      font-size: 0.85rem;
      max-height: 250px;
      overflow-y: auto;
      white-space: pre-wrap;
    }

    .seat-selector {
      display: flex;
      gap: 0.5rem;
      margin-bottom: 1rem;
    }

    .speed-selector {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.9rem;
      color: var(--text-muted);
    }

    .speed-btn {
      padding: 0.3rem 0.6rem;
      font-size: 0.8rem;
    }
  </style>
</head>
<body>

  <header>
    <div>
      <h1>Shellcade Test Playback</h1>
      <div class="meta-text">Script: ${path.basename(scriptPath)} | Total Steps: ${shots.length}</div>
    </div>
    <div class="meta-text" id="status-display">Shot 1 of ${shots.length}</div>
  </header>

  <div class="container">
    <!-- Left Panel: Terminal & Controls -->
    <div class="panel terminal-container">
      <!-- Seat Selector if multiple seats -->
      <div class="seat-selector" id="seat-selector-container">
        ${seats.map(s => `<button class="seat-btn ${s === 0 ? 'active' : ''}" onclick="selectSeat(${s})">Seat ${s}</button>`).join('')}
      </div>

      <div class="terminal-window">
        <div class="terminal-header">
          <div class="dot dot-red"></div>
          <div class="dot dot-yellow"></div>
          <div class="dot dot-green"></div>
          <div class="terminal-title" id="terminal-title">80x24 Canvas</div>
        </div>
        <div class="terminal-screen" id="screen"></div>
      </div>

      <div class="controls">
        <div class="slider-row">
          <span id="curr-time">Step 0</span>
          <input type="range" class="timeline" id="timeline" min="0" max="${shots.length - 1}" value="0" oninput="seekTo(this.value)">
          <span id="max-time">Step ${shots.length - 1}</span>
        </div>

        <div class="btn-row">
          <button onclick="prevStep()">Prev</button>
          <button id="play-btn" onclick="togglePlay()">Play</button>
          <button onclick="nextStep()">Next</button>
        </div>

        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div class="speed-selector">
            <span>Speed:</span>
            <button class="speed-btn active" id="speed-1x" onclick="setSpeed(1, '1x')">1x</button>
            <button class="speed-btn" id="speed-2x" onclick="setSpeed(2, '2x')">2x</button>
            <button class="speed-btn" id="speed-05x" onclick="setSpeed(0.5, '0.5x')">0.5x</button>
          </div>
          <button onclick="toggleViewMode()" id="view-mode-btn">View Code / Text</button>
        </div>
      </div>
    </div>

    <!-- Right Panel: Info & Script -->
    <div class="panel info-panel">
      <div>
        <div class="info-title">Active Step Details</div>
        <div style="margin-top: 1rem;">
          <p><strong>Shot Name:</strong> <span id="detail-name">-</span></p>
          <p style="margin-top: 0.5rem;"><strong>Format:</strong> <span id="detail-format">-</span></p>
        </div>
      </div>

      <div>
        <div class="info-title">Smoke Script</div>
        <div class="script-view" style="margin-top: 1rem;">${escapeHtml(scriptContent)}</div>
      </div>
    </div>
  </div>

  <script>
    const shots = ${JSON.stringify(shots)};
    let currentStep = 0;
    let currentSeat = 0;
    let isPlaying = false;
    let playInterval = null;
    let speedMultiplier = 1;
    let viewMode = 'ansi'; // 'ansi' or 'txt'

    function escapeHtml(unsafe) {
      return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
    }

    function ansiToHtml(ansiStr) {
      const parts = ansiStr.split(/\\x1b\\[([0-9;]*)m/);
      let html = '';
      
      let bold = false;
      let dim = false;
      let underline = false;
      let reverse = false;
      let fg = null;
      let bg = null;
      
      function getSpan() {
        let styles = [];
        
        let currentFg = fg;
        let currentBg = bg;
        if (reverse) {
          currentFg = bg || 'rgb(0,0,0)';
          currentBg = fg || 'rgb(255,255,255)';
        }
        
        if (bold) styles.push('font-weight: bold');
        if (dim) styles.push('opacity: 0.6');
        if (underline) styles.push('text-decoration: underline');
        if (currentFg) styles.push('color: ' + currentFg);
        if (currentBg) styles.push('background-color: ' + currentBg);
        
        return styles.length ? '<span style="' + styles.join(';') + '">' : '';
      }
      
      for (let i = 0; i < parts.length; i++) {
        if (i % 2 === 0) {
          let text = parts[i]
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/ /g, '&nbsp;');
          
          const span = getSpan();
          if (span) {
            html += span + text + '</span>';
          } else {
            html += text;
          }
        } else {
          const codes = parts[i].split(';').map(Number);
          if (codes.length === 0 || (codes.length === 1 && codes[0] === 0)) {
            bold = false;
            dim = false;
            underline = false;
            reverse = false;
            fg = null;
            bg = null;
          } else {
            for (let j = 0; j < codes.length; j++) {
              const code = codes[j];
              if (code === 0) {
                bold = false;
                dim = false;
                underline = false;
                reverse = false;
                fg = null;
                bg = null;
              } else if (code === 1) {
                bold = true;
              } else if (code === 2) {
                dim = true;
              } else if (code === 4) {
                underline = true;
              } else if (code === 7) {
                reverse = true;
              } else if (code === 38 && codes[j+1] === 2) {
                fg = 'rgb(' + codes[j+2] + ',' + codes[j+3] + ',' + codes[j+4] + ')';
                j += 4;
              } else if (code === 48 && codes[j+1] === 2) {
                bg = 'rgb(' + codes[j+2] + ',' + codes[j+3] + ',' + codes[j+4] + ')';
                j += 4;
              }
            }
          }
        }
      }
      return html;
    }

    function renderCurrentFrame() {
      const shot = shots[currentStep];
      if (!shot) return;

      // Update metadata
      document.getElementById('status-display').innerText = 'Shot ' + (currentStep + 1) + ' of ' + shots.length;
      document.getElementById('curr-time').innerText = 'Step ' + currentStep;
      document.getElementById('timeline').value = currentStep;
      document.getElementById('detail-name').innerText = shot.name;
      document.getElementById('detail-format').innerText = shot.collapsed ? 'Collapsed (Identical for all seats)' : 'Seat-specific';

      // Pick frame
      let seatData = shot.seats[currentSeat] || shot.seats[0]; // fallback to seat 0
      if (!seatData) {
        // Find first available seat
        const availableSeats = Object.keys(shot.seats);
        if (availableSeats.length > 0) {
          seatData = shot.seats[availableSeats[0]];
        }
      }

      const screen = document.getElementById('screen');
      if (!seatData) {
        screen.innerHTML = 'No data for this seat.';
        return;
      }

      if (viewMode === 'ansi') {
        screen.innerHTML = ansiToHtml(seatData.ansi);
      } else {
        screen.innerHTML = escapeHtml(seatData.txt);
      }
    }

    function seekTo(val) {
      currentStep = parseInt(val, 10);
      renderCurrentFrame();
    }

    function selectSeat(seatIdx) {
      currentSeat = seatIdx;
      document.querySelectorAll('.seat-btn').forEach(btn => {
        btn.classList.remove('active');
      });
      event.target.classList.add('active');
      renderCurrentFrame();
    }

    function prevStep() {
      if (currentStep > 0) {
        currentStep--;
        renderCurrentFrame();
      }
    }

    function nextStep() {
      if (currentStep < shots.length - 1) {
        currentStep++;
        renderCurrentFrame();
      } else {
        togglePlay(); // stop at end
      }
    }

    function togglePlay() {
      const btn = document.getElementById('play-btn');
      if (isPlaying) {
        clearInterval(playInterval);
        isPlaying = false;
        btn.innerText = 'Play';
      } else {
        isPlaying = true;
        btn.innerText = 'Pause';
        playInterval = setInterval(() => {
          nextStep();
        }, 300 / speedMultiplier);
      }
    }

    function setSpeed(multiplier, id) {
      speedMultiplier = multiplier;
      document.querySelectorAll('.speed-selector button').forEach(btn => {
        btn.classList.remove('active');
      });
      document.getElementById('speed-' + id).classList.add('active');
      if (isPlaying) {
        togglePlay(); // restart interval with new speed
        togglePlay();
      }
    }

    function toggleViewMode() {
      const btn = document.getElementById('view-mode-btn');
      if (viewMode === 'ansi') {
        viewMode = 'txt';
        btn.innerText = 'View ANSI Frame';
      } else {
        viewMode = 'ansi';
        btn.innerText = 'View Code / Text';
      }
      renderCurrentFrame();
    }

    // Init
    renderCurrentFrame();
  </script>
</body>
</html>`;

  const dashboardPath = path.join(outDir, 'playback.html');
  fs.writeFileSync(dashboardPath, htmlContent, 'utf8');
  console.log(`[Harness] Playback dashboard successfully created at: ${dashboardPath}`);
  return dashboardPath;
}

function main() {
  const options = parseArgs();
  if (options.help) {
    printHelp();
    return;
  }

  console.log('--- Starting Test Harness ---');
  
  // 1. Build guest WASM binary
  const buildSuccess = buildWasm();
  if (!buildSuccess) {
    console.error('[Harness] Build phase failed. Aborting tests.');
    process.exit(1);
  }

  // 2. Run native smoke tests to generate ANSI/TXT recordings
  const testSuccess = runSmokeTests(options.script, options.out);
  if (!testSuccess) {
    console.error('[Harness] Smoke test phase failed.');
    process.exit(1);
  }

  // 3. Compile the generated shots into a beautiful HTML playback dashboard
  const htmlPath = generatePlaybackHtml(options.out, options.script);
  console.log(`\n[Harness] Test run completed successfully.`);
  console.log(`[Harness] Open ${htmlPath} to view the frame-by-step playback!`);
}

main();
