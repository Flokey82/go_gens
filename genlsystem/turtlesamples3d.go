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

func Pyramid3d(fname string, n int) error {
	// Adapted from:
	// https://github.com/yalue/l_system_3d/blob/master/config.txt
	path := Lindenmayer([]string{"P"}, map[string][]string{
		"P": {"(", "(", "P", "N", "P", "N", "-", "N", "P", "N", "-", "N", "P", ")", "U", "N", "D", "P", ")"},
		"N": {"N", "N"},
	}, n)

	turtle := NewTurtle3d(map[string]func(*Turtle3d){
		"F": func(t *Turtle3d) {
			t.Draw(1)
		},
		"N": func(t *Turtle3d) {
			t.Move(1)
		},
		"-": func(t *Turtle3d) {
			t.Rotate(90)
		},
		"(": func(t *Turtle3d) { // push position and angle
			t.Save()
		},
		")": func(t *Turtle3d) { // pop position and angle
			t.Restore()
		},
		"U": func(t *Turtle3d) {
			// Face "upwards" to move to the upper pyramid.
			t.Rotate(45)
			t.Pitch(45)
		},
		"D": func(t *Turtle3d) {
			// Undo the "U" rotation in preparation for drawing the upper pyramid.
			t.Pitch(-45)
			t.Rotate(-45)
		},
		"P": func(t *Turtle3d) {
			// Keep track of our start position, we'll return here at the end.
			t.Save()

			// Bottom left -> bottom right, up edge from bottom right
			t.Draw(1)
			t.Save()
			t.Rotate(135)
			t.Pitch(45)
			t.Draw(1)
			t.Restore()

			// Bottom right -> up right, up edge from up right
			t.Rotate(90)
			t.Draw(1)
			t.Save()
			t.Rotate(135)
			t.Pitch(45)
			t.Draw(1)
			t.Restore()

			// up right -> up left, up edge from up left
			t.Rotate(90)
			t.Draw(1)
			t.Save()
			t.Rotate(135)
			t.Pitch(45)
			t.Draw(1)
			t.Restore()

			// up left -> bottom left, up edge from bottom left
			t.Rotate(90)
			t.Draw(1)
			t.Rotate(135)
			t.Pitch(45)
			t.Draw(1)

			// Return to bottom left, facing right.
			t.Restore()
		},
		">": func(t *Turtle3d) {
			t.Roll(23)
		},
		"<": func(t *Turtle3d) {
			t.Roll(-23)
		},
		"|": func(t *Turtle3d) { // turn around
			t.Rotate(180)
		},
	})

	return turtle.Go(fname, path)
}
