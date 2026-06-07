// game — a shellcade game. Run it right now: go run .
package main

import (
	"context"
	"fmt"
	"math"
	"strconv"
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
	return &room{
		services: svc,
	}
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
	Snake2Head  kit.Color
	Snake2Tail  kit.Color
	Food        kit.Color
	Obstacle    kit.Color
	Footer      kit.Color
	Key         kit.Color
	ModalBorder kit.Color
	Hazard      kit.Color
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
			Snake2Head:  kit.RGB(0xff, 0x00, 0x7f), // Neon Pink
			Snake2Tail:  kit.RGB(0xff, 0xa5, 0x00), // Neon Orange
			Food:        kit.RGB(0xff, 0x00, 0x7f), // Neon Pink
			Obstacle:    kit.RGB(0xff, 0x8c, 0x00), // Neon Orange/Amber
			Footer:      kit.RGB(0x00, 0xe5, 0xff), // Light Cyan
			Key:         kit.RGB(0xff, 0x00, 0x7f), // Neon Pink
			ModalBorder: kit.RGB(0xff, 0x00, 0x55), // Pink/Red
			Hazard:      kit.RGB(0xff, 0x00, 0xff), // Neon Magenta
		},
		{
			Name:        "Ocean",
			Border:      kit.RGB(0x00, 0x00, 0xcd), // Deep Blue
			Header:      kit.RGB(0xe0, 0xff, 0xff), // Light Cyan
			Dot:         kit.RGB(0x11, 0x22, 0x44), // Dark Navy
			SnakeHead:   kit.RGB(0x00, 0xff, 0xcc), // Teal/Cyan
			SnakeTail:   kit.RGB(0x00, 0x66, 0xff), // Blue
			Snake2Head:  kit.RGB(0xff, 0x7f, 0x50), // Coral
			Snake2Tail:  kit.RGB(0xff, 0xd7, 0x00), // Gold
			Food:        kit.RGB(0xff, 0x7f, 0x50), // Coral
			Obstacle:    kit.RGB(0xff, 0xd7, 0x00), // Gold
			Footer:      kit.RGB(0x00, 0xcd, 0xcd), // Teal
			Key:         kit.RGB(0xff, 0x7f, 0x50), // Coral
			ModalBorder: kit.RGB(0x00, 0xbf, 0xff), // Sky Blue
			Hazard:      kit.RGB(0xff, 0x55, 0x00), // Bright Orange-Red
		},
		{
			Name:        "Sunset",
			Border:      kit.RGB(0xdc, 0x14, 0x3c), // Crimson Red
			Header:      kit.RGB(0xff, 0xd7, 0x00), // Gold
			Dot:         kit.RGB(0x44, 0x22, 0x22), // Dark Rust
			SnakeHead:   kit.RGB(0xff, 0xa5, 0x00), // Yellow-Orange
			SnakeTail:   kit.RGB(0x8b, 0x00, 0x00), // Deep Red
			Snake2Head:  kit.RGB(0xda, 0x70, 0xd6), // Neon Purple/Orchid
			Snake2Tail:  kit.RGB(0xdc, 0x14, 0x3c), // Crimson
			Food:        kit.RGB(0xda, 0x70, 0xd6), // Neon Purple
			Obstacle:    kit.RGB(0xff, 0x00, 0xff), // Magenta
			Footer:      kit.RGB(0xff, 0xc0, 0xcb), // Light Orange/Pink
			Key:         kit.RGB(0xff, 0xd7, 0x00), // Gold
			ModalBorder: kit.RGB(0xff, 0x8c, 0x00), // Dark Orange
			Hazard:      kit.RGB(0xff, 0x24, 0x00), // Scarlet/Fiery Red
		},
		{
			Name:        "Matrix",
			Border:      kit.RGB(0x00, 0x64, 0x00), // Dark Green
			Header:      kit.RGB(0x00, 0xff, 0x00), // Neon Green
			Dot:         kit.RGB(0x00, 0x22, 0x00), // Very Dark Green
			SnakeHead:   kit.RGB(0xcc, 0xff, 0xcc), // White-Green
			SnakeTail:   kit.RGB(0x32, 0xcd, 0x32), // Lime Green
			Snake2Head:  kit.RGB(0xff, 0xff, 0x00), // Yellow
			Snake2Tail:  kit.RGB(0x00, 0x64, 0x00), // Forest Green
			Food:        kit.RGB(0xad, 0xff, 0x2f), // Matrix Yellow
			Obstacle:    kit.RGB(0x22, 0x8b, 0x22), // Forest Green
			Footer:      kit.RGB(0x98, 0xfb, 0x98), // Pale Green
			Key:         kit.RGB(0x00, 0xff, 0x00), // Neon Green
			ModalBorder: kit.RGB(0x7f, 0xff, 0x00), // Light Lime
			Hazard:      kit.RGB(0x00, 0xff, 0xff), // Neon Cyan (contrasts green)
		},
		{
			Name:        "Vaporwave",
			Border:      kit.RGB(0xda, 0x70, 0xd6), // Pastel Purple
			Header:      kit.RGB(0xff, 0xff, 0x00), // Bright Yellow
			Dot:         kit.RGB(0x33, 0x11, 0x44), // Pastel Violet
			SnakeHead:   kit.RGB(0xff, 0x69, 0xb4), // Hot Pink
			SnakeTail:   kit.RGB(0x00, 0xff, 0xff), // Cyan
			Snake2Head:  kit.RGB(0xee, 0x82, 0xee), // Soft Purple
			Snake2Tail:  kit.RGB(0xff, 0xff, 0x00), // Yellow
			Food:        kit.RGB(0xff, 0xa5, 0x00), // Neon Orange
			Obstacle:    kit.RGB(0x4b, 0x00, 0x82), // Deep Indigo
			Footer:      kit.RGB(0xee, 0x82, 0xee), // Soft Purple
			Key:         kit.RGB(0x00, 0xff, 0xff), // Cyan
			ModalBorder: kit.RGB(0xff, 0x00, 0xff), // Magenta
			Hazard:      kit.RGB(0x00, 0xff, 0x7f), // Spring Green
		},
	}
}

type GameMode int

const (
	ModeClassic GameMode = iota
	ModeHazard
	ModeMaze
	ModePortal
)

type Hazard struct {
	Pos        Point
	Dir        Point
	MinX, MaxX int
	MinY, MaxY int
}

