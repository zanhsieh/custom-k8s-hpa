//usr/bin/env go run $0 "$@"; exit
package regression

import (
	//"reflect"
	"testing"

	"github.com/gonum/matrix/mat64"
)

func TestRound(t *testing.T) {
	testCases := []struct {
		val	float64
		roundOn	float64
		places	int
		success	bool
		expected float64
	}{
		{3.556, .5, 2, true, 3.56},
	}
	for i := range testCases {
		tc := &testCases[i]
		resultVal := Round(tc.val, tc.roundOn, tc.places)
		if resultVal != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, resultVal)
		}
	}
}

func TestSolve(t *testing.T) {
	// classical case for polynomial regression
	// see: https://rosettacode.org/wiki/Polynomial_regression#Library_gonum.2Fmatrix
	var (
		x = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		y = []float64{1, 6, 17, 34, 57, 86, 121, 162, 209, 262, 321}
		expected = []float64{1, 2, 3}
		degree = 2
	)

	c := Solve(x, y, degree)
	e := mat64.NewDense(len(expected), 1, expected)
	if !mat64.EqualApprox(c, e, 0.5) {
		t.Errorf("expected \n%v\n got \n%v\n", mat64.Formatted(e), mat64.Formatted(c))
	}
}
//func TestVandermonde(t *testing.T) {}
//func TestPredict(t *testing.T) {}
