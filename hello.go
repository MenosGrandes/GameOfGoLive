package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
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

type Cell struct {
	WorldPosition Pair
	IsAlive       bool
}

func (c Cell) getNeighbours(w World) Pair {

	pos := c.WorldPosition
	var neighbours = [8]Pair{
		pos.subPair(Pair{1, 1}),
		pos.subPair(Pair{0, 1}),
		pos.subPair(Pair{1, 0}),
		pos.subPair(Pair{0, -1}),
		pos.addPair(Pair{1, 1}),
		pos.addPair(Pair{0, 1}),
		pos.addPair(Pair{1, 0}),
		pos.addPair(Pair{0, -1}),
	}

	aliveCount := Pair{}
	log.Println("For current Pos:")
	log.Println(pos)
	for _, e := range neighbours {
		index := Convert2d1d(e)
		if index < MAP_SIZE_X && index >= 0 {
			log.Println("nei pos")
			log.Println(e)
			log.Println(index)
			if w.Cells[index].IsAlive == true {
				aliveCount.addValueXMut(1)
			} else {
				aliveCount.addValueYMut(1)
			}
		}

	}

	log.Println("aliveCount")
	log.Println(aliveCount)
	return aliveCount
}
func (c *Cell) update(w World) {
	aliveDeath := c.getNeighbours(w)

	log.Println("!@!@")
	log.Println(c.WorldPosition)
	log.Println(aliveDeath)
}

func (c Cell) draw(s tcell.Screen) {
	// 3 x 2 is a squere
	corners := worldToScreenPosition(c.WorldPosition)
	upperLeftCorner := corners[0]
	lowerRightCorner := corners[1]

	cellAliveStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)
	cellDeadStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorYellow)
	if c.IsAlive {
		drawBox(s, upperLeftCorner.X, upperLeftCorner.Y, lowerRightCorner.X, lowerRightCorner.Y, cellAliveStyle, "")
	} else {
		drawBox(s, upperLeftCorner.X, upperLeftCorner.Y, lowerRightCorner.X, lowerRightCorner.Y, cellDeadStyle, "")
	}
}

// 12 x 12 size?
const (
	MAP_SIZE_X = 3
	MAP_SIZE_Y = 3
	MAP_SIZE   = MAP_SIZE_X * MAP_SIZE_Y
)

func Convert1d2(i int, max int) Pair {
	return Pair{i % max, i / max}
}
func Convert2d1d(i Pair) int {
	return i.X*MAP_SIZE_X + i.Y
}

type World struct {
	Cells [MAP_SIZE]Cell
}

func (w *World) init() {
	for i := 0; i < MAP_SIZE; i++ {
		w.Cells[i].WorldPosition = Convert1d2(i, MAP_SIZE_X)
	}
	w.Cells[4].IsAlive = true
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
	for _, e := range w.Cells {
		e.update(w)
	}
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

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}

	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}

func updateScreen(s tcell.Screen, world World) {
	s.Clear()
	w, h := s.Size()
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)
	drawBox(s, 1, 1, w-1, h-1, boxStyle, "Click and drag to draw a box")
	// drawBox(s, 1, 1, 3, 2, boxStyle, "1")
	// drawBox(s, 4, 1, 6, 2, boxStyle, "2")
	// drawBox(s, 7, 1, 9, 2, boxStyle, "3")
	//
	// drawBox(s, 1, 3, 3, 4, boxStyle, "4")
	// drawBox(s, 4, 3, 6, 4, boxStyle, "5")
	// drawBox(s, 7, 3, 9, 4, boxStyle, "5")
	world.update(s)
}
func main() {
	world := World{}
	world.init()
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

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
	ox, oy := -1, -1
	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			updateScreen(s, world)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				s.Clear()
			}
		case *tcell.EventMouse:
			x, y := ev.Position()

			switch ev.Buttons() {
			case tcell.Button1, tcell.Button2:
				if ox < 0 {
					ox, oy = x, y // record location when click started
				}

			case tcell.ButtonNone:
				if ox >= 0 {
					label := fmt.Sprintf("%d,%d to %d,%d", ox, oy, x, y)
					drawBox(s, ox, oy, x, y, boxStyle, label)
					ox, oy = -1, -1
				}
			}
		}
	}
}
