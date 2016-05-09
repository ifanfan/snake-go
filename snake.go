package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

var (
	// FoodColor food color
	FoodColor = termbox.ColorYellow

	//WallColor wall color
	WallColor = termbox.ColorBlue

	// SnakeColor snake color
	SnakeColor = termbox.ColorWhite
)

// Point is a base data
type Point struct {
	x int
	y int
}

//Transform vector transform
func (p *Point) Transform(v Vector) {
	p.x += v.x
	p.y += v.y
}

// Equal if point equal
func (p *Point) Equal(point Point) bool {
	return p.x == point.x && p.y == point.y
}

// ToDraw for render into term
type ToDraw interface {
	Color() termbox.Attribute
	Points() []Point
}

// ToText for render text
type ToText interface {
	Cursor() Point
	Color() termbox.Attribute
	Runes() []rune
}

// Wall is edge
type Wall struct {
	points []Point
	color  termbox.Attribute
	P1     Point
	P2     Point
}

// NewWall create wall
func NewWall(p1, p2 Point) *Wall {
	w := new(Wall)
	w.color = WallColor
	points := make([]Point, 0, (p2.x+p2.y)*2)
	for i := p1.x; i <= p2.x; i++ {
		for j := p1.y; j <= p2.y; j++ {
			if i == p1.x || i == p2.x || j == p1.y || j == p2.y {
				points = append(points, Point{i, j})
			}
		}
	}
	w.P1 = p1
	w.P2 = p2
	w.points = points
	return w
}

// Color for render
func (p *Wall) Color() termbox.Attribute {
	return p.color
}

// Points return all points
func (p *Wall) Points() []Point {
	return p.points
}

// Food for snake to eat
type Food struct {
	point Point
	color termbox.Attribute
}

// NewFood create new food
func NewFood(x, y int) *Food {
	f := new(Food)
	f.point.x = x
	f.point.y = y
	f.color = FoodColor
	return f
}

// Color for render
func (p *Food) Color() termbox.Attribute {
	return p.color
}

// Points return all points
func (p *Food) Points() []Point {
	return []Point{p.point}
}

// Vector vector
type Vector Point

var (
	// Up direction
	Up = Vector{0, -1}
	// Down direction
	Down = Vector{0, 1}
	// Left direction
	Left = Vector{-1, 0}
	// Right direction
	Right = Vector{1, 0}
)

// VcIsZero if the Vector is zero
func VcIsZero(p Vector) bool {
	return p.x == 0 && p.y == 0
}

//VcAdd Add vector add
func VcAdd(v0, v1 Vector) Vector {
	v0.x += v1.x
	v0.y += v1.y
	return v0
}

// VcEqual is vector equal
func VcEqual(v0, v1 Vector) bool {
	return v0.x == v1.y && v0.y == v1.y
}

// Snake game role
type Snake struct {
	color  termbox.Attribute
	points []Point
	vector Vector
}

// NewSnake create new snake
func NewSnake(x, y int) *Snake {
	s := new(Snake)
	s.color = termbox.ColorWhite
	s.points = make([]Point, 1, 100)
	s.points[0] = Point{x, y}
	s.vector = Right
	return s
}

// Color for render
func (p *Snake) Color() termbox.Attribute {
	return p.color
}

// Points return points
func (p *Snake) Points() []Point {
	return p.points
}

// Head return snake head
func (p *Snake) Head() Point {
	return p.points[0]
}

// SetVector set step vector
func (p *Snake) SetVector(v Vector) {
	if !VcIsZero(v) && !VcIsZero(VcAdd(p.vector, v)) {
		p.vector = v
	}
}

// Step snake  step by direction
func (p *Snake) Step() {
	copy(p.points[1:], p.points[:len(p.points)-1])
	p.points[0].Transform(p.vector)
}

// Eat snake eat
func (p *Snake) Eat() {
	point := p.points[0]
	point.Transform(p.vector)
	p.points = append(p.points, Point{})
	copy(p.points[1:], p.points[:len(p.points)-1])
	p.points[0] = point
}

// Scores snake length
type Scores struct {
	num    int
	cursor Point
	color  termbox.Attribute
}

// NewScores create score
func NewScores(x, y int) *Scores {
	s := new(Scores)
	s.cursor = Point{x, y}
	s.color = termbox.ColorWhite
	s.num = 0
	return s
}

// Cursor return text position
func (p *Scores) Cursor() Point {
	return p.cursor
}

// Color return color
func (p *Scores) Color() termbox.Attribute {
	return p.color
}

// Runes return text []rune
func (p *Scores) Runes() []rune {
	return []rune(fmt.Sprintf("NUM: %d", p.num))
}

// Inc number ++
func (p *Scores) Inc() {
	p.num++
}

// Render ToDraw data into term
func Render(ds ...ToDraw) {
	for _, d := range ds {
		color := d.Color()
		points := d.Points()
		for _, point := range points {
			// point width is 2
			termbox.SetCell(point.x*2, point.y, ' ', termbox.ColorDefault, color)
			termbox.SetCell(point.x*2+1, point.y, ' ', termbox.ColorDefault, color)
		}
	}
	//termbox.Flush()
}

// RenderT render text into term
func RenderT(ts ...ToText) {
	for _, t := range ts {
		rs := t.Runes()
		color := t.Color()
		p := t.Cursor()
		for i, r := range rs {
			termbox.SetCell(p.x+i, p.y, r, color, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

// Clear screen
func Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func intersect(p Point, ps []Point) bool {
	for _, exp := range ps {
		if p.Equal(exp) {
			return true
		}
	}
	return false
}

func randomFood(p1, p2 Point, excludes []Point) *Food {
	rand.Seed(time.Now().UnixNano())
	var x, y int
	for {
		x = rand.Intn(p2.x-p1.x-1) + p1.x + 1
		y = rand.Intn(p2.y-p1.y-1) + p1.y + 1
		if !intersect(Point{x, y}, excludes) {
			break
		}
	}
	return NewFood(x, y)
}

func canEat(s *Snake, f *Food) bool {
	p := s.Head()
	p.Transform(s.vector)
	return f.point.Equal(p)
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eq := make(chan termbox.Event)
	go func() {
		for {
			eq <- termbox.PollEvent()
		}
	}()

	tick := time.NewTicker(200 * time.Millisecond)
	defer tick.Stop()

	vector := Right

	width, height := termbox.Size()
	width = width / 2 / 2
	p1 := Point{2, 2}
	p2 := Point{width - 3, height - 3}

	scores := NewScores(4, 1)
	wall := NewWall(p1, p2)
	snake := NewSnake(wall.P2.x/2, wall.P2.y/2)
	food := randomFood(wall.P1, wall.P2, snake.Points())
	Render(wall, snake, food)
	RenderT(scores)

loop:
	for {
		select {
		case ev := <-eq:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyEsc:
					break loop
				case termbox.KeyArrowUp:
					vector = Up
				case termbox.KeyArrowDown:
					vector = Down
				case termbox.KeyArrowLeft:
					vector = Left
				case termbox.KeyArrowRight:
					vector = Right
				}
				snake.SetVector(vector)
			}
		case <-tick.C:
			if canEat(snake, food) {
				snake.Eat()
				scores.Inc()
				food = randomFood(wall.P1, wall.P2, snake.Points())
			} else {
				snake.Step()
				if intersect(snake.Head(), snake.Points()[1:]) || intersect(snake.Head(), wall.Points()) {
					snake.color = termbox.ColorRed
					tick.Stop()
				}
			}
			Clear()
			Render(wall, snake, food)
			RenderT(scores)
		}
	}
}
