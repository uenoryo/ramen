package ramen

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/uenoryo/ramen/slack"
)

var (
	ErrMissingBotBame    = errors.New("error missing bot name")
	ErrMissingRemindTime = errors.New("error missing remind time")
)

type Config struct {
	slack.Config `yaml:"slack"`
}

type Ramen struct {
	client *slack.Client
}

func New(cnf Config) *Ramen {
	slackCnf := slack.Config{
		BotName:        cnf.BotName,
		BotDisplayName: cnf.BotDisplayName,
		Token:          cnf.Token,
	}
	client := slack.New(slackCnf)
	ramen := &Ramen{
		client: client,
	}
	ramen.client.OnReceiveMessage = ramen.receiveAndReply
	return ramen
}

// Run (๑•̀ㅂ•́)و ｸﾞｯ
func (rmn Ramen) Run() error {
	rmn.client.Connect()
	rmn.client.Run()

	return nil
}

func (rmn Ramen) receiveAndReply(msg *slack.Message) {
	log.Println(msg)
	_, date, time, content, err := rmn.analysis(msg.Text)
	switch err {
	case nil:
		break
	case ErrMissingBotBame:
		// 反応しない
		return
	case ErrMissingRemindTime:
		// TODO: エラーを返す
		return
	default:
		// TODO: エラーを返す
		return
	}

	rmn.client.Post(msg.Channel, fmt.Sprintf("%s %s に「%s」をリマインドしますね！", date, time, content))
}

func (rmn Ramen) analysis(text string) (to, date, time, content string, err error) {
	text = strings.TrimSpace(text)

	// 1個目の要素がBotNameかどうか
	sps := strings.Split(text, " ")
	if len(sps) == 0 || !rmn.isBotName(sps[0]) {
		err = ErrMissingBotBame
		return
	}
	to = sps[0]

	// BotNameを削除して詰める
	text = strings.Replace(text, sps[0], "", 1)
	text = strings.TrimSpace(text)

	// 2個目の要素が日付または時間かどうか
	sps = strings.Split(text, " ")

	switch {
	case len(sps) == 0:
		err = ErrMissingRemindTime
		return
	case rmn.isDate(sps[0]):
		date = sps[0]
	case rmn.isTime(sps[0]):
		time = sps[0]
	default:
		err = ErrMissingRemindTime
		return
	}

	// 日付または時間を削除して詰める
	text = strings.Replace(text, sps[0], "", 1)
	text = strings.TrimSpace(text)

	sps = strings.Split(text, " ")

	// 2個目が日付だった場合は3個目が時間であるはず
	if date != "" {
		switch {
		case len(sps) == 0:
			err = ErrMissingRemindTime
			return
		case rmn.isTime(sps[0]):
			time = sps[0]
		default:
			err = ErrMissingRemindTime
			return
		}

		// 時間を削除して詰める
		text = strings.Replace(text, sps[0], "", 1)
		text = strings.TrimSpace(text)
	}

	// 残りは文章
	content = text
	return
}

func (rmn Ramen) isBotName(str string) bool {
	return str == fmt.Sprintf("@%s", rmn.client.BotName)
}

func (rmn Ramen) isDate(str string) bool {
	sps := strings.Split(str, "/")
	if len(sps) != 2 {
		return false
	}
	if len(sps[0]) != 1 && len(sps[0]) != 2 {
		return false
	}
	if len(sps[1]) != 1 && len(sps[1]) != 2 {
		return false
	}
	return true
}

func (rmn Ramen) isTime(str string) bool {
	sps := strings.Split(str, ":")
	if len(sps) != 2 {
		return false
	}
	if len(sps[0]) != 1 && len(sps[0]) != 2 {
		return false
	}
	if len(sps[1]) != 1 && len(sps[1]) != 2 {
		return false
	}
	return true
}
