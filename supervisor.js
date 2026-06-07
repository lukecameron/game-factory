const fs = require('fs');
const crypto = require('crypto');
const path = require('path');
const { spawnSync } = require('child_process');

const WATCHED_FILES = [
  'supervisor.js',
  'run-supervisor.sh',
  'tools/test-harness.js',
  'ITERATION-PROMPT.md',
  'REFLECTION-PROMPT.md'
];

const STATE_FILE = '.factory-state.json';

function getFileHash(filepath) {
  try {
    const content = fs.readFileSync(filepath);
    return crypto.createHash('sha256').update(content).digest('hex');
  } catch (e) {
    return null;
  }
}

function getAllHashes() {
  const hashes = {};
  for (const file of WATCHED_FILES) {
    hashes[file] = getFileHash(file);
  }
  return hashes;
}

function hashesChanged(oldHashes, newHashes) {
  for (const file of WATCHED_FILES) {
    if (oldHashes[file] !== newHashes[file]) {
      console.log(`[Supervisor] File changed: ${file}`);
      return true;
    }
  }
  return false;
}

function loadState() {
  if (fs.existsSync(STATE_FILE)) {
    try {
      return JSON.parse(fs.readFileSync(STATE_FILE, 'utf8'));
    } catch (e) {
      console.error('[Supervisor] Failed to parse state file, using defaults', e);
    }
  }
  return {
    current_iteration: 1,
    reflect_interval: 10,
    stage: 'iteration',
    delay_between_iterations_seconds: 300,
    last_updated: new Date().toISOString()
  };
}

function saveState(state) {
  state.last_updated = new Date().toISOString();
  fs.writeFileSync(STATE_FILE, JSON.stringify(state, null, 2), 'utf8');
}

function runAgent(promptFile) {
  console.log(`[Supervisor] Reading prompt from ${promptFile}...`);
  let prompt;
  try {
    prompt = fs.readFileSync(promptFile, 'utf8');
  } catch (e) {
    console.error(`[Supervisor] Failed to read prompt file ${promptFile}:`, e);
    return false;
  }
  
  console.log(`[Supervisor] Spawning agy CLI...`);
  // Ensure local tools and Go/TinyGo are in PATH for agy execution
  const localBin = '/home/luke/.local/bin';
  const goBin = '/home/luke/.local/go/bin';
  const tinygoBin = '/home/luke/.local/tinygo/bin';
  const customPath = `${localBin}:${goBin}:${tinygoBin}:${process.env.PATH}`;

  const result = spawnSync('agy', [
    '--dangerously-skip-permissions',
    '-p',
    prompt
  ], {
    stdio: 'inherit',
    env: { ...process.env, PATH: customPath, GOTOOLCHAIN: 'local' }
  });
  
  return result.status === 0;
}

function runGit(args) {
  const result = spawnSync('git', args, { stdio: 'inherit' });
  return result.status === 0;
}

function commitIteration(iteration) {
  console.log(`[Supervisor] Committing iteration ${iteration}...`);
  runGit(['add', '.']);
  runGit(['commit', '-m', `iteration ${iteration}: completed task`]);
  // We can push to the current branch
  runGit(['push', 'origin', 'main']);
}

function commitReflection(iteration) {
  console.log(`[Supervisor] Committing reflection after iteration ${iteration}...`);
  runGit(['add', '.']);
  runGit(['commit', '-m', `reflection after iteration ${iteration}: optimized factory`]);
  runGit(['push', 'origin', 'main']);
}

async function main() {
  const state = loadState();
  console.log(`[Supervisor] Loaded State:`, state);

  while (true) {
    const initialHashes = getAllHashes();

    if (state.stage === 'iteration') {
      console.log(`\n--- Starting Iteration ${state.current_iteration} ---`);
      
      const success = runAgent('ITERATION-PROMPT.md');
      if (!success) {
        console.error(`[Supervisor] Iteration ${state.current_iteration} failed. Pausing loop.`);
        process.exit(1);
      }
      
      commitIteration(state.current_iteration);
      
      // Reload state from disk to capture any agent updates (e.g., tuned delay or reflect_interval)
      const diskState = loadState();
      state.delay_between_iterations_seconds = diskState.delay_between_iterations_seconds;
      state.reflect_interval = diskState.reflect_interval;

      if (state.current_iteration % state.reflect_interval === 0) {
        state.stage = 'reflection';
      } else {
        state.current_iteration += 1;
      }
      saveState(state);
      
    } else if (state.stage === 'reflection') {
      console.log(`\n--- Starting Reflection ---`);
      
      const success = runAgent('REFLECTION-PROMPT.md');
      if (!success) {
        console.error(`[Supervisor] Reflection failed. Pausing loop.`);
        process.exit(1);
      }
      
      commitReflection(state.current_iteration);
      
      // Reload state from disk to capture any reflection agent updates
      const diskState = loadState();
      state.delay_between_iterations_seconds = diskState.delay_between_iterations_seconds;
      state.reflect_interval = diskState.reflect_interval;

      state.stage = 'iteration';
      state.current_iteration += 1;
      saveState(state);
    }

    // Check for self-updates
    const currentHashes = getAllHashes();
    if (hashesChanged(initialHashes, currentHashes)) {
      console.log(`[Supervisor] Self-update detected in watched files. Rebooting...`);
      process.exit(3);
    }

    const delayMs = state.delay_between_iterations_seconds * 1000;
    console.log(`[Supervisor] Step complete. Cooling down for ${state.delay_between_iterations_seconds} seconds...`);
    await new Promise(resolve => setTimeout(resolve, delayMs));
  }
}

main().catch(err => {
  console.error('[Supervisor] Fatal error in main loop:', err);
  process.exit(1);
});
