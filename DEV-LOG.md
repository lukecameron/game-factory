# Game Factory - Development Log

This log registers the technical design, architectural changes, and task completions implemented during each iteration.

---

## [2026-06-07] Task 2: Multi-segment Snake Movement & Input Safety
- **Status**: Completed.
- **Details**:
  - Replaced the single-point entity (`entityPos`) with a slice representing a multi-segment snake body (`snake []Point`) initialized to length 4 at startup in [main.go](file:///home/luke/dev/game-factory/game/main.go).
  - Implemented body segment movement logic where each segment shifts to take the position of the preceding segment on each clock tick in `tick()`.
  - Added an input safety lock (`lastMovedDir`) to prevent 180-degree self-collisions when players double-tap opposite steering directions within a single clock cycle.
  - Designed a premium truecolor neon gradient style for the snake body that transitions smoothly from Lime Green (`RGB(0x39, 0xff, 0x14)`) at the head to Cyan (`RGB(0x00, 0xe5, 0xff)`) at the tail.
  - Ensured background dot grid rendering skips any cell occupied by any portion of the snake body.

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
