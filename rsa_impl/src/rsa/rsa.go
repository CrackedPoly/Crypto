package rsa

import (
	"crypto/rand"
	"math/big"
)

type RSA struct {
	P *big.Int
	Q *big.Int
	N *big.Int
	E *big.Int
	D *big.Int
}

func NewRSA(keyLen int) RSA {
	operand := big.NewInt(0)
	one := big.NewInt(1)
	var p, q, e *big.Int
	n := new(big.Int)
	fn := new(big.Int)
	d := new(big.Int)
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
	n = n.Mul(p, q)
	fn = fn.Mul(operand.Sub(p, one), operand.Sub(q, one))

	for {
		e, err = rand.Int(rand.Reader, fn)
		if err == nil && one.Cmp(operand.GCD(nil, nil, e, fn)) == 0 {
			break
		}
	}
	d = d.ModInverse(e, fn)
	return RSA{
		P: p,
		Q: q,
		N: n,
		E: e,
		D: d,
	}
}

func NewSign(nString string, dString string) RSA {
	n := new(big.Int)
	d := new(big.Int)
	n.SetString(nString, 16)
	d.SetString(dString, 16)
	return RSA{
		P: nil,
		Q: nil,
		N: n,
		E: nil,
		D: d,
	}
}

func NewCheck(nString string, eString string) RSA {
	n := new(big.Int)
	e := new(big.Int)
	n.SetString(nString, 16)
	e.SetString(eString, 16)
	return RSA{
		P: nil,
		Q: nil,
		N: n,
		E: e,
		D: nil,
	}
}

func (r *RSA) Encrypt(plaintext string) string {
	p := new(big.Int)
	p.SetString(plaintext, 16)
	p.Exp(p, r.E, r.N)
	cipher := p.Text(16)
	return cipher
}

func (r *RSA) Decrypt(cipher string) string {
	p := new(big.Int)
	p.SetString(cipher, 16)
	p.Exp(p, r.D, r.N)
	plaintext := p.Text(16)
	return plaintext
}
