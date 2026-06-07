# game — a shellcade game

## Develop (instant, no wasm)

    go run .                 # play in this terminal; Esc leaves
    go run . -seats 2        # hot-seat multiplayer; Ctrl-T switches seats
    go run . -seed 42        # reproducible runs

## Build the wasm artifact (~4s)

    tinygo build -opt=1 -no-debug -gc=leaking \
        -o game.wasm -target wasip1 -buildmode=c-shared .

Then verify with the shellcade developer kit: shellcade-kit check game.wasm — and play the
real artifact with shellcade-kit play game.wasm.

## Learn more

- GUIDE.md in github.com/shellcade/kit/v2 — the authoring guide
- ABI.md — the contract your game targets
- github.com/shellcade/games — published example games using every SDK feature
