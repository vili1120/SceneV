package lang

import "fmt"

type Node interface {
	String() string
	GetPosStart() Position
	GetPosEnd() Position
}

type NumberNode struct {
	Tok      Token
	PosStart Position
	PosEnd   Position
}

func (nn NumberNode) String() string {
  return fmt.Sprintf("%v", nn.Tok.String())
}

func (nn *NumberNode) GetPosStart() Position {
	return nn.PosStart
}

func (nn *NumberNode) GetPosEnd() Position {
	return nn.PosEnd
}

func (nn *NumberNode) SetPos() *NumberNode {
  nn.PosStart = nn.Tok.PosStart
  nn.PosEnd = nn.Tok.PosEnd
  return nn
}

type IfNode struct {
	Cases    [][]Node
  ElseCase Node
	PosStart Position
	PosEnd   Position
}

func (in IfNode) String() string {
  return ""
}

func (in *IfNode) GetPosStart() Position {
	return in.PosStart
}

func (in *IfNode) GetPosEnd() Position {
	return in.PosEnd
}

func (in *IfNode) SetPos() *IfNode {
  in.PosStart = in.Cases[0][0].GetPosStart()
  if in.ElseCase != nil {
    in.PosEnd = in.ElseCase.GetPosEnd()
  } else {
    in.PosEnd = in.Cases[len(in.Cases)-1][0].GetPosEnd()
  }
  return in
}

type ForNode struct {
	VarNameTok Token
  StartVal Node
  EndVal Node
  StepVal Node
  BodyNode Node
	PosStart Position
	PosEnd   Position
}

func (fn ForNode) String() string {
  return ""
}

func (fn *ForNode) GetPosStart() Position {
	return fn.PosStart
}

func (fn *ForNode) GetPosEnd() Position {
	return fn.PosEnd
}

func (fn *ForNode) SetPos() *ForNode {
  fn.PosStart = fn.VarNameTok.PosStart
  fn.PosEnd = fn.BodyNode.GetPosEnd()
  return fn
}

type WhileNode struct {
	Cond     Node
  BodyNode Node
	PosStart Position
	PosEnd   Position
}

func (wn WhileNode) String() string {
  return ""
}

func (wn *WhileNode) GetPosStart() Position {
	return wn.PosStart
}

func (wn *WhileNode) GetPosEnd() Position {
	return wn.PosEnd
}

func (wn *WhileNode) SetPos() *WhileNode {
  wn.PosStart = wn.Cond.GetPosStart()
  wn.PosEnd = wn.BodyNode.GetPosEnd()
  return wn
}

type VarAccessNode struct {
	VarName  Token
	PosStart Position
	PosEnd   Position
}

func (van VarAccessNode) String() string {
  return ""
}

func (van *VarAccessNode) GetPosStart() Position {
	return van.PosStart
}

func (van *VarAccessNode) GetPosEnd() Position {
	return van.PosEnd
}

func (van *VarAccessNode) SetPos() *VarAccessNode {
  van.PosStart = van.VarName.PosStart
  van.PosEnd = van.VarName.PosEnd
  return van
}

type VarAssignNode struct {
	VarName     Token
  ValueNode   Node
	PosStart    Position
	PosEnd      Position
}

func (van VarAssignNode) String() string {
  return ""
}

func (van *VarAssignNode) GetPosStart() Position {
	return van.PosStart
}

func (van *VarAssignNode) GetPosEnd() Position {
	return van.PosEnd
}

func (van *VarAssignNode) SetPos() *VarAssignNode {
  van.PosStart = van.VarName.PosStart
  van.PosEnd = van.ValueNode.GetPosEnd()
  return van
}

type BinOpNode struct {
	LeftNode  Node
	OpTok     Token
	RightNode Node
	PosStart  Position
	PosEnd    Position
}

func (bon BinOpNode) String() string {
  return fmt.Sprintf("(%v, %v, %v)", bon.LeftNode.String(), bon.OpTok.String(), bon.RightNode.String())
}

func (bon *BinOpNode) GetPosStart() Position {
	return bon.PosStart
}

func (bon *BinOpNode) GetPosEnd() Position {
	return bon.PosEnd
}

func (bon *BinOpNode) SetPos() *BinOpNode {
  bon.PosStart = bon.LeftNode.GetPosStart()
  bon.PosEnd = bon.RightNode.GetPosEnd()
  return bon
}

type UnaryOpNode struct {
	OpTok    Token
	Node     Node
	PosStart Position
	PosEnd   Position
}

func (uop UnaryOpNode) String() string {
  return fmt.Sprintf("(%v, %v)", uop.OpTok, uop.Node)
}

func (uop *UnaryOpNode) GetPosStart() Position {
	return uop.PosStart
}

func (uop *UnaryOpNode) GetPosEnd() Position {
	return uop.PosEnd
}

func (uop *UnaryOpNode) SetPos() *UnaryOpNode {
  uop.PosStart = uop.OpTok.PosStart
  uop.PosEnd = uop.Node.GetPosEnd()
  return uop
}

type FuncDefNode struct {
  VarNameTok Token
  ArgNameToks []Token
  BodyNode Node
  PosStart Position
  PosEnd Position
}

func (fdn FuncDefNode) String() string {
  return ""
}

func (fdn *FuncDefNode) GetPosStart() Position {
	return fdn.PosStart
}

func (fdn *FuncDefNode) GetPosEnd() Position {
	return fdn.PosEnd
}

func (fdn *FuncDefNode) SetPos() *FuncDefNode {
  emptyTok := Token{}
  if fdn.VarNameTok != emptyTok {
    fdn.PosStart = fdn.VarNameTok.PosStart
  } else if len(fdn.ArgNameToks) > 0 {
    fdn.PosStart = fdn.ArgNameToks[0].PosStart
  } else {
    fdn.PosStart = fdn.BodyNode.GetPosStart()
  }
  fdn.PosEnd = fdn.BodyNode.GetPosEnd()
  return fdn
}

type CallNode struct {
  NodeToCall Node
  ArgNodes []Node
  PosStart Position
  PosEnd Position
}

func (cn CallNode) String() string {
  return ""
}

func (cn *CallNode) GetPosStart() Position {
	return cn.PosStart
}

func (cn *CallNode) GetPosEnd() Position {
	return cn.PosEnd
}

func (cn *CallNode) SetPos() *CallNode {
  cn.PosStart = cn.NodeToCall.GetPosStart()
  if len(cn.ArgNodes) > 0 {
    cn.PosEnd = cn.ArgNodes[len(cn.ArgNodes)-1].GetPosEnd()
  } else {
    cn.PosEnd = cn.NodeToCall.GetPosEnd()
  }
  return cn
}
