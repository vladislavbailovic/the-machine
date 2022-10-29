package debug

import (
	"fmt"
	"strings"
)

type Resolution uint8

const (
	Byte Resolution = 0
	Uint Resolution = iota
)

type Representation uint8

const (
	Binary  Representation = 0
	Hex     Representation = iota
	Decimal Representation = iota
)

type RenderingDirection uint8

const (
	Horizontal RenderingDirection = 0
	Vertical   RenderingDirection = iota
)

type Formatter struct {
	Numbers   Representation
	OutputAs  Resolution
	Rendering RenderingDirection
}

func (x Formatter) GetFormat() (string, string) {
	posFormat := "%4d"
	valFormat := "%#02x"
	switch x.Numbers {
	case Binary:
		switch x.OutputAs {
		case Byte:
			posFormat = "%10d"
			valFormat = "%#08b"
		case Uint:
			posFormat = "%18d"
			valFormat = "%#016b"
		}
	case Decimal:
		switch x.OutputAs {
		case Byte:
			posFormat = "%3d"
			valFormat = "%3d"
		case Uint:
			posFormat = "%5d"
			valFormat = "%05d"
		}

	}
	return posFormat, valFormat
}

func (x Formatter) Stitch(first []string, rest ...[]string) string {
	switch x.Rendering {
	case Vertical:
		return x.stitchCols(first, rest...)
	case Horizontal:
		return x.stitchRows(first, rest...)
	default:
		return fmt.Sprintf("ERROR: unknown rendering direction: %d", x.Rendering)
	}
}

func (x Formatter) stitchRows(first []string, rest ...[]string) string {
	out := make([]string, len(rest)+1)
	out[0] = strings.Join(first, " ")
	separator := strings.Repeat("-", len(out[0]))
	for idx, item := range rest {
		out[idx+1] = strings.Join(item, " ")
	}
	return strings.Join(out, fmt.Sprintf("\n%s\n", separator))
}

func (x Formatter) stitchCols(first []string, rest ...[]string) string {
	cols := make([]string, len(rest)+1)
	rows := make([]string, len(first))

	for rowIdx, item := range first {
		cols[0] = item
		ln := len(item)
		for colIdx, col := range rest {
			if rowIdx < len(col) {
				cols[colIdx+1] = col[rowIdx]
				ln = len(col[rowIdx])
			} else {
				cols[colIdx+1] = strings.Repeat(" ", ln)
			}
		}
		rows[rowIdx] = strings.Join(cols, " | ")
	}
	return strings.Join(rows, "\n")
}
