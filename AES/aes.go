package AES

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var sbox = [256]byte{
	0x63, 0x7c, 0x77, 0x7b, 0xf2, 0x6b, 0x6f, 0xc5, 0x30, 0x01, 0x67, 0x2b, 0xfe, 0xd7, 0xab, 0x76,
	0xca, 0x82, 0xc9, 0x7d, 0xfa, 0x59, 0x47, 0xf0, 0xad, 0xd4, 0xa2, 0xaf, 0x9c, 0xa4, 0x72, 0xc0,
	0xb7, 0xfd, 0x93, 0x26, 0x36, 0x3f, 0xf7, 0xcc, 0x34, 0xa5, 0xe5, 0xf1, 0x71, 0xd8, 0x31, 0x15,
	0x04, 0xc7, 0x23, 0xc3, 0x18, 0x96, 0x05, 0x9a, 0x07, 0x12, 0x80, 0xe2, 0xeb, 0x27, 0xb2, 0x75,
	0x09, 0x83, 0x2c, 0x1a, 0x1b, 0x6e, 0x5a, 0xa0, 0x52, 0x3b, 0xd6, 0xb3, 0x29, 0xe3, 0x2f, 0x84,
	0x53, 0xd1, 0x00, 0xed, 0x20, 0xfc, 0xb1, 0x5b, 0x6a, 0xcb, 0xbe, 0x39, 0x4a, 0x4c, 0x58, 0xcf,
	0xd0, 0xef, 0xaa, 0xfb, 0x43, 0x4d, 0x33, 0x85, 0x45, 0xf9, 0x02, 0x7f, 0x50, 0x3c, 0x9f, 0xa8,
	0x51, 0xa3, 0x40, 0x8f, 0x92, 0x9d, 0x38, 0xf5, 0xbc, 0xb6, 0xda, 0x21, 0x10, 0xff, 0xf3, 0xd2,
	0xcd, 0x0c, 0x13, 0xec, 0x5f, 0x97, 0x44, 0x17, 0xc4, 0xa7, 0x7e, 0x3d, 0x64, 0x5d, 0x19, 0x73,
	0x60, 0x81, 0x4f, 0xdc, 0x22, 0x2a, 0x90, 0x88, 0x46, 0xee, 0xb8, 0x14, 0xde, 0x5e, 0x0b, 0xdb,
	0xe0, 0x32, 0x3a, 0x0a, 0x49, 0x06, 0x24, 0x5c, 0xc2, 0xd3, 0xac, 0x62, 0x91, 0x95, 0xe4, 0x79,
	0xe7, 0xc8, 0x37, 0x6d, 0x8d, 0xd5, 0x4e, 0xa9, 0x6c, 0x56, 0xf4, 0xea, 0x65, 0x7a, 0xae, 0x08,
	0xba, 0x78, 0x25, 0x2e, 0x1c, 0xa6, 0xb4, 0xc6, 0xe8, 0xdd, 0x74, 0x1f, 0x4b, 0xbd, 0x8b, 0x8a,
	0x70, 0x3e, 0xb5, 0x66, 0x48, 0x03, 0xf6, 0x0e, 0x61, 0x35, 0x57, 0xb9, 0x86, 0xc1, 0x1d, 0x9e,
	0xe1, 0xf8, 0x98, 0x11, 0x69, 0xd9, 0x8e, 0x94, 0x9b, 0x1e, 0x87, 0xe9, 0xce, 0x55, 0x28, 0xdf,
	0x8c, 0xa1, 0x89, 0x0d, 0xbf, 0xe6, 0x42, 0x68, 0x41, 0x99, 0x2d, 0x0f, 0xb0, 0x54, 0xbb, 0x16,
}

