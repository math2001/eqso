package main

import (
	"fmt"
	"log"
	"strings"
)

const (
	digits = "1234567890"
	// Null is no operator
	Null = Symbol(iota)
	// Add operator
	Add = Symbol(iota)
	// Mul is the multiply operator
	Mul = Symbol(iota)
	// Div is the divide operator
	Div = Symbol(iota)
	// Exp is the exponential operator
	Exp = Symbol(iota)
	// Open bracket
	Open = Symbol(iota)
	// Close bracket
	Close = Symbol(iota)
)

// Symbol represents a special symbol: add, multiply, divide, and brackets
type Symbol int

func (o Symbol) String() string {
	switch o {
	case Null:
		return "{null}"
	case Add:
		return "{add}"
	case Mul:
		return "{mul}"
	case Div:
		return "{div}"
	case Exp:
		return "{exp}"
	case Open:
		return "("
	case Close:
		return ")"
	default:
		return fmt.Sprintf("{unknown: %d}", o)
	}
}

// Node is a number in the expression. 4 is considered to be a number, just as
// is the whole (1 - 4) for example
type Node struct {
	A, B     interface{} // either int or *Node
	Operator Symbol
}

func (n Node) String() string {
	return fmt.Sprintf("Node{%v %s %v}", n.A, n.Operator, n.B)
}

// Eval evaluates the node's value
func (n *Node) Eval() (int, error) {
	var a, b int
	var err error
	if node, isnode := n.A.(*Node); isnode {
		a, err = node.Eval()
		if err != nil {
			return 0, err
		}
	} else {
		r := n.A.(int)
		a = r
	}
	if node, isnode := n.B.(*Node); isnode {
		b, err = node.Eval()
		if err != nil {
			return 0, err
		}
	} else {
		r := n.B.(int)
		b = r
	}

	if n.Operator == Add {
		return a + b, nil
	} else if n.Operator == Mul {
		return a * b, nil
	} else if n.Operator == Div {
		return a / b, nil
	} else if n.Operator == Null {
		return a, nil
	}
	return 0, fmt.Errorf("Invalid operator %v", n.Operator)
}

// Expression represents a mathematical expression. There are only 2 types of
// operators: + and *, plus brackets
type Expression []interface{}

func contains(str rune, all string) bool {
	for _, s := range all {
		if str == s {
			return true
		}
	}
	return false
}

func main() {
	// expr, err := ToExpression("(1+2)*3")
	expr, err := Tokenize(strings.NewReader("(1+2*3)*4"))
	if err != nil {
		log.Fatal(err)
	}
	tree, err := Parse(expr)
	if err != nil {
		log.Fatal(err)
	}
	res, err := tree.Eval()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Result:", res)
}
