// game — a shellcade game. Run it right now: go run .
package main

import (
	"fmt"
	"math"
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
		Leaderboard: &kit.LeaderboardSpec{
			MetricLabel: "Score",
			Direction:   kit.HigherBetter,
			Aggregation: kit.BestResult,
			Format:      kit.Integer,
		},
	}
}

func (Game) NewRoom(cfg kit.RoomConfig, svc kit.Services) kit.Handler {
	return &room{}
}

type Point struct {
	X, Y int
}

type ScorePopup struct {
	X, Y      int
	Text      string
	Color     kit.Color
	CreatedAt time.Time
}

type Palette struct {
	Name        string
	Border      kit.Color
	Header      kit.Color
	Dot         kit.Color
	SnakeHead   kit.Color
	SnakeTail   kit.Color
	Food        kit.Color
	Obstacle    kit.Color
	Footer      kit.Color
	Key         kit.Color
	ModalBorder kit.Color
}

func getPalettes() []Palette {
	return []Palette{
		{
			Name:        "Cyberpunk",
			Border:      kit.RGB(0x8a, 0x2b, 0xe2), // Neon Violet
			Header:      kit.RGB(0x00, 0xff, 0xff), // Aqua/Cyan
			Dot:         kit.RGB(0x33, 0x33, 0x33), // Dark Gray
			SnakeHead:   kit.RGB(0x39, 0xff, 0x14), // Lime Green
			SnakeTail:   kit.RGB(0x00, 0xe5, 0xff), // Cyan
			Food:        kit.RGB(0xff, 0x00, 0x7f), // Neon Pink
			Obstacle:    kit.RGB(0xff, 0x8c, 0x00), // Neon Orange/Amber
			Footer:      kit.RGB(0x00, 0xe5, 0xff), // Light Cyan
			Key:         kit.RGB(0xff, 0x00, 0x7f), // Neon Pink
			ModalBorder: kit.RGB(0xff, 0x00, 0x55), // Pink/Red
		},
		{
			Name:        "Ocean",
			Border:      kit.RGB(0x00, 0x00, 0xcd), // Deep Blue
			Header:      kit.RGB(0xe0, 0xff, 0xff), // Light Cyan
			Dot:         kit.RGB(0x11, 0x22, 0x44), // Dark Navy
			SnakeHead:   kit.RGB(0x00, 0xff, 0xcc), // Teal/Cyan
			SnakeTail:   kit.RGB(0x00, 0x66, 0xff), // Blue
			Food:        kit.RGB(0xff, 0x7f, 0x50), // Coral
			Obstacle:    kit.RGB(0xff, 0xd7, 0x00), // Gold
			Footer:      kit.RGB(0x00, 0xcd, 0xcd), // Teal
			Key:         kit.RGB(0xff, 0x7f, 0x50), // Coral
			ModalBorder: kit.RGB(0x00, 0xbf, 0xff), // Sky Blue
		},
		{
			Name:        "Sunset",
			Border:      kit.RGB(0xdc, 0x14, 0x3c), // Crimson Red
			Header:      kit.RGB(0xff, 0xd7, 0x00), // Gold
			Dot:         kit.RGB(0x44, 0x22, 0x22), // Dark Rust
			SnakeHead:   kit.RGB(0xff, 0xa5, 0x00), // Yellow-Orange
			SnakeTail:   kit.RGB(0x8b, 0x00, 0x00), // Deep Red
			Food:        kit.RGB(0xda, 0x70, 0xd6), // Neon Purple
			Obstacle:    kit.RGB(0xff, 0x00, 0xff), // Magenta
			Footer:      kit.RGB(0xff, 0xc0, 0xcb), // Light Orange/Pink
			Key:         kit.RGB(0xff, 0xd7, 0x00), // Gold
			ModalBorder: kit.RGB(0xff, 0x8c, 0x00), // Dark Orange
		},
		{
			Name:        "Matrix",
			Border:      kit.RGB(0x00, 0x64, 0x00), // Dark Green
			Header:      kit.RGB(0x00, 0xff, 0x00), // Neon Green
			Dot:         kit.RGB(0x00, 0x22, 0x00), // Very Dark Green
			SnakeHead:   kit.RGB(0xcc, 0xff, 0xcc), // White-Green
			SnakeTail:   kit.RGB(0x32, 0xcd, 0x32), // Lime Green
			Food:        kit.RGB(0xad, 0xff, 0x2f), // Matrix Yellow
			Obstacle:    kit.RGB(0x22, 0x8b, 0x22), // Forest Green
			Footer:      kit.RGB(0x98, 0xfb, 0x98), // Pale Green
			Key:         kit.RGB(0x00, 0xff, 0x00), // Neon Green
			ModalBorder: kit.RGB(0x7f, 0xff, 0x00), // Light Lime
		},
		{
			Name:        "Vaporwave",
			Border:      kit.RGB(0xda, 0x70, 0xd6), // Pastel Purple
			Header:      kit.RGB(0xff, 0xff, 0x00), // Bright Yellow
			Dot:         kit.RGB(0x33, 0x11, 0x44), // Pastel Violet
			SnakeHead:   kit.RGB(0xff, 0x69, 0xb4), // Hot Pink
			SnakeTail:   kit.RGB(0x00, 0xff, 0xff), // Cyan
			Food:        kit.RGB(0xff, 0xa5, 0x00), // Neon Orange
			Obstacle:    kit.RGB(0x4b, 0x00, 0x82), // Deep Indigo
			Footer:      kit.RGB(0xee, 0x82, 0xee), // Soft Purple
			Key:         kit.RGB(0x00, 0xff, 0xff), // Cyan
			ModalBorder: kit.RGB(0xff, 0x00, 0xff), // Magenta
		},
	}
}

