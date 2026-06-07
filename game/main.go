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
		Name:             "Neon Snake",
		ShortDescription: "A retro-inspired premium Neon Snake game built for TUI.",
		MinPlayers:       1,
		MaxPlayers:       4,
		HeartbeatMS:      50, // 20 ticks per second
	}
}

func (Game) NewRoom(cfg kit.RoomConfig, svc kit.Services) kit.Handler {
	return &room{}
}

type Point struct {
	X, Y int
}

// room is one live room. ALL state lives here (and only here) — the host can
// snapshot and restore it, so key anything durable by Player.AccountID.
type room struct {
	kit.Base
	frame       *kit.Frame
	lastTick    time.Time
	tickRate    time.Duration
	score       int
	highScore   int
	gameStarted bool
	gameOver    bool

	snake        []Point
	entityDir    Point
	lastMovedDir Point
	food         Point
	obstacles    []Point

	startedAt time.Time
}

func (rm *room) OnStart(r kit.Room) {
	r.SetInputContext(kit.CtxNav)
	rm.lastTick = r.Now()
	rm.startedAt = r.Now()
	rm.tickRate = 150 * time.Millisecond
	rm.snake = []Point{
		{X: 19, Y: 9},
		{X: 18, Y: 9},
		{X: 17, Y: 9},
		{X: 16, Y: 9},
	}
	rm.entityDir = Point{X: 1, Y: 0}
	rm.lastMovedDir = Point{X: 1, Y: 0}
	rm.gameStarted = true
	rm.score = 0
	rm.gameOver = false

	// Generate initial food & obstacles
	rm.food = rm.randomFreePoint(r, 0)
	rm.obstacles = []Point{}
	for i := 0; i < 3; i++ {
		rm.obstacles = append(rm.obstacles, rm.randomFreePoint(r, 5))
	}
}

func (rm *room) OnJoin(r kit.Room, p kit.Player) { rm.render(r) }

func (rm *room) OnInput(r kit.Room, p kit.Player, in kit.Input) {
	action := kit.Resolve(in, kit.CtxNav)

	// Custom support for WASD controls
	if in.Kind == kit.InputRune {
		switch in.Rune {
		case 'w', 'W':
			action = kit.ActUp
		case 's', 'S':
			action = kit.ActDown
		case 'a', 'A':
			action = kit.ActLeft
		case 'd', 'D':
			action = kit.ActRight
		}
	}

	switch action {
	case kit.ActUp:
		if rm.lastMovedDir.Y != 1 && !rm.gameOver {
			rm.entityDir = Point{X: 0, Y: -1}
		}
	case kit.ActDown:
		if rm.lastMovedDir.Y != -1 && !rm.gameOver {
			rm.entityDir = Point{X: 0, Y: 1}
		}
	case kit.ActLeft:
		if rm.lastMovedDir.X != 1 && !rm.gameOver {
			rm.entityDir = Point{X: -1, Y: 0}
		}
	case kit.ActRight:
		if rm.lastMovedDir.X != -1 && !rm.gameOver {
			rm.entityDir = Point{X: 1, Y: 0}
		}
	case kit.ActConfirm:
		if rm.gameOver {
			rm.reset(r)
		} else {
			rm.gameStarted = !rm.gameStarted
		}
	}
	rm.render(r)
}

// OnWake is the host heartbeat — the ONLY time your code runs without input.
// Drive every animation, countdown, and timeout from CallContext time here.
func (rm *room) OnWake(r kit.Room) {
	now := r.Now()
	if rm.lastTick.IsZero() {
		rm.lastTick = now
	}
	if rm.startedAt.IsZero() {
		rm.startedAt = now
	}

	// Advance game state based on tickRate
	if rm.gameStarted && !rm.gameOver && now.Sub(rm.lastTick) >= rm.tickRate {
		rm.lastTick = now
		rm.tick(r)
	}

	rm.render(r)
}

func (rm *room) reset(r kit.Room) {
	rm.snake = []Point{
		{X: 19, Y: 9},
		{X: 18, Y: 9},
		{X: 17, Y: 9},
		{X: 16, Y: 9},
	}
	rm.entityDir = Point{X: 1, Y: 0}
	rm.lastMovedDir = Point{X: 1, Y: 0}
	rm.score = 0
	rm.gameOver = false
	rm.gameStarted = true
	rm.lastTick = r.Now()
	rm.startedAt = r.Now()

	// Regenerate food and obstacles
	rm.food = rm.randomFreePoint(r, 0)
	rm.obstacles = []Point{}
	for i := 0; i < 3; i++ {
		rm.obstacles = append(rm.obstacles, rm.randomFreePoint(r, 5))
	}
}

