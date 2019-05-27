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
	ErrInvalidRemindTime = errors.New("error invalid remind date time")
	ErrRemindTimeIsPast  = errors.New("error remind time is past")

	ReminderCheckIntervalSec = 30 * time.Second
)

type Config struct {
	slack.Config `yaml:"slack"`
}

type Ramen struct {
	client  *slack.Client
	storage storage.Storage
	nowFunc func() time.Time
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
	ramen.nowFunc = func() time.Time {
		return time.Now()
	}
	return ramen
}

// Run (๑•̀ㅂ•́)و ｸﾞｯ
func (rmn *Ramen) Run() error {
	if err := rmn.storage.Load(); err != nil {
		return errors.Wrap(err, "load storage failed.")
	}

	rmn.Watch()
	rmn.client.FetchUsers()
	rmn.client.Connect()
	rmn.client.Run()

	return nil
}

func (rmn *Ramen) Watch() {
	ticker := time.NewTicker(ReminderCheckIntervalSec)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := rmn.remindIfExists(); err != nil {
				log.Println("[ERROR] remind failed, message:", err.Error())
			}
		}
	}
}

func (rmn *Ramen) remindIfExists() error {
	now := rmn.nowFunc()
	for _, record := range rmn.storage.Data() {
		if now.Equal(record.RemindAt) || now.Before(record.RemindAt) {
			rmn.storage.Delete(record.ID)
			time.Sleep(2 * time.Second)
		}
	}
	return nil
}

func (rmn *Ramen) receiveAndReply(msg *slack.Message) {
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

	remindAt, err := rmn.strToTime(remindDate, remindTime)
	switch {
	case err == ErrInvalidRemindTime:
		rmn.client.Post(msg.Channel, "時間がよくわかりません...")
		return
	case err == ErrRemindTimeIsPast:
		rmn.client.Post(msg.Channel, "過ぎ去った時間には戻れません")
		return
	case err != nil:
		rmn.client.Post(msg.Channel, "エラーが発生しました: "+err.Error())
		return
	}

	record := &storage.Record{
		ID:        rmn.genID(remindAt.String() + msg.Text),
		UserID:    msg.User,
		Content:   content,
		CreatedAt: time.Now(),
		RemindAt:  remindAt,
	}
	if err := rmn.storage.Save(record); err != nil {
		log.Println("保存時にエラーが発生", err.Error())
	}

	rmn.client.Post(msg.Channel, fmt.Sprintf("<@%s> りょーかいです！ %s にリマインドしますね！", record.UserID, remindAt.Format("2006/01/02 15:04")))
}

func (rmn *Ramen) analysis(text string) (to, date, time, content string, err error) {
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

func (rmn *Ramen) isBotName(str string) bool {
	return str == fmt.Sprintf("@%s", rmn.client.BotName)
}

func (rmn *Ramen) isDate(str string) bool {
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

func (rmn *Ramen) isTime(str string) bool {
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

func (rmn *Ramen) strToTime(remindDate, remindTime string) (time.Time, error) {
	var (
		now        = rmn.nowFunc()
		remindYear = now.Format("2006")
	)
	if remindDate == "" {
		remindDate = now.Format("01/02")
	}

	dtStr := fmt.Sprintf("%s %s %s JST", remindYear, remindDate, remindTime)
	remindAt, err := time.Parse("2006 1/2 15:4 MST", dtStr)
	if err != nil {
		return time.Time{}, ErrInvalidRemindTime
	}
	remindMonth := remindAt.Format("1")

	if remindAt.Equal(now) || remindAt.Before(now) {
		// NOTE: 「年が替わる2ヶ月前」 且つ 「1月, 2月」を指定していた場合は来年を指していると判定する仕様
		isNearNextYaer := now.AddDate(0, 2, 0).Format("2006") != now.Format("2006")
		if isNearNextYaer && (remindMonth == "1" || remindMonth == "2") {
			remindAt = remindAt.AddDate(1, 0, 0)
		} else {
			return time.Time{}, ErrRemindTimeIsPast
		}
	}
	return remindAt, nil
}

func (rmn *Ramen) genID(str string) string {
	sha := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", sha[:15])
}
