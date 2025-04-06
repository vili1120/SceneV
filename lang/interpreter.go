package lang

import (
	"fmt"
	"reflect"
)

type RTResult struct {
  value any
  error *Error
}

func (rtr *RTResult) Register(res RTResult) any {
  if res.error != nil {
    rtr.error = res.error
  }
  return res.value
}

func (rtr *RTResult) Success(value any) RTResult{
  rtr.value = value
  return *rtr
}

func (rtr *RTResult) Failure(error Error) RTResult{
  rtr.error = &error
  return *rtr
}

type Interpreter struct {
  node Node
}

func (i *Interpreter) Visit(node Node, context Context) RTResult {
	nodeType := reflect.TypeOf(node)
	if nodeType.Kind() == reflect.Ptr {
		nodeType = nodeType.Elem()
	}
	MethodName := fmt.Sprintf("Visit%v", nodeType.Name())
	method := reflect.ValueOf(i).MethodByName(MethodName)
	if !method.IsValid() {
		panic("No " + MethodName + " method defined")
	}

	results := method.Call([]reflect.Value{reflect.ValueOf(node), reflect.ValueOf(context)})
	if len(results) == 0 {
		panic("Visit method did not return anything")
	}

	rtResult, ok := results[0].Interface().(RTResult)
	if !ok {
		panic(fmt.Sprintf("Visit method returned unexpected type: %T", results[0].Interface()))
	}

	return rtResult
}


func (i *Interpreter) VisitNumberNode(node *NumberNode, context Context) RTResult {
  res := RTResult{}
  return res.Success(
    NewNumber(node.Tok.value).SetContext(&context).SetPos(&node.PosStart, &node.PosEnd),
  )
}

func (i *Interpreter) VisitVarAccessNode(node *VarAccessNode, context Context) RTResult {
  res := RTResult{}
  var_name := node.VarName.value
  value := context.SymbolTable.Get(var_name.(string))
  if value == nil {
    return res.Failure(*RTError(
      node.PosStart, node.PosEnd,
      fmt.Sprintf("'%v' is not defined", var_name),
      context,
    ))
  }

  value = value.(*Number).Copy().SetPos(&node.PosStart, &node.PosEnd)
  return res.Success(value)
}

func (i *Interpreter) VisitVarAssignNode(node *VarAssignNode, context Context) RTResult {
  res := RTResult{}
  var_name := node.VarName.value
  value := res.Register(i.Visit(node.ValueNode, context))
  if res.error != nil { return res }
  context.SymbolTable.Set(var_name.(string), value)
  return res.Success(value)
}

func (i *Interpreter) VisitBinOpNode(node *BinOpNode, context Context) RTResult{
  res := RTResult{}
  left := res.Register(i.Visit(node.LeftNode, context))
  if res.error != nil {return res}
  right := res.Register(i.Visit(node.RightNode, context))
  if res.error != nil {return res}

	leftNum, ok1 := left.(*Number)
	rightNum, ok2 := right.(*Number)

	if !ok1 || !ok2 {
		panic("Operands must be numbers")
	}

  var result *Number
  var err *Error

  if node.OpTok.type_ == PLUS {
    result, err = leftNum.Add(rightNum)
  } else if node.OpTok.type_ == MINUS {
    result, err = leftNum.Sub(rightNum)
  } else if node.OpTok.type_ == MUL {
    result, err = leftNum.Mul(rightNum)
  } else if node.OpTok.type_ == DIV {
    result, err = leftNum.Div(rightNum)
  } else if node.OpTok.type_ == POW {
    result, err = leftNum.Pow(rightNum)
  }
  if err != nil {
    return res.Failure(*err)
  } else {
    return res.Success(result.SetPos(&node.PosStart, &node.PosEnd))
  }
}

func (i *Interpreter) VisitUnaryOpNode(node *UnaryOpNode, context Context) RTResult {
  res := RTResult{}
  number := res.Register(i.Visit(node.Node, context))
  if res.error != nil {return res}
  num, ok := number.(*Number)
  if !ok {
    panic("Operand must a number")
  }

  var err *Error = nil

  if node.OpTok.type_ == MINUS {
    num, err = num.Mul(NewNumber(-1))
  }

  if err != nil {
    return res.Failure(*err)
  } else {
    return res.Success(num.SetPos(&node.PosStart, &node.PosEnd))
  }
}

