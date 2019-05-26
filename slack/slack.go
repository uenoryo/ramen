package slack

import (
	"log"

	libslack "github.com/nlopes/slack"
)

type Client struct {
	client *libslack.Client
	rtm    *libslack.RTM
}

func New(token string) *Client {
	client := libslack.New(token)
	return &Client{
		client: client,
		rtm:    client.NewRTM(),
	}
}

func (cli *Client) Connect() {
	go cli.rtm.ManageConnection()
}

func (cli *Client) Run() {
	for {
		select {
		case msg := <-cli.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *libslack.ConnectedEvent:
				log.Println("connected")

			case *libslack.MessageEvent:
				user := ev.User
				text := ev.Text
				channel := ev.Channel

				if ev.Username == "testbot" {
					continue
				}

				cli.client.PostMessage(
					channel,
					libslack.MsgOptionText(text, false),
					libslack.MsgOptionUsername("testbot"),
				)
			}
		}
	}
}
