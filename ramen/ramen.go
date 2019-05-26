package ramen

import (
	"fmt"

	"github.com/uenoryo/ramen/slack"
)

type Config struct {
	slack.Config
}

type Ramen struct {
	client *slack.Client
}

func New(cnf Config) *Ramen {
	slackCnf := slack.Config{
		BotName: cnf.BotName,
		Token:   cnf.Token,
	}
	client := slack.New(slackCnf)
	client.OnReceiveMessage = func(msg *slack.Message) {
		client.Post(msg.Channel, fmt.Sprintf("「%s」受信したよ", msg.Text))
	}

	return &Ramen{
		client: client,
	}
}

// Run (๑•̀ㅂ•́)و ｸﾞｯ
func (rmn Ramen) Run() error {
	rmn.client.Connect()
	rmn.client.Run()

	return nil
}
