// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package polynomial

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

type Op uint8

const (
	ADD Op = iota
	SUB
	MUL
)

// probabilisticCheck checks if c == Op(a, b) using Scwhartz Zippel
func probabilisticCheck(a, b, c Polynomial, op Op) bool {

	var r fr.Element
	r.SetRandom()
	ar := a.Eval(&r)
	br := b.Eval(&r)
	cr := c.Eval(&r)

	switch op {
	case ADD:
		ar.Add(&ar, &br)
		return ar.Equal(&cr)
	case SUB:
		ar.Sub(&ar, &br)
		return ar.Equal(&cr)
	case MUL:
		ar.Mul(&ar, &br)
		return ar.Equal(&cr)
	default:
		panic("operation not supported")
	}
}

func randomPolynomial(size int) Polynomial {
	res := NewPolynomial(uint64(size))
	for i := 0; i < size; i++ {
		res[i].SetRandom()
	}
	return res
}

func TestPolynomialOperands(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10

	properties := gopter.NewProperties(parameters)

	properties.Property("(Div by x-a) check operands which are not receivers are not modified", prop.ForAll(
		func() bool {
			a := randomPolynomial(4)
			b := randomPolynomial(8)
			bc := b.Copy()
			var r fr.Element
			r.SetRandom()
			a.DividePolyByXminusA(&b, &r)
			return bc.Equal(&b)
		},
	))

	properties.Property("(ADD) check operands which are not receivers are not modified", prop.ForAll(
		func() bool {
			res := true
			{
				a := randomPolynomial(4)
				b := randomPolynomial(8)
				ac := a.Copy()
				bc := b.Copy()
				var c Polynomial
				c.Add(&a, &b)
				res = res && ac.Equal(&a) && bc.Equal(&b)
			}
			{
				a := randomPolynomial(4)
				b := randomPolynomial(8)
				ac := a.Copy()
				b.Add(&a, &b)
				res = res && ac.Equal(&a)
			}
			{
				a := randomPolynomial(4)
				b := randomPolynomial(8)
				bc := b.Copy()
				a.Add(&a, &b)
				res = res && bc.Equal(&b)
			}
			return res
		},
	))

	properties.Property("(ADD) check operands which are not receivers are not modified", prop.ForAll(
		func() bool {
			res := true
			{
				a := randomPolynomial(4)
				b := randomPolynomial(8)
				ac := a.Copy()
				bc := b.Copy()
				var c Polynomial
				c.Sub(&a, &b)
				res = res && ac.Equal(&a) && bc.Equal(&b)
			}
			{
				a := randomPolynomial(4)
				b := randomPolynomial(8)
				ac := a.Copy()
				b.Sub(&a, &b)
				res = res && ac.Equal(&a)
			}
			{
				a := randomPolynomial(4)
				b := randomPolynomial(8)
				bc := b.Copy()
				a.Sub(&a, &b)
				res = res && bc.Equal(&b)
			}
			return res
		},
	))

	properties.Property("(MUL) check operands which are not receivers are not modified", prop.ForAll(
		func() bool {
			res := true
			{
				a := randomPolynomial(4)
				b := randomPolynomial(8)
				ac := a.Copy()
				bc := b.Copy()
				var c Polynomial
				c.Mul(&a, &b)
				res = res && ac.Equal(&a) && bc.Equal(&b)
			}
			{
				a := randomPolynomial(4)
				b := randomPolynomial(8)
				ac := a.Copy()
				b.Mul(&a, &b)
				res = res && ac.Equal(&a)
			}
			{
				a := randomPolynomial(4)
				b := randomPolynomial(8)
				bc := b.Copy()
				a.Mul(&a, &b)
				res = res && bc.Equal(&b)
			}
			return res
		},
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

func TestPolynomialOps(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10

	properties := gopter.NewProperties(parameters)

	// size of polynomials [a,b,c] -> a = op(b, c)
	configs := [][3]int{
		{4, 8, 16},
		{16, 16, 16},
		{16, 8, 4},
	}

	properties.Property("p(a)+c=(p+c)(a)", prop.ForAll(
		func() bool {
			res := true
			for _, conf := range configs {
				a := randomPolynomial(conf[0])
				b := randomPolynomial(1)
				c := randomPolynomial(conf[2])
				c.AddConstant(&a, b[0])
				res = res && probabilisticCheck(a, b, c, ADD)
			}
			return res
		},
	))

	properties.Property("p(a)-c=(p-c)(a)", prop.ForAll(
		func() bool {
			res := true
			for _, conf := range configs {
				a := randomPolynomial(conf[0])
				b := randomPolynomial(1)
				c := randomPolynomial(conf[2])
				c.SubConstant(&a, b[0])
				res = res && probabilisticCheck(a, b, c, SUB)
			}
			return res
		},
	))

	properties.Property("p(a)+q(a)=(p+q)(a)", prop.ForAll(
		func() bool {
			res := true
			for _, conf := range configs {
				a := randomPolynomial(conf[0])
				b := randomPolynomial(conf[1])
				c := NewPolynomial(uint64(conf[2]))
				c.Add(&a, &b)
				res = res && probabilisticCheck(a, b, c, ADD)
			}
			return res
		},
	))

	properties.Property("p(a)-q(a)=(p-q)(a)", prop.ForAll(
		func() bool {
			res := true
			for _, conf := range configs {
				a := randomPolynomial(conf[0])
				b := randomPolynomial(conf[1])
				c := NewPolynomial(uint64(conf[2]))
				c.Sub(&a, &b)
				res = res && probabilisticCheck(a, b, c, SUB)
			}
			return res
		},
	))

	properties.Property("p(a)*q(a)=(p*q)(a)", prop.ForAll(
		func() bool {
			res := true
			for _, conf := range configs {
				a := randomPolynomial(conf[0])
				b := randomPolynomial(conf[1])
				c := NewPolynomial(uint64(conf[2]))
				c.Mul(&a, &b)
				res = res && probabilisticCheck(a, b, c, MUL)
			}
			return res
		},
	))

	properties.Property("(x-a)*q(a)=p(a)", prop.ForAll(
		func() bool {
			res := true
			for _, conf := range configs {

				p := randomPolynomial(conf[0])
				q := NewPolynomial(uint64(conf[2]))

				var r fr.Element
				r.SetRandom()
				pr := p.Eval(&r)
				p.SubConstant(&p, pr)

				q.DividePolyByXminusA(&p, &pr)

				xsubr := NewPolynomial(2)
				xsubr[0].Neg(&r)
				xsubr[1].SetOne()
				res = res && probabilisticCheck(xsubr, q, p, MUL)

			}

			p := randomPolynomial(16)
			var q Polynomial

			var r fr.Element
			r.SetRandom()
			pr := p.Eval(&r)
			p.SubConstant(&p, pr)

			q.DividePolyByXminusA(&p, &pr)

			xsubr := NewPolynomial(2)
			xsubr[0].Neg(&r)
			xsubr[1].SetOne()
			res = res && probabilisticCheck(xsubr, q, p, MUL)

			return res
		},
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

func TestPolynomialEval(t *testing.T) {

	// build polynomial
	f := make(Polynomial, 20)
	for i := 0; i < 20; i++ {
		f[i].SetOne()
	}

	// random value
	var point fr.Element
	point.SetRandom()

	// compute manually f(val)
	var expectedEval, one, den fr.Element
	var expo big.Int
	one.SetOne()
	expo.SetUint64(20)
	expectedEval.Exp(point, &expo).
		Sub(&expectedEval, &one)
	den.Sub(&point, &one)
	expectedEval.Div(&expectedEval, &den)

	// compute purported evaluation
	purportedEval := f.Eval(&point)

	// check
	if !purportedEval.Equal(&expectedEval) {
		t.Fatal("polynomial evaluation failed")
	}
}