func (rm *room) randomFreePoint(r kit.Room, avoidHeadRange int) Point {
	// Attempt up to 100 times to find a free spot
	for attempt := 0; attempt < 100; attempt++ {
		x := r.Rand().Intn(39)
		y := r.Rand().Intn(18)
		p := Point{X: x, Y: y}

		// Check snake collision
		inSnake := false
		for _, sp := range rm.snake {
			if sp == p {
				inSnake = true
				break
			}
		}
		if inSnake {
			continue
		}

		// Check food collision
		if p == rm.food {
			continue
		}

		// Check obstacles collision
		inObstacle := false
		for _, op := range rm.obstacles {
			if op == p {
				inObstacle = true
				break
			}
		}
		if inObstacle {
			continue
		}

		// Avoid snake head range if requested
		if avoidHeadRange > 0 && len(rm.snake) > 0 {
			head := rm.snake[0]
			distX := head.X - p.X
			distY := head.Y - p.Y
			if distX < 0 {
				distX = -distX
			}
			if distY < 0 {
				distY = -distY
			}
			if distX+distY <= avoidHeadRange {
				continue
			}
		}

		return p
	}
	// Fallback
	return Point{X: 10, Y: 10}
}

func (rm *room) tick(r kit.Room) {
	if len(rm.snake) == 0 {
		return
	}

	// Save tail segment position before we shift, in case we eat food and grow
	tail := rm.snake[len(rm.snake)-1]

	// Move the body: shift all elements down
	for i := len(rm.snake) - 1; i > 0; i-- {
		rm.snake[i] = rm.snake[i-1]
	}

	// Update the head position
	rm.snake[0].X += rm.entityDir.X
	rm.snake[0].Y += rm.entityDir.Y

	// Wrap around boundaries (Playfield width is 39 grid units, height is 18)
	if rm.snake[0].X < 0 {
		rm.snake[0].X = 38
	} else if rm.snake[0].X > 38 {
		rm.snake[0].X = 0
	}
	if rm.snake[0].Y < 0 {
		rm.snake[0].Y = 17
	} else if rm.snake[0].Y > 17 {
		rm.snake[0].Y = 0
	}

	// Update last moved direction
	rm.lastMovedDir = rm.entityDir

	// 1. Check self-collision
	for _, sp := range rm.snake[1:] {
		if rm.snake[0] == sp {
			rm.gameOver = true
			return
		}
	}

	// 2. Check obstacle-collision
	for _, op := range rm.obstacles {
		if rm.snake[0] == op {
			rm.gameOver = true
			return
		}
	}

	// 3. Check food-collision
	if rm.snake[0] == rm.food {
		// Grow snake by restoring tail
		rm.snake = append(rm.snake, tail)
		rm.score += 10
		if rm.score > rm.highScore {
			rm.highScore = rm.score
		}
		// Respawn food
		rm.food = rm.randomFreePoint(r, 0)
		// Spawn a new obstacle
		rm.obstacles = append(rm.obstacles, rm.randomFreePoint(r, 4))
	}
}

