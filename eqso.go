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
	// Multiply operator
	Multiply = Symbol(iota)
	// Divide operator
	Divide = Symbol(iota)
)

// Symbol represents a special symbol: add, multiply, divide, and brackets
type Symbol int

func (o Symbol) String() string {
	if o == Add {
		return "{add}"
	} else if o == Multiply {
		return "{mul}"
	} else if o == Divide {
		return "{div}"
	}
	return fmt.Sprintf("{unknown: %d}", o)
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

// ToExpression a string to a mathematical expression
func ToExpression(s string) (Expression, error) {
	var expr Expression
	var magnitude strings.Builder
	var positive = true
	s = s + " "
	for _, c := range s {
		if contains(c, digits) {
			magnitude.WriteRune(c)
		} else {
			if magnitude.Len() > 0 {
				// add the add operator between every real number
				if len(expr) > 0 {
					if _, ok := expr[len(expr)-1].(Real); ok {
						expr = append(expr, Add)
					}
				}
				m, err := strconv.Atoi(magnitude.String())
				if err != nil {
					return nil, err
				}
				expr = append(expr, Real{positive, m})
				magnitude.Reset()
				positive = true
			}
			if c == '-' {
				positive = !positive
			}
		}
	}
	return expr, nil
}

func main() {
}
