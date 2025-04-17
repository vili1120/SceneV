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

func (i *Interpreter) VisitStringNode(node *StringNode, context Context) RTResult {
  res := RTResult{}
  return res.Success(
    NewString(node.Tok.value.(string)).SetContext(&context).SetPos(&node.PosStart, &node.PosEnd),
  )
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
  //switch value.(type) {
  //  case *Number:
  //    value = value.(*Number).Copy().SetPos(&node.PosStart, &node.PosEnd)
  //  case *Value:
  //    value = value.(*Value).Copy().SetPos(&node.PosStart, &node.PosEnd)
  //  case *Function:
  //    value = value.(*Function).Copy().SetPos(&node.PosStart, &node.PosEnd)
  //    return res.Success(value)
  //}
  value = value.Copy().SetPos(&node.PosStart, &node.PosEnd)
  return res.Success(value)
}

func (i *Interpreter) VisitVarAssignNode(node *VarAssignNode, context Context) RTResult {
  res := RTResult{}
  var_name := node.VarName.value
  value := res.Register(i.Visit(node.ValueNode, context))
  if res.error != nil { return res }
  context.SymbolTable.Set(var_name.(string), value.(Val))
  return res.Success(nil)
}

func (i *Interpreter) VisitIfNode(node *IfNode, context Context) RTResult {
  res := RTResult{}

  for _, Case := range node.Cases {
    condition := Case[0]
    expr := Case[1]

    cond_val := res.Register(i.Visit(condition, context))
    if res.error != nil { return res }
    
    switch v := cond_val.(type) {
      case Val:
        if v.IsTrue() {
          expr_val := res.Register(i.Visit(expr, context))
          if res.error != nil { return res }
          return res.Success(expr_val)
        }
    }
  }
  if node.ElseCase != nil {
    else_val := res.Register(i.Visit(node.ElseCase, context))
    if res.error != nil { return res }
    return res.Success(else_val)
  }
  return res.Success(nil)
}

func (i *Interpreter) VisitForNode(node *ForNode, context Context) RTResult {
  res := RTResult{}

  startVal := res.Register(i.Visit(node.StartVal, context))
  if res.error != nil { return res }
  startNum := startVal.(*Number)

  endVal := res.Register(i.Visit(node.EndVal, context))
  if res.error != nil { return res }
  endNum := endVal.(*Number)
  
  var stepNum *Number
  if node.StepVal != nil {
    stepVal := res.Register(i.Visit(node.StepVal, context))
    if res.error != nil { return res }
    stepNum = stepVal.(*Number)
  } else {
    stepNum = NewNumber(1).(*Number)
  }

  IVal := startNum.value.(int)
  StepVal := stepNum.value.(int)
  EndVal := endNum.value.(int)
  
  var condition func() bool
  if StepVal >= 0 {
    condition = func() bool { return IVal <= EndVal }
  } else {
    condition = func() bool { return IVal >= EndVal }
  }
  for condition() {
    context.SymbolTable.Set(node.VarNameTok.value.(string), NewNumber(IVal))
    IVal += StepVal

    res.Register(i.Visit(node.BodyNode, context))
    if res.error != nil { return res }
  }
  return res.Success(nil)
}

func (i *Interpreter) VisitWhileNode(node *WhileNode, context Context) RTResult {
  res := RTResult{}

  for true{
    condition := res.Register(i.Visit(node.Cond, context)).(*Number)
    if res.error != nil { return res }

    if !condition.IsTrue() { break }

    res.Register(i.Visit(node.BodyNode, context))
    if res.error != nil { return res }
  }
  return res.Success(nil)
}

func (i *Interpreter) VisitFuncDefNode(node *FuncDefNode, context Context) RTResult {
  res := RTResult{}
  funcName := node.VarNameTok.value
  
  body := node.BodyNode
  var arg_names []string
  for _, v := range node.ArgNameToks {
    arg_names = append(arg_names, v.value.(string))
  }

  funcValue := NewFunction(funcName.(string), body, arg_names).SetContext(&context).SetPos(&node.PosStart, &node.PosEnd)
  
  emptyTok := Token{value: ""}
  if node.VarNameTok != emptyTok {
    context.SymbolTable.Set(funcName.(string), funcValue)
  }

  return res.Success(funcValue)
}

func (i *Interpreter) VisitCallNode(node *CallNode, context Context) RTResult {
  res := RTResult{}
  args := []Val{}

  valueToCall := res.Register(i.Visit(node.NodeToCall, context))
  if res.error != nil { return res }
  var CallVal *Function
  if fn, ok := valueToCall.(*Function); ok {
    CallVal = fn.SetPos(&node.PosStart, &node.PosEnd).(*Function)
  }

  for _, argNode := range node.ArgNodes {
    args = append(args, res.Register(i.Visit(argNode, context)).(Val))
    if res.error != nil { return res }
  }
  returnVal := res.Register(CallVal.Execute(args))
  if res.error != nil { return res }
  return res.Success(returnVal)
}

func (i *Interpreter) VisitBinOpNode(node *BinOpNode, context Context) RTResult{
  res := RTResult{}
  left := res.Register(i.Visit(node.LeftNode, context))
  if res.error != nil {return res}
  right := res.Register(i.Visit(node.RightNode, context))
  if res.error != nil {return res}

	leftNum, ok1 := left.(Val)
	rightNum, ok2 := right.(Val)

	if !ok1 || !ok2 {
		panic("Operands must be values")
	}

  var result Val
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
  } else if node.OpTok.type_ == EE {
    result, err = leftNum.CompEQ(rightNum)
  } else if node.OpTok.type_ == NE {
    result, err = leftNum.CompNE(rightNum)
  } else if node.OpTok.type_ == LT {
    result, err = leftNum.CompLT(rightNum)
  } else if node.OpTok.type_ == GT {
    result, err = leftNum.CompGT(rightNum)
  } else if node.OpTok.type_ == LTE {
    result, err = leftNum.CompLTE(rightNum)
  } else if node.OpTok.type_ == GTE {
    result, err = leftNum.CompGTE(rightNum)
  } else if  node.OpTok.Matches(KEYWORD, "and") {
    result, err = leftNum.And(rightNum)
  } else if  node.OpTok.Matches(KEYWORD, "or") {
    result, err = leftNum.Or(rightNum)
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
  num, ok := number.(Val)
  if !ok {
    panic("Operand must be a number")
  }

  var err *Error = nil

  if node.OpTok.type_ == MINUS {
    num, err = num.Mul(NewNumber(-1))
  } else if node.OpTok.Matches(KEYWORD, "not") {
    num, err = num.Not()
  }

  if err != nil {
    return res.Failure(*err)
  } else {
    return res.Success(num.SetPos(&node.PosStart, &node.PosEnd))
  }
}