// room is one live room. ALL state lives here (and only here).
type room struct {
	kit.Base
	services    kit.Services
	pb1         int
	pb2         int
	newPB1      bool
	newPB2      bool
	frame       *kit.Frame
	lastTick    time.Time
	tickRate    time.Duration
	score1      int
	score2      int
	highScore   int
	gameStarted bool
	gameOver    bool

	snake1        []Point
	entityDir1    Point
	lastMovedDir1 Point

	snake2        []Point
	entityDir2    Point
	lastMovedDir2 Point

	crashed1 bool
	crashed2 bool

	food      Point
	obstacles []Point

	startedAt time.Time

	// Task 4 Aesthetics
	themeIndex      int
	popups          []ScorePopup
	lastCollisionAt time.Time

	// Task 5 Multiplayer & Active player tracking
	activePlayer    kit.Player
	activePlayerSet bool

	gameMode GameMode
	hazards  []Hazard

	// Task 8 Power-ups
	powerUpPos       Point
	powerUpType      string // "SHIELD" or "FREEZE"
	powerUpActive    bool
	powerUpSpawnedAt time.Time

	p1PowerUpType   string
	p1PowerUpExpiry time.Time

	p2PowerUpType   string
	p2PowerUpExpiry time.Time

	tickCount int
	lastWake  time.Time
	p2IsBot   bool

	portalA Point
	portalB Point
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

func (rm *room) initHazards() {
	rm.hazards = []Hazard{}
	if rm.gameMode == ModeClassic {
		return
	}
	if rm.gameMode == ModeHazard {
		rm.hazards = []Hazard{
			{Pos: Point{X: 5, Y: 3}, Dir: Point{X: 1, Y: 0}, MinX: 5, MaxX: 33, MinY: 3, MaxY: 3},
			{Pos: Point{X: 33, Y: 14}, Dir: Point{X: -1, Y: 0}, MinX: 5, MaxX: 33, MinY: 14, MaxY: 14},
			{Pos: Point{X: 8, Y: 4}, Dir: Point{X: 0, Y: 1}, MinX: 8, MaxX: 8, MinY: 4, MaxY: 13},
			{Pos: Point{X: 30, Y: 13}, Dir: Point{X: 0, Y: -1}, MinX: 30, MaxX: 30, MinY: 4, MaxY: 13},
		}
	} else if rm.gameMode == ModeMaze {
		rm.hazards = []Hazard{
			{Pos: Point{X: 5, Y: 7}, Dir: Point{X: 1, Y: 0}, MinX: 5, MaxX: 33, MinY: 7, MaxY: 7},
			{Pos: Point{X: 33, Y: 10}, Dir: Point{X: -1, Y: 0}, MinX: 5, MaxX: 33, MinY: 10, MaxY: 10},
			{Pos: Point{X: 5, Y: 2}, Dir: Point{X: 0, Y: 1}, MinX: 5, MaxX: 5, MinY: 2, MaxY: 15},
			{Pos: Point{X: 33, Y: 15}, Dir: Point{X: 0, Y: -1}, MinX: 33, MaxX: 33, MinY: 2, MaxY: 15},
		}
	}
}

func (rm *room) isMazeWall(p Point) bool {
	if rm.gameMode != ModeMaze {
		return false
	}
	// Wall 1 & 3: X from 8 to 14, Y = 4 or 13
	if p.X >= 8 && p.X <= 14 && (p.Y == 4 || p.Y == 13) {
		return true
	}
	// Wall 2 & 4: X from 24 to 30, Y = 4 or 13
	if p.X >= 24 && p.X <= 30 && (p.Y == 4 || p.Y == 13) {
		return true
	}
	// Wall 5 & 6: X = 19, Y from 2 to 5 or 12 to 15
	if p.X == 19 && ((p.Y >= 2 && p.Y <= 5) || (p.Y >= 12 && p.Y <= 15)) {
		return true
	}
	return false
}

func (rm *room) getModeName() string {
	switch rm.gameMode {
	case ModeClassic:
		return "CLASSIC"
	case ModeHazard:
		return "HAZARDS"
	case ModeMaze:
		return "MAZE"
	case ModePortal:
		return "PORTALS"
	default:
		return "CLASSIC"
	}
}

func (rm *room) OnStart(r kit.Room) {
	r.SetInputContext(kit.CtxNav)
	rm.lastTick = r.Now()
	rm.startedAt = r.Now()
	rm.tickRate = 150 * time.Millisecond
	rm.snake1 = []Point{
		{X: 10, Y: 9},
		{X: 9, Y: 9},
		{X: 8, Y: 9},
		{X: 7, Y: 9},
	}
	rm.entityDir1 = Point{X: 1, Y: 0}
	rm.lastMovedDir1 = Point{X: 1, Y: 0}

	rm.snake2 = []Point{
		{X: 28, Y: 9},
		{X: 29, Y: 9},
		{X: 30, Y: 9},
		{X: 31, Y: 9},
	}
	rm.entityDir2 = Point{X: -1, Y: 0}
	rm.lastMovedDir2 = Point{X: -1, Y: 0}

	rm.crashed1 = false
	rm.crashed2 = false

	rm.gameStarted = true
	rm.score1 = 0
	rm.score2 = 0
	rm.gameOver = false
	rm.themeIndex = 0
	rm.gameMode = ModeClassic
	rm.popups = []ScorePopup{}
	rm.lastCollisionAt = time.Time{}
	rm.activePlayer = kit.Player{}
	rm.activePlayerSet = false

	rm.powerUpActive = false
	rm.p1PowerUpType = ""
	rm.p1PowerUpExpiry = time.Time{}
	rm.p2PowerUpType = ""
	rm.p2PowerUpExpiry = time.Time{}
	rm.tickCount = 0
	rm.lastWake = r.Now()
	rm.p2IsBot = false
	rm.portalA = Point{X: 9, Y: 9}
	rm.portalB = Point{X: 29, Y: 9}

	// Initialize hazards
	rm.initHazards()

	// Generate initial food & obstacles
	rm.food = rm.randomFreePoint(r, 0)
	rm.obstacles = []Point{}
	for i := 0; i < 3; i++ {
		rm.obstacles = append(rm.obstacles, rm.randomFreePoint(r, 5))
	}
}

func (rm *room) OnJoin(r kit.Room, p kit.Player) {
	rm.updateActivePlayer(r)
	if len(r.Members()) >= 2 {
		rm.p2IsBot = false
	}
	rm.loadPersonalBests(r)
	rm.render(r)
}

func (rm *room) OnInput(r kit.Room, p kit.Player, in kit.Input) {
	rm.activePlayer = p
	rm.activePlayerSet = true

	// Identify player index
	members := r.Members()
	isPlayer1 := true
	isPlayer2 := false
	if len(members) >= 2 {
		if p.AccountID == members[0].AccountID {
			isPlayer1 = true
			isPlayer2 = false
		} else if p.AccountID == members[1].AccountID {
			isPlayer1 = false
			isPlayer2 = true
		} else {
			isPlayer1 = false
			isPlayer2 = false
		}
	}

	// Handle steering keys concurrently
	if !rm.gameOver && rm.gameStarted {
		if in.Kind == kit.InputRune {
			// Player 1 controls Snake 1 using WASD
			if isPlayer1 {
				switch in.Rune {
				case 'w', 'W':
					if rm.lastMovedDir1.Y != 1 {
						rm.entityDir1 = Point{X: 0, Y: -1}
					}
				case 's', 'S':
					if rm.lastMovedDir1.Y != -1 {
						rm.entityDir1 = Point{X: 0, Y: 1}
					}
				case 'a', 'A':
					if rm.lastMovedDir1.X != 1 {
						rm.entityDir1 = Point{X: -1, Y: 0}
					}
				case 'd', 'D':
					if rm.lastMovedDir1.X != -1 {
						rm.entityDir1 = Point{X: 1, Y: 0}
					}
				}
			}
		} else if in.Kind == kit.InputKey {
			// Player 2 controls Snake 2 using Arrow keys.
			// In single-player co-op (len(members) < 2), Player 1 controls Snake 2 using Arrow keys.
			canControlSnake2 := isPlayer2 || (isPlayer1 && len(members) < 2)
			if canControlSnake2 {
				if rm.p2IsBot {
					rm.p2IsBot = false
				}
				switch in.Key {
				case kit.KeyUp:
					if rm.lastMovedDir2.Y != 1 {
						rm.entityDir2 = Point{X: 0, Y: -1}
					}
				case kit.KeyDown:
					if rm.lastMovedDir2.Y != -1 {
						rm.entityDir2 = Point{X: 0, Y: 1}
					}
				case kit.KeyLeft:
					if rm.lastMovedDir2.X != 1 {
						rm.entityDir2 = Point{X: -1, Y: 0}
					}
				case kit.KeyRight:
					if rm.lastMovedDir2.X != -1 {
						rm.entityDir2 = Point{X: 1, Y: 0}
					}
				}
			}
		}
	}

	action := kit.Resolve(in, kit.CtxNav)
	switch action {
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

	// Switch Mode support
	if in.Kind == kit.InputRune && (in.Rune == 'm' || in.Rune == 'M') {
		rm.gameMode = (rm.gameMode + 1) % 4
		rm.reset(r)
	}

	// Switch Bot support
	if in.Kind == kit.InputRune && (in.Rune == 'b' || in.Rune == 'B') {
		rm.p2IsBot = !rm.p2IsBot
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
	if rm.lastWake.IsZero() {
		rm.lastWake = now
	}

	elapsedSinceLastWake := now.Sub(rm.lastWake)
	rm.lastWake = now

	// If paused or game over, shift power-up timers forward by the elapsed real time
	if (!rm.gameStarted || rm.gameOver) && elapsedSinceLastWake > 0 {
		if !rm.p1PowerUpExpiry.IsZero() {
			rm.p1PowerUpExpiry = rm.p1PowerUpExpiry.Add(elapsedSinceLastWake)
		}
		if !rm.p2PowerUpExpiry.IsZero() {
			rm.p2PowerUpExpiry = rm.p2PowerUpExpiry.Add(elapsedSinceLastWake)
		}
		if rm.powerUpActive && !rm.powerUpSpawnedAt.IsZero() {
			rm.powerUpSpawnedAt = rm.powerUpSpawnedAt.Add(elapsedSinceLastWake)
		}
	}

	// Advance game state based on tickRate
	if rm.gameStarted && !rm.gameOver && now.Sub(rm.lastTick) >= rm.tickRate {
		rm.lastTick = now
		rm.tick(r)
	}

	rm.render(r)
}

func (rm *room) loadPersonalBests(r kit.Room) {
	if rm.services.Accounts == nil {
		return
	}
	members := r.Members()
	if len(members) >= 1 {
		p1Acct := rm.services.Accounts.For(members[0])
		if p1Acct != nil {
			store := p1Acct.Store()
			if store != nil {
				val, ok, err := store.Get(context.Background(), "personal_best")
				if err == nil && ok {
					rm.pb1, _ = strconv.Atoi(string(val))
				} else {
					rm.pb1 = 0
				}
			}
		}
	}
	if len(members) >= 2 {
		p2Acct := rm.services.Accounts.For(members[1])
		if p2Acct != nil {
			store := p2Acct.Store()
			if store != nil {
				val, ok, err := store.Get(context.Background(), "personal_best")
				if err == nil && ok {
					rm.pb2, _ = strconv.Atoi(string(val))
				} else {
					rm.pb2 = 0
				}
			}
		}
	}
}

func (rm *room) savePersonalBests(r kit.Room) {
	if rm.services.Accounts == nil {
		return
	}
	members := r.Members()
	if len(members) >= 2 {
		// Multiplayer
		if rm.score1 > rm.pb1 {
			rm.pb1 = rm.score1
			rm.newPB1 = true
			p1Acct := rm.services.Accounts.For(members[0])
			if p1Acct != nil {
				store := p1Acct.Store()
				if store != nil {
					_ = store.Set(context.Background(), "personal_best", []byte(strconv.Itoa(rm.score1)), kit.MergeMax)
				}
			}
		}
		if rm.score2 > rm.pb2 {
			rm.pb2 = rm.score2
			rm.newPB2 = true
			p2Acct := rm.services.Accounts.For(members[1])
			if p2Acct != nil {
				store := p2Acct.Store()
				if store != nil {
					_ = store.Set(context.Background(), "personal_best", []byte(strconv.Itoa(rm.score2)), kit.MergeMax)
				}
			}
		}
	} else if len(members) == 1 {
		// Single player / Co-op
		maxScore := rm.score1
		if rm.score2 > maxScore {
			maxScore = rm.score2
		}
		if maxScore > rm.pb1 {
			rm.pb1 = maxScore
			rm.newPB1 = true
			p1Acct := rm.services.Accounts.For(members[0])
			if p1Acct != nil {
				store := p1Acct.Store()
				if store != nil {
					_ = store.Set(context.Background(), "personal_best", []byte(strconv.Itoa(maxScore)), kit.MergeMax)
				}
			}
		}
	}
}

func (rm *room) reset(r kit.Room) {
	rm.snake1 = []Point{
		{X: 10, Y: 9},
		{X: 9, Y: 9},
		{X: 8, Y: 9},
		{X: 7, Y: 9},
	}
	rm.entityDir1 = Point{X: 1, Y: 0}
	rm.lastMovedDir1 = Point{X: 1, Y: 0}

	rm.snake2 = []Point{
		{X: 28, Y: 9},
		{X: 29, Y: 9},
		{X: 30, Y: 9},
		{X: 31, Y: 9},
	}
	rm.entityDir2 = Point{X: -1, Y: 0}
	rm.lastMovedDir2 = Point{X: -1, Y: 0}

	rm.crashed1 = false
	rm.crashed2 = false

	rm.score1 = 0
	rm.score2 = 0
	rm.gameOver = false
	rm.gameStarted = true
	rm.lastTick = r.Now()
	rm.startedAt = r.Now()
	rm.tickRate = 150 * time.Millisecond
	rm.popups = []ScorePopup{}
	rm.lastCollisionAt = time.Time{}

	rm.powerUpActive = false
	rm.p1PowerUpType = ""
	rm.p1PowerUpExpiry = time.Time{}
	rm.p2PowerUpType = ""
	rm.p2PowerUpExpiry = time.Time{}
	rm.tickCount = 0
	rm.lastWake = r.Now()
	rm.portalA = Point{X: 9, Y: 9}
	rm.portalB = Point{X: 29, Y: 9}

	// Initialize hazards
	rm.initHazards()

	// Regenerate food and obstacles
	rm.food = rm.randomFreePoint(r, 0)
	rm.obstacles = []Point{}
	for i := 0; i < 3; i++ {
		rm.obstacles = append(rm.obstacles, rm.randomFreePoint(r, 5))
	}

	rm.newPB1 = false
	rm.newPB2 = false
	rm.loadPersonalBests(r)
}

func (rm *room) randomFreePoint(r kit.Room, avoidHeadRange int) Point {
	// Attempt up to 100 times to find a free spot
	for attempt := 0; attempt < 100; attempt++ {
		x := r.Rand().Intn(39)
		y := r.Rand().Intn(18)
		p := Point{X: x, Y: y}

		// Check maze wall collision
		if rm.isMazeWall(p) {
			continue
		}

		// Check snake 1 collision
		inSnake1 := false
		for _, sp := range rm.snake1 {
			if sp == p {
				inSnake1 = true
				break
			}
		}
		if inSnake1 {
			continue
		}

		// Check snake 2 collision
		inSnake2 := false
		for _, sp := range rm.snake2 {
			if sp == p {
				inSnake2 = true
				break
			}
		}
		if inSnake2 {
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

		// Check power-up collision
		if rm.powerUpActive && p == rm.powerUpPos {
			continue
		}

		// Check patrolling hazards collision
		inHazard := false
		for _, hp := range rm.hazards {
			if hp.Pos == p {
				inHazard = true
				break
			}
		}
		if inHazard {
			continue
		}

		// Check portals collision
		if rm.gameMode == ModePortal && (p == rm.portalA || p == rm.portalB) {
			continue
		}

		// Avoid snake head range if requested
		if avoidHeadRange > 0 {
			headAvoid := false
			if len(rm.snake1) > 0 {
				head1 := rm.snake1[0]
				distX := head1.X - p.X
				distY := head1.Y - p.Y
				if distX < 0 {
					distX = -distX
				}
				if distY < 0 {
					distY = -distY
				}
				if distX+distY <= avoidHeadRange {
					headAvoid = true
				}
			}
			if len(rm.snake2) > 0 {
				head2 := rm.snake2[0]
				distX := head2.X - p.X
				distY := head2.Y - p.Y
				if distX < 0 {
					distX = -distX
				}
				if distY < 0 {
					distY = -distY
				}
				if distX+distY <= avoidHeadRange {
					headAvoid = true
				}
			}
			if headAvoid {
				continue
			}
		}

		return p
	}
	// Fallback
	return Point{X: 10, Y: 10}
}

func (rm *room) tick(r kit.Room) {
	if len(rm.snake1) == 0 || len(rm.snake2) == 0 {
		return
	}
	rm.tickCount++

	now := r.Now()
	p1ShieldActive := !rm.p1PowerUpExpiry.IsZero() && now.Before(rm.p1PowerUpExpiry) && rm.p1PowerUpType == "SHIELD"
	p1FreezeActive := !rm.p1PowerUpExpiry.IsZero() && now.Before(rm.p1PowerUpExpiry) && rm.p1PowerUpType == "FREEZE"

	p2ShieldActive := !rm.p2PowerUpExpiry.IsZero() && now.Before(rm.p2PowerUpExpiry) && rm.p2PowerUpType == "SHIELD"
	p2FreezeActive := !rm.p2PowerUpExpiry.IsZero() && now.Before(rm.p2PowerUpExpiry) && rm.p2PowerUpType == "FREEZE"

	// Save tail segment positions
	tail1 := rm.snake1[len(rm.snake1)-1]
	tail2 := rm.snake2[len(rm.snake2)-1]
	oldHead1 := rm.snake1[0]
	oldHead2 := rm.snake2[0]

	move1 := true
	if p2FreezeActive && rm.tickCount%2 == 0 {
		move1 = false
	}

	move2 := true
	if p1FreezeActive && rm.tickCount%2 == 0 {
		move2 = false
	}

	if move1 {
		// Move Snake 1 body
		for i := len(rm.snake1) - 1; i > 0; i-- {
			rm.snake1[i] = rm.snake1[i-1]
		}
		// Update Snake 1 head position
		head1 := rm.snake1[0]
		if rm.gameMode == ModePortal && head1 == rm.portalA {
			rm.snake1[0].X = (rm.portalB.X + rm.entityDir1.X + 39) % 39
			rm.snake1[0].Y = (rm.portalB.Y + rm.entityDir1.Y + 18) % 18
		} else if rm.gameMode == ModePortal && head1 == rm.portalB {
			rm.snake1[0].X = (rm.portalA.X + rm.entityDir1.X + 39) % 39
			rm.snake1[0].Y = (rm.portalA.Y + rm.entityDir1.Y + 18) % 18
		} else {
			rm.snake1[0].X += rm.entityDir1.X
			rm.snake1[0].Y += rm.entityDir1.Y

			// Wrap boundaries for Snake 1
			if rm.snake1[0].X < 0 {
				rm.snake1[0].X = 38
			} else if rm.snake1[0].X > 38 {
				rm.snake1[0].X = 0
			}
			if rm.snake1[0].Y < 0 {
				rm.snake1[0].Y = 17
			} else if rm.snake1[0].Y > 17 {
				rm.snake1[0].Y = 0
			}
		}
		rm.lastMovedDir1 = rm.entityDir1
	}

	if move2 {
		if rm.p2IsBot {
			oppositeDir := Point{X: -rm.lastMovedDir2.X, Y: -rm.lastMovedDir2.Y}
			target := rm.food
			if rm.powerUpActive {
				target = rm.powerUpPos
			}
			dir, ok := rm.findShortestDir(rm.snake2[0], target, oppositeDir, p1FreezeActive, p2FreezeActive)
			if ok {
				rm.entityDir2 = dir
			}
		}
		// Move Snake 2 body
		for i := len(rm.snake2) - 1; i > 0; i-- {
			rm.snake2[i] = rm.snake2[i-1]
		}
		// Update Snake 2 head position
		head2 := rm.snake2[0]
		if rm.gameMode == ModePortal && head2 == rm.portalA {
			rm.snake2[0].X = (rm.portalB.X + rm.entityDir2.X + 39) % 39
			rm.snake2[0].Y = (rm.portalB.Y + rm.entityDir2.Y + 18) % 18
		} else if rm.gameMode == ModePortal && head2 == rm.portalB {
			rm.snake2[0].X = (rm.portalA.X + rm.entityDir2.X + 39) % 39
			rm.snake2[0].Y = (rm.portalA.Y + rm.entityDir2.Y + 18) % 18
		} else {
			rm.snake2[0].X += rm.entityDir2.X
			rm.snake2[0].Y += rm.entityDir2.Y

			// Wrap boundaries for Snake 2
			if rm.snake2[0].X < 0 {
				rm.snake2[0].X = 38
			} else if rm.snake2[0].X > 38 {
				rm.snake2[0].X = 0
			}
			if rm.snake2[0].Y < 0 {
				rm.snake2[0].Y = 17
			} else if rm.snake2[0].Y > 17 {
				rm.snake2[0].Y = 0
			}
		}
		rm.lastMovedDir2 = rm.entityDir2
	}

	// Store old hazard positions for swap collision check
	oldHazards := make([]Point, len(rm.hazards))
	for i, h := range rm.hazards {
		oldHazards[i] = h.Pos
	}

	// Update Patrolling Hazards (frozen if freeze is active)
	hazardsMove := !p1FreezeActive && !p2FreezeActive
	if hazardsMove {
		for i := range rm.hazards {
			h := &rm.hazards[i]
			nextX := h.Pos.X + h.Dir.X
			nextY := h.Pos.Y + h.Dir.Y

			// Check if out of bounds
			if nextX < h.MinX || nextX > h.MaxX || nextY < h.MinY || nextY > h.MaxY {
				h.Dir.X = -h.Dir.X
				h.Dir.Y = -h.Dir.Y
				nextX = h.Pos.X + h.Dir.X
				nextY = h.Pos.Y + h.Dir.Y
			}
			h.Pos.X = nextX
			h.Pos.Y = nextY
		}
	}

	// Check collisions for Snake 1
	c1Self := false
	for _, sp := range rm.snake1[1:] {
		if rm.snake1[0] == sp {
			c1Self = true
		}
	}
	c1Obstacle := false
	for _, op := range rm.obstacles {
		if rm.snake1[0] == op {
			c1Obstacle = true
		}
	}
	c1Maze := rm.isMazeWall(rm.snake1[0])
	c1Hazard := false
	for i, hp := range rm.hazards {
		if rm.snake1[0] == hp.Pos || (rm.snake1[0] == oldHazards[i] && oldHead1 == hp.Pos) {
			c1Hazard = true
		}
	}
	c1Snake2 := false
	for _, sp := range rm.snake2 {
		if rm.snake1[0] == sp {
			c1Snake2 = true
		}
	}

	if (c1Self || c1Obstacle || c1Maze || c1Hazard || c1Snake2) && !p1ShieldActive {
		rm.crashed1 = true
	}

	// Check collisions for Snake 2
	c2Self := false
	for _, sp := range rm.snake2[1:] {
		if rm.snake2[0] == sp {
			c2Self = true
		}
	}
	c2Obstacle := false
	for _, op := range rm.obstacles {
		if rm.snake2[0] == op {
			c2Obstacle = true
		}
	}
	c2Maze := rm.isMazeWall(rm.snake2[0])
	c2Hazard := false
	for i, hp := range rm.hazards {
		if rm.snake2[0] == hp.Pos || (rm.snake2[0] == oldHazards[i] && oldHead2 == hp.Pos) {
			c2Hazard = true
		}
	}
	c2Snake1 := false
	for _, sp := range rm.snake1 {
		if rm.snake2[0] == sp {
			c2Snake1 = true
		}
	}

	if (c2Self || c2Obstacle || c2Maze || c2Hazard || c2Snake1) && !p2ShieldActive {
		rm.crashed2 = true
	}

	// If either crashed, it's Game Over!
	if rm.crashed1 || rm.crashed2 {
		rm.gameOver = true
		rm.lastCollisionAt = r.Now()
		rm.savePersonalBests(r)

		// Post results to the leaderboard
		members := r.Members()
		if len(members) >= 2 {
			p1Result := kit.PlayerResult{
				Player: members[0],
				Metric: rm.score1,
				Status: kit.StatusFinished,
			}
			p2Result := kit.PlayerResult{
				Player: members[1],
				Metric: rm.score2,
				Status: kit.StatusFinished,
			}
			r.Post(kit.Result{
				Rankings: []kit.PlayerResult{p1Result, p2Result},
			})
		} else {
			// Single player
			r.Post(kit.Result{
				Rankings: []kit.PlayerResult{
					{
						Player: rm.getActivePlayer(r),
						Metric: rm.score1,
						Status: kit.StatusFinished,
					},
				},
			})
		}
		return
	}

	// Check power-up collision
	if rm.powerUpActive {
		// Check if 8 seconds expired
		if now.Sub(rm.powerUpSpawnedAt) >= 8*time.Second {
			rm.powerUpActive = false
		} else {
			// Check Snake 1 head
			if rm.snake1[0] == rm.powerUpPos {
				rm.p1PowerUpType = rm.powerUpType
				rm.p1PowerUpExpiry = now.Add(6 * time.Second)
				rm.powerUpActive = false

				// Trigger score popup / text popup
				palettes := getPalettes()
				theme := palettes[rm.themeIndex]
				popupText := "+" + rm.powerUpType
				rm.popups = append(rm.popups, ScorePopup{
					X:         rm.powerUpPos.X,
					Y:         rm.powerUpPos.Y,
					Text:      popupText,
					Color:     theme.SnakeHead,
					CreatedAt: now,
				})
			} else if rm.snake2[0] == rm.powerUpPos {
				// Check Snake 2 head
				rm.p2PowerUpType = rm.powerUpType
				rm.p2PowerUpExpiry = now.Add(6 * time.Second)
				rm.powerUpActive = false

				// Trigger score popup / text popup
				palettes := getPalettes()
				theme := palettes[rm.themeIndex]
				popupText := "+" + rm.powerUpType
				rm.popups = append(rm.popups, ScorePopup{
					X:         rm.powerUpPos.X,
					Y:         rm.powerUpPos.Y,
					Text:      popupText,
					Color:     theme.Snake2Head,
					CreatedAt: now,
				})
			}
		}
	}

	// Check food-collision for Snake 1
	if rm.snake1[0] == rm.food {
		rm.snake1 = append(rm.snake1, tail1)
		rm.score1 += 10
		if rm.score1 > rm.highScore {
			rm.highScore = rm.score1
		}
		rm.onFoodEaten(r, rm.food, 1)
	} else if rm.snake2[0] == rm.food {
		// Check food-collision for Snake 2
		rm.snake2 = append(rm.snake2, tail2)
		rm.score2 += 10
		if rm.score2 > rm.highScore {
			rm.highScore = rm.score2
		}
		rm.onFoodEaten(r, rm.food, 2)
	}
}

func (rm *room) onFoodEaten(r kit.Room, foodPos Point, snakeNum int) {
	totalScore := rm.score1 + rm.score2
	speedMs := 150 - (totalScore/10)*5
	if speedMs < 60 {
		speedMs = 60
	}
	rm.tickRate = time.Duration(speedMs) * time.Millisecond

	palettes := getPalettes()
	theme := palettes[rm.themeIndex]
	var popupColor kit.Color
	if snakeNum == 1 {
		popupColor = theme.SnakeHead
	} else {
		popupColor = theme.Snake2Head
	}
	rm.popups = append(rm.popups, ScorePopup{
		X:         foodPos.X,
		Y:         foodPos.Y,
		Text:      "+10",
		Color:     popupColor,
		CreatedAt: r.Now(),
	})

	rm.food = rm.randomFreePoint(r, 0)
	rm.obstacles = append(rm.obstacles, rm.randomFreePoint(r, 4))

	// Spawn a power-up if not already active
	if !rm.powerUpActive {
		roll := r.Rand().Intn(100)
		if roll < 100 {
			pType := "SHIELD"
			if r.Rand().Intn(2) == 0 {
				pType = "FREEZE"
			}
			rm.powerUpPos = rm.randomFreePoint(r, 2)
			rm.powerUpType = pType
			rm.powerUpActive = true
			rm.powerUpSpawnedAt = r.Now()
		}
	}
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}
	pad := (width - len(text)) / 2
	return fmt.Sprintf("%*s%s", pad, "", text)
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

	members := r.Members()

	// 2. Draw Header Content
	f.Text(1, 2, "▲▼ NEON DUEL ▲▼", headerStyle)

	if len(members) >= 2 {
		f.Text(1, 19, "P1:", headerStyle)
		f.Text(1, 22, fmt.Sprintf("%04d", rm.score1), valueStyle)
		f.Text(1, 27, "P2:", headerStyle)
		f.Text(1, 30, fmt.Sprintf("%04d", rm.score2), valueStyle)
	} else {
		f.Text(1, 19, "S1:", headerStyle)
		f.Text(1, 22, fmt.Sprintf("%04d", rm.score1), valueStyle)
		f.Text(1, 27, "S2:", headerStyle)
		f.Text(1, 30, fmt.Sprintf("%04d", rm.score2), valueStyle)
	}

	f.Text(1, 36, "HI:", headerStyle)
	f.Text(1, 39, fmt.Sprintf("%04d", rm.highScore), valueStyle)

	f.Text(1, 54, "THEME:", headerStyle)
	f.Text(1, 60, theme.Name, valueStyle)

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
	for _, sp := range rm.snake1 {
		occupied[sp] = true
	}
	for _, sp := range rm.snake2 {
		occupied[sp] = true
	}
	occupied[rm.food] = true
	for _, op := range rm.obstacles {
		occupied[op] = true
	}
	if rm.powerUpActive {
		occupied[rm.powerUpPos] = true
	}
	if rm.gameMode == ModePortal {
		occupied[rm.portalA] = true
		occupied[rm.portalB] = true
	}

	p1ShieldActive := !rm.p1PowerUpExpiry.IsZero() && now.Before(rm.p1PowerUpExpiry) && rm.p1PowerUpType == "SHIELD"
	p1FreezeActive := !rm.p1PowerUpExpiry.IsZero() && now.Before(rm.p1PowerUpExpiry) && rm.p1PowerUpType == "FREEZE"

	p2ShieldActive := !rm.p2PowerUpExpiry.IsZero() && now.Before(rm.p2PowerUpExpiry) && rm.p2PowerUpType == "SHIELD"
	p2FreezeActive := !rm.p2PowerUpExpiry.IsZero() && now.Before(rm.p2PowerUpExpiry) && rm.p2PowerUpType == "FREEZE"

	// 3. Draw Grid Dots
	for y := 0; y < 18; y++ {
		for x := 0; x < 39; x++ {
			p := Point{X: x, Y: y}
			if occupied[p] || rm.isMazeWall(p) {
				continue
			}
			f.SetRune(3+y, 1+x*2, '·', dotStyle)
		}
	}

	// 3.5 Draw Maze Walls if in ModeMaze
	if rm.gameMode == ModeMaze {
		wallStyle := kit.Style{FG: theme.Border, Attr: kit.AttrBold}
		for y := 0; y < 18; y++ {
			for x := 0; x < 39; x++ {
				p := Point{X: x, Y: y}
				if rm.isMazeWall(p) {
					f.SetWide(3+y, 1+x*2, '▒', wallStyle)
				}
			}
		}
	}

	// 3.6 Draw Portals if in ModePortal
	if rm.gameMode == ModePortal {
		elapsed := now.Sub(rm.startedAt)
		pulseA := 0.75 + 0.25*math.Sin(float64(elapsed.Milliseconds())*0.01)
		pulseB := 0.75 + 0.25*math.Sin(float64(elapsed.Milliseconds())*0.01 + math.Pi)

		// Portal A: neon blue/cyan
		colorA := brightenColor(kit.RGB(0x00, 0xd2, 0xff), pulseA)
		styleA := kit.Style{FG: colorA, Attr: kit.AttrBold}

		// Portal B: neon orange/gold
		colorB := brightenColor(kit.RGB(0xff, 0x9d, 0x00), pulseB)
		styleB := kit.Style{FG: colorB, Attr: kit.AttrBold}

		f.SetWide(3+rm.portalA.Y, 1+rm.portalA.X*2, '◎', styleA)
		f.SetWide(3+rm.portalB.Y, 1+rm.portalB.X*2, '◎', styleB)
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

	// 4.5 Draw Power-Up on field if active
	if rm.powerUpActive {
		var powerUpStyle kit.Style
		var powerUpGlyph rune
		if rm.powerUpType == "SHIELD" {
			shieldPulse := 0.75 + 0.25*math.Sin(float64(elapsed.Milliseconds())*0.01)
			powerUpStyle = kit.Style{FG: brightenColor(kit.RGB(0xff, 0xd7, 0x00), shieldPulse), Attr: kit.AttrBold}
			powerUpGlyph = '🛡'
		} else {
			freezePulse := 0.75 + 0.25*math.Sin(float64(elapsed.Milliseconds())*0.01)
			powerUpStyle = kit.Style{FG: brightenColor(kit.RGB(0x00, 0xff, 0xff), freezePulse), Attr: kit.AttrBold}
			powerUpGlyph = '❄'
		}
		f.SetWide(3+rm.powerUpPos.Y, 1+rm.powerUpPos.X*2, powerUpGlyph, powerUpStyle)
	}

	// 5. Draw Obstacles (Neon triangles with warning flash when head is close)
	for _, op := range rm.obstacles {
		var dist1, dist2 int
		if len(rm.snake1) > 0 {
			head1 := rm.snake1[0]
			dx := head1.X - op.X
			dy := head1.Y - op.Y
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			dist1 = dx + dy
		} else {
			dist1 = 999
		}
		if len(rm.snake2) > 0 {
			head2 := rm.snake2[0]
			dx := head2.X - op.X
			dy := head2.Y - op.Y
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			dist2 = dx + dy
		} else {
			dist2 = 999
		}

		minDist := dist1
		if dist2 < minDist {
			minDist = dist2
		}

		var obstacleStyle kit.Style
		if minDist <= 2 && rm.gameStarted && !rm.gameOver {
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

	// 5.5 Draw Patrolling Hazards (Neon diamond with pulse and head warning flash)
	for _, hp := range rm.hazards {
		var dist1, dist2 int
		if len(rm.snake1) > 0 {
			head1 := rm.snake1[0]
			dx := head1.X - hp.Pos.X
			dy := head1.Y - hp.Pos.Y
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			dist1 = dx + dy
		} else {
			dist1 = 999
		}
		if len(rm.snake2) > 0 {
			head2 := rm.snake2[0]
			dx := head2.X - hp.Pos.X
			dy := head2.Y - hp.Pos.Y
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			dist2 = dx + dy
		} else {
			dist2 = 999
		}

		minDist := dist1
		if dist2 < minDist {
			minDist = dist2
		}

		var hazardStyle kit.Style
		if minDist <= 2 && rm.gameStarted && !rm.gameOver {
			flashCycle := (elapsed.Milliseconds() / 100) % 2
			if flashCycle == 0 {
				hazardStyle = kit.Style{FG: kit.RGB(0xff, 0xff, 0xff), Attr: kit.AttrBold}
			} else {
				hazardStyle = kit.Style{FG: kit.RGB(0xff, 0x00, 0x55), Attr: kit.AttrBold}
			}
		} else {
			hazardPulse := 0.75 + 0.25*math.Sin(float64(elapsed.Milliseconds())*0.008)
			hazardStyle = kit.Style{FG: brightenColor(theme.Hazard, hazardPulse), Attr: kit.AttrBold}
		}

		f.SetWide(3+hp.Pos.Y, 1+hp.Pos.X*2, '❖', hazardStyle)
	}

	// 6. Draw Snake 1 (gradient from SnakeHead to SnakeTail, flowing dynamically)
	n1 := len(rm.snake1)
	timeShift := float64(now.Sub(rm.startedAt).Milliseconds()%2000) / 2000.0
	for i := n1 - 1; i >= 0; i-- {
		p := rm.snake1[i]
		var segmentStyle kit.Style
		if p1ShieldActive {
			shieldPulse := 0.75 + 0.25*math.Sin(float64(now.Sub(rm.startedAt).Milliseconds()-int64(i*50))*0.01)
			segmentStyle = kit.Style{FG: brightenColor(kit.RGB(0xff, 0xd7, 0x00), shieldPulse), Attr: kit.AttrBold}
		} else if p2FreezeActive {
			segmentStyle = kit.Style{FG: interpolateColor(theme.SnakeHead, kit.RGB(0x00, 0xbf, 0xff), 0.5), Attr: kit.AttrDim}
		} else {
			if i == 0 {
				headPulse := 0.85 + 0.15*math.Sin(float64(now.Sub(rm.startedAt).Milliseconds())*0.006)
				segmentStyle = kit.Style{FG: brightenColor(theme.SnakeHead, headPulse)}
			} else {
				ratio := float64(i) / float64(n1-1)
				shiftedRatio := ratio + timeShift
				if shiftedRatio > 1.0 {
					shiftedRatio -= 1.0
				}
				segmentStyle = kit.Style{FG: interpolateColor(theme.SnakeHead, theme.SnakeTail, shiftedRatio)}
			}
		}
		f.SetWide(3+p.Y, 1+p.X*2, '█', segmentStyle)
	}

	// Draw Snake 2 (gradient from Snake2Head to Snake2Tail, flowing dynamically)
	n2 := len(rm.snake2)
	for i := n2 - 1; i >= 0; i-- {
		p := rm.snake2[i]
		var segmentStyle kit.Style
		if p2ShieldActive {
			shieldPulse := 0.75 + 0.25*math.Sin(float64(now.Sub(rm.startedAt).Milliseconds()-int64(i*50))*0.01)
			segmentStyle = kit.Style{FG: brightenColor(kit.RGB(0xff, 0xd7, 0x00), shieldPulse), Attr: kit.AttrBold}
		} else if p1FreezeActive {
			segmentStyle = kit.Style{FG: interpolateColor(theme.Snake2Head, kit.RGB(0x00, 0xbf, 0xff), 0.5), Attr: kit.AttrDim}
		} else {
			if i == 0 {
				headPulse := 0.85 + 0.15*math.Sin(float64(now.Sub(rm.startedAt).Milliseconds())*0.006)
				segmentStyle = kit.Style{FG: brightenColor(theme.Snake2Head, headPulse)}
			} else {
				ratio := float64(i) / float64(n2-1)
				shiftedRatio := ratio + timeShift
				if shiftedRatio > 1.0 {
					shiftedRatio -= 1.0
				}
				segmentStyle = kit.Style{FG: interpolateColor(theme.Snake2Head, theme.Snake2Tail, shiftedRatio)}
			}
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

	// 7.5 Draw Divider on row 21 with Mode and player text
	var dividerText string
	p1Status := ""
	if p1ShieldActive {
		remSec := int(math.Ceil(rm.p1PowerUpExpiry.Sub(now).Seconds()))
		p1Status = fmt.Sprintf(" [🛡 %ds]", remSec)
	} else if p1FreezeActive {
		remSec := int(math.Ceil(rm.p1PowerUpExpiry.Sub(now).Seconds()))
		p1Status = fmt.Sprintf(" [❄ %ds]", remSec)
	}

	p2Status := ""
	if p2ShieldActive {
		remSec := int(math.Ceil(rm.p2PowerUpExpiry.Sub(now).Seconds()))
		p2Status = fmt.Sprintf(" [🛡 %ds]", remSec)
	} else if p2FreezeActive {
		remSec := int(math.Ceil(rm.p2PowerUpExpiry.Sub(now).Seconds()))
		p2Status = fmt.Sprintf(" [❄ %ds]", remSec)
	}

	if len(members) >= 2 {
		dividerText = fmt.Sprintf(" MODE: %s │ P1: %s%s VS P2: %s%s ", rm.getModeName(), members[0].Handle, p1Status, members[1].Handle, p2Status)
	} else {
		botText := "CO-OP"
		if rm.p2IsBot {
			botText = "AI-BOT"
		}
		dividerText = fmt.Sprintf(" MODE: %s │ PLAYSTYLE: %s %s%s ", rm.getModeName(), botText, p1Status, p2Status)
	}
	rm.drawDividerWithText(f, 21, dividerText, now)

	// 8. Draw Footer Content
	f.Text(22, 2, "CONTROLS:", footerStyle)
	col := 12
	col = f.Text(22, col, " [", footerStyle)
	if len(members) >= 2 {
		col = f.Text(22, col, "P1:WASD P2:Arrows", keyStyle)
	} else {
		if rm.p2IsBot {
			col = f.Text(22, col, "WASD/Bot", keyStyle)
		} else {
			col = f.Text(22, col, "WASD/Arrows", keyStyle)
		}
	}
	col = f.Text(22, col, "] Move", footerStyle)

	col = f.Text(22, col+1, " [", footerStyle)
	col = f.Text(22, col, "T", keyStyle)
	col = f.Text(22, col, "]Theme", footerStyle)

	col = f.Text(22, col+1, " [", footerStyle)
	col = f.Text(22, col, "M", keyStyle)
	col = f.Text(22, col, "]Mode", footerStyle)

	if len(members) < 2 {
		col = f.Text(22, col+1, " [", footerStyle)
		col = f.Text(22, col, "B", keyStyle)
		col = f.Text(22, col, "]Bot", footerStyle)
	}

	col = f.Text(22, col+1, " [", footerStyle)
	col = f.Text(22, col, "Spc", keyStyle)
	col = f.Text(22, col, "]Pause", footerStyle)

	col = f.Text(22, col+1, " [", footerStyle)
	col = f.Text(22, col, "Esc", keyStyle)
	col = f.Text(22, col, "]Quit", footerStyle)

	// 9. Draw Game Over Overlay
	if rm.gameOver {
		modalStyle := kit.Style{FG: theme.ModalBorder, Attr: kit.AttrBold}
		textStyle := kit.Style{FG: kit.RGB(0xff, 0xff, 0xff), Attr: kit.AttrBold}
		subTextStyle := kit.Style{FG: theme.Border}

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

		var winnerMsg string
		if rm.crashed1 && rm.crashed2 {
			winnerMsg = "DRAW / MUTUAL CRASH"
		} else if rm.crashed1 {
			if len(members) >= 2 {
				winnerMsg = fmt.Sprintf("%s WINS!", members[1].Handle)
			} else {
				winnerMsg = "SNAKE 2 WINS!"
			}
		} else {
			if len(members) >= 2 {
				winnerMsg = fmt.Sprintf("%s WINS!", members[0].Handle)
			} else {
				winnerMsg = "SNAKE 1 WINS!"
			}
		}

		scoreText := fmt.Sprintf("S1: %04d  VS  S2: %04d", rm.score1, rm.score2)

		f.Text(11, 21, centerText(winnerMsg, 38), textStyle)
		f.Text(12, 21, centerText(scoreText, 38), textStyle)

		f.Text(13, 25, "Press [SPACE] to Restart", subTextStyle)
	}

	// Send viewport to each member
	for _, p := range r.Members() {
		pb := rm.pb1
		isNewPB := rm.newPB1
		if len(members) >= 2 && p.AccountID == members[1].AccountID {
			pb = rm.pb2
			isNewPB = rm.newPB2
		}

		currentPB := pb
		isCurrentlyNewPB := isNewPB

		if len(members) >= 2 {
			if p.AccountID == members[0].AccountID {
				if rm.score1 > pb {
					currentPB = rm.score1
					isCurrentlyNewPB = true
				}
			} else {
				if rm.score2 > pb {
					currentPB = rm.score2
					isCurrentlyNewPB = true
				}
			}
		} else {
			maxScore := rm.score1
			if rm.score2 > maxScore {
				maxScore = rm.score2
			}
			if maxScore > pb {
				currentPB = maxScore
				isCurrentlyNewPB = true
			}
		}

		var pbText string
		var pbValStyle kit.Style
		if isCurrentlyNewPB {
			elapsed := now.Sub(rm.startedAt)
			pulse := (elapsed.Milliseconds() / 200) % 2
			if pulse == 0 {
				pbText = "*NEW*"
				pbValStyle = kit.Style{FG: kit.RGB(0xff, 0xff, 0x00), Attr: kit.AttrBold}
			} else {
				pbText = fmt.Sprintf("%04d", currentPB)
				pbValStyle = kit.Style{FG: kit.RGB(0xff, 0x00, 0x55), Attr: kit.AttrBold}
			}
		} else {
			pbText = fmt.Sprintf("%04d", currentPB)
			pbValStyle = valueStyle
		}

		f.Text(1, 45, "PB:", headerStyle)
		f.Text(1, 48, pbText, pbValStyle)

		if rm.gameOver {
			var pbMsg string
			var pbMsgStyle kit.Style
			if isCurrentlyNewPB {
				pbMsg = "✨ NEW PERSONAL BEST! ✨"
				pulse := (now.Sub(rm.startedAt).Milliseconds() / 150) % 3
				switch pulse {
				case 0:
					pbMsgStyle = kit.Style{FG: kit.RGB(0xff, 0xd7, 0x00), Attr: kit.AttrBold}
				case 1:
					pbMsgStyle = kit.Style{FG: kit.RGB(0x00, 0xff, 0xff), Attr: kit.AttrBold}
				default:
					pbMsgStyle = kit.Style{FG: kit.RGB(0xff, 0x00, 0x7f), Attr: kit.AttrBold}
				}
			} else {
				pbMsg = fmt.Sprintf("Personal Best: %04d", pb)
				pbMsgStyle = kit.Style{FG: theme.Footer}
			}
			f.Text(9, 21, centerText(pbMsg, 38), pbMsgStyle)
		}

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

func (rm *room) findShortestDir(start Point, target Point, oppositeDir Point, p1FreezeActive, p2FreezeActive bool) (Point, bool) {
	// BFS pathfinding queue
	type QueueNode struct {
		pos  Point
		path []Point
	}

	queue := []QueueNode{
		{pos: start, path: []Point{}},
	}
	visited := make(map[Point]bool)
	visited[start] = true

	isSafeBFS := func(p Point, isStart bool, prevDir Point) bool {
		if rm.isMazeWall(p) {
			return false
		}
		for _, op := range rm.obstacles {
			if op == p {
				return false
			}
		}
		for _, sp := range rm.snake1 {
			if sp == p {
				return false
			}
		}
		for _, sp := range rm.snake2 {
			if sp == p {
				return false
			}
		}
		for _, hp := range rm.hazards {
			if hp.Pos == p {
				return false
			}
			if !p1FreezeActive && !p2FreezeActive {
				nextHX := hp.Pos.X + hp.Dir.X
				nextHY := hp.Pos.Y + hp.Dir.Y
				if nextHX < hp.MinX || nextHX > hp.MaxX || nextHY < hp.MinY || nextHY > hp.MaxY {
					nextHX = hp.Pos.X - hp.Dir.X
					nextHY = hp.Pos.Y - hp.Dir.Y
				}
				if p.X == nextHX && p.Y == nextHY {
					return false
				}
			}
		}
		return true
	}

	dirs := []Point{
		{X: 0, Y: -1}, // Up
		{X: 0, Y: 1},  // Down
		{X: -1, Y: 0}, // Left
		{X: 1, Y: 0},  // Right
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.pos == target {
			if len(curr.path) > 0 {
				return curr.path[0], true
			}
			return Point{X: 0, Y: 0}, false
		}

		for _, d := range dirs {
			if curr.pos == start && d == oppositeDir {
				continue
			}

			var nextP Point
			if rm.gameMode == ModePortal && curr.pos == rm.portalA {
				nextP = Point{
					X: (rm.portalB.X + d.X + 39) % 39,
					Y: (rm.portalB.Y + d.Y + 18) % 18,
				}
			} else if rm.gameMode == ModePortal && curr.pos == rm.portalB {
				nextP = Point{
					X: (rm.portalA.X + d.X + 39) % 39,
					Y: (rm.portalA.Y + d.Y + 18) % 18,
				}
			} else {
				nextP = Point{
					X: (curr.pos.X + d.X + 39) % 39,
					Y: (curr.pos.Y + d.Y + 18) % 18,
				}
			}

			if !visited[nextP] {
				if nextP == target || isSafeBFS(nextP, curr.pos == start, d) {
					visited[nextP] = true
					newPath := make([]Point, len(curr.path)+1)
					copy(newPath, curr.path)
					newPath[len(curr.path)] = d
					queue = append(queue, QueueNode{pos: nextP, path: newPath})
				}
			}
		}
	}

	// Fallback to survival space maximization if no path to target found
	var bestDir Point
	bestSpace := -1
	hasSafeMove := false

	for _, d := range dirs {
		if d == oppositeDir {
			continue
		}
		var nextP Point
		if rm.gameMode == ModePortal && start == rm.portalA {
			nextP = Point{
				X: (rm.portalB.X + d.X + 39) % 39,
				Y: (rm.portalB.Y + d.Y + 18) % 18,
			}
		} else if rm.gameMode == ModePortal && start == rm.portalB {
			nextP = Point{
				X: (rm.portalA.X + d.X + 39) % 39,
				Y: (rm.portalA.Y + d.Y + 18) % 18,
			}
		} else {
			nextP = Point{
				X: (start.X + d.X + 39) % 39,
				Y: (start.Y + d.Y + 18) % 18,
			}
		}

		if isSafeBFS(nextP, true, d) {
			hasSafeMove = true
			space := rm.countReachableSpace(nextP, p1FreezeActive, p2FreezeActive)
			if space > bestSpace {
				bestSpace = space
				bestDir = d
			}
		}
	}

	if hasSafeMove {
		return bestDir, true
	}

	return Point{X: 0, Y: 0}, false
}

func (rm *room) countReachableSpace(start Point, p1FreezeActive, p2FreezeActive bool) int {
	visited := make(map[Point]bool)
	visited[start] = true
	queue := []Point{start}
	count := 0
	maxDepth := 50

	isSafeBFS := func(p Point) bool {
		if rm.isMazeWall(p) {
			return false
		}
		for _, op := range rm.obstacles {
			if op == p {
				return false
			}
		}
		for _, sp := range rm.snake1 {
			if sp == p {
				return false
			}
		}
		for _, sp := range rm.snake2 {
			if sp == p {
				return false
			}
		}
		for _, hp := range rm.hazards {
			if hp.Pos == p {
				return false
			}
			if !p1FreezeActive && !p2FreezeActive {
				nextHX := hp.Pos.X + hp.Dir.X
				nextHY := hp.Pos.Y + hp.Dir.Y
				if nextHX < hp.MinX || nextHX > hp.MaxX || nextHY < hp.MinY || nextHY > hp.MaxY {
					nextHX = hp.Pos.X - hp.Dir.X
					nextHY = hp.Pos.Y - hp.Dir.Y
				}
				if p.X == nextHX && p.Y == nextHY {
					return false
				}
			}
		}
		return true
	}

	dirs := []Point{
		{X: 0, Y: -1},
		{X: 0, Y: 1},
		{X: -1, Y: 0},
		{X: 1, Y: 0},
	}

	for len(queue) > 0 && count < maxDepth {
		curr := queue[0]
		queue = queue[1:]
		count++

		for _, d := range dirs {
			var nextP Point
			if rm.gameMode == ModePortal && curr == rm.portalA {
				nextP = Point{
					X: (rm.portalB.X + d.X + 39) % 39,
					Y: (rm.portalB.Y + d.Y + 18) % 18,
				}
			} else if rm.gameMode == ModePortal && curr == rm.portalB {
				nextP = Point{
					X: (rm.portalA.X + d.X + 39) % 39,
					Y: (rm.portalA.Y + d.Y + 18) % 18,
				}
			} else {
				nextP = Point{
					X: (curr.X + d.X + 39) % 39,
					Y: (curr.Y + d.Y + 18) % 18,
				}
			}
			if !visited[nextP] && isSafeBFS(nextP) {
				visited[nextP] = true
				queue = append(queue, nextP)
			}
		}
	}
	return count
}
