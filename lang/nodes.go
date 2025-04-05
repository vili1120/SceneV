package lang

import "fmt"

type Node interface {
  String() string
}

type NumberNode struct {
  Tok Token
}

func (nn NumberNode) String() string {
  return fmt.Sprintf("%v", nn.Tok.String())
}

type BinOpNode struct {
  LeftNode Node
  OpTok Token
  RightNode Node
}

func (bon BinOpNode) String() string {
  return fmt.Sprintf("(%v, %v, %v)", bon.LeftNode.String(), bon.OpTok.String(), bon.RightNode.String())
}

type UnaryOpNode struct {
  OpTok Token
  Node Node
}

func (uop UnaryOpNode) String() string {
  return fmt.Sprintf("(%v, %v)", uop.OpTok, uop.Node)
}