var inv_sbox = [256]byte{
	0x52, 0x09, 0x6a, 0xd5, 0x30, 0x36, 0xa5, 0x38, 0xbf, 0x40, 0xa3, 0x9e, 0x81, 0xf3, 0xd7, 0xfb,
	0x7c, 0xe3, 0x39, 0x82, 0x9b, 0x2f, 0xff, 0x87, 0x34, 0x8e, 0x43, 0x44, 0xc4, 0xde, 0xe9, 0xcb,
	0x54, 0x7b, 0x94, 0x32, 0xa6, 0xc2, 0x23, 0x3d, 0xee, 0x4c, 0x95, 0x0b, 0x42, 0xfa, 0xc3, 0x4e,
	0x08, 0x2e, 0xa1, 0x66, 0x28, 0xd9, 0x24, 0xb2, 0x76, 0x5b, 0xa2, 0x49, 0x6d, 0x8b, 0xd1, 0x25,
	0x72, 0xf8, 0xf6, 0x64, 0x86, 0x68, 0x98, 0x16, 0xd4, 0xa4, 0x5c, 0xcc, 0x5d, 0x65, 0xb6, 0x92,
	0x6c, 0x70, 0x48, 0x50, 0xfd, 0xed, 0xb9, 0xda, 0x5e, 0x15, 0x46, 0x57, 0xa7, 0x8d, 0x9d, 0x84,
	0x90, 0xd8, 0xab, 0x00, 0x8c, 0xbc, 0xd3, 0x0a, 0xf7, 0xe4, 0x58, 0x05, 0xb8, 0xb3, 0x45, 0x06,
	0xd0, 0x2c, 0x1e, 0x8f, 0xca, 0x3f, 0x0f, 0x02, 0xc1, 0xaf, 0xbd, 0x03, 0x01, 0x13, 0x8a, 0x6b,
	0x3a, 0x91, 0x11, 0x41, 0x4f, 0x67, 0xdc, 0xea, 0x97, 0xf2, 0xcf, 0xce, 0xf0, 0xb4, 0xe6, 0x73,
	0x96, 0xac, 0x74, 0x22, 0xe7, 0xad, 0x35, 0x85, 0xe2, 0xf9, 0x37, 0xe8, 0x1c, 0x75, 0xdf, 0x6e,
	0x47, 0xf1, 0x1a, 0x71, 0x1d, 0x29, 0xc5, 0x89, 0x6f, 0xb7, 0x62, 0x0e, 0xaa, 0x18, 0xbe, 0x1b,
	0xfc, 0x56, 0x3e, 0x4b, 0xc6, 0xd2, 0x79, 0x20, 0x9a, 0xdb, 0xc0, 0xfe, 0x78, 0xcd, 0x5a, 0xf4,
	0x1f, 0xdd, 0xa8, 0x33, 0x88, 0x07, 0xc7, 0x31, 0xb1, 0x12, 0x10, 0x59, 0x27, 0x80, 0xec, 0x5f,
	0x60, 0x51, 0x7f, 0xa9, 0x19, 0xb5, 0x4a, 0x0d, 0x2d, 0xe5, 0x7a, 0x9f, 0x93, 0xc9, 0x9c, 0xef,
	0xa0, 0xe0, 0x3b, 0x4d, 0xae, 0x2a, 0xf5, 0xb0, 0xc8, 0xeb, 0xbb, 0x3c, 0x83, 0x53, 0x99, 0x61,
	0x17, 0x2b, 0x04, 0x7e, 0xba, 0x77, 0xd6, 0x26, 0xe1, 0x69, 0x14, 0x63, 0x55, 0x21, 0x0c, 0x7d,
}

var rcon = [10]uint32{
	0x01000000, 0x02000000, 0x04000000, 0x08000000, 0x10000000,
	0x20000000, 0x40000000, 0x80000000, 0x1b000000, 0x36000000,
}

type AES struct {
	// number of rounds
	nr int
	// number of word in the key
	nk int
	// number(word) of block
	nb int
	// length(bit) of block
	len int
	// key
	key []byte
}

type paddingFunc func(*[]byte, int)

func NewAES(key []byte) (*AES, error) {
	var nk, nr int
	switch len(key) {
	case 16:
		nk = 4
		nr = 10
	case 24:
		nk = 6
		nr = 12
	case 32:
		nk = 8
		nr = 14
	default:
		return nil, errors.New("Key length doesn't match.")
	}
	return &AES{
		nr:  nr,
		nk:  nk,
		nb:  4,
		len: 16,
		key: key,
	}, nil
}

func ReadHex(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil
	}
	for i, _ := range data {
		if data[i] >= 'a' && data[i] <= 'f' {
			data[i] = data[i] - 'a' + 'A'
		}
	}
	msg := make([]byte, len(data)/2)
	for i, _ := range msg {
		if data[2*i] >= 'A' && data[2*i] <= 'F' {
			msg[i] = (data[2*i] - 'A' + 10) << 4
		} else if data[2*i] >= '0' && data[2*i] <= '9' {
			msg[i] = (data[2*i] - '0') << 4
		}

		if data[2*i+1] >= 'A' && data[2*i+1] <= 'F' {
			msg[i] ^= data[2*i+1] - 'A' + 10
		} else if data[2*i+1] >= '0' && data[2*i+1] <= '9' {
			msg[i] ^= data[2*i+1] - '0'
		}
	}
	return msg
}

