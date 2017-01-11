//usr/bin/env go run $0 "$@"; exit
package regression

import (
	"fmt"
	"math"

	"github.com/gonum/matrix/mat64"
)

const (
	debug = false
)

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func Solve(x, y []float64, degree int) *mat64.Dense {

	a := Vandermonde(x, degree)
	if debug {
		fmt.Printf("a=>\n%v\n", mat64.Formatted(a))
	}

	b := mat64.NewDense(len(y), 1, y)
	if debug {
		fmt.Printf("b=>\n%v\n", mat64.Formatted(b))
	}

	c := mat64.NewDense(degree+1, 1, nil)

	qr := new(mat64.QR)
	qr.Factorize(a)

	err := c.SolveQR(qr, false, b)
	if err != nil {
		fmt.Println(err)
		return mat64.NewDense(degree+1, 1, nil)
	} else {
		return c
	}
}

func Vandermonde(a []float64, degree int) *mat64.Dense {
	x := mat64.NewDense(len(a), degree+1, nil)
	for i := range a {
		for j, p := 0, 1.; j <= degree; j, p = j+1, p*a[i] {
			x.Set(i, j, p)
		}
	}
	return x
}

func Predict(thisP float64, c *mat64.Dense, degree int) float64 {
	p := []float64{thisP}
	vp := Vandermonde(p, degree)
	if debug {
		fmt.Printf("vp=>\n%v\n", mat64.Formatted(vp))
	}

	result := mat64.NewDense(1, 1, nil)
	result.Mul(vp, c)
	return result.At(0, 0)
}
