package lang

import (
	"fmt"
	"strings"
)

func stringWithArrows(text string, posStart Position, posEnd Position) string {
	result := ""

	// Find the start of the line
	idxStart := strings.LastIndex(text[:posStart.idx], "\n")
	if idxStart == -1 {
		idxStart = 0
	} else {
		idxStart += 1
	}

	// Find the end of the line
	idxEnd := strings.Index(text[idxStart:], "\n")
	if idxEnd == -1 {
		idxEnd = len(text)
	} else {
		idxEnd += idxStart
	}

	line := text[idxStart:idxEnd]

	// Column positions
	colStart := posStart.col
	colEnd := posEnd.col

	// Ensure colEnd is not before colStart
	if colEnd <= colStart {
		colEnd = colStart + 1
	}

	// Build result string with arrows
	result += line + "\n"
	result += strings.Repeat(" ", colStart) + strings.Repeat("^", colEnd-colStart) + "\n"

	return strings.ReplaceAll(result, "\t", "")
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

type Error struct {
  PosStart Position
  PosEnd Position
  ErrorName string
  Details string
}

func (e Error) AsString() string {
  result := fmt.Sprintf("%v: %v", e.ErrorName, e.Details)
  result += "\n"+fmt.Sprintf("(File: %v, Line: %v)", e.PosStart.fn, e.PosEnd.ln + 1)
  result += "\n\n" + stringWithArrows(e.PosStart.ftxt, e.PosStart, e.PosEnd)
  return result
}

func IllegalCharError(pos_start, pos_end Position, details string) *Error {
  return &Error{pos_start, pos_end, "Illegal Character", details}
}

func InvalidSyntaxError(pos_start, pos_end Position, details string) *Error {
  return &Error{pos_start, pos_end, "Invalid Syntax", details}
}