func WriteHex(filename string, msg []byte) error {
	//hexString := make([]byte, 2*len(msg))
	//for i, _ := range msg{
	//	if (msg[i]>>4) >= 10 {
	//		hexString[2*i] = (msg[i]>>4) + 'A'
	//	}else if (msg[i]>>4) <= 9 {
	//		hexString[2*i+1] = (msg[i]>>4) + '0'
	//	}
	//
	//	if ((msg[i]<<4)>>4) >= 10 {
	//		hexString[2*i] = ((msg[i]<<4)>>4) + 'A'
	//	}else if ((msg[i]<<4)>>4) <= 9 {
	//		hexString[2*i+1] = ((msg[i]<<4)>>4) + '0'
	//	}
	//}

	file, err := os.OpenFile(filename, os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Open file err =", err)
		return err
	}
	defer file.Close()

	var test string
	test = fmt.Sprintf("%02X", msg)
	_, err1 := file.WriteString(test)
	if err1 != nil {
		return err1
	}
	//for i, _ := range hexString{
	//	_, err1 := fmt.Fprintf(file, "%s", string(hexString[i]))
	//	if err1 != nil {
	//		return err1
	//	}
	//}

	//file, err := os.OpenFile(filename, os.O_APPEND, 0666)
	//if err != nil{
	//	fmt.Println("Open file err =", err)
	//	return err
	//}
	//defer file.Close()
	//for i, _ := range msg{
	//	_, err1 := fmt.Fprintf(file, "%02X", msg[i])
	//	if err1 != nil {
	//		return err1
	//	}
	//}
	return nil
}

//
func DumpWords(msg string, in []uint32) {
	fmt.Printf("%s", msg)
	for i, v := range in {
		if i%4 == 0 {
			fmt.Printf("\nw[%02d]: %.8X ", i/4, v)
		} else {
			fmt.Printf("%.8X ", v)
		}
	}
	fmt.Println("\n")
}

//
func DumpBytes(msg string, in []byte) {
	fmt.Printf("%s", msg)
	for i, v := range in {
		if i%16 == 0 {
			fmt.Printf("\nblock[%d]: %02X", i/16, v)
		} else {
			if i%4 == 0 {
				fmt.Printf(" %02X", v)
			} else {
				fmt.Printf("%02X", v)
			}
		}
	}
	fmt.Println("\n")
}

// LittleEndian or BigEndian matters.
func (a *AES) keyExpansion() []uint32 {
	var w []uint32
	var j int
	for i := 0; i < 4; i++ {
		w = append(w, binary.BigEndian.Uint32(a.key[4*i:4*i+4]))
	}
	for i := 4; (j < 10) || (i%4 != 0); i++ {
		if i%4 != 0 {
			w = append(w, w[i-4]^w[i-1])
		} else {
			tempW := make([]byte, 4)
			binary.BigEndian.PutUint32(tempW, w[i-1])
			RotWord(tempW)
			a.subBytes(tempW)
			w = append(w, w[i-4]^rcon[j]^binary.BigEndian.Uint32(tempW))
			j++
		}
	}

	//DumpWords("keyExpansion:", w)
	return w
}

func (a *AES) EncryptECB(in []byte, pad paddingFunc) {
	pad(&in, a.len)
	roundKeys := a.keyExpansion()
	for i := 0; i < len(in); i += 16 {
		a.encryptBlock(in[i:i+16], roundKeys)
	}
	//fmt.Printf("AES-%d ECB encrypted cipher:\n", a.len*8)
	//DumpBytes("", in)
}

func (a *AES) DecryptECB(in []byte) {
	//pad(&in, a.len)
	roundKeys := a.keyExpansion()
	for i := 0; i < len(in); i += 16 {
		a.decryptBlock(in[i:i+16], roundKeys)
	}
	//fmt.Printf("AES-%d ECB decrypted plaintext:", a.len*8)
	//DumpBytes("", in)
}

