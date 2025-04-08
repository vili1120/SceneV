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

func (n *Number) Copy() *Number {
  copy := NewNumber(n.value)
  copy.SetPos(n.PosStart, n.PosEnd)
  copy.SetContext(n.Context)
  return copy
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

func (n *Number) CompEQ(other *Number) (*Number, *Error) {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(BoolToInt(v1 == v2)).SetContext(n.Context), nil
		case float64:
			return NewNumber(BoolToInt(float64(v1) == v2)).SetContext(n.Context), nil
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(BoolToInt(v1 == float64(v2))).SetContext(n.Context), nil
		case float64:
			return NewNumber(BoolToInt(v1 == v2)).SetContext(n.Context), nil
		}
	}
	return nil, nil
}

func (n *Number) CompNE(other *Number) (*Number, *Error) {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(BoolToInt(v1 != v2)).SetContext(n.Context), nil
		case float64:
			return NewNumber(BoolToInt(float64(v1) != v2)).SetContext(n.Context), nil
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(BoolToInt(v1 != float64(v2))).SetContext(n.Context), nil
		case float64:
			return NewNumber(BoolToInt(v1 != v2)).SetContext(n.Context), nil
		}
	}
	return nil, nil
}

func (n *Number) CompLT(other *Number) (*Number, *Error) {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(BoolToInt(v1 < v2)).SetContext(n.Context), nil
		case float64:
			return NewNumber(BoolToInt(float64(v1) < v2)).SetContext(n.Context), nil
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(BoolToInt(v1 < float64(v2))).SetContext(n.Context), nil
		case float64:
			return NewNumber(BoolToInt(v1 < v2)).SetContext(n.Context), nil
		}
	}
	return nil, nil
}

func (n *Number) CompGT(other *Number) (*Number, *Error) {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(BoolToInt(v1 > v2)).SetContext(n.Context), nil
		case float64:
			return NewNumber(BoolToInt(float64(v1) > v2)).SetContext(n.Context), nil
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(BoolToInt(v1 > float64(v2))).SetContext(n.Context), nil
		case float64:
			return NewNumber(BoolToInt(v1 > v2)).SetContext(n.Context), nil
		}
	}
	return nil, nil
}

func (n *Number) CompLTE(other *Number) (*Number, *Error) {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(BoolToInt(v1 <= v2)).SetContext(n.Context), nil
		case float64:
			return NewNumber(BoolToInt(float64(v1) <= v2)).SetContext(n.Context), nil
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(BoolToInt(v1 <= float64(v2))).SetContext(n.Context), nil
		case float64:
			return NewNumber(BoolToInt(v1 <= v2)).SetContext(n.Context), nil
		}
	}
	return nil, nil
}

func (n *Number) CompGTE(other *Number) (*Number, *Error) {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(BoolToInt(v1 >= v2)).SetContext(n.Context), nil
		case float64:
			return NewNumber(BoolToInt(float64(v1) >= v2)).SetContext(n.Context), nil
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(BoolToInt(v1 >= float64(v2))).SetContext(n.Context), nil
		case float64:
			return NewNumber(BoolToInt(v1 >= v2)).SetContext(n.Context), nil
		}
	}
	return nil, nil
}

func (n *Number) And(other *Number) (*Number, *Error) {
	return NewNumber(BoolToInt(NumToBool(n.value) && NumToBool(other.value))).SetContext(n.Context), nil
}

func (n *Number) Or(other *Number) (*Number, *Error) {
	return NewNumber(BoolToInt(NumToBool(n.value) || NumToBool(other.value))).SetContext(n.Context), nil
}

func (n *Number) Not() (*Number, *Error) {
  if n.value == 0 {
    return NewNumber(1).SetContext(n.Context), nil
  }
  return NewNumber(0).SetContext(n.Context), nil
}

func (n *Number) IsTrue() bool {
  return NumToBool(n.value)
}

func (n Number) String() string {
  return fmt.Sprintf("%v", n.value)
}

func BoolToInt(val bool) int {
  if val == true {
    return 1
  } else {
    return 0
  }
}

func NumToBool(num any) bool {
	switch v := num.(type) {
	case int:
		return v != 0
	case float64:
		return v != 0.0
	default:
		return false
	}
}
