package main

import (
	"log"
	"reflect"
	"testing"
)

// factorial returns n! It doesn't support negative numbers (which would raise
// a Math Error)
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

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
		"1+2*3": r{
			tree: &Node{Real{true, 1}, &Node{Real{true, 2}, Real{true, 3}, Mul}, Add},
		},
		"(1+2)*3": r{
			tree: &Node{&Node{Real{true, 1}, Real{true, 2}, Add}, Real{true, 3}, Mul},
		},
		"(1+2)+3": r{
			tree: &Node{&Node{Real{true, 1}, Real{true, 2}, Add}, Real{true, 3}, Add},
		},
		"(10*(1+3))+1": r{
			tree: &Node{&Node{10, &Node{1, 3, Add}, Mul}, 1, Add},
		},
		"1+2": r{
			tree: &Node{Real{true, 1}, Real{true, 2}, Add},
		},
		"1": r{
			tree: &Node{Real{true, 1}, nil, Null},
		},
		"10+59*32/4": r{
			tree: &Node{R(10), &Node{&Node{R(59), R(32), Mul}, R(4), Div}, Add},
		},
		"(10)": r{
			tree: &Node{R(10), nil, Null},
		},
		"20(10 + 2)": r{
			tree: &Node{R(20), &Node{R(10), R(2), Add}, Mul},
		},
		"-20(10 + 2)*-3": r{
			tree: &Node{
				&Node{
					&Node{
						&Node{R(-1), R(20), Mul},
						&Node{R(10), R(2), Add},
						Mul,
					},
					R(-1),
					Mul,
				},
				R(3),
				Mul,
			},
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
func TestEval(t *testing.T) {
	type r struct {
		res int
		err error
	}
	var argresult = map[string]r{
		"1+2": r{
			res: 3,
		},
		"1+2+3+4+5+6+7+8+9+10":    r{55, nil},
		"1*2*3*4*5*6*7*8*9*10":    r{factorial(10), nil},
		"12*43+32*-35":            r{12*43 + 32*-35, nil},
		"(10+8)*28/6":             r{(10 + 8) * 28 / 6, nil},
		"(10*(22+4)-10/(4/2))+11": r{(10*(22+4) - 10/(4/2)) + 11, nil},
	}
	for arg, expected := range argresult {
		expr, err := ToExpression(arg)
		if err != nil {
			log.Fatalf("This shouldn't happen: %s", err)
		}
		tree, err := Parse(expr)
		if err != nil {
			log.Fatalf("This shouldn't happen either: %s", err)
		}
		res, err := tree.Eval()
		if res != expected.res || err != expected.err {
			t.Errorf("Different result/err for %v:\nshould have (%v, %v)\ngot         (%v, %v)",
				arg, expected.res, expected.err, res, err)
		}
	}
}
