package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	errMissingClosing   = fmt.Errorf("missing closing bracket")
	errUnmatchedClosing = fmt.Errorf("extra closing bracket found")
)

func addIfLastIsTerm(expr Expression, s Symbol) Expression {
	if len(expr) == 0 {
		return expr
	}
	_, isterm := expr[len(expr)-1].(int)
	// a closing bracket is considered to be a term, since what is inside the
	// bracket *actually* is a term
	isclose := expr[len(expr)-1] == Close
	if isterm || isclose {
		return append(expr, s)
	}
	return expr
}

// Tokenize a string to a list of token (Symbols or integers). It doesn't know
// anything about the validity of the expression
func Tokenize(reader io.Reader) (Expression, error) {
	var (
		expr      Expression
		magnitude strings.Builder
		positive  = true
		r         = bufio.NewReader(reader)
		stop      = false
		brackets  = 0 // brackets counter, to make sure that they are valid
	)
	for {
		if stop {
			break
		}
		ru, _, err := r.ReadRune()
		if err == io.EOF {
			// we don't break now, because we want to last number of symbol to
			// be added to the expression
			stop = true
		} else if err != nil {
			return nil, fmt.Errorf("couldn't read rune: %s", err)
		}
		if contains(ru, digits) {
			magnitude.WriteRune(ru)
			continue
		}
		if magnitude.Len() > 0 {
			// here, we have a rune different than a digit, therefore we must
			// add the digits that we have stored in magnitude
			// But before, we need to add a + (adding a negative number is the
			// same thing as sustracting a positive one)
			expr = addIfLastIsTerm(expr, Add)
			m, err := strconv.Atoi(magnitude.String())
			if err != nil {
				return nil, fmt.Errorf("couldn't convert magnitude to string: %s", err)
			}
			if positive {
				expr = append(expr, m)
			} else {
				expr = append(expr, -m)
			}

			magnitude.Reset()
			positive = true
		}
		if ru == '^' {
			expr = append(expr, Exp)
		}
		if ru == '-' {
			expr = addIfLastIsTerm(expr, Add)
			positive = false
		}
		if ru == '+' {
			expr = append(expr, Add)
		}
		if ru == '(' {
			if !positive {
				// here, we have something like '... - ( ...'
				// which is '... -1 * ( ...'
				expr = append(expr, -1)
			}
			expr = addIfLastIsTerm(expr, Mul)
			expr = append(expr, Open)
			positive = true // reset to default
			brackets++
		}
		if ru == ')' {
			expr = append(expr, Close)
			brackets--
			if brackets < 0 {
				return nil, errUnmatchedClosing
			}
		}
		if ru == '/' {
			expr = append(expr, Div)
		}
		if ru == '*' {
			expr = append(expr, Mul)
		}
	}
	if brackets > 0 {
		return nil, errMissingClosing
	}
	return expr, nil
}
