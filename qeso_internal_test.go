package main

import (
	"reflect"
	"testing"
)

func TestToExpression(t *testing.T) {
	type r struct {
		expr Expression
		err  error
	}
	var argresult = map[string]r{
		"1+2": r{
			expr: Expression{Real{true, 1}, Add, Real{true, 2}},
		},
		"1+2+3": r{
			expr: Expression{Real{true, 1}, Add, Real{true, 2}, Add, Real{true, 3}},
		},
		"1+2-3": r{
			expr: Expression{Real{true, 1}, Add, Real{true, 2}, Add, Real{false, 3}},
		},
		"1 + 2  - 3": r{
			expr: Expression{Real{true, 1}, Add, Real{true, 2}, Add, Real{false, 3}},
		},
		"1 - - 2": r{
			expr: Expression{Real{true, 1}, Add, Real{true, 2}},
		},
		"1 - - - 2": r{
			expr: Expression{Real{true, 1}, Add, Real{false, 2}},
		},
		"12 + 43": r{
			expr: Expression{Real{true, 12}, Add, Real{true, 43}},
		},
	}
	for arg, expected := range argresult {
		expr, err := ToExpression(arg)
		if !reflect.DeepEqual(expr, expected.expr) || err != expected.err {
			t.Errorf("Different result/err for %v: should have (%v, %v), got (%v, %v)",
				arg, expected.expr, expected.err, expr, err)
		}
	}
}
