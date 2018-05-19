package main

import (
	"fmt"
	"math"
	"math/rand"

	"time"

	"github.com/oskca/gopherjs-canvas"
	"github.com/oskca/gopherjs-dom"
)

// A Ship is the player.
type Ship struct {
	x, y, dx, dy, dir float64
}

// An Asteroid floats through space and can kill the player.
type Asteroid struct {
	x, y, dx, dy, dir, spin, r float64
	path                       []float64
}

// A Missile is fired from the Ship to destroy Asteroids.
type Missile struct {
	x, y, dx, dy, fuse float64
}

// KeysPressed contains the current keys being pressed.
type KeysPressed struct {
	left, right, thrust, missile bool
}

const (
	fps           = 30
	maxSpeed      = 6.0
	twoPi         = math.Pi * 2
	keyA          = 65
	keyS          = 83
	keyK          = 75
	keyL          = 76
	keyLeftArrow  = 37
	keyRightArrow = 39
	keyUpArrow    = 38
	keySpace      = 32
)

var (
	ship         *Ship
	cnvs         *canvas.Canvas
	ctx          *canvas.Context2D
	keysPressed  *KeysPressed
	canvasWidth  float64
	canvasHeight float64
	rng          *rand.Rand
	asteroids    []Asteroid
	missiles     []Missile
	missileTimer int
	score        int
)

func randomRange(min, max float64) float64 {
	return min + rng.Float64()*(max-min)
}

func randomSign() float64 {
	if rng.Float32() < 0.5 {
		return -1
	}

	return 1
}

func clamp(lo, hi, x float64) float64 {
	return math.Max(lo, math.Min(hi, x))
}

func (ship Ship) draw(ctx *canvas.Context2D) {
	ctx.Save()
	ctx.Translate(ship.x, ship.y)
	ctx.Rotate(ship.dir)
	ctx.StrokeStyle = "#00ADD8"
	ctx.LineWidth = 1
	ctx.BeginPath()
	ctx.MoveTo(0, -1)
	ctx.MoveTo(0, -10)
	ctx.LineTo(-7, 10)
	ctx.LineTo(0, 8)
	ctx.LineTo(7, 10)
	ctx.LineTo(0, -10)
	ctx.Stroke()
	ctx.Restore()
}

func (asteroid Asteroid) draw(ctx *canvas.Context2D) {
	ctx.Save()
	ctx.Translate(asteroid.x, asteroid.y)
	ctx.Rotate(asteroid.dir)
	ctx.StrokeStyle = "#5DC9E2"
	ctx.LineWidth = 1
	ctx.BeginPath()

	lenFloat := float64(len(asteroid.path))
	theta := lenFloat * (twoPi / lenFloat)
	off := asteroid.path[len(asteroid.path)-1]

	ctx.MoveTo(asteroid.r*off*math.Cos(theta), asteroid.r*off*math.Sin(theta))

	for i, off := range asteroid.path {
		theta := float64(i+1) * (twoPi / lenFloat)

		ctx.LineTo(asteroid.r*off*math.Cos(theta), asteroid.r*off*math.Sin(theta))
	}

	ctx.Stroke()
	ctx.Restore()
}

func (missile Missile) draw(ctx *canvas.Context2D) {
	ctx.Save()
	ctx.StrokeStyle = "#00ADD8"
	ctx.BeginPath()
	ctx.Arc(missile.x, missile.y, 1.0, 0.0, twoPi, false)
	ctx.Stroke()
	ctx.Restore()
}

func drawScore() {
	scoreText := fmt.Sprintf("Score: %d", score)

	ctx.Save()
	ctx.FillStyle = "#00ADD8"
	ctx.Font = "16px silkscreen"
	ctx.FillText(scoreText, 20, 20, 100.0)
	ctx.Restore()
}

func updateShip() {
	leftAdjust := 0.0
	if keysPressed.left {
		leftAdjust = (-math.Pi / 16.0)
	}

	rightAdjust := 0.0
	if keysPressed.right {
		rightAdjust = (math.Pi / 16.0)
	}

	ship.dir = ship.dir + leftAdjust + rightAdjust

	thrust := 0.0
	if keysPressed.thrust {
		thrust = 0.5
	}

	dx := clamp(-maxSpeed, maxSpeed, ship.dx+thrust*math.Cos(ship.dir-math.Pi/2.0))
	dy := clamp(-maxSpeed, maxSpeed, ship.dy+thrust*math.Sin(ship.dir-math.Pi/2.0))

	ship.dx = dx
	ship.dy = dy

	ship.x = math.Mod(canvasWidth+ship.x+ship.dx, canvasWidth)
	ship.y = math.Mod(canvasHeight+ship.y+ship.dy, canvasHeight)
}

