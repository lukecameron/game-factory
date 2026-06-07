# Game Factory - Testing & Feel Log

This log documents test coverage, gameplay feel, input responsiveness, and visual verification results for each iteration.

---

## [2026-06-07] Task 9: AI Bot Option for Snake 2 in Single-Player/Co-op Mode Test
- **Status**: Passed successfully.
- **Coverage**:
  - Verification of AI Bot toggle functionality: Player 1 triggers the bot toggle key (`b`) which transitions the active playstyle.
  - Verification of BFS-based guide movement: Once enabled, Snake 2 steers itself UP towards the active food coordinate on successive ticks.
  - Verification of Manual override behavior: Player 2 inputs a steering override key (`left` keypress), which immediately deactivates the AI Bot and places the snake back under manual control.
  - Verification of UI constraints: TUI status panels, divider details, and bottom navigation controls adapt cleanly within 80-column limits when bot states change.
- **Playback URL**: [playback.html](file:///home/luke/dev/game-factory/test-logs/shots/playback.html)
- **Gameplay Feel Observations**:
  - Toggling between co-op manual play and CPU bot mode feels extremely fluid and responsive.
  - The BFS pathfinding allows the bot to make smart decisions, hunting down food efficiently without crashing.
  - The instant override is a superb UX feature; players can let the bot handle simple straight paths, but grab the controls immediately when they want to make fine-tuned maneuvers, which feels incredibly smooth and empowering.

## [2026-06-07] Task 8: Power-ups (Shield & Freeze) and Pause State Timer Tracking Test
- **Status**: Passed successfully.
- **Coverage**:
  - Verification of concurrent steering controls for dual-snake setup.
  - Cycle of multiple game modes (Classic, Hazards, Maze).
  - Food eating logic resulting in +10 score popups, length growth, speed scaling, and deterministic Shield power-up spawning on the board.
  - Steering Snake 2 downwards and then rightwards to avoid colliding with boundaries/obstacles.
  - Verification of power-up collection: Snake 2 successfully collects the Shield power-up `🛡`.
  - Active power-up visualization rendering test: verified that the floating `+SHIELD` popup displays at collection, the snake body glows gold with a custom sinusoidal wave pulse style, and the Row 21 divider displays `[🛡 6s]` and subsequently counts down correctly.
  - Verified pause/resume mechanics do not drain power-up active durations.
- **Playback URL**: [playback.html](file:///home/luke/dev/game-factory/test-logs/shots/playback.html)
- **Gameplay Feel Observations**:
  - Collecting power-ups feels extremely satisfying and adds a strategic layer to navigation, especially in maze and hazard modes.
  - The custom gold pulse wave animation on the snake segments when shielded looks highly premium and makes the invincibility state feel very powerful and responsive.
  - The floating text popups (`+SHIELD`, `+FREEZE`, `+10`) rise and fade smoothly, offering juicy visual feedback.
  - The pause time adjustment keeps the game feeling completely fair; players can pause without losing their hard-earned power-ups.

## [2026-06-07] Task 7: Dual-Snake Multiplayer Concurrent Steering and Key Bindings Test
- **Status**: Passed successfully.
- **Coverage**:
  - Initial load rendering in classic mode with 2 seats (multiplayer layout).
  - Theme switching verification for seat 0 (Player 1) and seat 1 (Player 2).
  - Normal gameplay automatic movement leading to mutual head-on collision and Game Over display.
  - Game restart using Space key from Player 1 (seat 0).
  - Concurrent steering input verification:
    - Player 1 (seat 0) sends WASD input 's' to steer Snake 1 downwards.
    - Player 2 (seat 1) sends Arrow Key input 'up' to steer Snake 2 upwards.
    - Advancing time by 300ms to verify both snakes react concurrently and move safely in their new directions.
  - Game mode switching to **Hazards**, then **Maze**, and cycling back to **Classic** while verifying correct layouts and hazard movement.
- **Playback URL**: [playback.html](file:///home/luke/dev/game-factory/test-logs/shots/playback.html)
- **Gameplay Feel Observations**:
  - The steering controls feel incredibly responsive. Restricting Player 1 to WASD and Player 2 to Arrow keys ensures no overlapping keystroke conflicts, making concurrent local play on a shared keyboard extremely fun and smooth.
  - The bottom TUI footer automatically changes labels (e.g. `[P1:WASD P2:Arrows]` vs `[S1:WASD S2:Arrows]`) which makes it very clear how to play, fitting perfectly within the 80x24 TUI screen limits.
  - Dual-control co-op mode (in single-player) offers a unique and tricky brain-teaser challenge, as you navigate both snakes concurrently trying to coordinate their paths without crashing into each other or obstacles.

## [2026-06-07] Task 6: Game Modes, Patrolling Hazards, and Custom Maze Layout Test
- **Status**: Passed successfully.
- **Coverage**:
  - Initial load rendering in classic mode.
  - Theme switching and multiplayer active seat updates.
  - Normal gameplay movement and food eating.
  - Classic static obstacle collision and Game Over restart.
  - Mode switching to **Hazards**:
    - Pressing 'm' switches to Hazards mode, initializing 4 patrolling hazards with theme-colored diamond `❖` glyphs.
    - Advanced time by 300ms to verify hazard movement along their patrolled segments.
  - Mode switching to **Maze**:
    - Pressing 'm' switches to Maze mode, rendering symmetric neon maze walls `▒` and patrolling hazards at safe paths.
    - Advanced time by 300ms to verify hazard movement in the maze layout.
  - Return to **Classic** mode:
    - Pressing 'm' cycles back to Classic mode.
- **Playback URL**: [playback.html](file:///home/luke/dev/game-factory/test-logs/shots/playback.html)
- **Gameplay Feel Observations**:
  - Cycling modes with a single key press ('m') restarts the game instantly and is extremely responsive, making it fun to try different layouts.
  - The patrolling hazards look extremely premium with their slow sinusoidal breathing/pulsing colors, and the proximity-based warning flashes (white/red when <= 2 tiles) add intense visual feedback.
  - The maze mode feels like a completely different game, requiring tight navigation in corridors while avoiding patrolling security drones.
  - The swap-collision protection is incredibly robust; passing through hazards head-on correctly registers a collision instantly.

## [2026-06-07] Task 5: Leaderboards, Dynamic Speed Scaling, and Hot-Seat Multiplayer UI Test
- **Status**: Passed successfully.
- **Coverage**:
  - Initial load rendering in multiplayer mode (2 seats joined).
  - Theme switching verification:
    - Pressing 't' cycles theme from **Cyberpunk** to **Ocean**.
  - Active seat switching verification:
    - Switching active seat to `seat1` updates the active seat indicator in the header to `SEAT: seat1` and the bottom-divider text to `ACTIVE SEAT: seat1`.
    - Subsequent inputs from `seat1` (pressing 't') cycle theme to **Sunset**.
    - Switching seat back to `seat0` updates header/divider to `SEAT: seat0` and `ACTIVE SEAT: seat0`.
    - Input from `seat0` cycles theme to **Matrix**.
  - Snake movement and food eating verification.
  - Collision detection:
    - Crashing into the obstacle registers Game Over.
    - Verified that the Game Over modal renders the final score along with the player seat that died: `FINAL SCORE: 0010 (seat0)`.
    - Verifies that leaderboard score is correctly posted on game over.
  - Restart confirmation restores status to PLAYING and resets speed to base (150ms).
- **Playback URL**: [playback.html](file:///home/luke/dev/game-factory/test-logs/shots/playback.html)
- **Gameplay Feel Observations**:
  - The multiplayer seat switching is immediate and visual, and the highlighted seat in both the header and the bottom divider makes hot-seat feel completely integrated.
  - Spacing of the header elements is clean and has perfect 1-cell gaps, resolving any overlaps or text truncation.
  - Dynamic speed scaling (speeding up as you eat more food) makes the game feel progressively intense, adding a real arcade challenge.
  - Resetting speed on restart keeps the restart loop very fair and satisfying.

---

## [2026-06-07] Task 4: Premium Aesthetics, Themes, and Animations Test
- **Status**: Passed successfully.
- **Coverage**:
  - Initial load rendering.
  - Theme switching verification:
    - Pressing 't' cycles the theme from **Cyberpunk Neon** to **Neon Ocean**.
    - Pressing 't' again cycles the theme to **Sunset Glow**.
    - Captured separate shots (`theme_neon_ocean` and `theme_sunset_glow`) to verify color transitions, layouts, and header elements.
  - Food Eating & Score Popup verification:
    - Moving to food cell and eating it creates the "+10" popup at the exact food coordinate, rises vertically, and disappears correctly.
  - Proximity Warning & Danger Obstacles verification:
    - Approaching an obstacle within 2 cells activates the warning flashing animation to alert the player.
  - Game Over Collision & Impact Flash verification:
    - Crashing into the obstacle triggers a full screen-border impact flash for 300ms, halts gameplay ticks, and brings up the themed GAME OVER menu.
  - Restart confirmation restores status to PLAYING.
- **Playback URL**: [playback.html](file:///home/luke/dev/game-factory/test-logs/shots/playback.html)
- **Gameplay Feel Observations**:
  - The flowing border lights and glowing gradient snake segment colors make the playfield feel exceptionally alive and arcade-like.
  - The twinkling/sparkling food shape cycles add a delightful layer of visual polish.
  - The warning flash on close obstacles is highly practical and visually dramatic, warning players before they crash.
  - The floating "+10" popup gives direct, juicy feedback upon eating food.
  - The red/white impact border flash when crashing makes collisions feel impactful and satisfyingly punchy.

## [2026-06-07] Task 3: Collision Detection & Game Mechanics Test
- **Status**: Passed successfully.
- **Coverage**:
  - Initial 4-segment snake layout and random placement of food (seed 1 spawns at `X: 23, Y: 15`) and 3 obstacles.
  - Forward rightward movement (4 ticks / 600ms) to align head with food's X column.
  - Downward steering input (`down` key) and advance (6 ticks / 900ms) to successfully eat the food.
  - Growth verification (snake grows from length 4 to 5 segments), scoring verification (score increases to `10`, high score is `10`), and dynamic obstacle/food generation verification.
  - Leftward steering input (`left` key) and advance (12 ticks / 1800ms) to position the snake head under the obstacle at `X: 11, Y: 11`.
  - Upward steering input (`up` key) and advance (4 ticks / 600ms) to trigger a head-on collision with the obstacle.
  - Game Over verification: ensures game tick halts, movements are blocked, and the premium centered "GAME OVER" modal is rendered correctly.
  - Restart mechanism verification: space key press triggers a call to `reset()`, resetting snake length to 4, clearing score, regenerating 3 new obstacles/food, and returning status to `PLAYING`.
- **Playback URL**: [playback.html](file:///home/luke/dev/game-factory/test-logs/shots/playback.html)
- **Gameplay Feel Observations**:
  - The collision check feels extremely immediate and accurate, locking the game state instantly when hit.
  - The pulsing glowing food star and the dynamic orange triangle obstacles look premium, cohesive, and neon-arcade-like.
  - The centered Game Over modal has a beautiful double-line border and flashing text, which gives a strong, polished retro-arcade feel upon defeat.
  - The instant restart via the SPACE bar provides a highly responsive retry loop, making it very satisfying to jump back in.

## [2026-06-07] Task 2: Multi-segment Snake Movement & Input Safety Test
- **Status**: Passed successfully.
- **Coverage**:
  - Initial 4-segment snake layout rendering.
  - Automatic rightward movement (3 ticks / 450ms).
  - 180-degree safety lock verification (attempting to turn left while moving right is blocked).
  - Downward steering input handling (Arrow Down key press, followed by 300ms advance).
  - 180-degree safety lock verification (attempting to turn up while moving down is blocked).
  - Leftward steering input handling (Arrow Left key press, followed by 300ms advance to confirm turning left from moving down).
  - Pausing/Unpausing mechanism (Space key press, followed by 300ms advance to ensure movement stopped).
- **Playback URL**: [playback.html](file:///home/luke/dev/game-factory/test-logs/shots/playback.html)
- **Gameplay Feel Observations**:
  - The multi-segment body movement trailing the head functions smoothly, looking very natural.
  - The color gradient transition down the snake body is visually gorgeous and adds a highly premium, glowing neon quality.
  - The input turn lock makes the controls feel robust and fair, preventing accidental self-destruction from rapid consecutive key presses.

## [2026-06-07] Task 1: Core Loop & Board Layout Test
- **Status**: Passed successfully.
- **Coverage**:
  - Initial screen layout rendering.
  - Automatic rightward movement (3 ticks / 450ms).
  - Downward steering input handling (Arrow Down key press, followed by 300ms advance).
  - Pausing/Unpausing mechanism (Space key press, followed by 300ms advance to ensure movement stopped).
- **Playback URL**: [playback.html](file:///home/luke/dev/game-factory/test-logs/shots/playback.html)
- **Gameplay Feel Observations**:
  - The double-width grid mapping resolves the typical 2:1 terminal cell aspect ratio distortion, making horizontal/vertical movement feel visually consistent.
  - The 150ms tick rate is a solid pace for game navigation.
  - The pulsing truecolor status text and the styled dot grid playfield create a premium arcade feel.

## [2026-06-07] Task 0: Scaffolding Smoke Test
- **Status**: Verified successfully.
- **Coverage**:
  - Initial frame load (`wake`).
  - Keypress input verification (`rune: ' '` for space press).
  - One-shot wake delay verification (`advance: 2.5s` to test game state reset).
- **Playback URL**: [playback.html](file:///home/luke/dev/game-factory/test-logs/shots/playback.html)
- **Gameplay Feel Observations**:
  - The UI is highly responsive; the frame transitions upon SPACE press occur instantaneously.
  - The TUI spacing constraints (80x24) are respected, and the reset transition after the 2-second delay operates smoothly via the wake mechanism.

## [2026-06-07] Pre-Iteration
- **Status**: Setup verification pending.
- **Details**: Toolchain verified natively. Awaiting Phase 3 game scaffolding to execute first smoke test.
