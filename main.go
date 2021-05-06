package main

import (
	_ "./AES"
	_ "encoding/binary"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
)

func main(){
	var plainfile, keyfile, vifile, cipherfile, mode string
	var app = &cli.App{
		Name: 	"应用密码学实践-2019141440070-罗鉴",
		Usage: 	"AES加密与解密",
		EnableBashCompletion: true,
		Commands: cli.Commands{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "plainfile",
				Usage: "明文文件的位置",
				Aliases: []string{"p"},
				Destination: &plainfile,
			},
			&cli.StringFlag{
				Name: "keyfile",
				Usage: "密文文件的位置",
				Aliases: []string{"k"},
				Destination: &keyfile,
			},
			&cli.StringFlag{
				Name: "vifile",
				Usage: "初始向量文件的位置",
				Aliases: []string{"v"},
				Destination: &vifile,
			},
			&cli.StringFlag{
				Name: "cipherfile",
				Usage: "输出密文文件的位置",
				Aliases: []string{"c"},
				Destination: &cipherfile,
			},
			&cli.StringFlag{
				Name: "mode",
				Usage: "选择模式，有ECB、CBC、CFB、OFB四种",
				Aliases: []string{"m"},
				Destination: &mode,
			},
		},
		Action: func(ctx *cli.Context) (err error) {

			return nil
		},

	}

	sort.Sort(cli.FlagsByName(app.Flags))
	app.Run(os.Args)

}
