# Game Factory - Testing & Feel Log

This log documents test coverage, gameplay feel, input responsiveness, and visual verification results for each iteration.

---

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
