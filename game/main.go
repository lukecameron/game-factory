// game — a shellcade game. Run it right now: go run .
package main

import (
	"fmt"
	"time"

	kit "github.com/shellcade/kit/v2"
)

func main() { kit.Main(Game{}) }

// Game is the registry entry: metadata + a per-room behavior factory.
type Game struct{}

func (Game) Meta() kit.GameMeta {
	return kit.GameMeta{
		Slug:             "game",
		Name:             "game",
		ShortDescription: "Describe your game in one line.",
		MinPlayers:       1,
		MaxPlayers:       4,
	}
}

func (Game) NewRoom(cfg kit.RoomConfig, svc kit.Services) kit.Handler {
	return &room{}
}

// room is one live room. ALL state lives here (and only here) — the host can
// snapshot and restore it, so key anything durable by Player.AccountID.
type room struct {
	kit.Base
	presses  int
	deadline time.Time // a wake-driven one-shot: see OnWake
}

func (rm *room) OnStart(r kit.Room) {
	r.SetInputContext(kit.CtxNav)
}

func (rm *room) OnJoin(r kit.Room, p kit.Player) { rm.render(r) }

func (rm *room) OnInput(r kit.Room, p kit.Player, in kit.Input) {
	switch kit.Resolve(in, kit.CtxNav) {
	case kit.ActConfirm:
		rm.presses++
		// One-shot timer, the wake way: store a deadline, check it in OnWake.
		rm.deadline = r.Now().Add(2 * time.Second)
	}
	rm.render(r)
}

// OnWake is the host heartbeat — the ONLY time your code runs without input.
// Drive every animation, countdown, and timeout from CallContext time here.
func (rm *room) OnWake(r kit.Room) {
	if !rm.deadline.IsZero() && r.Now().After(rm.deadline) {
		rm.deadline = time.Time{}
		rm.presses = 0 // the timeout fired: reset
	}
	rm.render(r)
}

func (rm *room) render(r kit.Room) {
	f := kit.NewFrame() // frames are POINTERS, always (see ABI.md §6)
	title := kit.Style{FG: kit.Cyan, Attr: kit.AttrBold}
	dim := kit.Style{FG: kit.DimGray}

	f.Text(2, 4, "*** game ***", title)
	f.Text(10, 4, fmt.Sprintf("SPACE pressed %d times", rm.presses), kit.Style{FG: kit.White})
	if !rm.deadline.IsZero() {
		left := rm.deadline.Sub(r.Now()).Round(100 * time.Millisecond)
		f.Text(12, 4, fmt.Sprintf("resetting in %s...", left), kit.Style{FG: kit.Yellow})
	}
	f.Text(kit.Rows-1, 2, "SPACE press   Esc leave", dim)

	for _, p := range r.Members() {
		r.Send(p, f)
	}
}
