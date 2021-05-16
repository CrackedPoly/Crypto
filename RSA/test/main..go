package main

import "fmt"

var plainfile, nfile, efile, dfile, cipherfile string

func main() {
	//var app = &cli.App{
	//	Name:                 "应用密码学实践-2019141440070-罗鉴",
	//	Usage:                "RSA加解密和数字签名",
	//	EnableBashCompletion: true,
	//	Commands:           cli.Commands{
	//		{
	//			Name:	"genkey",
	//			Usage:  "产生RSA算法中的p、q、n、e、d",
	//			Action: genkeyAction,
	//		},
	//		{
	//			Name: 	"encrypt",
	//			Usage: 	"加密指定文件中的信息",
	//			Action: encryptAction,
	//			Flags: []cli.Flag{
	//				&cli.StringFlag{
	//					Name:        "plainfile",
	//					Usage:       "指定明文文件的位置和名称",
	//					Aliases:     []string{"p"},
	//					Destination: &plainfile,
	//					Required: true,
	//				},
	//				&cli.StringFlag{
	//					Name:        "nfile",
	//					Usage:       "指定存放整数 n 的文件的位置和名称",
	//					Aliases:     []string{"n"},
	//					Destination: &nfile,
	//					Required: true,
	//				},
	//				&cli.StringFlag{
	//					Name:        "efile",
	//					Usage:       "在数据加密时，指定存放整数 e 的文件的位置和名称",
	//					Aliases:     []string{"e"},
	//					Destination: &efile,
	//					Required: true,
	//				},
	//				&cli.StringFlag{
	//					Name:        "cipherfile",
	//					Usage:       "指定密文文件的位置和名称",
	//					Aliases:     []string{"c"},
	//					Destination: &cipherfile,
	//					Required: true,
	//				},
	//			},
	//		},
	//		{
	//			Name: 	"sign",
	//			Usage: 	"对指定文件进行签名",
	//			Action: signAction,
	//			Flags: []cli.Flag{
	//				&cli.StringFlag{
	//					Name:        "plainfile",
	//					Usage:       "指定明文文件的位置和名称",
	//					Aliases:     []string{"p"},
	//					Destination: &plainfile,
	//					Required: true,
	//				},
	//				&cli.StringFlag{
	//					Name:        "nfile",
	//					Usage:       "指定存放整数 n 的文件的位置和名称",
	//					Aliases:     []string{"n"},
	//					Destination: &nfile,
	//					Required: true,
	//				},
	//				&cli.StringFlag{
	//					Name:        "dfile",
	//					Usage:       "在数据签名时，指定存放整数 e 的文件的位置和名称",
	//					Aliases:     []string{"d"},
	//					Destination: &efile,
	//					Required: true,
	//				},
	//				&cli.StringFlag{
	//					Name:        "cipherfile",
	//					Usage:       "指定签名文件的位置和名称",
	//					Aliases:     []string{"c"},
	//					Destination: &cipherfile,
	//					Required: true,
	//				},
	//			},
	//		},
	//	},
	//	Action: func(ctx *cli.Context) (err error) {
	//		_ = ctx.App.Command("help").Action(ctx)
	//		_ = afterAction(ctx)
	//		return
	//	},
	//}
	//
	//err := app.Run(os.Args)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//rsa := RSA.NewCheck("73299B42DBD959CDB3FB176BD1", "10001")
	//fmt.Println(rsa.Encrypt("63727970746F677261706879"))
	s := ReadHex("rsa_n.txt")
	fmt.Println(s)
}
