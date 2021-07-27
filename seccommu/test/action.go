package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/UnknwoonUser/Crypto/seccommu"
	"github.com/otokaze/go-kit/log"
	"github.com/peterh/liner"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	log2 "log"
	"strings"
)

func newcerAction(ctx *cli.Context) error {
	err := seccommu.NewCer(childcer, childpvk, rootcer, rootpvk, cn)

	fmt.Printf("Generate \"%s\", \"%s\" successfully.\n", childpvk, childcer)
	return err
}

func sendAction(ctx *cli.Context) error {
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

func getkeyAction(ctx *cli.Context) error {
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

func getMessageAction(ctx *cli.Context) error {
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

func verifyAction(ctx *cli.Context) error {
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

func diffAction(ctx *cli.Context) error {
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

func afterAction(ctx *cli.Context) (err error) {
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
