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

type ForNode struct {
	VarNameTok Token
  StartVal Node
  EndVal Node
  StepVal Node
  BodyNode Node
	PosStart Position
	PosEnd   Position
}

func (in ForNode) String() string {
  return ""
}

func (in *ForNode) GetPosStart() Position {
	return in.PosStart
}

func (in *ForNode) GetPosEnd() Position {
	return in.PosEnd
}

type WhileNode struct {
	Cond     Node
  BodyNode Node
	PosStart Position
	PosEnd   Position
}

func (in WhileNode) String() string {
  return ""
}

func (in *WhileNode) GetPosStart() Position {
	return in.PosStart
}

func (in *WhileNode) GetPosEnd() Position {
	return in.PosEnd
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

