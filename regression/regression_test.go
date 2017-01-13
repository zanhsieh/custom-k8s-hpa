//usr/bin/env go run $0 "$@"; exit
package regression

import (
	"math"
	"reflect"
	"testing"

	"github.com/gonum/matrix/mat64"
)

func TestRound(t *testing.T) {
	testCases := []struct {
		input    []interface{}
		success  bool
		expected float64
	}{
		{[]interface{}{3.556, .5, 2}, true, 3.56},
	}
	for i := range testCases {
		tc := &testCases[i]
		arr := tc.input
		var args []reflect.Value
		for _, x := range arr {
			args = append(args, reflect.ValueOf(x))
		}
		fun := reflect.ValueOf(Round)
		result := fun.Call(args)
		resultVal := result[0].Interface().(float64)
		if resultVal != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, resultVal)
		}
	}
}

func TestSolve(t *testing.T) {
	// classical case for polynomial regression
	// see: https://rosettacode.org/wiki/Polynomial_regression#Library_gonum.2Fmatrix
	var (
		x        = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		y        = []float64{1, 6, 17, 34, 57, 86, 121, 162, 209, 262, 321}
		expected = []float64{1, 2, 3}
		degree   = 2
	)

	c := Solve(x, y, degree)
	e := mat64.NewDense(len(expected), 1, expected)
	if !mat64.EqualApprox(c, e, 0.5) {
		t.Errorf("expected \n%v\n got \n%v\n", mat64.Formatted(e), mat64.Formatted(c))
	}
}

func TestVandermonde(t *testing.T) {
	var (
		x      = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		c0     = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		c1     = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		c2     = []float64{0, 1, 4, 9, 16, 25, 36, 49, 64, 81, 100}
		degree = 2
	)
	a := Vandermonde(x, degree)
	e := mat64.NewDense(len(x), degree+1, nil)
	e.SetCol(0, c0)
	e.SetCol(1, c1)
	e.SetCol(2, c2)
	if !mat64.EqualApprox(a, e, 0.5) {
		t.Errorf("expected \n%v\n got \n%v\n", mat64.Formatted(e), mat64.Formatted(a))
	}
}

func TestPredict(t *testing.T) {
	var (
		thisP  = 11.0
		cf     = []float64{1, 2, 3}
		degree = 2
	)
	c := mat64.NewDense(len(cf), 1, cf)
	result := Predict(thisP, c, degree)
	expected := 0.0
	for i := range cf {
		expected = expected + cf[i]*math.Pow(thisP, float64(i))
	}
	if result != expected {
		t.Errorf("expected \n%v\n got \n%v\n", expected, result)
	}
}
