package main

import (
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
			expr: makeexpr(10, Add, -1, Mul, Open, 5, -1, Mul, 3, Close),
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
