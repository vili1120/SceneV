package lang

import (
	"fmt"
	"math"
)

//type Val interface {
//  SetPos(pos_start, pos_end *Position) *Val
//  SetContext(context *Context) *Val
//  Copy() *Val
//  Add(other *Val) (*Val, *Error)
//  Sub(other *Val) (*Val, *Error)
//  Mul(other *Val) (*Val, *Error)
//  Div(other *Val) (*Val, *Error)
//  Pow(other *Val) (*Val, *Error)
//  CompEQ(other *Val) (*Val, *Error)
//  CompNE(other *Val) (*Val, *Error)
//  CompLT(other *Val) (*Val, *Error)
//  CompGT(other *Val) (*Val, *Error)
//  CompLTE(other *Val) (*Val, *Error)
//  CompGTE(other *Val) (*Val, *Error)
//  And(other *Val) (*Val, *Error)
//  Or(other *Val) (*Val, *Error)
//  Not() (Val, *Error)
//}

type Val interface {
  SetPos(*Position, *Position) Val
  SetContext(*Context) Val
  Copy() Val
  Add(Val) (Val, *Error)
  Sub(Val) (Val, *Error)
  Mul(Val) (Val, *Error)
  Div(Val) (Val, *Error)
  Pow(Val) (Val, *Error)
  CompEQ(Val) (Val, *Error)
  CompNE(Val) (Val, *Error)
  CompLT(Val) (Val, *Error)
  CompGT(Val) (Val, *Error)
  CompLTE(Val) (Val, *Error)
  CompGTE(Val) (Val, *Error)
  And(Val) (Val, *Error)
  Or(Val) (Val, *Error) 
  Not() (Val, *Error)
  IsTrue() bool
  String() string
}

type Value struct {
  PosStart *Position
  PosEnd *Position
  Context *Context
}

func (v *Value) SetPos(pos_start, pos_end *Position) Val {
  v.PosStart = pos_start
  v.PosEnd = pos_end
  return v
}

func (v *Value) Copy() Val {
  panic("No copy method defined")
  return nil
}

func (v *Value) SetContext(context *Context) Val {
  v.Context = context
  return v
}

func (v *Value) Add(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) Sub(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) Mul(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) Div(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) Pow(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) CompEQ(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) CompNE(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) CompLT(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) CompGT(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) CompLTE(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) CompGTE(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) And(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) Or(other Val) (Val, *Error) {
  return nil, v.IllegalOperation(other)
}

func (v *Value) Not() (Val, *Error) {
  return nil, v.IllegalOperation(v)
}

func (v *Value) Execute(args []any) RTResult {
  res := RTResult{}
  return res.Failure(*v.IllegalOperation(nil))
}

func (v *Value) IsTrue() bool {
  return false
}

func (v *Value) IllegalOperation(other any) *Error {
  if other == nil { other = v }
  switch o := other.(type) {
    case *Number:
      return RTError(
        *v.PosStart, *o.PosEnd,
        "Illegal operation",
        *v.Context,
      )
  }
  return nil
}

func (v *Value) String() string {
  return ""
}

///////////////////////////////////////////////////////////////////////////////

func NewString(value string) Val {
  s := &StringVal{
    value: value,
  }
  s.SetPos(nil, nil)
  s.SetContext(nil)
  return s
}

type StringVal struct {
  Value
  value string
  Context *Context
}

func (s *StringVal) SetPos(PosStart, PosEnd *Position) Val {
  s.PosStart = PosStart
  s.PosEnd = PosEnd
  return s
}

func (s *StringVal) SetContext(context *Context) Val {
  s.Context = context
  return s
}

func (s *StringVal) copy() Val {
  copy := NewString(s.value)
  copy.SetPos(s.PosStart, s.PosEnd)
  copy.SetContext(s.Context)
  return copy
}

func (s *StringVal) Add(other Val) (Val, *Error) {
  switch o := other.(type) {
    case *StringVal:
      return NewString(s.value+o.value).SetContext(s.Context), nil
  }
  return nil, s.IllegalOperation(other)
}

func (s *StringVal) Div(other Val) (Val, *Error) {
  switch o := other.(type) {
    case *Number:
      runes := []rune(s.value)
      return NewString(string(runes[o.value.(int)])).SetContext(s.Context), nil
  }
  return nil, s.IllegalOperation(other)
}

