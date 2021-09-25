package action

import (
	"context"
	"flag"
	"fmt"
	"github.com/CrackedPoly/Crypto/rsa_impl/src/rsa"
	"github.com/CrackedPoly/Crypto/utils"
	"github.com/otokaze/go-kit/log"
	"github.com/peterh/liner"
	"github.com/urfave/cli/v2"
	"strings"
)

// GenkeyAction
// @Description:
// @param ctx
// @return error
func GenkeyAction(ctx *cli.Context) error {
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
	utils.WriteStringHex("p.txt", pString)
	utils.WriteStringHex("q.txt", qString)
	utils.WriteStringHex("n.txt", nString)
	utils.WriteStringHex("e.txt", eString)
	utils.WriteStringHex("d.txt", dString)
	return nil
}

func EncryptAction(ctx *cli.Context) error {
	plainfile := ctx.String("p")
	efile := ctx.String("e")
	cipherfile := ctx.String("c")

	plaintext := utils.ReadStringHex(plainfile)
	s := utils.ReadStringHex(efile)
	fmt.Println("public key: ", s)
	tmp := strings.Fields(s)
	Rsa := rsa.NewCheck(tmp[0], tmp[1])
	cipher := Rsa.Encrypt(plaintext)
	fmt.Println("cipher: ", cipher)
	utils.WriteStringHex(cipherfile, cipher)
	return nil
}

func SignAction(ctx *cli.Context) error {
	plainfile := ctx.String("p")
	efile := ctx.String("e")
	cipherfile := ctx.String("c")

	plaintext := utils.ReadStringHex(plainfile)
	s := utils.ReadStringHex(efile)
	fmt.Println("private key: ", s)
	tmp := strings.Fields(s)
	Rsa := rsa.NewSign(tmp[0], tmp[1])
	cipher := Rsa.Decrypt(plaintext)
	fmt.Println("signed message: ", cipher)
	utils.WriteStringHex(cipherfile, cipher)
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
