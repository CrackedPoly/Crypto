package seccommu

import (
	"crypto"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
)

func rand128() []byte {
	var r *big.Int
	r, _ = rand.Int(rand.Reader, big.NewInt(128))
	tmp := make([]byte, 128)
	return r.FillBytes(tmp)
}

func parsePubKey(cerfile string) *rsa.PublicKey {
	bytes, _ := ioutil.ReadFile(cerfile)
	cer, err := x509.ParseCertificate(bytes)
	if err != nil {
		log.Fatal(err)
	}
	pubKey := cer.PublicKey.(*rsa.PublicKey)
	return pubKey
}

func parsePriKey(pvkfile string) *rsa.PrivateKey {
	cmd := exec.Command("openssl", "rsa", "-inform", "pvk", "-in", pvkfile, "-outform", "pem", "-out", pvkfile+".pem")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()
	//out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := ioutil.ReadFile(pvkfile + ".pem")
	block, _ := pem.Decode(bytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		log.Fatal("failed to decode PEM block containing private key")
	}
	priKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	_ = os.Remove(pvkfile + ".pem")
	return priKey
}

func NewCer(childCer string, childKey string, rootCer string, rootKey string, cn string) error {
	cmd := exec.Command("./tools/makecert.exe", "-n", "CN="+cn, "-iv", rootKey, "-ic", rootCer, "-sv", childKey, childCer)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	err := cmd.Run()
	//out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func EncryptAES(key []byte, plain []byte) []byte {
	cipher := make([]byte, len(plain))
	block, _ := aes.NewCipher(key)
	block.Encrypt(cipher, plain)
	// 写入文件
	return cipher
}

func DecryptAES(key []byte, cipher []byte) []byte {
	plain := make([]byte, len(cipher))
	block, _ := aes.NewCipher(key)
	block.Decrypt(plain, cipher)
	return plain
}

func EncryptRSA(pubKey *rsa.PublicKey, plain []byte) []byte {
	cipher, _ := rsa.EncryptPKCS1v15(rand.Reader, pubKey, plain)
	// 写入文件
	return cipher
}

func DecryptRSA(priKey *rsa.PrivateKey, cipher []byte) []byte {
	plain, _ := rsa.DecryptPKCS1v15(rand.Reader, priKey, cipher)
	return plain
}

func Hash(data []byte) [32]byte {
	return sha256.Sum256(data)
}

func Sign(priKey *rsa.PrivateKey, data []byte) []byte {
	signature, _ := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, data)
	return signature
	// 写入文件
}
