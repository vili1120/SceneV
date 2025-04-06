package lang

import (
	"fmt"
	"math"
)

func NewNumber(value any) *Number {
  n := &Number{
    value: value,
  }
  n.SetPos(nil, nil)
  n.SetContext(nil)
  return n
}

type Number struct {
  value any
  PosStart, PosEnd *Position
  Context *Context
}

func (n *Number) SetPos(pos_start, pos_end *Position) *Number {
  n.PosStart = pos_start
  n.PosEnd = pos_end
  return n
}

func (n *Number) SetContext(context *Context) *Number {
  n.Context = context
  return n
}

func (n *Number) Add(other *Number) (*Number, *Error) {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 + v2).SetContext(n.Context), nil
		case float64:
			return NewNumber(float64(v1) + v2).SetContext(n.Context), nil
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 + float64(v2)).SetContext(n.Context), nil
		case float64:
			return NewNumber(v1 + v2).SetContext(n.Context), nil
		}
	}
	return nil, nil
}

func (n *Number) Sub(other *Number) (*Number, *Error) {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 - v2).SetContext(n.Context), nil
		case float64:
			return NewNumber(float64(v1) - v2).SetContext(n.Context), nil
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 - float64(v2)).SetContext(n.Context), nil
		case float64:
			return NewNumber(v1 - v2).SetContext(n.Context), nil
		}
	}
	return nil, nil
}

func (n *Number) Mul(other *Number) (*Number, *Error) {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 * v2).SetContext(n.Context), nil
		case float64:
			return NewNumber(float64(v1) * v2).SetContext(n.Context), nil
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 * float64(v2)).SetContext(n.Context), nil
		case float64:
			return NewNumber(v1 * v2).SetContext(n.Context), nil
		}
	}
	return nil, nil
}

func (n *Number) Div(other *Number) (*Number, *Error) {
  if other.value == 0 || other.value == 0.0 {
    return nil, RTError(
      *other.PosStart, *other.PosEnd,
      "Division by zero",
      *n.Context,
    )
  }
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(float64(v1) / float64(v2)).SetContext(n.Context), nil
		case float64:
			return NewNumber(float64(v1) / v2).SetContext(n.Context), nil
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 / float64(v2)).SetContext(n.Context), nil
		case float64:
			return NewNumber(v1 / v2).SetContext(n.Context), nil
		}
	}
	return nil, nil
}

func (n *Number) Pow(other *Number) (*Number, *Error) {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(math.Pow(float64(v1), float64(v2))).SetContext(n.Context), nil
		case float64:
			return NewNumber(math.Pow(float64(v1), v2)).SetContext(n.Context), nil
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(math.Pow(v1, float64(v2))).SetContext(n.Context), nil
		case float64:
			return NewNumber(math.Pow(v1, v2)).SetContext(n.Context), nil
		}
	}
	return nil, nil
}


func (n Number) String() string {
  return fmt.Sprintf("%v", n.value)
}
