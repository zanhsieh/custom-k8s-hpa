//usr/bin/env go run $0 "$@"; exit
package regression

import (
	"os"
	"reflect"
	"testing"
)

type RoundInputVal struct {
	Val	float64
	RoundOn	float64
	Places	int
}

func TestRound(t *testing.T) {
	testCases := []struct {
		input	[]RoundInputVal
		success	bool
		expected float64
	}{
		{[]RoundInputVal{3.556, .5, 2}, true, 3.6},
	}
	for i := range testCases {
		tc := &testCases[i]
		arr := tc.input[0]
		var args []reflect.Value
		for _, x := range arr {
			args = append(args, reflect.ValueOf(x))
		}
		fun := reflect.ValueOf(regression.Round)
		result := fun.Call(args)
		resultVal := result[0].Interface().(float64)
		if resultVal != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, resultVal)
		}
	}
}

//func TestSolve(t *testing.T) {}
//func TestVandermonde(t *testing.T) {}
//func TestPredict(t *testing.T) {}
