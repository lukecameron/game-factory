# Game Factory - Iteration Instructions

You are an agentic coding assistant running in an automated game development loop. Your goal is to implement a single gameplay feature or change, test it agentically, and update the project logs.

---

## Your Step-by-Step Workflow

### 1. Load Current Context
- Read the files:
  - [ROADMAP.md](file:///home/luke/dev/game-factory/ROADMAP.md) (identify the first unchecked task).
  - [DEV-LOG.md](file:///home/luke/dev/game-factory/DEV-LOG.md) (understand previous implementation details).
  - [TEST-LOG.md](file:///home/luke/dev/game-factory/TEST-LOG.md) (understand previous testing observations and feedback).
- Select exactly **one** task to complete.

### 2. Implement the Feature
- Make code modifications **only** inside the `game/` subdirectory.
- Preserve the 80x24 TUI constraints, focus on keyboard-driven controls, and maintain code simplicity (refer to [NORTH-STAR.md](file:///home/luke/dev/game-factory/NORTH-STAR.md) guidelines).

### 3. Design and Run Tests
- Write/update the smoke test yaml script at `game/smoke.yaml` to cover your implementation. Make sure it uses valid step definitions (e.g., `rune`, `key`, `text`, `advance`, `wake`, `shot`).
- Run the test harness from the root workspace:
  ```bash
  node tools/test-harness.js --script game/smoke.yaml --out test-logs/shots
  ```
- Review the terminal output and fix any compilation or smoke test errors.

### 4. Document Your Work
- **Roadmap**: Update [ROADMAP.md](file:///home/luke/dev/game-factory/ROADMAP.md) by checking off the completed task. If you discovered new bugs or next steps, append them as new unchecked tasks.
- **Dev Log**: Add a structured entry in [DEV-LOG.md](file:///home/luke/dev/game-factory/DEV-LOG.md) documenting:
  - The technical details of the code changes.
  - File links for edited files.
- **Test Log**: Add an entry in [TEST-LOG.md](file:///home/luke/dev/game-factory/TEST-LOG.md) documenting:
  - What steps are covered in your smoke test.
  - Playback URL: Create a link to the generated HTML dashboard (e.g. `[playback.html](file:///home/luke/dev/game-factory/test-logs/shots/playback.html)`).
  - Gameplay feel observations (responsiveness, animation pacing, scoring feel).

### 5. Final Check
Ensure the WASM binary was successfully compiled at `game/game.wasm`. Once done, complete your execution.