func (a *AES) EncryptCBC(in []byte, iv []byte, pad paddingFunc) {
	pad(&in, a.len)
	roundKeys := a.keyExpansion()
	ivTmp := make([]byte, len(iv))
	copy(ivTmp, iv)

	for i := 0; i < len(in); i += 16 {
		Xor(in[i:i+16], ivTmp)
		a.encryptBlock(in[i:i+16], roundKeys)
		copy(ivTmp, in[i:i+16])
	}
	//fmt.Printf("AES-%d CBC encrypted cipher:\n", a.len*8)
	//DumpBytes("", in)
}

func (a *AES) DecryptCBC(in []byte, iv []byte) {
	roundKeys := a.keyExpansion()
	ivTmp := make([]byte, len(iv))
	copy(ivTmp, iv)
	reg := make([]byte, len(iv))

	for i := 0; i < len(in); i += 16 {
		copy(reg, in[i:i+16])
		a.decryptBlock(in[i:i+16], roundKeys)
		Xor(in[i:i+16], ivTmp)
		copy(ivTmp, reg)
	}
	//fmt.Printf("AES-%d CBC decrypted plaintext:", a.len*8)
	//DumpBytes("", in)
}

func (a *AES) EncryptCFB32(in []byte, iv []byte, pad paddingFunc) {
	pad(&in, 4)
	roundKeys := a.keyExpansion()
	ivTmp := make([]byte, len(iv))
	copy(ivTmp, iv)

	for i := 0; i < len(in); i += 4 {
		a.encryptBlock(ivTmp, roundKeys)
		Xor(in[i:i+4], ivTmp[0:4])
		ivTmp = append(ivTmp[4:], in[i:i+4]...)
	}
	//fmt.Printf("AES-%d CFB32 encrypted cipher:\n", a.len*8)
	//DumpBytes("", in)
}

func (a *AES) DecryptCFB32(in []byte, iv []byte) {
	roundKeys := a.keyExpansion()
	ivTmp := make([]byte, len(iv))
	copy(ivTmp, iv)
	cipherTmp := make([]byte, len(in))
	copy(cipherTmp, in)

	for i := 0; i < len(in); i += 4 {
		a.encryptBlock(ivTmp, roundKeys)
		Xor(in[i:i+4], ivTmp[0:4])
		ivTmp = append(ivTmp[4:], cipherTmp[i:i+4]...)
	}
	//fmt.Printf("AES-%d CFB32 decrypted plaintext:", a.len*8)
	//DumpBytes("", in)
}

func (a *AES) EncryptOFB32(in []byte, iv []byte, pad paddingFunc) {
	pad(&in, 4)
	roundKeys := a.keyExpansion()
	ivTmp := make([]byte, len(iv))
	copy(ivTmp, iv)

	for i := 0; i < len(in); i += 4 {
		a.encryptBlock(ivTmp, roundKeys)
		Xor(in[i:i+4], ivTmp[0:4])
		ivTmp = append(ivTmp[4:], ivTmp[0:4]...)
	}
	//fmt.Printf("AES-%d OFB32 encrypted cipher:\n", a.len*8)
	//DumpBytes("", in)
}

func (a *AES) DecryptOFB32(in []byte, iv []byte) {
	roundKeys := a.keyExpansion()
	ivTmp := make([]byte, len(iv))
	copy(ivTmp, iv)

	for i := 0; i < len(in); i += 4 {
		a.encryptBlock(ivTmp, roundKeys)
		Xor(in[i:i+4], ivTmp[0:4])
		iv = append(ivTmp[4:], ivTmp[0:4]...)
	}
	//fmt.Printf("AES-%d OFB32 decrypted plaintext:", a.len*8)
	//DumpBytes("", in)
}

//
func PaddingZeros(in *[]byte, blockLen int) {
	for len(*in)%blockLen != 0 {
		*in = append(*in, 0x00)
	}
}

//
func (a *AES) subBytes(state []byte) {
	for i, v := range state {
		state[i] = sbox[v]
	}
}

//
func (a *AES) invSubBytes(state []byte) {
	for i, v := range state {
		state[i] = inv_sbox[v]
	}
}

//
func (a *AES) shiftRow(in []byte, i int, n int) {
	in[i], in[i+4*1], in[i+4*2], in[i+4*3] = in[i+4*(n%4)], in[i+4*((n+1)%4)], in[i+4*((n+2)%4)], in[i+4*((n+3)%4)]
}

func RotWord(in []byte) {
	in[0], in[1], in[2], in[3] = in[1], in[2], in[3], in[0]
}

