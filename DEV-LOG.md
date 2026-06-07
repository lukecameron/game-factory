# Game Factory - Development Log

This log registers the technical design, architectural changes, and task completions implemented during each iteration.

---

## [2026-06-07] Task 9: AI Bot Option for Snake 2 in Single-Player/Co-op Mode
- **Status**: Completed.
- **Details**:
  - Implemented an AI Bot option for Snake 2 enabled/disabled via the `B` / `b` key in single-player/co-op mode.
  - Added BFS pathfinding to navigate Snake 2 towards the active target (either food or power-up).
  - Implemented a fallback survival algorithm that uses a BFS flood-fill (up to 50 depth) to identify the safest immediate move when a direct path to the target is blocked.
  - Configured manual override detection: pressing any Arrow Key automatically deactivates the bot, letting the human player assume control of Snake 2 instantly.
  - Redesigned footer controls and playstyle divider text to fit within TUI column limits (e.g. `[WASD/Bot]` or `[WASD/Arrows]`, `PLAYSTYLE: AI-BOT`).
  - Edited: [main.go](file:///home/luke/dev/game-factory/game/main.go), [smoke.yaml](file:///home/luke/dev/game-factory/game/smoke.yaml).

## [2026-06-07] Task 8: Power-ups (Shield & Freeze) and Pause State Timer Tracking
- **Status**: Completed.
- **Details**:
  - Implemented Shield and Freeze power-up behavior. Collectibles spawn deterministically on food eaten and expire after 8 seconds if not collected.
  - Picking up a Shield gives the player temporary protection against obstacles, maze walls, hazards, self-collision, and collision with the other snake.
  - Picking up a Freeze slows down the opponent snake (moves on alternate ticks) and completely freezes all patrolling hazards.
  - Implemented pause-aware and game-over-aware timer tracking: in [main.go](file:///home/luke/dev/game-factory/game/main.go), the elapsed real time is added back to power-up spawned times and expiration deadlines when the game is paused or over, ensuring power-ups do not drain while the game is paused.
  - Added power-up status rendering inside the Row 21 divider line showing type and remaining seconds (e.g. `[🛡 6s]` or `[❄ 4s]`), along with custom snake body gradients (pulsing gold for Shield, cyan-tinted dim for Freeze).
  - Edited: [main.go](file:///home/luke/dev/game-factory/game/main.go), [smoke.yaml](file:///home/luke/dev/game-factory/game/smoke.yaml).

## [2026-06-07] Task 7: Dual-Snake Multiplayer Concurrent Steering and Key Bindings
- **Status**: Completed.
- **Details**:
  - Implemented concurrent dual-snake input handling: Player 1 controls Snake 1 using WASD keys, while Player 2 controls Snake 2 using Arrow keys.
  - In single-player mode, the single active seat controls both Snake 1 (WASD) and Snake 2 (Arrow keys) concurrently as a dual-control co-op challenge.
  - In multiplayer mode, the controls are strictly separated (Player 1 only receives WASD inputs; Player 2 only receives Arrow key inputs) to allow simultaneous play without key binding conflicts.
  - Redesigned and optimized the bottom TUI footer layout to display specific control labels dynamically based on play mode (e.g. `[P1:WASD P2:Arrows]` in multiplayer, and `[S1:WASD S2:Arrows]` in single-player/co-op mode), while maintaining the strict 80-column TUI boundary constraints.
  - Edited: [main.go](file:///home/luke/dev/game-factory/game/main.go), [smoke.yaml](file:///home/luke/dev/game-factory/game/smoke.yaml).

## [2026-06-07] Task 6: Game Modes, Patrolling Hazards, and Custom Maze Layout
- **Status**: Completed.
- **Details**:
  - Implemented 3 game modes: `CLASSIC`, `HAZARDS`, and `MAZE`.
  - Added mode cycling support: pressing the `M` key cycles modes and resets the game session immediately.
  - Implemented static maze walls for `MAZE` mode: designed a beautiful, symmetric, and highly playable neon maze barrier layout.
  - Added dynamic patrolling hazards: security drones/hazards (`❖` glyph) that move back and forth along predetermined paths.
  - Custom hazard styles: designed custom truecolor hazard palette entries matching all 5 themes. Hazards pulse dynamically and flash white/red warning colors when the snake's head is within 2 cells.
  - Collision checks: implemented comprehensive hazard and maze wall collision checks, including a swap-collision safety check to prevent passing through moving hazards.
  - UI layout improvements: updated the row 21 divider to display the active game mode, and restructured the footer controls to avoid overlaps and respect TUI constraints.
  - Edited: [main.go](file:///home/luke/dev/game-factory/game/main.go), [smoke.yaml](file:///home/luke/dev/game-factory/game/smoke.yaml).

## [2026-06-07] Task 5: Leaderboards, Dynamic Speed Scaling, and Hot-Seat Multiplayer UI
- **Status**: Completed.
- **Details**:
  - Implemented static leaderboard registration in GameMeta. Posts the final score to the leaderboard on collision (game over) using `r.Post`.
  - Added active seat tracking using `activePlayer` on the room struct, updating it during `OnJoin` and `OnInput` calls.
  - Revamped the top header layout dynamically. In multiplayer mode (`len(r.Members()) > 1`), a dedicated `SEAT: <handle>` indicator appears next to the score/theme info, with precise mathematical padding to avoid any cell overlaps or text collisions.
  - Implemented a centered bottom-divider status text: When in multiplayer mode, the bottom double-line divider displays ` ACTIVE SEAT: <handle> ` centered inside the divider line.
  - Updated the centered Game Over modal: Displays the active seat handle along with the final score, e.g. `FINAL SCORE: 0010 (seat0)`.
  - Added dynamic difficulty/speed scaling: Eaten food increases the tick rate (decreases tick duration by 5ms for every 10 points), clamped at 60ms min.
  - Added tick rate resetting: Resetting the game properly restores the tick rate back to 150ms.
  - Edited: [main.go](file:///home/luke/dev/game-factory/game/main.go), [smoke.yaml](file:///home/luke/dev/game-factory/game/smoke.yaml).

## [2026-06-07] Task 4: Premium Aesthetics, Custom Palettes, & Micro-animations
- **Status**: Completed.
- **Details**:
  - Implemented 5 vibrant, handcrafted truecolor themes:
    - **Cyberpunk Neon**: Neon Violet, Aqua, Lime Green, Pink, Orange.
    - **Neon Ocean**: Deep Blue, Teal, Coral, Gold, Sky Blue.
    - **Sunset Glow**: Crimson, Gold, Orange, Deep Red, Orchid.
    - **Matrix Grid**: Dark Green, Neon Green, Forest Green, Lime, Yellow-Green.
    - **Vaporwave**: Pastel Purple, Yellow, Violet, Hot Pink, Cyan, Magenta.
  - Added theme switching support: Player can cycle through themes dynamically by pressing the `T` or `t` key during the game.
  - Implemented micro-animations:
    - **Border Laser Pulse**: A dynamic bright laser segment flows around the double-line borders.
    - **Twinkling Sparkle Food**: Food item twinkles by cycling through different star/sparkle glyphs (`★`, `☆`, `✦`, `✧`) and pulsing its truecolor intensity.
    - **Proximity Warning Flash**: Obstacles near the snake's head (distance <= 2 cells) flash rapidly between the theme color and white/red warning colors.
    - **Floating Score Popups**: Eaten food triggers a floating `+10` text effect that rises and dims over 1 second.
    - **Collision Flash Screen Effect**: Playfield border flashes white/pink for 300ms on impact.
  - Edited: [main.go](file:///home/luke/dev/game-factory/game/main.go), [smoke.yaml](file:///home/luke/dev/game-factory/game/smoke.yaml).

## [2026-06-07] Task 3: Collision Detection, Spawning, Scoring, and Game Over Restart Mechanics
- **Status**: Completed.
- **Details**:
  - Added food spawning and eating logic: Spawns food at random unoccupied cells. Eating food increases the score by 10, updates the high score, grows the snake by 1 segment, and triggers new food and obstacle spawning.
  - Implemented obstacle generation: Generates 3 initial obstacles at startup, and generates one new obstacle dynamically every time food is eaten, using deterministic random seeding (`r.Rand()`).
  - Added full collision detection:
    - **Self-collision**: Checks if the snake head collides with any other segment in the snake body.
    - **Obstacle-collision**: Checks if the snake head collides with any spawned obstacles.
    - **Game Over state**: Pauses game ticks, blocks movements, and displays a flashing game over modal centered on the screen.
  - Implemented restart mechanism: Pressing `Space` (ActConfirm) when the game is over calls a reset routine that clears obstacles/food/snake and starts a new session.
  - Edited: [main.go](file:///home/luke/dev/game-factory/game/main.go).

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
