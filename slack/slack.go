package slack

import (
	"fmt"

	libslack "github.com/nlopes/slack"
)

type Config struct {
	Token          string `yaml:"token"`
	BotName        string `yaml:"bot_name"`
	BotDisplayName string `yaml:"bot_display_name"`
}

type Client struct {
	BotName          string
	BotDisplayName   string
	client           *libslack.Client
	rtm              *libslack.RTM
	OnConnected      func()
	OnReceiveMessage func(*Message)
}

type Message struct {
	*libslack.MessageEvent
}

func New(cnf Config) *Client {
	client := libslack.New(cnf.Token)
	return &Client{
		BotName:          cnf.BotName,
		BotDisplayName:   cnf.BotDisplayName,
		client:           client,
		rtm:              client.NewRTM(),
		OnConnected:      onConnectedDefault,
		OnReceiveMessage: onReceiveMessageDefault,
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
				cli.OnConnected()

			case *libslack.MessageEvent:
				if ev.Username == cli.BotDisplayName {
					continue
				}

				cli.OnReceiveMessage(&Message{ev})
			}
		}
	}
}

func (cli *Client) Post(channel, text string) {
	cli.client.PostMessage(
		channel,
		libslack.MsgOptionText(text, false),
		libslack.MsgOptionUsername(cli.BotDisplayName),
	)
}

func onConnectedDefault() {
	fmt.Println("connected")
}

func onReceiveMessageDefault(msg *Message) {
	fmt.Printf("received message:%s\n", msg.Text)
}
