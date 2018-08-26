package main

import (
	"fmt"
	"log"
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
	case Open:
		return "("
	case Close:
		return ")"
	default:
		return fmt.Sprintf("{unknown: %d}", o)
	}
}

// Real represents any real number
type Real struct {
	Positive  bool
	Magnitude int
}

func (r Real) String() string {
	var p string
	if r.Positive {
		p = "+"
	} else {
		p = "-"
	}
	return fmt.Sprintf("%s%d", p, r.Magnitude)
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
	expr, err := ToExpression("(1+2*3)*4")
	if err != nil {
		log.Fatal(err)
	}
	tree, err := Parse(expr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Got tree:", tree)
}