func (s *StringVal) CompEQ(other Val) (Val, *Error) {
  switch o := other.(type) {
    case *StringVal:
      if s.value == o.value {
        return NewNumber(1).SetContext(s.Context), nil
      } else {
        return NewNumber(0).SetContext(s.Context), nil
      }
  }
  return nil, s.IllegalOperation(other)
}

func (s *StringVal) CompNE(other Val) (Val, *Error) {
  switch o := other.(type) {
    case *StringVal:
      if s.value != o.value {
        return NewNumber(1).SetContext(s.Context), nil
      } else {
        return NewNumber(0).SetContext(s.Context), nil
      }
  }
  return nil, s.IllegalOperation(other)
}

func (s *StringVal) IsTrue() bool {
  if len(s.value) > 0 {
    return true
  } else {
    return false
  }
}

func (s StringVal) String() string {
  return fmt.Sprintf("\"%v\"", s.value)
}

///////////////////////////////////////////////////////////////////////////////

func NewNumber(value any) Val {
  n := &Number{
    value: value,
  }
  n.SetPos(nil, nil)
  n.SetContext(nil)
  return n
}

type Number struct {
  Value
  value any
  Context *Context
}

func (n *Number) SetPos(pos_start, pos_end *Position) Val {
  n.PosStart = pos_start
  n.PosEnd = pos_end
  return n
}

func (n *Number) Copy() Val {
  copy := NewNumber(n.value)
  copy.SetPos(n.PosStart, n.PosEnd)
  copy.SetContext(n.Context)
  return copy
}

func (n *Number) SetContext(context *Context) Val {
  n.Context = context
  return n
}

