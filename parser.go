package main

import (
	"fmt"
)

var errNotFound = fmt.Errorf("function never returned true")

// indexof runs fn on every element of the slice after 'after', and if it
// returns true, it returns this (index, element)
func indexof(expr Expression, fn func(int, interface{}) (bool, error), after int) (int, interface{}, error) {
	for i, e := range expr {
		ok, err := fn(i, e)
		if err != nil {
			return 0, nil, err
		}
		if ok {
			return i, e, nil
		}
	}
	return 0, nil, errNotFound
}

// // A sub parse of the parse function
// func parsebrackets(expr Expression) (Expression, error) {

// }

// It returns an expression containing one Node
// The reason being that it's recursive (so, it calls itself with expression
// with multiple Nodes/unparsed tokens)
func parse(expr Expression) (Expression, error) {
	// look for brackets
	i, _, err := indexof(expr, func(i int, e interface{}) (bool, error) {
		return e == Open, nil
	}, 0)
	if err == errNotFound {
		// look for * or /
		i, _, err := indexof(expr, func(i int, e interface{}) (bool, error) {
			return e == Mul || e == Div, nil
		}, 0)
		if err == errNotFound {
			// look for +
			i, _, err := indexof(expr, func(i int, e interface{}) (bool, error) {
				return e == Add, nil
			}, 0)
			if err == errNotFound {
				if len(expr) != 1 {
					return nil, fmt.Errorf("invalid expression after parsing %v", expr)
				}
				switch expr[0].(type) {
				case Real:
					return Expression{&Node{expr[0], nil, Null}}, nil
				case *Node:
					return expr, nil
				}
				return nil, fmt.Errorf("Got 'empty' expression of %d elements: %v", len(expr), expr)
			} else if err != nil {
				return nil, err
			}
			expr = append(
				append(expr[:i-1], &Node{expr[i-1], expr[i+1], expr[i].(Symbol)}),
				expr[i+2:]...)
			return parse(expr)
		} else if err != nil {
			return nil, err
		}
		// the operands and then the operator
		expr = append(append(expr[:i-1], &Node{expr[i-1], expr[i+1], expr[i].(Symbol)}), expr[i+2:]...)
		return parse(expr)
	} else if err != nil {
		return nil, err
	}
	j, _, err := indexof(expr, func(i int, e interface{}) (bool, error) {
		return e == Close, nil
	}, i)
	if err != nil {
		return nil, err
	}
	sub, err := parse(expr[i+1 : j])
	if err != nil {
		return nil, err
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
