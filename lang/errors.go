package lang

import (
	"fmt"
	"strings"
)

// ======= Helper Function =======
func stringWithArrows(text string, posStart Position, posEnd Position) string {
	result := ""

	idxStart := strings.LastIndex(text[:posStart.idx], "\n")
	if idxStart == -1 {
		idxStart = 0
	} else {
		idxStart += 1
	}

	idxEnd := strings.Index(text[idxStart:], "\n")
	if idxEnd == -1 {
		idxEnd = len(text)
	} else {
		idxEnd += idxStart
	}

	line := text[idxStart:idxEnd]
	colStart := posStart.col
	colEnd := posEnd.col

	if colEnd <= colStart {
		colEnd = colStart + 1
	}

	result += line + "\n"
	result += strings.Repeat(" ", colStart) + strings.Repeat("^", colEnd-colStart) + "\n"

	return strings.ReplaceAll(result, "\t", "")
}

// ======= Base Error =======
type Error struct {
	PosStart  Position
	PosEnd    Position
	ErrorName string
	Details   string
  Type string
  Context *Context
}

func (e *Error) AsString() string {
  if e.Type == "rterror" && e.Context != nil {
    result := e.GenerateTraceback()
	  result += fmt.Sprintf("%v: %v", e.ErrorName, e.Details)
	  result += "\n\n" + stringWithArrows(e.PosStart.ftxt, e.PosStart, e.PosEnd)
	  return result
  }
	result := fmt.Sprintf("%v: %v", e.ErrorName, e.Details)
	result += "\n" + fmt.Sprintf("(File: %v, Line: %v)", e.PosStart.fn, e.PosEnd.ln+1)
	result += "\n\n" + stringWithArrows(e.PosStart.ftxt, e.PosStart, e.PosEnd)
	return result
}

func (e *Error) GenerateTraceback() string {
  result := ""
  pos := e.PosStart
  ctx := e.Context

  for ctx != nil {
    result = fmt.Sprintf("  (File %v, line %v, in %v)\n", pos.fn, pos.ln+1, ctx.DisplayName) + result
    if ctx.ParentEntryPos == nil {
      break
    }
    pos = *ctx.ParentEntryPos
    ctx = ctx.Parent
  }
  return "Traceback (most recent call last):\n" + result
}

func IllegalCharError(posStart, posEnd Position, details string) *Error {
	return &Error{
		PosStart:  posStart,
		PosEnd:    posEnd,
		ErrorName: "Illegal Character",
		Details:   details,
    Type: "",
	}
}

func ExpectedCharError(posStart, posEnd Position, details string) *Error {
	return &Error{
		PosStart:  posStart,
		PosEnd:    posEnd,
		ErrorName: "Expected Character",
		Details:   details,
    Type: "",
	}
}

func InvalidSyntaxError(posStart, posEnd Position, details string) *Error {
	return &Error{
		PosStart:  posStart,
		PosEnd:    posEnd,
		ErrorName: "Invalid Syntax",
		Details:   details,
    Type: "",
	}
}

func RTError(posStart, posEnd Position, details string, context Context) *Error {
	return &Error{
		PosStart:  posStart,
		PosEnd:    posEnd,
		ErrorName: "Illegal Character",
		Details:   details,
    Type: "rterror",
    Context: &context,
	}
}
