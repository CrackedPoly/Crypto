package action

import (
	"context"
	"flag"
	"fmt"
	"github.com/CrackedPoly/Crypto/seccommu/src/seccommu"
	"github.com/otokaze/go-kit/log"
	"github.com/peterh/liner"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	log2 "log"
	"strings"
)

func NewcerAction(ctx *cli.Context) error {
	childcer := ctx.String("childcer")
	childpvk := ctx.String("childpvk")
	rootcer := ctx.String("rootcer")
	rootpvk := ctx.String("rootpvk")
	cn := ctx.String("cn")

	err := seccommu.NewCer(childcer, childpvk, rootcer, rootpvk, cn)

	fmt.Printf("Generate \"%s\", \"%s\" successfully.\n", childpvk, childcer)
	return err
}

func SendAction(ctx *cli.Context) error {
	cer := ctx.String("cer")
	keyout := ctx.String("keyout")
	msg := ctx.String("msg")
	msgout := ctx.String("msgout")
	pvk := ctx.String("pvk")
	signedout := ctx.String("signedout")

	key := seccommu.Rand16()
	pubKey := seccommu.ParsePubKey(cer)
	seccommu.EncryptRSA(pubKey, key, keyout)

	plaintext, err := ioutil.ReadFile(msg)
	if err != nil {
		return err
	}
	_, err1 := seccommu.EncryptFile(key, plaintext, msgout)
	if err1 != nil {
		return err1
	}

	priKey := seccommu.ParsePriKey(pvk)
	seccommu.Sign(priKey, seccommu.Hash(plaintext), signedout)

	fmt.Printf("Encrypted key in \"%s\"\n", keyout)
	fmt.Printf("Encrypted message in \"%s\"\n", msgout)
	fmt.Printf("Signed hash of message in \"%s\"\n", signedout)

	return nil
}

func GetKeyAction(ctx *cli.Context) error {
	keyout := ctx.String("keyout")
	pvk := ctx.String("pvk")
	keyrec := ctx.String("keyrec")

	priKey := seccommu.ParsePriKey(pvk)
	cipher, err := ioutil.ReadFile(keyout)
	if err != nil {
		return err
	}
	key := seccommu.DecryptRSA(priKey, cipher)

	fmt.Printf("key: %x, recovered in %s successfully.\n", key, keyrec)

	err1 := ioutil.WriteFile(keyrec, key, 0777)
	if err1 != nil {
		return err1
	}
	return nil
}

func GetMessageAction(ctx *cli.Context) error {
	keyrec := ctx.String("keyrec")
	msgout := ctx.String("msgout")
	msgrec := ctx.String("msgrec")

	key, err := ioutil.ReadFile(keyrec)
	if err != nil {
		return err
	}

	cipher, err1 := ioutil.ReadFile(msgout)
	if err1 != nil {
		return err1
	}

	_, err2 := seccommu.DecryptFile(key, cipher, msgrec)
	if err2 != nil {
		return err2
	}

	fmt.Printf("Message recovered in %s successfully.\n", msgrec)
	return nil
}

func VerifyAction(ctx *cli.Context) error {
	signedout := ctx.String("signedout")
	cer := ctx.String("cer")
	msgrec := ctx.String("msgrec")

	sig, err := ioutil.ReadFile(signedout)
	if err != nil {
		return err
	}

	pubKey := seccommu.ParsePubKey(cer)

	msg, err1 := ioutil.ReadFile(msgrec)
	if err1 != nil {
		return err1
	}
	hashed := seccommu.Hash(msg)

	err2 := seccommu.Verify(pubKey, hashed, sig)
	if err2 != nil {
		log2.Fatal(err2)
	}
	println("true")
	return nil
}

func DiffAction(ctx *cli.Context) error {
	msg := ctx.String("msg")
	msgrec := ctx.String("msgrec")

	equal, err := seccommu.Equal(msg, msgrec)
	if err != nil {
		return err
	}
	if equal {
		println("success")
	} else {
		println("failure")
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
