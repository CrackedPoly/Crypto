package main

import (
	"context"
	"flag"
	"fmt"
	rsa "github.com/UnknwoonUser/Crypto/RSA"
	"github.com/otokaze/go-kit/log"
	"github.com/peterh/liner"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"strings"
)

func genkeyAction(ctx *cli.Context) error {
	Rsa := rsa.NewRSA(256)
	pString := Rsa.P.Text(16)
	qString := Rsa.Q.Text(16)
	nString := Rsa.N.Text(16)
	eString := Rsa.E.Text(16)
	dString := Rsa.D.Text(16)
	fmt.Println("p: ", pString)
	fmt.Println("q: ", qString)
	fmt.Println("n: ", nString)
	fmt.Println("e: ", eString)
	fmt.Println("d: ", dString)
	WriteHex("p.txt", pString)
	WriteHex("q.txt", qString)
	WriteHex("n.txt", nString)
	WriteHex("e.txt", eString)
	WriteHex("d.txt", dString)
	return nil
}

func encryptAction(ctx *cli.Context) error {
	plaintext := ReadHex(plainfile)
	eString := ReadHex(efile)
	nString := ReadHex(nfile)
	Rsa := rsa.NewCheck(nString, eString)
	cipher := Rsa.Encrypt(plaintext)
	fmt.Println("cipher: ", cipher)
	WriteHex(cipherfile, cipher)
	return nil
}

func signAction(ctx *cli.Context) error {
	plaintext := ReadHex(plainfile)
	dString := ReadHex(dfile)
	nString := ReadHex(nfile)
	Rsa := rsa.NewSign(nString, dString)
	cipher := Rsa.Decrypt(plaintext)
	fmt.Println("signed cipher: ", cipher)
	WriteHex(cipherfile, cipher)
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

func ReadHex(filename string) string {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error in reading", filename, err)
	}
	return string(f)
}

func WriteHex(filename string, msg string) {
	err := ioutil.WriteFile(filename, []byte(msg), 0666)
	if err != nil {
		fmt.Println("Error in writing", filename, err)
	}
}