func (n *Number) Add(other Val) (Val, *Error) {
  switch o := other.(type) {
    case *Number:
	    switch v1 := n.value.(type) {
	    case int:
	    	switch v2 := o.value.(type) {
	    	case int:
	    		return NewNumber(v1 + v2).SetContext(n.Context), nil
	    	case float64:
	    		return NewNumber(float64(v1) + v2).SetContext(n.Context), nil
	    	}
	    case float64:
	    	switch v2 := o.value.(type) {
	    	case int:
	    		return NewNumber(v1 + float64(v2)).SetContext(n.Context), nil
	    	case float64:
	    		return NewNumber(v1 + v2).SetContext(n.Context), nil
	    	}
	    }
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) Sub(other Val) (Val, *Error) {
  switch o := other.(type) {
    case *Number:
	    switch v1 := n.value.(type) {
	    case int:
	    	switch v2 := o.value.(type) {
	    	case int:
	    		return NewNumber(v1 - v2).SetContext(n.Context), nil
	    	case float64:
	    		return NewNumber(float64(v1) - v2).SetContext(n.Context), nil
	    	}
	    case float64:
	    	switch v2 := o.value.(type) {
	    	case int:
	    		return NewNumber(v1 - float64(v2)).SetContext(n.Context), nil
	    	case float64:
	    		return NewNumber(v1 - v2).SetContext(n.Context), nil
	    	}
	    }
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) Mul(other Val) (Val, *Error) {
  switch o := other.(type) {
    case *Number:
	    switch v1 := n.value.(type) {
	    case int:
	    	switch v2 := o.value.(type) {
	    	case int:
	    		return NewNumber(v1 * v2).SetContext(n.Context), nil
	    	case float64:
	    		return NewNumber(float64(v1) * v2).SetContext(n.Context), nil
	    	}
	    case float64:
	    	switch v2 := o.value.(type) {
	    	case int:
	    		return NewNumber(v1 * float64(v2)).SetContext(n.Context), nil
	    	case float64:
	    		return NewNumber(v1 * v2).SetContext(n.Context), nil
	    	}
	    }
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) Div(other Val) (Val, *Error) {
  switch o := other.(type) {
    case *Number:
      if o.value == 0 || o.value == 0.0 {
        return nil, RTError(
          *o.PosStart, *o.PosEnd,
          "Division by zero",
          *n.Context,
        )
      }
	    switch v1 := n.value.(type) {
	    case int:
	    	switch v2 := o.value.(type) {
	    	case int:
	    		return NewNumber(float64(v1) / float64(v2)).SetContext(n.Context), nil
	    	case float64:
	    		return NewNumber(float64(v1) / v2).SetContext(n.Context), nil
	    	}
	    case float64:
	    	switch v2 := o.value.(type) {
	    	case int:
	    		return NewNumber(v1 / float64(v2)).SetContext(n.Context), nil
	    	case float64:
	    		return NewNumber(v1 / v2).SetContext(n.Context), nil
	    	}
	    }
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) Pow(other Val) (Val, *Error) {
  switch o := other.(type) {
    case *Number:
	    switch v1 := n.value.(type) {
	    case int:
	    	switch v2 := o.value.(type) {
	    	case int:
	    		return NewNumber(math.Pow(float64(v1), float64(v2))).SetContext(n.Context), nil
	    	case float64:
	    		return NewNumber(math.Pow(float64(v1), v2)).SetContext(n.Context), nil
	    	}
	    case float64:
	    	switch v2 := o.value.(type) {
	    	case int:
	    		return NewNumber(math.Pow(v1, float64(v2))).SetContext(n.Context), nil
	    	case float64:
	    		return NewNumber(math.Pow(v1, v2)).SetContext(n.Context), nil
	    	}
	    }
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) CompEQ(other Val) (Val, *Error) {
  switch other := other.(type) {
    case *Number:
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
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) CompNE(other Val) (Val, *Error) {
  switch other := other.(type) {
    case *Number:
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
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) CompLT(other Val) (Val, *Error) {
  switch other := other.(type) {
    case *Number:
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
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) CompGT(other Val) (Val, *Error) {
  switch other := other.(type) {
    case *Number:
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
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) CompLTE(other Val) (Val, *Error) {
  switch other := other.(type) {
    case *Number:
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
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) CompGTE(other Val) (Val, *Error) {
  switch other := other.(type) {
    case *Number:
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
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) And(other Val) (Val, *Error) {
  switch other := other.(type) {
    case *Number:
	  return NewNumber(BoolToInt(NumToBool(n.value) && NumToBool(other.value))).SetContext(n.Context), nil
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) Or(other Val) (Val, *Error) {
  switch other := other.(type) {
    case *Number:
	  return NewNumber(BoolToInt(NumToBool(n.value) || NumToBool(other.value))).SetContext(n.Context), nil
  }
  return nil, n.Value.IllegalOperation(other)
}

func (n *Number) Not() (Val, *Error) {
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

///////////////////////////////////////////////////////////////////////////

func NewFunction(name string, body Node, argNames []string) *Function {
  var Name string
  if name == "" {
    Name = "<anonymous>"
  } else {
    Name = name
  }
  f := &Function{
    Name: Name,
    BodyNode: body,
    ArgNames: argNames,
  }
  f.SetPos(nil, nil)
  f.SetContext(nil)
  return f
}

type Function struct {
  Value
  Name string
  BodyNode Node
  ArgNames []string
}

func (f *Function) SetPos(pos_start, pos_end *Position) Val {
  f.PosStart = pos_start
  f.PosEnd = pos_end
  return f
}

func (f *Function) SetContext(context *Context) Val {
  f.Context = context
  return f
}

func (f *Function) Execute(args []Val) (RTResult){
  res := RTResult{}
  interpreter := Interpreter{}

  newCtx := Context{DisplayName: f.Name, Parent: f.Context, ParentEntryPos: f.PosStart}
  newCtx.SymbolTable = NewSymbolTable(newCtx.Parent.SymbolTable)

  if len(args) > len(f.ArgNames) {
    return res.Failure(*RTError(
      *f.PosStart, *f.PosEnd,
      fmt.Sprintf("%v too few args passed into '%v'", len(f.ArgNames)-len(args), f.Name),
      *f.Context,
    ))
  }
  if len(args) < len(f.ArgNames) {
    return res.Failure(*RTError(
      *f.PosStart, *f.PosEnd,
      fmt.Sprintf("%v too many args passed into '%v'", len(f.ArgNames)-len(args), f.Name),
      *f.Context,
    ))
  }

  for i := range len(args) {
    argName := f.ArgNames[i]
    argVal := args[i]
    argVal.SetContext(&newCtx)
    newCtx.SymbolTable.Set(argName, argVal)
  }
  Val := res.Register(interpreter.Visit(f.BodyNode, newCtx))
  if res.error != nil { return res }
  return res.Success(Val)
}

func (f *Function) Copy() Val {
  copy := Function{Name: f.Name, BodyNode: f.BodyNode, ArgNames: f.ArgNames}
  copy.SetContext(f.Context)
  copy.SetPos(f.PosStart, f.PosEnd)
  return &copy
}

func (f Function) String() string {
  return fmt.Sprintf("<function %v>", f.Name)
}