// room is one live room. ALL state lives here (and only here).
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

	// Task 4 Aesthetics
	themeIndex      int
	popups          []ScorePopup
	lastCollisionAt time.Time

	// Task 5 Multiplayer & Active player tracking
	activePlayer    kit.Player
	activePlayerSet bool
}

func (rm *room) updateActivePlayer(r kit.Room) {
	members := r.Members()
	if len(members) == 0 {
		return
	}
	// If active player is not set or not in the room anymore, choose the first available member
	found := false
	if rm.activePlayerSet {
		for _, m := range members {
			if m.AccountID == rm.activePlayer.AccountID {
				found = true
				break
			}
		}
	}
	if !found {
		rm.activePlayer = members[0]
		rm.activePlayerSet = true
	}
}

func (rm *room) getActivePlayer(r kit.Room) kit.Player {
	rm.updateActivePlayer(r)
	return rm.activePlayer
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
	rm.themeIndex = 0
	rm.popups = []ScorePopup{}
	rm.lastCollisionAt = time.Time{}
	rm.activePlayer = kit.Player{}
	rm.activePlayerSet = false

	// Generate initial food & obstacles
	rm.food = rm.randomFreePoint(r, 0)
	rm.obstacles = []Point{}
	for i := 0; i < 3; i++ {
		rm.obstacles = append(rm.obstacles, rm.randomFreePoint(r, 5))
	}
}

func (rm *room) OnJoin(r kit.Room, p kit.Player) {
	rm.updateActivePlayer(r)
	rm.render(r)
}

func (rm *room) OnInput(r kit.Room, p kit.Player, in kit.Input) {
	rm.activePlayer = p
	rm.activePlayerSet = true
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

	// Switch Theme support
	if in.Kind == kit.InputRune && (in.Rune == 't' || in.Rune == 'T') {
		palettes := getPalettes()
		rm.themeIndex = (rm.themeIndex + 1) % len(palettes)
	}

	rm.render(r)
}

