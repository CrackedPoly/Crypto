package action

import (
	"context"
	"errors"
	"flag"
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
		if err != nil {
			log.Error("%s", err)
			return err
		}
		cipher := _aes.EncryptECB(plain, utils.PKCS7Padding)
		utils.WriteBytesHex(cipherfile, cipher)
	case "CBC":
		plain := utils.ReadBytesHex(plainfile)
		key := utils.ReadBytesHex(keyfile)
		iv := utils.ReadBytesHex(vifile)
		_aes, err := aes.NewAES(key)
		if err != nil {
			log.Error("%s", err)
			return err
		}
		cipher := _aes.EncryptCBC(plain, iv[0:16], utils.PKCS7Padding)
		utils.WriteBytesHex(cipherfile, cipher)
	case "CFB":
		plain := utils.ReadBytesHex(plainfile)
		key := utils.ReadBytesHex(keyfile)
		iv := utils.ReadBytesHex(vifile)
		_aes, err := aes.NewAES(key)
		if err != nil {
			log.Error("%s", err)
			return err
		}
		cipher := _aes.EncryptCFB(plain, iv[0:16], 16)
		utils.WriteBytesHex(cipherfile, cipher)
	case "OFB":
		plain := utils.ReadBytesHex(plainfile)
		key := utils.ReadBytesHex(keyfile)
		iv := utils.ReadBytesHex(vifile)
		_aes, err := aes.NewAES(key)
		if err != nil {
			log.Error("%s", err)
			return err
		}
		cipher := _aes.EncryptOFB(plain, iv[0:16])
		utils.WriteBytesHex(cipherfile, cipher)
	default:
		log.Error("invalid mode")
		return errors.New("invalid mode")
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
		if err != nil {
			log.Error("%s", err)
			return err
		}
		plain := _aes.DecryptECB(cipher, utils.PKCS7Unpadding)
		utils.WriteBytesHex(plainfile, plain)
	case "CBC":
		cipher := utils.ReadBytesHex(cipherfile)
		key := utils.ReadBytesHex(keyfile)
		iv := utils.ReadBytesHex(vifile)
		_aes, err := aes.NewAES(key)
		if err != nil {
			log.Error("%s", err)
			return err
		}
		plain := _aes.DecryptCBC(cipher, iv[0:16], utils.PKCS7Unpadding)
		utils.WriteBytesHex(plainfile, plain)
	case "CFB":
		cipher := utils.ReadBytesHex(cipherfile)
		key := utils.ReadBytesHex(keyfile)
		iv := utils.ReadBytesHex(vifile)
		_aes, err := aes.NewAES(key)
		if err != nil {
			log.Error("%s", err)
			return err
		}
		plain := _aes.DecryptCFB(cipher, iv[0:16], 16)
		utils.WriteBytesHex(plainfile, plain)
	case "OFB":
		cipher := utils.ReadBytesHex(cipherfile)
		key := utils.ReadBytesHex(keyfile)
		iv := utils.ReadBytesHex(vifile)
		_aes, err := aes.NewAES(key)
		if err == nil {
			log.Error("%s", err)
			return err
		}
		plain := _aes.DecryptOFB(cipher, iv[0:16])
		utils.WriteBytesHex(plainfile, plain)
	default:
		log.Error("invalid mode")
		return errors.New("invalid mode")
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
