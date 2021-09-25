package main

import (
	"github.com/CrackedPoly/Crypto/seccommu/src/action"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	var (
		app = &cli.App{
			Name:                 "应用密码学综合实践-2019141440070-罗鉴",
			Usage:                "安全传输",
			EnableBashCompletion: true,
			Commands: cli.Commands{
				{
					Name:   "newcer",
					Usage:  "申请证书",
					Action: action.NewcerAction,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "childcer",
							Usage:    "指定子证书文件的位置和名称",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "childpvk",
							Usage:    "指定子证书密钥文件的位置和名称",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "rootcer",
							Usage:    "指定根证书文件的位置和名称",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "rootpvk",
							Usage:    "指定根证书密钥文件的位置和名称",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "CN",
							Usage:    "指定证书使用者的名称",
							Required: true,
						},
					},
				},
				{
					Name:   "send",
					Usage:  "生成随机数之后，加密消息，对消息签名，公钥加密随机数",
					Action: action.SendAction,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "msg",
							Usage:    "指定消息文件的位置",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "pvk",
							Usage:    "指定用于签名的私钥文件的位置",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "cer",
							Usage:    "指定用于加密随机数的证书文件的位置",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "msgout",
							Usage:    "指定存有加密消息的文件",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "keyout",
							Usage:    "指定存有加密随机数的文件",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "signedout",
							Usage:    "指定存有对散列值进行签名的文件",
							Required: true,
						},
					},
				},
				{
					Name:   "getkey",
					Usage:  "使用乙的私钥解出用于加密的随机数",
					Action: action.GetKeyAction,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "pvk",
							Usage:    "指定解密方的私钥",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "keyout",
							Usage:    "指定存有加密随机数的文件",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "keyrec",
							Usage:    "指定恢复随机数的文件",
							Required: true,
						},
					},
				},
				{
					Name:   "getmsg",
					Usage:  "使用恢复出来的随机数对消息密文进行解密",
					Action: action.GetMessageAction,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "msgout",
							Usage:    "指定存有加密消息的文件",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "msgrec",
							Usage:    "指定保存恢复消息的文件",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "keyrec",
							Usage:    "指定保存恢复随机数的文件",
							Required: true,
						},
					},
				},
				{
					Name:   "verify",
					Usage:  "验证甲的数字签名",
					Action: action.VerifyAction,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "signedout",
							Usage:    "指定存有对散列值进行签名的文件",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "msgrec",
							Usage:    "指定保存恢复消息的文件",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "cer",
							Usage:    "指定用于验证签名的证书文件",
							Required: true,
						},
					},
				},
				{
					Name:   "diff",
					Usage:  "检验原始消息和恢复消息是否一样",
					Action: action.DiffAction,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "msg",
							Usage:    "指定消息文件的位置",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "msgrec",
							Usage:    "指定保存恢复消息的文件",
							Required: true,
						},
					},
				},
			},
			Action: func(ctx *cli.Context) (err error) {
				_ = ctx.App.Command("help").Action(ctx)
				_ = action.AfterAction(ctx)
				return
			},
		}
	)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
