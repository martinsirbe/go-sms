package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/martinsirbe/go-sms/pkg/sms"
)

const (
	appName        = "go-sms"
	appDescription = "Short Message Service (SMS) text message sender using AWS Simple Notification Service."
)

func main() {
	app := cli.App(appName, appDescription)

	senderID := app.String(cli.StringOpt{
		Name: "sender-id",
		Desc: "The sender ID which will appear on the receiver's device. (Optional, " +
			"if provided will override sender ID provided via configuration file.)",
		EnvVar: "SENDER_ID",
	})
	receiver := app.String(cli.StringOpt{
		Name:   "receiver",
		Desc:   "The receiver mobile phone number (in the E.164 format). (Mandatory)",
		EnvVar: "RECEIVER",
	})
	message := app.String(cli.StringOpt{
		Name:   "message",
		Desc:   "The text message you wish to send. (Mandatory)",
		EnvVar: "MESSAGE",
	})
	configPath := app.String(cli.StringOpt{
		Name:   "config-path",
		Desc:   "The path to the configurations file. (Mandatory)",
		EnvVar: "GO_SMS_CONFIG_PATH",
	})

	app.Action = func() {
		configFile, err := os.ReadFile(*configPath)
		if err != nil {
			log.WithError(err).Fatal("failed to load go-sms config.yaml")
		}

		var config sms.Config
		if err = yaml.Unmarshal(configFile, &config); err != nil {
			log.WithError(err).Fatal("failed to unmarshal go-sms config.yaml")
		}

		sender := sms.New(config)
		if senderID != nil {
			sender.WithSenderID(*senderID)
		}

		id, err := sender.Send(*message, *receiver)
		if err != nil {
			log.WithError(err).Fatal("failed to send the text message")
		}

		log.Infof("successfully sent a text message to - %s, message id - %s", *receiver, *id)
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Panicf("app failed to run")
	}
}
