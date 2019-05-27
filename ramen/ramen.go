package ramen

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/uenoryo/ramen/slack"
	"github.com/uenoryo/ramen/storage"
)

var (
	ErrMissingBotBame    = errors.New("error missing bot name")
	ErrMissingRemindTime = errors.New("error missing remind time")
)

type Config struct {
	slack.Config `yaml:"slack"`
}

type Ramen struct {
	client  *slack.Client
	storage storage.Storage
}

func New(cnf Config) *Ramen {
	slackCnf := slack.Config{
		BotName:        cnf.BotName,
		BotDisplayName: cnf.BotDisplayName,
		Token:          cnf.Token,
	}
	client := slack.New(slackCnf)

	storage := storage.NewFileStorage()
	ramen := &Ramen{
		client:  client,
		storage: storage,
	}
	ramen.client.OnReceiveMessage = ramen.receiveAndReply
	return ramen
}

// Run (๑•̀ㅂ•́)و ｸﾞｯ
func (rmn Ramen) Run() error {
	if err := rmn.storage.Load(); err != nil {
		return errors.Wrap(err, "load storage failed.")
	}

	rmn.client.FetchUsers()
	rmn.client.Connect()
	rmn.client.Run()

	return nil
}

func (rmn Ramen) receiveAndReply(msg *slack.Message) {
	log.Println(msg)
	_, remindDate, remindTime, content, err := rmn.analysis(msg.Text)
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

	record := &storage.Record{
		ID:        rmn.genID(msg.Text),
		UserID:    msg.User,
		Content:   content,
		CreatedAt: time.Now(),
		RemindAt:  time.Now(),
	}
	if err := rmn.storage.Save(record); err != nil {
		log.Println("保存時にエラーが発生", err.Error())
	}

	rmn.client.Post(msg.Channel, fmt.Sprintf("<@%s> %s %s に「%s」をリマインドしますね！", record.UserID, remindDate, remindTime, content))
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

func (rmn Ramen) genID(str string) string {
	sha := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", sha[:15])
}
