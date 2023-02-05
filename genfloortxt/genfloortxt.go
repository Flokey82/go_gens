// Package genfloortxt renders a primitive floor plan read from an ASCII file to the console.
// NOTE: In future I plan to add some procedural generation bits to it.
package genfloortxt

import (
	"bufio"
	"io"
	"strings"
)

// cp437 mapping from char code to unicode replacement rune.
// var cp437 = []rune("\x00☺☻♥♦♣♠•◘○◙♂♀♪♬☼►◄↕‼¶§▬↨↑↓→←∟↔▲▼ !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⌂ÇüéâäàåçêëèïîìÄÅÉæÆôöòûùÿÖÜ¢£¥₧ƒáíóúñÑªº¿⌐¬½¼¡«»░▒▓│┤╡╢╖╕╣║╗╝╜╛┐└┴┬├─┼╞╟╚╔╩╦╠═╬╧╨╤╥╙╘╒╓╫╪┘┌█▄▌▐▀αßΓπΣσµτΦΘΩδ∞φε∩≡±≥≤⌠⌡÷≈°∙·√ⁿ²■\u00A0")

const (
	CellWall   = '#'
	CellWindow = 'W'
	CellDoor   = 'D'
)

// Plan represents a parsed floorplan.
type Plan struct {
	cells  [][]byte
	Height int
	Width  int
}

// ReadPlan reads a floor plan from the given reader.
func ReadPlan(r io.Reader) *Plan {
	p := &Plan{}
	br := bufio.NewReader(r)
	var maxLen int
	for {
		// Read each line and transform it into a byte array.
		line, _, err := br.ReadLine()
		if err != nil {
			break
		}
		ln := make([]byte, len(line))
		copy(ln, line)
		p.cells = append(p.cells, ln)
		if len(ln) > maxLen {
			maxLen = len(ln)
		}
	}

	// Store the dimensions of the floor plan.
	p.Height = len(p.cells)
	p.Width = maxLen

	// Pad the cells with empty space if needed.
	for i, l := range p.cells {
		if len(l) < maxLen {
			p.cells[i] = append(p.cells[i], make([]byte, maxLen-len(l))...)
		}
	}
	return p
}

// Render 'renders' the floor plan to an array of strings.
func (p *Plan) Render() (lines []string) {
	// Iterate over the cells and render them.
	for y, xRange := range p.cells {
		var sb strings.Builder
		for x := range xRange {
			// Based on the cellVal there are different ways to render it.
			sb.WriteRune(p.renderCell(x, y))
		}
		lines = append(lines, sb.String())
	}
	return lines
}

// renderCell returns the rune for the given coordinates.
func (p *Plan) renderCell(x, y int) rune {
	// Get the values of neighboring cells.
	var nc, ec, sc, wc byte
	if y > 0 {
		nc = p.cells[y-1][x]
	}
	if y < p.Height-1 {
		sc = p.cells[y+1][x]
	}
	if x > 0 {
		wc = p.cells[y][x-1]
	}
	if x < p.Width-1 {
		ec = p.cells[y][x+1]
	}

	// Based on the cellVal there are different ways to render it.
	switch cell := p.cells[y][x]; cell {
	case CellWall:
		// Encode the type of the cell based on the neighboring cells.
		t := encodeType(nc, ec, sc, wc, cell)
		switch t {
		case 10: // 1010
			return '═'
		case 3: // 0011
			return '╚'
		case 6: // 0110
			return '╔'
		case 7: // 0111
			return '╠'
		case 13: // 1101
			return '╣'
		case 5: // 0101
			return '║'
		case 9: // 1001
			return '╝'
		case 12: // 1100
			return '╗'
		case 11: // 1011
			return '╩'
		case 14: // 1110
			return '╦'
		case 15: // 1111
			return '╬'
		default:
			return '█'
		}
	case CellWindow:
		return '░'
	case CellDoor:
		return '▒'
	}
	return ' '
}

// encodeType flips bits 0, 1, 2, 3 if 'nc', 'ec', 'sc', or 'wc' are == 'cell' type.
// Every combination of matched neighbor types represent a unique number.
func encodeType(nc, ec, sc, wc, cell byte) (res byte) {
	if nc == cell {
		res |= 1
	}
	if ec == cell {
		res |= 1 << 1
	}
	if sc == cell {
		res |= 1 << 2
	}
	if wc == cell {
		res |= 1 << 3
	}
	return res
}
