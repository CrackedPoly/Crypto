package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var childcer, childpvk, rootcer, rootpvk, cn string
var msg, pvk, cer, msgout, keyout, signedout string
var keyrec, msgrec string

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
					Action: newcerAction,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "childcer",
							Usage:       "指定子证书文件的位置和名称",
							Destination: &childcer,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "childpvk",
							Usage:       "指定子证书密钥文件的位置和名称",
							Destination: &childpvk,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "rootcer",
							Usage:       "指定根证书文件的位置和名称",
							Destination: &rootcer,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "rootpvk",
							Usage:       "指定根证书密钥文件的位置和名称",
							Destination: &rootpvk,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "CN",
							Usage:       "指定证书使用者的名称",
							Destination: &cn,
							Required:    true,
						},
					},
				},
				{
					Name:   "send",
					Usage:  "生成随机数之后，加密消息，对消息签名，公钥加密随机数",
					Action: sendAction,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "msg",
							Usage:       "指定消息文件的位置",
							Destination: &msg,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "pvk",
							Usage:       "指定用于签名的私钥文件的位置",
							Destination: &pvk,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "cer",
							Usage:       "指定用于加密随机数的证书文件的位置",
							Destination: &cer,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "msgout",
							Usage:       "指定存有加密消息的文件",
							Destination: &msgout,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "keyout",
							Usage:       "指定存有加密随机数的文件",
							Destination: &keyout,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "signedout",
							Usage:       "指定存有对散列值进行签名的文件",
							Destination: &signedout,
							Required:    true,
						},
					},
				},
				{
					Name:   "getkey",
					Usage:  "使用乙的私钥解出用于加密的随机数",
					Action: getkeyAction,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "pvk",
							Usage:       "指定解密方的私钥",
							Destination: &pvk,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "keyout",
							Usage:       "指定存有加密随机数的文件",
							Destination: &keyout,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "keyrec",
							Usage:       "指定恢复随机数的文件",
							Destination: &keyrec,
							Required:    true,
						},
					},
				},
				{
					Name:   "getmsg",
					Usage:  "使用恢复出来的随机数对消息密文进行解密",
					Action: getMessageAction,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "msgout",
							Usage:       "指定存有加密消息的文件",
							Destination: &msgout,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "msgrec",
							Usage:       "指定保存恢复消息的文件",
							Destination: &msgrec,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "keyrec",
							Usage:       "指定保存恢复随机数的文件",
							Destination: &keyrec,
							Required:    true,
						},
					},
				},
				{
					Name:   "verify",
					Usage:  "验证甲的数字签名",
					Action: verifyAction,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "signedout",
							Usage:       "指定存有对散列值进行签名的文件",
							Destination: &signedout,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "msgrec",
							Usage:       "指定保存恢复消息的文件",
							Destination: &msgrec,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "cer",
							Usage:       "指定用于验证签名的证书文件",
							Destination: &cer,
							Required:    true,
						},
					},
				},
				{
					Name:   "diff",
					Usage:  "检验原始消息和恢复消息是否一样",
					Action: diffAction,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "msg",
							Usage:       "指定消息文件的位置",
							Destination: &msg,
							Required:    true,
						},
						&cli.StringFlag{
							Name:        "msgrec",
							Usage:       "指定保存恢复消息的文件",
							Destination: &msgrec,
							Required:    true,
						},
					},
				},
			},
			Action: func(ctx *cli.Context) (err error) {
				_ = ctx.App.Command("help").Action(ctx)
				_ = afterAction(ctx)
				return
			},
		}
	)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
