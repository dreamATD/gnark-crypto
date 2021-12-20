package fptower

func (z *E12) nSquare(n int) {
	for i := 0; i < n; i++ {
		z.CyclotomicSquare(z)
	}
}

func (z *E12) nSquareCompressed(n int) {
	for i := 0; i < n; i++ {
		z.CyclotomicSquareCompressed(z)
	}
}

// Expt set z to x^t in E12 and return z
func (z *E12) Expt(x *E12) *E12 {

	// Expt computation is derived from the addition chain:
	//
	//	_1000     = 1 << 3
	//	_1001     = 1 + _1000
	//	_1001000  = _1001 << 3
	//	_1010001  = _1001 + _1001000
	//	_10011001 = _1001000 + _1010001
	//	i67       = ((_10011001 << 5 + _1001) << 10 + _1010001) << 41
	//	return      1 + i67
	//
	// Operations: 62 squares 6 multiplies
	//
	// Generated by github.com/mmcloughlin/addchain v0.4.0.

	// Allocate Temporaries.
	var result, t0, t1 E12

	// Step 3: result = x^0x8
	result.CyclotomicSquare(x)
	result.nSquare(2)

	// Step 4: t0 = x^0x9
	t0.Mul(x, &result)

	// Step 7: t1 = x^0x48
	t1.CyclotomicSquare(&t0)
	t1.nSquare(2)

	// Step 8: result = x^0x51
	result.Mul(&t0, &t1)

	// Step 9: t1 = x^0x99
	t1.Mul(&t1, &result)

	// Step 14: t1 = x^0x1320
	t1.nSquare(5)

	// Step 15: t0 = x^0x1329
	t0.Mul(&t0, &t1)

	// Step 25: t0 = x^0x4ca400
	t0.nSquare(10)

	// Step 26: result = x^0x4ca451
	result.Mul(&result, &t0)

	// Step 67: result = x^0x9948a20000000000
	result.nSquareCompressed(41)
	result.Decompress(&result)

	// Step 68: result = x^0x9948a20000000001
	z.Mul(x, &result)

	return z
}

// MulBy014 multiplication by sparse element (c0, c1, 0, 0, c4)
func (z *E12) MulBy014(c0, c1, c4 *E2) *E12 {

	var a, b E6
	var d E2

	a.Set(&z.C0)
	a.MulBy01(c0, c1)

	b.Set(&z.C1)
	b.MulBy1(c4)
	d.Add(c1, c4)

	z.C1.Add(&z.C1, &z.C0)
	z.C1.MulBy01(c0, &d)
	z.C1.Sub(&z.C1, &a)
	z.C1.Sub(&z.C1, &b)
	z.C0.MulByNonResidue(&b)
	z.C0.Add(&z.C0, &a)

	return z
}

// Mul014By014 multiplication of sparse element (c0,c1,0,0,c4,0) by sparse element (d0,d1,0,0,d4,0)
func (z *E12) Mul014By014(d0, d1, d4, c0, c1, c4 *E2) *E12 {
	var tmp, x0, x1, x4, x04, x01, x14 E2
	x0.Mul(c0, d0)
	x1.Mul(c1, d1)
	x4.Mul(c4, d4)
	tmp.Add(c0, c4)
	x04.Add(d0, d4).
		Mul(&x04, &tmp).
		Sub(&x04, &x0).
		Sub(&x04, &x4)
	tmp.Add(c0, c1)
	x01.Add(d0, d1).
		Mul(&x01, &tmp).
		Sub(&x01, &x0).
		Sub(&x01, &x1)
	tmp.Add(c1, c4)
	x14.Add(d1, d4).
		Mul(&x14, &tmp).
		Sub(&x14, &x1).
		Sub(&x14, &x4)

	z.C0.B0.MulByNonResidue(&x4).
		Add(&z.C0.B0, &x0)
	z.C0.B1.Set(&x01)
	z.C0.B2.Set(&x1)
	z.C1.B0.SetZero()
	z.C1.B1.Set(&x04)
	z.C1.B2.Set(&x14)

	return z
}
