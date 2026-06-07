# Game Factory Roadmap

This file defines the tasks to be completed by the game factory. Each iteration, the agent will claim and complete exactly one task.

## Tasks

- [x] **Task 0**: Bootstrap the game scaffold using `shellcade-kit new game` and verify that the test harness compiles the WASM binary and runs the smoke tests correctly.
- [x] **Task 1**: Implement the core game loop and draw a basic 80x24 board structure.
- [x] **Task 2**: Implement player controls (input handling for keyboard keys) and render basic entity/player movements.
- [x] **Task 3**: Add collision detection and basic game mechanics (e.g. scoring, obstacle generation, win/loss criteria).
- [x] **Task 4**: Add color palettes, custom truecolor aesthetics, and micro-animations to enhance the game's premium feel.
- [x] **Task 5**: Final polishing, edge cases handling, and verification of multiplayer hot-seat seat switching if applicable.

## Future Ideas
- [x] **Task 6**: Add moving/patrolling hazards or complex maps (e.g. mazes) to the playfield.
- [x] **Task 7**: Implement dual-snake multiplayer where two players steer their own snakes concurrently using different key bindings.
- [x] **Task 8**: Implement Power-ups (Shield & Freeze) to dynamically alter gameplay mechanics and track them correctly during pause/game-over states.
- [x] **Task 9**: Implement an AI Bot player option for Snake 2 in single-player/co-op mode. Add a toggle key B to switch between CO-OP and BOT, guiding Snake 2 using BFS pathfinding and fallback survival flood fill, with manual steering overriding the bot instantly.
- [x] **Task 10**: Implement Portal Mode. Add a 4th game mode "PORTALS" featuring symmetric neon blue & orange portal rings. Teleport snakes through portals dynamically, ensuring body segments flow smoothly and the AI Bot navigates portal pathways using portal-aware BFS search.

## Future Ideas
- [x] **Task 11**: Add a customizable settings menu (e.g. custom snake skins/characters, toggle grid dots, or custom start tick rates) accessible via the `S` or `s` key.
- [x] **Task 12**: Implement retro terminal beep/audio cues using standard ASCII bell `\a` or distinct patterns for events like eating food, collecting power-ups, or crashing.