// OnWake is the host heartbeat.
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
	rm.tickRate = 150 * time.Millisecond
	rm.popups = []ScorePopup{}
	rm.lastCollisionAt = time.Time{}

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
			rm.lastCollisionAt = r.Now()
			r.Post(kit.Result{
				Rankings: []kit.PlayerResult{
					{
						Player: rm.getActivePlayer(r),
						Metric: rm.score,
						Status: kit.StatusFinished,
					},
				},
			})
			return
		}
	}

	// 2. Check obstacle-collision
	for _, op := range rm.obstacles {
		if rm.snake[0] == op {
			rm.gameOver = true
			rm.lastCollisionAt = r.Now()
			r.Post(kit.Result{
				Rankings: []kit.PlayerResult{
					{
						Player: rm.getActivePlayer(r),
						Metric: rm.score,
						Status: kit.StatusFinished,
					},
				},
			})
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

		// Dynamic difficulty speed up: decrease tick rate as score increases (clamp at 60ms)
		speedMs := 150 - (rm.score/10)*5
		if speedMs < 60 {
			speedMs = 60
		}
		rm.tickRate = time.Duration(speedMs) * time.Millisecond

		// Add score popup
		palettes := getPalettes()
		theme := palettes[rm.themeIndex]
		rm.popups = append(rm.popups, ScorePopup{
			X:         rm.food.X,
			Y:         rm.food.Y,
			Text:      "+10",
			Color:     theme.Food,
			CreatedAt: r.Now(),
		})

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

	now := r.Now()
	palettes := getPalettes()
	theme := palettes[rm.themeIndex]

	// Styles
	headerStyle := kit.Style{FG: theme.Header, Attr: kit.AttrBold}
	valueStyle := kit.Style{FG: kit.RGB(0xff, 0xff, 0xff)}
	dotStyle := kit.Style{FG: theme.Dot}
	footerStyle := kit.Style{FG: theme.Footer}
	keyStyle := kit.Style{FG: theme.Key, Attr: kit.AttrBold}

	// 1. Draw Border (Animates dynamically)
	rm.drawBorder(f, now)

	rm.updateActivePlayer(r)

	// 2. Draw Header Content
	f.Text(1, 2, "▲▼ NEON SNAKE ▲▼", headerStyle)
	f.Text(1, 19, "SCORE:", headerStyle)
	f.Text(1, 25, fmt.Sprintf("%04d", rm.score), valueStyle)

	f.Text(1, 30, "HIGH:", headerStyle)
	f.Text(1, 35, fmt.Sprintf("%04d", rm.highScore), valueStyle)

	// In multiplayer, show active player handle
	var themeStartCol int
	if len(r.Members()) > 1 && rm.activePlayerSet {
		f.Text(1, 40, "SEAT:", headerStyle)
		f.Text(1, 45, rm.activePlayer.Handle, valueStyle)
		themeStartCol = 52
	} else {
		themeStartCol = 40
	}
	f.Text(1, themeStartCol, "THEME:", headerStyle)
	f.Text(1, themeStartCol+7, theme.Name, valueStyle)

	// Pulsing status effect (right-aligned to column 78)
	var statusText string
	var statusStyle kit.Style
	if rm.gameOver {
		statusText = "GAME OVER"
		statusStyle = kit.Style{FG: kit.RGB(0xff, 0x00, 0x55), Attr: kit.AttrBold}
	} else if !rm.gameStarted {
		statusText = "PAUSED"
		statusStyle = kit.Style{FG: kit.RGB(0xff, 0xff, 0x00), Attr: kit.AttrBold}
	} else {
		statusText = "PLAYING"
		elapsed := now.Sub(rm.startedAt)
		pulse := (elapsed.Milliseconds() / 150) % 10
		if pulse < 5 {
			statusStyle = kit.Style{FG: theme.Key, Attr: kit.AttrBold}
		} else {
			statusStyle = kit.Style{FG: brightenColor(theme.Key, 0.6)}
		}
	}
	f.Text(1, 78-len(statusText), statusText, statusStyle)

	// Create lookup for coordinates to skip drawing grid dots
	occupied := make(map[Point]bool)
	for _, sp := range rm.snake {
		occupied[sp] = true
	}
	occupied[rm.food] = true
	for _, op := range rm.obstacles {
		occupied[op] = true
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

	// 4. Draw Food (Pulsing glowing neon star, rotating/twinkling glyph)
	elapsed := now.Sub(rm.startedAt)
	foodPulse := (elapsed.Milliseconds() / 150) % 6
	var foodColor kit.Color
	switch foodPulse {
	case 0, 5:
		foodColor = theme.Food
	case 1, 4:
		foodColor = brightenColor(theme.Food, 1.2)
	default:
		foodColor = brightenColor(theme.Food, 1.5)
	}
	foodStyle := kit.Style{FG: foodColor, Attr: kit.AttrBold}

	foodGlyphs := []rune{'★', '☆', '✦', '✧'}
	foodGlyphIdx := (elapsed.Milliseconds() / 250) % int64(len(foodGlyphs))
	foodGlyph := foodGlyphs[foodGlyphIdx]
	f.SetWide(3+rm.food.Y, 1+rm.food.X*2, foodGlyph, foodStyle)

	// 5. Draw Obstacles (Neon triangles with warning flash when head is close)
	for _, op := range rm.obstacles {
		head := rm.snake[0]
		distX := head.X - op.X
		if distX < 0 {
			distX = -distX
		}
		distY := head.Y - op.Y
		if distY < 0 {
			distY = -distY
		}
		dist := distX + distY

		var obstacleStyle kit.Style
		if dist <= 2 && rm.gameStarted && !rm.gameOver {
			flashCycle := (elapsed.Milliseconds() / 100) % 2
			if flashCycle == 0 {
				obstacleStyle = kit.Style{FG: kit.RGB(0xff, 0xff, 0xff), Attr: kit.AttrBold}
			} else {
				obstacleStyle = kit.Style{FG: kit.RGB(0xff, 0x00, 0x55), Attr: kit.AttrBold}
			}
		} else {
			obstacleStyle = kit.Style{FG: theme.Obstacle, Attr: kit.AttrBold}
		}

		f.SetWide(3+op.Y, 1+op.X*2, '▲', obstacleStyle)
	}

	// 6. Draw Snake (Neon gradient from SnakeHead to SnakeTail, flowing dynamically)
	n := len(rm.snake)
	timeShift := float64(now.Sub(rm.startedAt).Milliseconds()%2000) / 2000.0
	for i := n - 1; i >= 0; i-- {
		p := rm.snake[i]
		var segmentStyle kit.Style
		if i == 0 {
			headPulse := 0.85 + 0.15*math.Sin(float64(now.Sub(rm.startedAt).Milliseconds())*0.006)
			segmentStyle = kit.Style{FG: brightenColor(theme.SnakeHead, headPulse)}
		} else {
			ratio := float64(i) / float64(n-1)
			shiftedRatio := ratio + timeShift
			if shiftedRatio > 1.0 {
				shiftedRatio -= 1.0
			}
			segmentStyle = kit.Style{FG: interpolateColor(theme.SnakeHead, theme.SnakeTail, shiftedRatio)}
		}
		f.SetWide(3+p.Y, 1+p.X*2, '█', segmentStyle)
	}

	// 7. Draw Floating Score Popups
	var activePopups []ScorePopup
	for _, p := range rm.popups {
		age := now.Sub(p.CreatedAt)
		if age < 1000*time.Millisecond {
			yOffset := int(age.Milliseconds() / 200)
			drawY := 3 + p.Y - yOffset
			drawX := 1 + p.X*2
			if drawX > 75 {
				drawX = 75
			}
			if drawY > 2 && drawY < 21 {
				style := kit.Style{FG: p.Color, Attr: kit.AttrBold}
				if age > 500*time.Millisecond {
					style.Attr = kit.AttrDim
				}
				f.Text(drawY, drawX, p.Text, style)
			}
			activePopups = append(activePopups, p)
		}
	}
	rm.popups = activePopups

	// 7.5 Draw Divider on row 21 with active player text in multiplayer
	if len(r.Members()) > 1 && rm.activePlayerSet {
		activeText := fmt.Sprintf(" ACTIVE SEAT: %s ", rm.activePlayer.Handle)
		rm.drawDividerWithText(f, 21, activeText, now)
	}

	// 8. Draw Footer Content
	f.Text(22, 4, "CONTROLS:", footerStyle)
	col := 15
	col = f.Text(22, col, " [", footerStyle)
	col = f.Text(22, col, "WASD", keyStyle)
	col = f.Text(22, col, "/Arrows] Move", footerStyle)

	col = f.Text(22, 38, " [", footerStyle)
	col = f.Text(22, col, "T", keyStyle)
	col = f.Text(22, col, "] Theme", footerStyle)

	col = f.Text(22, 49, " [", footerStyle)
	col = f.Text(22, col, "Space", keyStyle)
	col = f.Text(22, col, "] Pause", footerStyle)

	col = f.Text(22, 63, " [", footerStyle)
	col = f.Text(22, col, "Esc", keyStyle)
	col = f.Text(22, col, "] Quit", footerStyle)

	// 9. Draw Game Over Overlay
	if rm.gameOver {
		modalStyle := kit.Style{FG: theme.ModalBorder, Attr: kit.AttrBold}
		textStyle := kit.Style{FG: kit.RGB(0xff, 0xff, 0xff), Attr: kit.AttrBold}
		subTextStyle := kit.Style{FG: theme.Border}
		infoStyle := kit.Style{FG: theme.Header}

		f.Text(8, 20, "╔══════════════════════════════════════╗", modalStyle)
		f.Text(9, 20, "║                                      ║", modalStyle)
		f.Text(10, 20, "║              GAME OVER               ║", modalStyle)
		f.Text(11, 20, "║                                      ║", modalStyle)
		f.Text(12, 20, "║                                      ║", modalStyle)
		f.Text(13, 20, "║                                      ║", modalStyle)
		f.Text(14, 20, "╚══════════════════════════════════════╝", modalStyle)

		gameOverPulse := (elapsed.Milliseconds() / 250) % 2
		var titleStyle kit.Style
		if gameOverPulse == 0 {
			titleStyle = kit.Style{FG: theme.ModalBorder, Attr: kit.AttrBold}
		} else {
			titleStyle = kit.Style{FG: theme.Key, Attr: kit.AttrBold}
		}
		f.Text(10, 35, "GAME OVER", titleStyle)

		var scoreText string
		if len(r.Members()) > 1 && rm.activePlayerSet {
			scoreText = fmt.Sprintf("FINAL SCORE: %04d (%s)", rm.score, rm.activePlayer.Handle)
		} else {
			scoreText = fmt.Sprintf("FINAL SCORE: %04d", rm.score)
		}
		f.Text(11, 24, scoreText, textStyle)
		f.Text(11, 24+len(scoreText)+2, "★", infoStyle)

		f.Text(13, 25, "Press [SPACE] to Restart", subTextStyle)
	}

	// Send viewport to each member
	for _, p := range r.Members() {
		r.Send(p, f)
	}
}

func (rm *room) getBorderStyle(r, c int, now time.Time) kit.Style {
	palettes := getPalettes()
	theme := palettes[rm.themeIndex]

	// Check if flashing from recent collision
	if !rm.lastCollisionAt.IsZero() && now.Sub(rm.lastCollisionAt) < 300*time.Millisecond {
		flashCycle := (now.Sub(rm.lastCollisionAt).Milliseconds() / 75) % 2
		if flashCycle == 0 {
			return kit.Style{FG: kit.RGB(0xff, 0xff, 0xff), Attr: kit.AttrBold}
		} else {
			return kit.Style{FG: kit.RGB(0xff, 0x00, 0x55), Attr: kit.AttrBold}
		}
	}

	baseStyle := kit.Style{FG: theme.Border}
	// Dynamic wave animation along the border
	timeShift := int(now.Sub(rm.startedAt).Milliseconds() / 80)
	pos := r + c - timeShift
	if pos < 0 {
		pos = -pos
	}
	wave := pos % 16
	if wave < 3 {
		return kit.Style{FG: brightenColor(theme.Border, 1.8), Attr: kit.AttrBold}
	}
	return baseStyle
}

func (rm *room) drawBorder(f *kit.Frame, now time.Time) {
	// Top border
	f.SetRune(0, 0, '╔', rm.getBorderStyle(0, 0, now))
	for c := 1; c < 79; c++ {
		f.SetRune(0, c, '═', rm.getBorderStyle(0, c, now))
	}
	f.SetRune(0, 79, '╗', rm.getBorderStyle(0, 79, now))

	// Row 2 divider
	f.SetRune(2, 0, '╠', rm.getBorderStyle(2, 0, now))
	for c := 1; c < 79; c++ {
		f.SetRune(2, c, '═', rm.getBorderStyle(2, c, now))
	}
	f.SetRune(2, 79, '╣', rm.getBorderStyle(2, 79, now))

	// Row 21 divider
	f.SetRune(21, 0, '╠', rm.getBorderStyle(21, 0, now))
	for c := 1; c < 79; c++ {
		f.SetRune(21, c, '═', rm.getBorderStyle(21, c, now))
	}
	f.SetRune(21, 79, '╣', rm.getBorderStyle(21, 79, now))

	// Bottom border
	f.SetRune(23, 0, '╚', rm.getBorderStyle(23, 0, now))
	for c := 1; c < 79; c++ {
		f.SetRune(23, c, '═', rm.getBorderStyle(23, c, now))
	}
	f.SetRune(23, 79, '╝', rm.getBorderStyle(23, 79, now))

	// Vertical borders
	for r := 1; r < 23; r++ {
		if r == 2 || r == 21 {
			continue
		}
		f.SetRune(r, 0, '║', rm.getBorderStyle(r, 0, now))
		f.SetRune(r, 79, '║', rm.getBorderStyle(r, 79, now))
	}
}

func (rm *room) drawDividerWithText(f *kit.Frame, row int, text string, now time.Time) {
	f.SetRune(row, 0, '╠', rm.getBorderStyle(row, 0, now))
	f.SetRune(row, 79, '╣', rm.getBorderStyle(row, 79, now))

	textLen := len(text)
	startCol := 40 - textLen/2
	endCol := startCol + textLen

	for c := 1; c < 79; c++ {
		if c >= startCol && c < endCol {
			char := rune(text[c-startCol])
			f.SetRune(row, c, char, kit.Style{FG: getPalettes()[rm.themeIndex].Header, Attr: kit.AttrBold})
		} else {
			f.SetRune(row, c, '═', rm.getBorderStyle(row, c, now))
		}
	}
}

func brightenColor(c kit.Color, factor float64) kit.Color {
	r, g, b := c.RGBVals()
	newR := float64(r) * factor
	newG := float64(g) * factor
	newB := float64(b) * factor
	if newR > 255 {
		newR = 255
	}
	if newG > 255 {
		newG = 255
	}
	if newB > 255 {
		newB = 255
	}
	return kit.RGB(uint8(newR), uint8(newG), uint8(newB))
}

func interpolateColor(c1, c2 kit.Color, t float64) kit.Color {
	r1, g1, b1 := c1.RGBVals()
	r2, g2, b2 := c2.RGBVals()
	r := uint8(float64(r1) + t*float64(int(r2)-int(r1)))
	g := uint8(float64(g1) + t*float64(int(g2)-int(g1)))
	b := uint8(float64(b1) + t*float64(int(b2)-int(b1)))
	return kit.RGB(r, g, b)
}