func (rm *room) render(r kit.Room) {
	if rm.frame == nil {
		rm.frame = kit.NewFrame()
	}
	f := rm.frame
	f.Clear()

	// Styles
	borderStyle := kit.Style{FG: kit.RGB(0x8a, 0x2b, 0xe2)}                         // Neon Violet
	headerStyle := kit.Style{FG: kit.RGB(0x00, 0xff, 0xff), Attr: kit.AttrBold}     // Aqua/Cyan
	valueStyle := kit.Style{FG: kit.RGB(0xff, 0xff, 0xff)}                          // White
	dotStyle := kit.Style{FG: kit.RGB(0x44, 0x44, 0x44)}                            // Dark Gray
	footerStyle := kit.Style{FG: kit.RGB(0x00, 0xe5, 0xff)}
	keyStyle := kit.Style{FG: kit.RGB(0xff, 0x00, 0x7f), Attr: kit.AttrBold}        // Neon Pink

	// 1. Draw Border
	rm.drawBorder(f, borderStyle)

	// 2. Draw Header Content
	f.Text(1, 3, "▲▼ NEON SNAKE ▲▼", headerStyle)
	f.Text(1, 30, "SCORE:", headerStyle)
	f.Text(1, 37, fmt.Sprintf("%04d", rm.score), valueStyle)

	f.Text(1, 46, "HIGH:", headerStyle)
	f.Text(1, 52, fmt.Sprintf("%04d", rm.highScore), valueStyle)

	// Pulsing status effect
	var statusText string
	var statusStyle kit.Style
	if rm.gameOver {
		statusText = "GAME OVER"
		statusStyle = kit.Style{FG: kit.RGB(0xff, 0x00, 0x55), Attr: kit.AttrBold} // Pink/Red
	} else if !rm.gameStarted {
		statusText = "PAUSED"
		statusStyle = kit.Style{FG: kit.RGB(0xff, 0xff, 0x00), Attr: kit.AttrBold} // Yellow
	} else {
		statusText = "PLAYING"
		elapsed := r.Now().Sub(rm.startedAt)
		pulse := (elapsed.Milliseconds() / 150) % 10
		if pulse < 5 {
			statusStyle = kit.Style{FG: kit.RGB(0xff, 0x00, 0x7f), Attr: kit.AttrBold} // Bright Pink
		} else {
			statusStyle = kit.Style{FG: kit.RGB(0xaa, 0x00, 0x55)} // Darker Pink
		}
	}
	f.Text(1, 62, "STATUS: "+statusText, statusStyle)

	// Create lookup for coordinates to skip drawing grid dots
	occupied := make(map[Point]bool)
	for _, p := range rm.snake {
		occupied[p] = true
	}
	occupied[rm.food] = true
	for _, p := range rm.obstacles {
		occupied[p] = true
	}

	// 3. Draw Grid Dots
	for y := 0; y < 18; y++ {
		for x := 0; x < 39; x++ {
			if occupied[Point{X: x, Y: y}] {
				continue
			}
			f.SetRune(3+y, 1+x*2, '·', dotStyle)
		}
	}

	// 4. Draw Food (Pulsing glowing neon star)
	elapsed := r.Now().Sub(rm.startedAt)
	foodPulse := (elapsed.Milliseconds() / 150) % 6
	var foodColor kit.Color
	switch foodPulse {
	case 0, 5:
		foodColor = kit.RGB(0xff, 0x00, 0x55) // Neon pink-red
	case 1, 4:
		foodColor = kit.RGB(0xff, 0x55, 0x7f) // Slightly lighter
	default:
		foodColor = kit.RGB(0xff, 0xaa, 0xcc) // Very light pink pulse
	}
	foodStyle := kit.Style{FG: foodColor, Attr: kit.AttrBold}
	f.SetWide(3+rm.food.Y, 1+rm.food.X*2, '★', foodStyle)

	// 5. Draw Obstacles (Neon Amber triangles)
	obstacleStyle := kit.Style{FG: kit.RGB(0xff, 0x8c, 0x00), Attr: kit.AttrBold} // Dark orange / amber
	for _, op := range rm.obstacles {
		f.SetWide(3+op.Y, 1+op.X*2, '▲', obstacleStyle)
	}

	// 6. Draw Snake (Neon gradient from Lime Green head to Cyan tail)
	n := len(rm.snake)
	for i := n - 1; i >= 0; i-- {
		p := rm.snake[i]
		var segmentStyle kit.Style
		if i == 0 {
			// Head: Bright Lime Green
			segmentStyle = kit.Style{FG: kit.RGB(0x39, 0xff, 0x14)}
		} else {
			// Gradient color: fade from Lime Green (0x39, 0xff, 0x14) to Cyan (0x00, 0xe5, 0xff)
			ratio := float64(i) / float64(n-1)
			rVal := uint8(0x39 - int(ratio*float64(0x39)))
			gVal := uint8(0xff - int(ratio*float64(0xff-0xe5)))
			bVal := uint8(0x14 + int(ratio*float64(0xff-0x14)))
			segmentStyle = kit.Style{FG: kit.RGB(rVal, gVal, bVal)}
		}
		f.SetWide(3+p.Y, 1+p.X*2, '█', segmentStyle)
	}

	// 7. Draw Footer Content
	f.Text(22, 4, "CONTROLS:", footerStyle)
	col := 15
	col = f.Text(22, col, " [", footerStyle)
	col = f.Text(22, col, "WASD", keyStyle)
	col = f.Text(22, col, "/Arrows] Move", footerStyle)

	col = f.Text(22, 45, " [", footerStyle)
	col = f.Text(22, col, "Space", keyStyle)
	col = f.Text(22, col, "] Pause", footerStyle)

	col = f.Text(22, 63, " [", footerStyle)
	col = f.Text(22, col, "Esc", keyStyle)
	col = f.Text(22, col, "] Quit", footerStyle)

	// 8. Draw Game Over Overlay
	if rm.gameOver {
		modalStyle := kit.Style{FG: kit.RGB(0xff, 0x00, 0x55), Attr: kit.AttrBold} // Neon Pink/Red border
		textStyle := kit.Style{FG: kit.RGB(0xff, 0xff, 0xff), Attr: kit.AttrBold}  // White bold
		subTextStyle := kit.Style{FG: kit.RGB(0x8a, 0x2b, 0xe2)}                  // Violet
		infoStyle := kit.Style{FG: kit.RGB(0x00, 0xff, 0xff)}                     // Aqua

		// Draw a box from row 8 to 14, col 20 to 60 (width 40, height 7)
		// Top border of modal
		f.Text(8, 20, "╔══════════════════════════════════════╗", modalStyle)
		f.Text(9, 20, "║                                      ║", modalStyle)
		f.Text(10, 20, "║              GAME OVER               ║", modalStyle)
		f.Text(11, 20, "║                                      ║", modalStyle)
		f.Text(12, 20, "║                                      ║", modalStyle)
		f.Text(13, 20, "║                                      ║", modalStyle)
		f.Text(14, 20, "╚══════════════════════════════════════╝", modalStyle)

		// Set the text content inside the box
		// Pulsing game over text
		gameOverPulse := (elapsed.Milliseconds() / 250) % 2
		var titleStyle kit.Style
		if gameOverPulse == 0 {
			titleStyle = kit.Style{FG: kit.RGB(0xff, 0x00, 0x55), Attr: kit.AttrBold}
		} else {
			titleStyle = kit.Style{FG: kit.RGB(0xff, 0xaa, 0x00), Attr: kit.AttrBold}
		}
		f.Text(10, 35, "GAME OVER", titleStyle)

		// Draw final score
		f.Text(11, 28, fmt.Sprintf("FINAL SCORE: %04d", rm.score), textStyle)
		f.Text(11, 47, "★", infoStyle) // Little star at the end

		// Draw restart instruction
		f.Text(13, 25, "Press [SPACE] to Restart", subTextStyle)
	}

	// Send viewport to each member
	for _, p := range r.Members() {
		r.Send(p, f)
	}
}

