package main

import (
	_ "encoding/binary"
	"fmt"
	Aes "github.com/UnknwoonUser/Crypto/AES"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
	"strings"
	"time"
	_ "time"
)

func main() {
	var plainfile, keyfile, vifile, cipherfile, mode string
	var app = &cli.App{
		Name:                 "应用密码学实践-2019141440070-罗鉴",
		Usage:                "AES加密与解密",
		EnableBashCompletion: true,
		Commands:             cli.Commands{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "plainfile",
				Usage:       "指定明文件的位置和名称",
				Aliases:     []string{"p"},
				Destination: &plainfile,
			},
			&cli.StringFlag{
				Name:        "keyfile",
				Usage:       "指定密钥文件的位置和名称",
				Aliases:     []string{"k"},
				Destination: &keyfile,
			},
			&cli.StringFlag{
				Name:        "vifile",
				Usage:       "指定初始化向量文件的位置和名称",
				Aliases:     []string{"v"},
				Destination: &vifile,
			},
			&cli.StringFlag{
				Name:        "cipherfile",
				Usage:       "指定密文文件的位置和名称",
				Aliases:     []string{"c"},
				Destination: &cipherfile,
			},
			&cli.StringFlag{
				Name:        "mode",
				Usage:       "指定加密的操作模式，有ECB、CBC、CFB、OFB四种",
				Aliases:     []string{"m"},
				Destination: &mode,
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			switch strings.ToTitle(mode) {
			case "ECB":
				msg := Aes.ReadHex(plainfile)
				key := Aes.ReadHex(keyfile)
				aes, err := Aes.NewAES(key)
				if err == nil {
					aes.EncryptECB(msg, Aes.PaddingZeros)
					err1 := Aes.WriteHex(cipherfile, msg)
					if err1 != nil {
						err1.Error()
					}
					//aes.DecryptECB(msg)
				}
			case "CBC":
				msg := Aes.ReadHex(plainfile)
				key := Aes.ReadHex(keyfile)
				iv := Aes.ReadHex(vifile)
				aes, err := Aes.NewAES(key)
				if err == nil {
					aes.EncryptCBC(msg, iv, Aes.PaddingZeros)
					err1 := Aes.WriteHex(cipherfile, msg)
					if err1 != nil {
						err1.Error()
					}
					//aes.DecryptCBC(msg, iv)
				}
			case "CFB":
				msg := Aes.ReadHex(plainfile)
				key := Aes.ReadHex(keyfile)
				iv := Aes.ReadHex(vifile)
				aes, err := Aes.NewAES(key)
				if err == nil {
					aes.EncryptCFB32(msg, iv, Aes.PaddingZeros)
					err1 := Aes.WriteHex(cipherfile, msg)
					if err1 != nil {
						err1.Error()
					}
					//aes.DecryptCFB32(msg, iv)
				}
			case "OFB":
				msg := Aes.ReadHex(plainfile)
				key := Aes.ReadHex(keyfile)
				iv := Aes.ReadHex(vifile)
				aes, err := Aes.NewAES(key)
				if err == nil {
					aes.EncryptOFB32(msg, iv, Aes.PaddingZeros)
					err1 := Aes.WriteHex(cipherfile, msg)
					if err1 != nil {
						err1.Error()
					}
					//aes.DecryptOFB32(msg, iv)
				}
			case "ECBTEST":
				msg := Aes.ReadHex(plainfile)
				key := Aes.ReadHex(keyfile)
				aes, err := Aes.NewAES(key)
				if err == nil {
					start := time.Now() // 获取当前时间
					for i := 0; i <10; i++ {
						aes.EncryptECB(msg, Aes.PaddingZeros)
					}
					for i := 0; i <10; i++ {
						aes.DecryptECB(msg)
					}
					elapsed := time.Since(start)
					fmt.Println("十次加解密时间：", elapsed)
					fmt.Println("速度：", 0.005/float64(elapsed) ,"MB/s")
					aes.DecryptECB(msg)
				}
			case "CBCTEST":
				msg := Aes.ReadHex(plainfile)
				key := Aes.ReadHex(keyfile)
				iv := Aes.ReadHex(vifile)
				aes, err := Aes.NewAES(key)
				if err == nil {
					start := time.Now() // 获取当前时间
					for i := 0; i <10; i++ {
						aes.EncryptCBC(msg, iv, Aes.PaddingZeros)
					}
					for i := 0; i <10; i++ {
						aes.DecryptCBC(msg, iv)
					}
					elapsed := time.Since(start)
					fmt.Println("十次加解密时间：", elapsed)
					fmt.Println("速度：", 0.005/float64(elapsed) ,"MB/s")
					aes.DecryptCBC(msg, iv)
				}
			case "CFBTEST":
				msg := Aes.ReadHex(plainfile)
				key := Aes.ReadHex(keyfile)
				iv := Aes.ReadHex(vifile)
				aes, err := Aes.NewAES(key)
				if err == nil {
					start := time.Now() // 获取当前时间
					for i := 0; i <10; i++ {
						aes.EncryptCFB32(msg, iv, Aes.PaddingZeros)
					}
					for i := 0; i <10; i++ {
						aes.DecryptCFB32(msg, iv)
					}
					elapsed := time.Since(start)
					fmt.Println("十次加解密时间：", elapsed)
					fmt.Println("速度：", 0.005/float64(elapsed) ,"MB/s")
					aes.DecryptCBC(msg, iv)
				}
			case "OFBTEST":
				msg := Aes.ReadHex(plainfile)
				key := Aes.ReadHex(keyfile)
				iv := Aes.ReadHex(vifile)
				aes, err := Aes.NewAES(key)
				if err == nil {
					start := time.Now() // 获取当前时间
					for i := 0; i <10; i++ {
						aes.EncryptOFB32(msg, iv, Aes.PaddingZeros)
					}
					for i := 0; i <10; i++ {
						aes.DecryptOFB32(msg, iv)
					}
					elapsed := time.Since(start)
					fmt.Println("十次加解密时间：", elapsed)
					fmt.Println("速度：", 0.005/float64(elapsed) ,"MB/s")
					aes.DecryptCBC(msg, iv)
				}
			default:
				fmt.Println("Invalid mode!")
			}
			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	app.Run(os.Args)

}
