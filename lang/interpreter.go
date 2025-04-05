package lang

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
  node Node
}

func (i Interpreter) Visit(node Node) any {
	nodeType := reflect.TypeOf(node)
	if nodeType.Kind() == reflect.Ptr {
		nodeType = nodeType.Elem()
	}
  MethodName := fmt.Sprintf("Visit%v", nodeType.Name())
  method := reflect.ValueOf(i).MethodByName(MethodName)
  if !method.IsValid() {
  	panic("No " + MethodName + " method defined")
  }
  
  results := method.Call([]reflect.Value{reflect.ValueOf(node)})
  return results[0].Interface()
}

func (i Interpreter) VisitNumberNode(node *NumberNode) *Number {
  return NewNumber(node.Tok.value).SetPos(&node.PosStart, &node.PosEnd)
}

func (i Interpreter) VisitBinOpNode(node *BinOpNode) *Number{
  left := i.Visit(node.LeftNode)
  right := i.Visit(node.RightNode)

	leftNum, ok1 := left.(*Number)
	rightNum, ok2 := right.(*Number)

	if !ok1 || !ok2 {
		panic("Operands must be numbers")
	}

  var result *Number

  if node.OpTok.type_ == PLUS {
    result = leftNum.Add(rightNum)
  } else if node.OpTok.type_ == MINUS {
    result = leftNum.Sub(rightNum)
  } else if node.OpTok.type_ == MUL {
    result = leftNum.Mul(rightNum)
  } else if node.OpTok.type_ == DIV {
    result = leftNum.Div(rightNum)
  }

  return result.SetPos(&node.PosStart, &node.PosEnd)
}

func (i Interpreter) VisitUnaryOpNode(node *UnaryOpNode) *Number {
  number := i.Visit(node.Node)
  num, ok := number.(*Number)
  if !ok {
    panic("Operand must a number")
  }

  if node.OpTok.type_ == MINUS {
    num = num.Mul(NewNumber(-1))
  }

  return num.SetPos(&node.PosStart, &node.PosEnd)
}

