package main

import (
	"fmt"
	"strconv"
	"strings"
)

const digits = "1234567890"
const ignores = " "
const (
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

// If the previous element in expr is a Real number (or a closed bracket, but
// this is technically the same thing), add val
func previousIsReal(expr Expression, s Symbol) Expression {
	if len(expr) > 0 {
		_, r := expr[len(expr)-1].(Real)
		c := expr[len(expr)-1] == Close
		if r || c {
			expr = append(expr, s)
		}
	}
	return expr
}

// ToExpression a string to a mathematical expression
func ToExpression(s string) (Expression, error) {
	var expr Expression
	var magnitude strings.Builder
	s = s + " "
	for _, c := range s {
		if contains(c, digits) {
			magnitude.WriteRune(c)
		} else {
			// add the add operator between every real number
			if magnitude.Len() > 0 {
				expr = previousIsReal(expr, Add)
				m, err := strconv.Atoi(magnitude.String())
				if err != nil {
					return nil, err
				}
				expr = append(expr, Real{true, m})
				magnitude.Reset()
			}
			if c == '-' {
				expr = previousIsReal(expr, Add)
				expr = append(expr, Real{false, 1}, Mul)
			}
			if c == '+' {
				expr = append(expr, Add)
			}
			if c == '*' {
				expr = append(expr, Mul)
			}
			if c == '/' {
				expr = append(expr, Div)
			}
			if c == '(' {
				expr = previousIsReal(expr, Mul)
				expr = append(expr, Open)
			}
			if c == ')' {
				expr = append(expr, Close)
			}
		}
	}
	return expr, nil
}

func main() {
}
