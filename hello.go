package main

import (
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Pair struct {
	X int
	Y int
}

func (self Pair) addPair(other Pair) Pair {
	r := self
	r.addPairMut(other)
	return r
}
func (self *Pair) addPairMut(other Pair) {
	self.X += other.X
	self.Y += other.Y
}
func (self Pair) mulPair(other Pair) Pair {
	r := self
	r.mulPairMut(other)
	return r
}
func (self *Pair) mulPairMut(other Pair) {
	self.X *= other.X
	self.Y *= other.Y
}

func (self *Pair) mulValMut(other int) {
	self.X *= other
	self.Y *= other
}
func (self Pair) mulVal(other int) Pair {
	r := self
	r.mulValMut(other)
	return r
}
func (self *Pair) addValueMut(other int) {
	self.X += other
	self.Y += other
}
func (self Pair) addValue(other int) Pair {
	r := self
	r.addValueMut(other)
	return r
}
func (self *Pair) addValueXMut(other int) {
	self.X += other
}
func (self *Pair) addValueYMut(other int) {
	self.Y += other
}
func (self *Pair) subPairMut(other Pair) {
	self.X += other.X
	self.Y += other.Y
}
func (self Pair) subPair(other Pair) Pair {
	r := self
	r.subPairMut(other)
	return r
}

type Rule struct{}

func (r Rule) modify(w World) World {
	newWorld := w

	for index, cell := range w.Cells {
		aliveDeath := cell.getNeighbours(w)
		// Dies if less than Two living nei
		if aliveDeath.X < 2 {
			newWorld.Cells[index].IsAlive = false

		} else if aliveDeath.X <= 3 {
			newWorld.Cells[index].IsAlive = true

		} else if aliveDeath.X > 3 {
			newWorld.Cells[index].IsAlive = false
		} else if aliveDeath.X == 3 {
			newWorld.Cells[index].IsAlive = true
		}
	}
	return newWorld
}

type Cell struct {
	WorldPosition Pair
	IsAlive       bool
}

func (c Cell) getNeighbours(w World) Pair {

	pos := c.WorldPosition
	var neighbours = [8]Pair{
		pos.addPair(Pair{-1, -1}),
		pos.addPair(Pair{0, -1}),
		pos.addPair(Pair{+1, -1}),
		pos.addPair(Pair{+1, 0}),
		pos.addPair(Pair{+1, +1}),
		pos.addPair(Pair{0, +1}),
		pos.addPair(Pair{-1, +1}),
		pos.addPair(Pair{-1, 0}),
	}

	aliveCount := Pair{}
	for _, e := range neighbours {
		index := w.Convert2d1d(e)
		if index >= 0 && index < len(w.Cells) {
			if w.Cells[index].IsAlive == true {
				aliveCount.addValueXMut(1)
			} else {
				aliveCount.addValueYMut(1)
			}
		}
	}

	return aliveCount
}

func (c Cell) draw(s tcell.Screen) {
	// 3 x 2 is a squere
	corners := worldToScreenPosition(c.WorldPosition)
	upperLeftCorner := corners[0]
	lowerRightCorner := corners[1]

	cellAliveStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	cellDeadStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)
	if c.IsAlive {
		drawBox(s, upperLeftCorner.X, upperLeftCorner.Y, lowerRightCorner.X, lowerRightCorner.Y, cellAliveStyle, "")
	} else {
		drawBox(s, upperLeftCorner.X, upperLeftCorner.Y, lowerRightCorner.X, lowerRightCorner.Y, cellDeadStyle, "")
	}
}

func Convert1d2d(i int, maxX int) Pair {
	return Pair{i % maxX, i / maxX}
}
func Convert2d1d(i Pair, maxX int) int {
	r := i.Y*maxX + i.X
	return r
}

type World struct {
	Cells   []Cell
	MapSize Pair
}

func (w *World) Convert1d2d(i int) Pair {
	return Convert1d2d(i, w.MapSize.X)
}
func (w *World) Convert2d1d(i Pair) int {
	return Convert2d1d(i, w.MapSize.X)
}

func (w *World) init(mapSize Pair) {
	w.MapSize = mapSize
	w.Cells = make([]Cell, w.MapSize.X*w.MapSize.Y)
	for i := 0; i < w.MapSize.X*w.MapSize.Y; i++ {
		w.Cells[i].WorldPosition = w.Convert1d2d(i)
		if i%2 == 0 {
			w.Cells[i].IsAlive = true
		}
	}
}

func worldToScreenPosition(pos Pair) [2]Pair {

	transpose := Pair{3, 2}
	begin := Pair{1, 1}
	diffPairLucRlc := Pair{2, 1}
	rlc := pos
	rlc.addPairMut(begin)
	rlc.addValueXMut(transpose.X * rlc.X)
	rlc.addValueYMut(transpose.Y * rlc.Y)

	luc := rlc
	luc.subPairMut(diffPairLucRlc)
	arr := [2]Pair{luc, rlc}
	return arr
}

func (w World) update(s tcell.Screen) {
	w.draw(s)
}
func (w World) draw(s tcell.Screen) {
	for _, e := range w.Cells {
		e.draw(s)
	}

}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}
	// // Draw borders
	// for col := x1; col <= x2; col++ {
	// 	s.SetContent(col, y1, tcell.RuneHLine, nil, style)
	// 	s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	// }
	// for row := y1 + 1; row < y2; row++ {
	// 	s.SetContent(x1, row, tcell.RuneVLine, nil, style)
	// 	s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	// }
	//
	// // Only draw corners if necessary
	// if y1 != y2 && x1 != x2 {
	// 	s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
	// 	s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
	// 	s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
	// 	s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	// }

	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}
func Run(s tcell.Screen, world World) {
	for {

		s.Clear()

		r := Rule{}
		world = r.modify(world)
		world.update(s)
		s.Show()

		time.Sleep(40 * time.Millisecond)
	}
}

func main() {

	world := World{}
	world.init(Pair{20, 20})
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)

	// Draw initial boxes

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Here's how to get the screen size when you need it.
	// xmax, ymax := s.Size()

	// Here's an example of how to inject a keystroke where it will
	// be picked up by the next PollEvent call.  Note that the
	// queue is LIFO, it has a limited length, and PostEvent() can
	// return an error.
	// s.PostEvent(tcell.NewEventKey(tcell.KeyRune, rune('a'), 0))

	// Event loop
	go Run(s, world)

	for {

		// // Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			}
		}
	}

}
