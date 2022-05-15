package genlsystem

func Hilbert3d(fname string, n int) error {
	path := Lindenmayer([]string{"X"}, map[string][]string{
		"X": {"^", "<", "X", "F", "^", "<", "X", "F", "X", "-", "F", "^", ">", ">", "X", "F", "X", "&", "F", "+", ">", ">", "X", "F", "X", "-", "F", ">", "X", "-", ">"},
	}, n)

	turtle := NewTurtle3d(map[string]func(*Turtle3d){
		"F": func(t *Turtle3d) {
			t.Draw(0.2)
		},
		"-": func(t *Turtle3d) {
			t.Rotate(-90)
		},
		"+": func(t *Turtle3d) {
			t.Rotate(90)
		},
		"^": func(t *Turtle3d) {
			t.Pitch(90)
		},
		"&": func(t *Turtle3d) {
			t.Pitch(-90)
		},
		">": func(t *Turtle3d) {
			t.Roll(90)
		},
		"<": func(t *Turtle3d) {
			t.Roll(-90)
		},
	})

	return turtle.Go(fname, path)
}

func Plant3d(fname string, n int) error {
	path := Lindenmayer([]string{"F"}, map[string][]string{
		"F": {"F", "F", "-", "[", "-", "F", "+", "F", "+", "F", "]", "+", "[", "+", "F", "-", "F", "-", "F", "]"},
		"-": {"-", ">"},
		"+": {"+", "<"},
	}, n)

	turtle := NewTurtle3d(map[string]func(*Turtle3d){
		"F": func(t *Turtle3d) {
			t.Draw(0.2)
		},
		"-": func(t *Turtle3d) {
			t.Rotate(-23)
		},
		"+": func(t *Turtle3d) {
			t.Rotate(23)
		},
		"^": func(t *Turtle3d) {
			t.Pitch(23)
		},
		"&": func(t *Turtle3d) {
			t.Pitch(-23)
		},
		">": func(t *Turtle3d) {
			t.Roll(23)
		},
		"<": func(t *Turtle3d) {
			t.Roll(-23)
		},
		"[": func(t *Turtle3d) { // push position and angle
			t.Save()
		},
		"]": func(t *Turtle3d) { // pop position and angle
			t.Restore()
		},
		"|": func(t *Turtle3d) { // turn around
			t.Rotate(180)
		},
	})

	return turtle.Go(fname, path)
}