func updateAsteroids() {
	for i, asteroid := range asteroids {
		asteroids[i].x = math.Mod(canvasWidth+asteroid.x+asteroid.dx, canvasWidth)
		asteroids[i].y = math.Mod(canvasHeight+asteroid.y+asteroid.dy, canvasHeight)

		dir := asteroid.dir + asteroid.spin
		if dir > twoPi {
			dir = dir - twoPi
		}

		asteroids[i].dir = dir
	}
}

func updateMissiles() {
	if missileTimer != 0 {
		missileTimer++

		if missileTimer >= fps/3 {
			missileTimer = 0
		}
	}

	if keysPressed.missile && missileTimer == 0 {
		dx := ship.dx + 8.0*math.Cos(ship.dir-math.Pi/2.0)
		dy := ship.dy + 8.0*math.Sin(ship.dir-math.Pi/2.0)
		missile := Missile{ship.x, ship.y, dx, dy, 50.0}
		missiles = append(missiles, missile)
		missileTimer = 1
	}

	k := 0
	for _, missile := range missiles {
		missile.fuse -= 1.0

		if missile.fuse != 0 {
			missile.x = math.Mod(canvasWidth+missile.x+missile.dx, canvasWidth)
			missile.y = math.Mod(canvasHeight+missile.y+missile.dy, canvasHeight)
			missiles[k] = missile
			k++
		}
	}

	missiles = missiles[:k]
}

func update() {
	updateShip()
	updateAsteroids()
	updateMissiles()
}

func draw() {
	ctx.FillStyle = "white"
	ctx.FillRect(0.0, 0.0, canvasWidth, canvasHeight)

	ship.draw(ctx)

	for _, asteroid := range asteroids {
		asteroid.draw(ctx)
	}

	for _, missile := range missiles {
		missile.draw(ctx)
	}

	// Score
	drawScore()
}

func gameLoop() {
	update()
	draw()
	dom.Window().Call("setTimeout", gameLoop, 1000/fps)
}

func updateKeysPressed(keyCode int, pressed bool) {
	switch keyCode {
	case keyA, keyLeftArrow: // a
		keysPressed.left = pressed
	case keyS, keyRightArrow: // s
		keysPressed.right = pressed
	case keyK, keyUpArrow: // k
		keysPressed.thrust = pressed
	case keyL, keySpace: // l
		keysPressed.missile = pressed
	}
}

func makeAsteroidPath() []float64 {
	path := make([]float64, 12)

	for i := range path {
		path[i] = randomRange(0.7, 1.1)
	}

	return path
}

func makeAsteroids() {
	asteroids = make([]Asteroid, 10)

	for i := range asteroids {
		path := makeAsteroidPath()
		x := randomRange(0, canvasWidth)
		y := randomRange(0, canvasHeight)
		dx := randomRange(1, 2) * randomSign()
		dy := randomRange(1, 2) * randomSign()
		spin := randomRange(-0.1, 0.1)
		asteroids[i] = Asteroid{x, y, dx, dy, 0, spin, 50, path}
	}
}

func main() {
	window := dom.Window()
	doc := dom.Document()
	cnvs := canvas.New(doc.GetElementById("canvas").Object)

	canvasWidth = float64(cnvs.Width)
	canvasHeight = float64(cnvs.Height)
	ctx = cnvs.GetContext2D()

	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	ship = &Ship{canvasWidth / 2.0, canvasHeight / 2.0, 0.0, 0.0, 0.0}
	missiles = make([]Missile, 0)
	missileTimer = 0
	keysPressed = &KeysPressed{false, false, false, false}
	makeAsteroids()

	window.AddEventListener(dom.EvtKeydown, func(event *dom.Event) {
		updateKeysPressed(event.KeyCode, true)
	})

	window.AddEventListener(dom.EvtKeyup, func(event *dom.Event) {
		updateKeysPressed(event.KeyCode, false)
	})

	gameLoop()
}
