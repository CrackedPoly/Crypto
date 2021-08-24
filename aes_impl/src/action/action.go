package action

import (
	"context"
	"flag"
	"fmt"
	"github.com/UnknwoonUser/Crypto/aes_impl/src/aes"
	"github.com/UnknwoonUser/Crypto/utils"
	"github.com/otokaze/go-kit/log"
	"github.com/peterh/liner"
	"github.com/urfave/cli/v2"
	"strings"
)

func EncryptAction(ctx *cli.Context) (err error) {
	mode := ctx.String("m")
	plainfile := ctx.String("p")
	keyfile := ctx.String("k")
	vifile := ctx.String("v")
	cipherfile := ctx.String("c")

	switch strings.ToTitle(mode) {
	case "ECB":
		plain := utils.ReadBytesHex(plainfile)
		key := utils.ReadBytesHex(keyfile)
		_aes, err := aes.NewAES(key)
		if err == nil {
			cipher := _aes.EncryptECB(plain, utils.PKCS7Padding)
			utils.WriteBytesHex(cipherfile, cipher)
			//aes.DecryptECB(msg)
		}
	case "CBC":
		plain := utils.ReadBytesHex(plainfile)
		key := utils.ReadBytesHex(keyfile)
		iv := utils.ReadBytesHex(vifile)
		_aes, err := aes.NewAES(key)
		if err == nil {
			cipher := _aes.EncryptCBC(plain, iv, utils.PKCS7Padding)
			utils.WriteBytesHex(cipherfile, cipher)
			//aes.DecryptCBC(msg, iv)
		}
	case "CFB":
		plain := utils.ReadBytesHex(plainfile)
		key := utils.ReadBytesHex(keyfile)
		iv := utils.ReadBytesHex(vifile)
		_aes, err := aes.NewAES(key)
		if err == nil {
			cipher := _aes.EncryptCFB(plain, iv, 16)
			utils.WriteBytesHex(cipherfile, cipher)
			//aes.DecryptCFB32(msg, iv)
		}
	case "OFB":
		plain := utils.ReadBytesHex(plainfile)
		key := utils.ReadBytesHex(keyfile)
		iv := utils.ReadBytesHex(vifile)
		_aes, err := aes.NewAES(key)
		if err == nil {
			cipher := _aes.EncryptOFB(plain, iv)
			utils.WriteBytesHex(cipherfile, cipher)
			//aes.DecryptOFB32(msg, iv)
		}
	default:
		fmt.Println("Invalid mode!")
	}
	return nil
}

func DecryptAction(ctx *cli.Context) (err error) {
	mode := ctx.String("m")
	plainfile := ctx.String("p")
	keyfile := ctx.String("k")
	vifile := ctx.String("v")
	cipherfile := ctx.String("c")

	switch strings.ToTitle(mode) {
	case "ECB":
		cipher := utils.ReadBytesHex(cipherfile)
		key := utils.ReadBytesHex(keyfile)
		_aes, err := aes.NewAES(key)
		if err == nil {
			plain := _aes.DecryptECB(cipher, utils.PKCS7Unpadding)
			utils.WriteBytesHex(plainfile, plain)
		}
	case "CBC":
		cipher := utils.ReadBytesHex(cipherfile)
		key := utils.ReadBytesHex(keyfile)
		iv := utils.ReadBytesHex(vifile)
		_aes, err := aes.NewAES(key)
		if err == nil {
			plain := _aes.DecryptCBC(cipher, iv, utils.PKCS7Unpadding)
			utils.WriteBytesHex(plainfile, plain)
		}
	case "CFB":
		cipher := utils.ReadBytesHex(cipherfile)
		key := utils.ReadBytesHex(keyfile)
		iv := utils.ReadBytesHex(vifile)
		_aes, err := aes.NewAES(key)
		if err == nil {
			plain := _aes.DecryptCFB(cipher, iv, 16)
			utils.WriteBytesHex(plainfile, plain)
		}
	case "OFB":
		cipher := utils.ReadBytesHex(cipherfile)
		key := utils.ReadBytesHex(keyfile)
		iv := utils.ReadBytesHex(vifile)
		_aes, err := aes.NewAES(key)
		if err == nil {
			plain := _aes.DecryptOFB(cipher, iv)
			utils.WriteBytesHex(plainfile, plain)
		}
	default:
		fmt.Println("Invalid mode!")
	}
	return nil
}

