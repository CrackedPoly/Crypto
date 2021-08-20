package seccommu

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
)

func Rand16() []byte {
	var r *big.Int
	r, _ = rand.Prime(rand.Reader, 128)
	tmp := make([]byte, 16)
	return r.FillBytes(tmp)
}

func ParsePubKey(cerfile string) *rsa.PublicKey {
	bytes, _ := ioutil.ReadFile(cerfile)
	cer, err := x509.ParseCertificate(bytes)
	if err != nil {
		log.Fatal(err)
	}
	pubKey := cer.PublicKey.(*rsa.PublicKey)
	return pubKey
}

func ParsePriKey(pvkfile string) *rsa.PrivateKey {
	cmd := exec.Command("openssl", "rsa", "-inform", "pvk", "-in", pvkfile, "-outform", "pem", "-out", pvkfile+".pem")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()
	//out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := ioutil.ReadFile(pvkfile + ".pem")
	block, _ := pem.Decode(bytes)
	if block == nil || block.Type != "rsa_impl PRIVATE KEY" {
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

func EncryptRSA(pubKey *rsa.PublicKey, plain []byte, filename string) []byte {
	cipher, _ := rsa.EncryptPKCS1v15(rand.Reader, pubKey, plain)
	// 写入文件
	_ = ioutil.WriteFile(filename, cipher, 0777)
	return cipher
}

func DecryptRSA(priKey *rsa.PrivateKey, cipher []byte) []byte {
	plain, _ := rsa.DecryptPKCS1v15(rand.Reader, priKey, cipher)
	return plain
}

func Hash(data []byte) [32]byte {
	return sha256.Sum256(data)
}

func Sign(priKey *rsa.PrivateKey, data [32]byte, filename string) []byte {
	signature, _ := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, data[:])
	// 写入文件
	_ = ioutil.WriteFile(filename, signature, 0777)
	return signature
}

func Verify(pubKey *rsa.PublicKey, hashed [32]byte, sig []byte) error {
	err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], sig)
	if err != nil {
		return err
	}
	return nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func EncryptFile(key []byte, plaintext []byte, outFilename string) (string, error) {
	of, err := os.Create(outFilename)
	if err != nil {
		return "", err
	}
	defer of.Close()

	iv := make([]byte, aes.BlockSize)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	padded := PKCS5Padding(plaintext, mode.BlockSize())
	ciphertext := make([]byte, len(padded))

	mode.CryptBlocks(ciphertext, padded)

	if _, err = of.Write(ciphertext); err != nil {
		return "", err
	}
	return outFilename, nil
}

func DecryptFile(key []byte, ciphertext []byte, outFilename string) (string, error) {
	of, err := os.Create(outFilename)
	if err != nil {
		return "", err
	}
	defer of.Close()

	iv := make([]byte, aes.BlockSize)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	padded := make([]byte, len(ciphertext))
	mode.CryptBlocks(padded, ciphertext)
	plaintext := PKCS5Unpadding(padded)

	if _, err := of.Write(plaintext); err != nil {
		return "", err
	}
	return outFilename, nil
}

func Equal(f1 string, f2 string) (bool, error) {
	file1, err := ioutil.ReadFile(f1)
	if err != nil {
		return false, err
	}
	file2, err := ioutil.ReadFile(f2)
	if err != nil {
		return false, err
	}
	return cmp.Equal(Hash(file1), Hash(file2)), nil
}
