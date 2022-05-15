package genlsystem

import (
	"image"
	"image/color"
	"math/rand"
)

func Hilbert(n int) image.Image {
	path := Lindenmayer([]string{"A"}, map[string][]string{
		"A": {"-", "B", "F", "+", "A", "F", "A", "+", "F", "B", "-"},
		"B": {"+", "A", "F", "-", "B", "F", "B", "-", "F", "A", "+"},
	}, n)

	turtle := NewTurtle(map[string]func(*Turtle){
		"F": func(t *Turtle) {
			t.Draw(5.0, 0)
		},
		"-": func(t *Turtle) {
			t.Turn(-90)
		},
		"+": func(t *Turtle) {
			t.Turn(90)
		},
	})

	return turtle.Go(path)
}

func Tree(n int) image.Image {
	segmentlength := 10.0

	green := color.RGBA{0x60, 0xFF, 0x00, 0xFF}
	black := color.RGBA{0x33, 0x33, 0x33, 0xFF}

	path := Lindenmayer([]string{"leaf"}, map[string][]string{
		"leaf":   {"branch", "[", "<", "leaf", "]", "<>", "branch", "[", ">", "leaf", "]"},
		"branch": {"trunk", "[", "<", "leaf", "]", "[", "<>", "branch", "]", "[", ">", "leaf", "]"},
		"trunk":  {"+", "trunk", "<>", "trunk", "-"},
		"+":      {"+", "+"},
		"-":      {"-", "-"},
	}, n)

	rmm := func(min, max float64) float64 {
		return rand.Float64()*(max-min) + min
	}

	turtle := NewTurtle(map[string]func(*Turtle){
		"leaf": func(t *Turtle) { // first gen, green leave
			t.SetColor(green)
			w := t.Width()
			t.SetWidth(10 * rand.Float64())
			t.Draw(segmentlength/2*rand.Float64(), 0)
			t.SetWidth(w)
		},
		"branch": func(t *Turtle) { // second gen, black branch
			t.SetColor(black)
			//t.SetWidth(1)
			t.Draw(segmentlength*rand.Float64(), 0)
		},
		"trunk": func(t *Turtle) { // second gen, black branch
			t.SetColor(black)
			//t.SetWidth(2)
			t.Draw(segmentlength*rand.Float64(), 0)
		},
		"+": func(t *Turtle) { // thicken
			t.SetWidth(t.Width() + 0.25)
		},
		"-": func(t *Turtle) { // thinning
			// w := t.Width() - 1
			// if w < 1 { w = 1 }
			t.SetWidth(t.Width() - 0.25)
		},
		"<": func(t *Turtle) { // turn left
			t.Turn(rmm(-45, 0))
		},
		">": func(t *Turtle) { // turn right
			t.Turn(rmm(0, 45))
		},
		"<>": func(t *Turtle) { // wiggle
			t.Turn(rmm(-2, 2))
		},
		"[": func(t *Turtle) { // push position and angle
			t.Save()
		},
		"]": func(t *Turtle) { // pop position and angle
			t.Restore()
		},
	})

	turtle.SetColor(black)
	turtle.SetWidth(1)
	turtle.Turn(-90)

	return turtle.Go(path)
}

func BinTree(n int) image.Image {
	segmentlenght := 1.0

	green := color.RGBA{0x33, 0xFF, 0x33, 0xFF}
	black := color.RGBA{0x00, 0x00, 0x00, 0xFF}

	path := Lindenmayer([]string{"0"}, map[string][]string{
		"1": {"1", "1"},
		"0": {"1", "[", "0", "]", "0"},
	}, n)

	turtle := NewTurtle(map[string]func(*Turtle){
		"0": func(t *Turtle) {
			// draw a line segment ending in a leaf
			t.SetColor(green)
			t.SetWidth(5)
			t.Draw(segmentlenght, 0)
		},
		"1": func(t *Turtle) {
			// draw a line segment
			t.SetColor(black)
			t.SetWidth(1)
			t.Draw(segmentlenght, 0)
		},
		"[": func(t *Turtle) {
			// push position and angle, turn left 45 degrees
			t.Save()
			t.Turn(-45)
		},
		"]": func(t *Turtle) {
			// pop position and angle, turn right 45 degrees
			t.Restore()
			t.Turn(45)
		},
	})

	turtle.Turn(-90)
	return turtle.Go(path)
}

func Plant(n int) image.Image {
	segmentlength := 4.0

	green := color.RGBA{0x60, 0xFF, 0x00, 0xFF}

	path := Lindenmayer([]string{"X"}, map[string][]string{
		"X": {"F", "-", "[", "[", "X", "]", "+", "X", "]", "+", "F", "[", "+", "F", "X", "]", "-", "X"},
		"F": {"F", "F"},
	}, n)

	turtle := NewTurtle(map[string]func(*Turtle){
		"F": func(t *Turtle) {
			// draw forward
			t.Draw(segmentlength*rand.Float64(), 0)
		},
		"-": func(t *Turtle) {
			// turn left 25°
			t.Turn(-20 + rand.Float64()*10)
		},
		"+": func(t *Turtle) {
			// turn right 25°
			t.Turn(20 + rand.Float64()*10)
		},
		"[": func(t *Turtle) {
			// push position and angle
			t.Save()
		},
		"]": func(t *Turtle) {
			// pop position and angle
			t.Restore()
		},
	})

	turtle.SetColor(green)
	turtle.SetWidth(1)
	turtle.Turn(-90)

	return turtle.Go(path)
}
