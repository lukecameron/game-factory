# Game Factory - Development Log

This log registers the technical design, architectural changes, and task completions implemented during each iteration.

---

## [2026-06-07] Task 1: Core Loop and TUI Board Layout
- **Status**: Completed.
- **Details**:
  - Implemented a 39x18 game playfield inside the 80x24 TUI constraints in [main.go](file:///home/luke/dev/game-factory/game/main.go).
  - Used double-width grid mapping (2 columns per grid unit) to handle terminal cell aspect ratio (2:1).
  - Implemented the core game loop in `OnWake` using comparisons against `r.Now()` (150ms tick rate).
  - Built a truecolor neon visual layout with a contiguous double-line border, pulsing status header, and dot grid overlay.
  - Added support for WASD steering controls alongside standard arrow keys in `OnInput`.

## [2026-06-07] Task 0: Scaffolding and Harness Verification
- **Status**: Completed verification of build toolchain and smoke test harness.
- **Details**:
  - Verified compilation of the WASM binary at [game.wasm](file:///home/luke/dev/game-factory/game/game.wasm).
  - Executed the custom test harness using [smoke.yaml](file:///home/luke/dev/game-factory/game/smoke.yaml).
  - Validated that the harness correctly generates snapshots and the interactive [playback.html](file:///home/luke/dev/game-factory/test-logs/shots/playback.html).

## [2026-06-07] Bootstrap (Pre-Iteration)
- **Status**: Completed toolchain setup.
- **Details**:
  - Installed Go 1.25.0 locally in `~/.local/go`.
  - Installed TinyGo 0.35.0 locally in `~/.local/tinygo`.
  - Downloaded `shellcade-kit` CLI binary (`v2.3.0`) and saved to `~/.local/bin`.
  - Created `run-supervisor.sh` and `supervisor.js` to coordinate the loop.
  - Implemented the custom testing harness `tools/test-harness.js` that compiles WASM and produces interactive HTML playbacks.
