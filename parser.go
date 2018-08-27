package main

import (
	"fmt"
)

var errNotFound = fmt.Errorf("function never returned true")

func parseAdd(expr Expression) (Expression, error) {
	var (
		i     int
		found bool
		e     interface{}
	)
	for i, e = range expr {
		if e == Add {
			found = true
			break
		}
	}
	if !found {
		switch expr[0].(type) {
		case int:
			return Expression{&Node{expr[0], nil, Null}}, nil
		case *Node:
			return expr, nil
		}
		return nil, fmt.Errorf("invalid type for first element: %T", expr[0])
	}
	expr = append(
		append(expr[:i-1], &Node{expr[i-1], expr[i+1], expr[i].(Symbol)}),
		expr[i+2:]...)
	// same things as parseMulDiv
	return parseAdd(expr)
}

func parseMulDiv(expr Expression) (Expression, error) {
	var (
		i     int
		found bool
		e     interface{}
	)
	for i, e = range expr {
		if e == Mul || e == Div {
			found = true
			break
		}
	}
	if !found {
		return parseAdd(expr)
	}
	// we add the first operand, the second operand, and the symbol
	expr = append(
		append(expr[:i-1], &Node{expr[i-1], expr[i+1], expr[i].(Symbol)}),
		expr[i+2:]...)
	// since we got here, this means that the tasks higher up are done
	// (brackets), so they don't need to run again
	return parseMulDiv(expr)
}

func parseExponentials(expr Expression) (Expression, error) {
	var (
		i     int
		e     interface{}
		found bool
	)
	for i, e = range expr {
		if e == Exp {
			found = true
			break
		}
	}
	if !found {
		return parseMulDiv(expr)
	}
	expr = append(
		append(expr[:i-1], &Node{expr[i-1], expr[i+1], expr[i].(Symbol)}),
		expr[i+2:]...)
	return parseExponentials(expr)
}

// It returns an expression containing one Node
// The reason being that it's recursive (so, it calls itself with expression
// with multiple Nodes/unparsed tokens)
func parse(expr Expression) (Expression, error) {
	// look for brackets
	var i, j int
	var e interface{}
	var found = false
	for i, e = range expr {
		if e == Open {
			found = true
			break
		}
	}
	if !found { // we don't have any brackets
		return parseExponentials(expr)
	}
	var opencount = 0
	found = false
	for j, e = range expr[i:] {
		if e == Open {
			opencount++
		}
		if e == Close {
			opencount--
			if opencount == 0 {
				found = true
				break
			}
		}
	}
	j += i
	if !found {
		// this means that we have found an opening bracket, but no closing
		// note that this shouldn't happen as the bracket count is checked in
		// the tokenizer
		return nil, fmt.Errorf("no matching bracket in %s", expr[i:])
	}
	sub, err := parse(expr[i+1 : j])
	if err != nil {
		return nil, fmt.Errorf("sub parsing %v: %s", expr[i+1:j], err)
	}
	// replace every element in the brackets with sub
	// expr = append(append(expr[:i+1], sub), expr[j:]...)
	expr = append(expr[:i], append(sub, expr[j+1:]...)...)
	return parse(expr)
}

// Parse transforms an expression into a tree of nodes
func Parse(expr Expression) (*Node, error) {
	expr, err := parse(expr)
	if err != nil {
		return nil, err
	}
	node, ok := expr[0].(*Node)
	if !ok {
		return nil, fmt.Errorf("invalid expression result: should have one *Node, got %T", expr[0])
	}
	return node, nil
}
