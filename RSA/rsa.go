package RSA

import (
	"crypto/rand"
	"math/big"
)

type RSA struct {
	p      *big.Int
	q      *big.Int
	n      *big.Int
	fn     *big.Int
	e      *big.Int
	d      *big.Int
	keyLen int
}

func NewRSA(keyLen int) RSA {
	operand := big.NewInt(0)
	one := big.NewInt(1)
	var p, q, n, fn, e, d *big.Int
	var err error
	for {
		p, err = rand.Prime(rand.Reader, keyLen)
		if err == nil {
			break
		}
	}
	for {
		q, err = rand.Prime(rand.Reader, keyLen)
		if err == nil {
			break
		}
	}
	n = operand.Mul(p, q)
	fn = operand.Mul(operand.Sub(p, one), operand.Sub(q, one))

	for {
		e, err = rand.Int(rand.Reader, fn)
		if err == nil && one.Cmp(operand.GCD(nil, nil, e, fn)) == 0 {
			break
		}
	}
	d = operand.ModInverse(e, fn)
	return RSA{
		p:      p,
		q:      q,
		n:      n,
		fn:     fn,
		e:      e,
		d:      d,
		keyLen: keyLen,
	}
}

func NewSign(n *big.Int, e *big.Int, d *big.Int) RSA {
	return RSA{
		p:      nil,
		q:      nil,
		n:      n,
		fn:     nil,
		e:      e,
		d:      d,
		keyLen: nil,
	}
}

func (r *RSA) Encrypt(plaintext string) string {
	p := new(big.Int)
	p.SetString(plaintext, 16)
	p.Exp(p, r.e, r.n)
	return p.Text(16)
}
