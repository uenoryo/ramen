package slack

import (
	"fmt"

	libslack "github.com/nlopes/slack"
)

type Config struct {
	Token   string
	BotName string
}

type Client struct {
	BotName          string
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
				if ev.Username == cli.BotName {
					continue
				}

				cli.OnReceiveMessage(&Message{ev})
			}
		}
	}
}

func onConnectedDefault() {
	fmt.Println("connected")
}

func onReceiveMessageDefault(msg *Message) {
	fmt.Printf("received message:%s\n", msg.Text)
}