//
func (a *AES) shiftRows(state []byte) {
	a.shiftRow(state, 1, 1)
	a.shiftRow(state, 2, 2)
	a.shiftRow(state, 3, 3)
}

//
func (a *AES) invShiftRows(state []byte) {
	a.shiftRow(state, 1, 3)
	a.shiftRow(state, 2, 2)
	a.shiftRow(state, 3, 1)
}

//
func xtime(in byte) byte {
	return (in << 1) ^ (((in >> 7) & 1) * 0x1b)
}

//
func xtimes(in byte, ts int) byte {
	for ts > 0 {
		in = xtime(in)
		ts--
	}
	return in
}

//
func mulBytes(x byte, y byte) byte {
	return (((y >> 0) & 0x01) * xtimes(x, 0)) ^
		(((y >> 1) & 0x01) * xtimes(x, 1)) ^
		(((y >> 2) & 0x01) * xtimes(x, 2)) ^
		(((y >> 3) & 0x01) * xtimes(x, 3)) ^
		(((y >> 4) & 0x01) * xtimes(x, 4)) ^
		(((y >> 5) & 0x01) * xtimes(x, 5)) ^
		(((y >> 6) & 0x01) * xtimes(x, 6)) ^
		(((y >> 7) & 0x01) * xtimes(x, 7))
}

//
func mulWords(x []byte, y []byte) {
	tmp := make([]byte, 4)
	copy(tmp, x)

	x[0] = mulBytes(tmp[0], y[3]) ^ mulBytes(tmp[1], y[0]) ^ mulBytes(tmp[2], y[1]) ^ mulBytes(tmp[3], y[2])
	x[1] = mulBytes(tmp[0], y[2]) ^ mulBytes(tmp[1], y[3]) ^ mulBytes(tmp[2], y[0]) ^ mulBytes(tmp[3], y[1])
	x[2] = mulBytes(tmp[0], y[1]) ^ mulBytes(tmp[1], y[2]) ^ mulBytes(tmp[2], y[3]) ^ mulBytes(tmp[3], y[0])
	x[3] = mulBytes(tmp[0], y[0]) ^ mulBytes(tmp[1], y[1]) ^ mulBytes(tmp[2], y[2]) ^ mulBytes(tmp[3], y[3])
}

//
func (a *AES) mixColumns(state []byte) {
	s := []byte{0x03, 0x01, 0x01, 0x02}
	for i := 0; i < len(state); i += 4 {
		mulWords(state[i:i+4], s)
	}
}

//
func (a *AES) invMixColumns(state []byte) {
	s := []byte{0x0b, 0x0d, 0x09, 0x0e}
	for i := 0; i < len(state); i += 4 {
		mulWords(state[i:i+4], s)
	}
}

//
func Xor(x []byte, y []byte) {
	if len(x) == len(y) {
		for i := 0; i < len(x); i++ {
			x[i] = x[i] ^ y[i]
		}
	}
}

//
func (a *AES) addRoundKey(state []byte, w []uint32) {
	tmp := make([]byte, a.len)
	for i := 0; i < len(w); i += 1 {
		binary.BigEndian.PutUint32(tmp[4*i:4*i+4], w[i])
	}
	Xor(state, tmp)
}

func (a *AES) encryptBlock(state []byte, roundKeys []uint32) {
	a.addRoundKey(state, roundKeys[0:4])
	for round := 1; round < a.nr; round++ {
		a.subBytes(state)  //ok
		a.shiftRows(state) //ok
		a.mixColumns(state)
		a.addRoundKey(state, roundKeys[4*round:4*round+4]) //ok
		//if round == 1 {
		//DumpBytes("第1轮加密结果：", state)
		//}
	}
	a.subBytes(state)
	a.shiftRows(state)
	a.addRoundKey(state, roundKeys[a.nr*4:a.nr*4+4])
}

func (a *AES) decryptBlock(state []byte, roundKeys []uint32) {
	a.addRoundKey(state, roundKeys[a.nr*4:a.nr*4+4])
	for round := a.nr - 1; round > 0; round-- {
		a.invShiftRows(state)
		a.invSubBytes(state)
		a.addRoundKey(state, roundKeys[4*round:4*round+4])
		a.invMixColumns(state)
	}
	a.invShiftRows(state)
	a.invSubBytes(state)
	a.addRoundKey(state, roundKeys[0:4])
}