func (rm *room) drawBorder(f *kit.Frame, borderStyle kit.Style) {
	// Top border
	f.SetRune(0, 0, '╔', borderStyle)
	for c := 1; c < 79; c++ {
		f.SetRune(0, c, '═', borderStyle)
	}
	f.SetRune(0, 79, '╗', borderStyle)

	// Row 2 divider
	f.SetRune(2, 0, '╠', borderStyle)
	for c := 1; c < 79; c++ {
		f.SetRune(2, c, '═', borderStyle)
	}
	f.SetRune(2, 79, '╣', borderStyle)

	// Row 21 divider
	f.SetRune(21, 0, '╠', borderStyle)
	for c := 1; c < 79; c++ {
		f.SetRune(21, c, '═', borderStyle)
	}
	f.SetRune(21, 79, '╣', borderStyle)

	// Bottom border
	f.SetRune(23, 0, '╚', borderStyle)
	for c := 1; c < 79; c++ {
		f.SetRune(23, c, '═', borderStyle)
	}
	f.SetRune(23, 79, '╝', borderStyle)

	// Vertical borders
	for r := 1; r < 23; r++ {
		if r == 2 || r == 21 {
			continue
		}
		f.SetRune(r, 0, '║', borderStyle)
		f.SetRune(r, 79, '║', borderStyle)
	}
}
