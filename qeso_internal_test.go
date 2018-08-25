package main

import (
	"log"
	"reflect"
	"testing"
)

func makeexpr(values ...interface{}) Expression {
	var expr Expression
	for _, val := range values {
		switch val.(type) {
		case Symbol:
			expr = append(expr, val.(Symbol))
		case int:
			var (
				pos bool
				mag int
			)
			pos = val.(int) >= 0
			if pos {
				mag = val.(int)
			} else {
				mag = -val.(int)
			}
			expr = append(expr, Real{
				Positive:  pos,
				Magnitude: mag,
			})
		}
	}
	return expr
}

func TestToExpression(t *testing.T) {
	type r struct {
		expr Expression
		err  error
	}
	var argresult = map[string]r{
		"1+2": r{
			expr: makeexpr(1, Add, 2),
		},
		"1+2+3": r{
			expr: makeexpr(1, Add, 2, Add, 3),
		},
		"1+2-3": r{
			expr: makeexpr(1, Add, 2, Add, -1, Mul, 3),
		},
		"1234+4321": r{
			expr: makeexpr(1234, Add, 4321),
		},
		"1234-4321": r{
			expr: makeexpr(1234, Add, -1, Mul, 4321),
		},
		"10 - (5 - 3)": r{
			expr: makeexpr(10, Add, -1, Mul, Open, 5, Add, -1, Mul, 3, Close),
		},
		"10 + (5 - 3)": r{
			expr: makeexpr(10, Add, Open, 5, Add, -1, Mul, 3, Close),
		},
		"10(5 - 3)": r{
			expr: makeexpr(10, Mul, Open, 5, Add, -1, Mul, 3, Close),
		},
		"10(5 - 3)(4 + 3)": r{
			expr: makeexpr(10, Mul, Open, 5, Add, -1, Mul, 3, Close, Mul,
				Open, 4, Add, 3, Close),
		},
		"10 / 3": r{
			expr: makeexpr(10, Div, 3),
		},
		"10 / (1 + 2)": r{
			expr: makeexpr(10, Div, Open, 1, Add, 2, Close),
		},
		"10 / (1 + 2(4 + 2))": r{
			expr: makeexpr(10, Div, Open, 1, Add, 2, Mul, Open, 4, Add, 2,
				Close, Close),
		},
		"10 / -3": r{
			expr: makeexpr(10, Div, -1, Mul, 3),
		},
	}
	for arg, expected := range argresult {
		expr, err := ToExpression(arg)
		if !reflect.DeepEqual(expr, expected.expr) || err != expected.err {
			t.Errorf("Different result/err for %v:\nshould have (%v, %v)\ngot         (%v, %v)",
				arg, expected.expr, expected.err, expr, err)
		}
	}
}

func TestParser(t *testing.T) {
	type r struct {
		tree *Node
		err  error
	}
	var argresult = map[string]r{
		"(1+2)": r{
			tree: &Node{1, 2, Add},
		},
		"1+2*3": r{
			tree: &Node{1, &Node{2, 3, Mul}, Add},
		},
		"1+2": r{
			tree: &Node{1, 2, Add},
		},
	}
	for arg, expected := range argresult {
		expr, err := ToExpression(arg)
		if err != nil {
			log.Fatalf("This shouldn't happen: %s", err)
		}
		tree, err := Parse(expr)
		if !reflect.DeepEqual(tree, expected.tree) || err != expected.err {
			t.Errorf("Different result/err for %v:\nshould have (%v, %v)\ngot         (%v, %v)",
				arg, expected.tree, expected.err, tree, err)
		}
	}
}