func AfterAction(ctx *cli.Context) (err error) {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	defer line.Close()
	line.SetCompleter(func(line string) (cs []string) {
		for _, c := range ctx.App.Commands {
			if strings.HasPrefix(c.Name, strings.ToLower(line)) {
				cs = append(cs, c.Name)
			}
		}
		return
	})
	for {
		var input string
		if input, err = line.Prompt("> "); err != nil {
			if err == liner.ErrPromptAborted {
				err = nil
				return
			}
			log.Error("line.Prompt() error(%v)", err)
			continue
		}
		line.AppendHistory(input)
		input = strings.TrimSpace(input)
		args := strings.Split(input, " ")
		if args[0] == "" {
			continue
		}
		var command *cli.Command
		if command = ctx.App.Command(args[0]); command == nil {
			log.Error("command(%s) not found!", args[0])
			continue
		}
		var fset = flag.NewFlagSet(args[0], flag.ContinueOnError)
		for _, f := range command.Flags {
			f.Apply(fset)
		}
		if !command.SkipFlagParsing {
			if err = fset.Parse(args[1:]); err != nil {
				if err == flag.ErrHelp {
					err = nil
					continue
				}
				log.Error("fs.Parse(%v) error(%v)", args[1:], err)
				continue
			}
		}
		var key interface{} = "args"
		nCtx := cli.NewContext(ctx.App, fset, ctx)
		nCtx.Context = context.WithValue(nCtx.Context, key, args[1:])
		nCtx.Command = command
		command.Action(nCtx)
	}
}

//case "ECBTEST":
//msg := Aes.ReadStringHex(plainfile)
//key := Aes.ReadStringHex(keyfile)
//aes, err := Aes.NewAES(key)
//if err == nil {
//start := time.Now() // 获取当前时间
//for i := 0; i < 10; i++ {
//aes.EncryptECB(msg, Aes.ZeroPadding)
//}
//for i := 0; i < 10; i++ {
//aes.DecryptECB(msg)
//}
//elapsed := time.Since(start)
//fmt.Println("十次加解密时间：", elapsed)
//fmt.Println("速度：", 0.005/float64(elapsed), "MB/s")
//aes.DecryptECB(msg)
//}
//case "CBCTEST":
//msg := Aes.ReadStringHex(plainfile)
//key := Aes.ReadStringHex(keyfile)
//iv := Aes.ReadStringHex(vifile)
//aes, err := Aes.NewAES(key)
//if err == nil {
//start := time.Now() // 获取当前时间
//for i := 0; i < 10; i++ {
//aes.EncryptCBC(msg, iv, Aes.ZeroPadding)
//}
//for i := 0; i < 10; i++ {
//aes.DecryptCBC(msg, iv)
//}
//elapsed := time.Since(start)
//fmt.Println("十次加解密时间：", elapsed)
//fmt.Println("速度：", 0.005/float64(elapsed), "MB/s")
//aes.DecryptCBC(msg, iv)
//}
//case "CFBTEST":
//msg := Aes.ReadStringHex(plainfile)
//key := Aes.ReadStringHex(keyfile)
//iv := Aes.ReadStringHex(vifile)
//aes, err := Aes.NewAES(key)
//if err == nil {
//start := time.Now() // 获取当前时间
//for i := 0; i < 10; i++ {
//aes.EncryptCFB32(msg, iv, Aes.ZeroPadding)
//}
//for i := 0; i < 10; i++ {
//aes.DecryptCFB32(msg, iv)
//}
//elapsed := time.Since(start)
//fmt.Println("十次加解密时间：", elapsed)
//fmt.Println("速度：", 0.005/float64(elapsed), "MB/s")
//aes.DecryptCBC(msg, iv)
//}
//case "OFBTEST":
//msg := Aes.ReadStringHex(plainfile)
//key := Aes.ReadStringHex(keyfile)
//iv := Aes.ReadStringHex(vifile)
//aes, err := Aes.NewAES(key)
//if err == nil {
//start := time.Now() // 获取当前时间
//for i := 0; i < 10; i++ {
//aes.EncryptOFB32(msg, iv, Aes.ZeroPadding)
//}
//for i := 0; i < 10; i++ {
//aes.DecryptOFB32(msg, iv)
//}
//elapsed := time.Since(start)
//fmt.Println("十次加解密时间：", elapsed)
//fmt.Println("速度：", 0.005/float64(elapsed), "MB/s")
//aes.DecryptCBC(msg, iv)
//}
