package ramen

import (
	"fmt"
	"testing"
	"time"

	"github.com/uenoryo/ramen/slack"
)

func Test_isDate(t *testing.T) {
	t.Parallel()

	type Test struct {
		Input  string
		Expect bool
	}

	cases := []Test{
		{
			Input:  "",
			Expect: false,
		},
		{
			Input:  "10",
			Expect: false,
		},
		{
			Input:  "10/",
			Expect: false,
		},
		{
			Input:  "10/0",
			Expect: true,
		},
		{
			Input:  "1/00",
			Expect: true,
		},
		{
			Input:  "1/0",
			Expect: true,
		},
		{
			Input:  "10/00",
			Expect: true,
		},
		{
			Input:  "110/00",
			Expect: false,
		},
		{
			Input:  "100/100",
			Expect: false,
		},
	}

	for _, test := range cases {
		ramen := &Ramen{}
		if g, w := ramen.isDate(test.Input), test.Expect; g != w {
			t.Errorf("error isDate. input %s, got %v, want %v", test.Input, g, w)
		}
	}
}

func Test_isTime(t *testing.T) {
	t.Parallel()

	type Test struct {
		Input  string
		Expect bool
	}

	cases := []Test{
		{
			Input:  "",
			Expect: false,
		},
		{
			Input:  "10",
			Expect: false,
		},
		{
			Input:  "10:",
			Expect: false,
		},
		{
			Input:  "10:0",
			Expect: true,
		},
		{
			Input:  "1:00",
			Expect: true,
		},
		{
			Input:  "1:0",
			Expect: true,
		},
		{
			Input:  "10:00",
			Expect: true,
		},
		{
			Input:  "110:00",
			Expect: false,
		},
		{
			Input:  "100:100",
			Expect: false,
		},
	}

	for _, test := range cases {
		ramen := &Ramen{}
		if g, w := ramen.isTime(test.Input), test.Expect; g != w {
			t.Errorf("error isTime. input %s, got %v, want %v", test.Input, g, w)
		}
	}
}

func Test_isBotName(t *testing.T) {
	t.Parallel()

	type Test struct {
		Input  string
		Expect bool
	}

	cases := []Test{
		{
			Input:  "@testbot",
			Expect: true,
		},
		{
			Input:  "testbot",
			Expect: false,
		},
	}

	for _, test := range cases {
		ramen := &Ramen{
			client: &slack.Client{
				BotName: "testbot",
			},
		}
		if g, w := ramen.isBotName(test.Input), test.Expect; g != w {
			t.Errorf("error isBotName. input %s, got %v, want %v", test.Input, g, w)
		}
	}
}

func Test_analysis(t *testing.T) {
	t.Parallel()

	type Expect struct {
		To      string
		Date    string
		Time    string
		Content string
		Error   error
	}

	type Test struct {
		Input  string
		Expect Expect
	}

	cases := []Test{
		{
			Input: "@testbot 12/10 10:00 eat ramen",
			Expect: Expect{
				To:      "@testbot",
				Date:    "12/10",
				Time:    "10:00",
				Content: "eat ramen",
				Error:   nil,
			},
		},
		{
			Input: "@testbot 10:00 eat ramen",
			Expect: Expect{
				To:      "@testbot",
				Date:    "",
				Time:    "10:00",
				Content: "eat ramen",
				Error:   nil,
			},
		},
		{
			Input: "    @testbot    12/10     10:00     eat ramen   ",
			Expect: Expect{
				To:      "@testbot",
				Date:    "12/10",
				Time:    "10:00",
				Content: "eat ramen",
				Error:   nil,
			},
		},
		{
			Input: "testbot 12/10 10:00 eat ramen", // @ がない
			Expect: Expect{
				Error: ErrMissingBotBame,
			},
		},
		{
			Input: "@testbot 12/100 10:00 eat ramen", // 日付が変
			Expect: Expect{
				Error: ErrMissingRemindTime,
			},
		},
		{
			Input: "@testbot 12/10 :00 eat ramen", // 時間が変
			Expect: Expect{
				Error: ErrMissingRemindTime,
			},
		},
	}

	for _, test := range cases {
		ramen := &Ramen{
			client: &slack.Client{
				BotName: "testbot",
			},
		}
		to, date, time, content, err := ramen.analysis(test.Input)
		if err != test.Expect.Error {
			t.Fatalf("error is not match. got %v, want %v", err, test.Expect.Error)
		}
		if test.Expect.Error != nil {
			continue
		}

		if to != test.Expect.To {
			t.Errorf("error result to. got %s, want %s", to, test.Expect.To)
		}
		if date != test.Expect.Date {
			t.Errorf("error result date. got %s, want %s", date, test.Expect.Date)
		}
		if time != test.Expect.Time {
			t.Errorf("error result time. got %s, want %s", time, test.Expect.Time)
		}
		if content != test.Expect.Content {
			t.Errorf("error result content. got %s, want %s", content, test.Expect.Content)
		}
	}
}

func Test_strToTime(t *testing.T) {
	t.Parallel()

	type Input struct {
		Date string
		Time string
	}

	type Test struct {
		Input  Input
		Expect time.Time
		Error  error
	}

	strToTime := func(str string) time.Time {
		dt, err := time.Parse("2006-01-02 15:04:05 MST", str)
		if err != nil {
			t.Fatal("error parse time", err.Error())
		}
		return dt
	}

	testNowFunc := func() time.Time {
		return strToTime("2019-11-10 12:00:00 JST")
	}

	cases := []Test{
		{
			Input: Input{
				Date: "12/15",
				Time: "10:15",
			},
			Expect: strToTime("2019-12-15 10:15:00 JST"),
		},
		{
			Input: Input{
				Date: "11/12",
				Time: "10:40",
			},
			Expect: strToTime("2019-11-12 10:40:00 JST"),
		},
		{
			Input: Input{
				Date: "2/5", // 年跨ぎ
				Time: "3:4",
			},
			Expect: strToTime("2020-02-05 03:04:00 JST"),
		},
		{
			Input: Input{
				Date: "10/10", // 過去
				Time: "11:59",
			},
			Error: ErrRemindTimeIsPast,
		},
		{
			Input: Input{
				Date: "1/1", // 省略形
				Time: "2:3",
			},
			Expect: strToTime("2020-01-01 02:03:00 JST"),
		},
		{
			Input: Input{
				Date: "", // 省略
				Time: "12:01",
			},
			Expect: strToTime("2019-11-10 12:01:00 JST"),
		},
		{
			Input: Input{
				Date: "", // 省略
				Time: "12:00",
			},
			Error: ErrRemindTimeIsPast,
		},
		{
			Input: Input{
				Date: "",
				Time: "",
			},
			Error: ErrInvalidRemindTime,
		},
		{
			Input: Input{
				Date: "12/32", // 存在しない日
				Time: "10:00",
			},
			Error: ErrInvalidRemindTime,
		},
		{
			Input: Input{
				Date: "12/1",
				Time: "24:00", // 存在しない時間
			},
			Error: ErrInvalidRemindTime,
		},
	}

	for _, test := range cases {
		t.Run(fmt.Sprintf("input %s %s", test.Input.Date, test.Input.Time), func(t *testing.T) {
			ramen := &Ramen{
				client: &slack.Client{
					BotName: "testbot",
				},
				nowFunc: testNowFunc,
			}

			res, err := ramen.strToTime(test.Input.Date, test.Input.Time)
			if err != test.Error {
				t.Fatalf("error is not match. got %v, want %v", err, test.Error)
			}
			if test.Error != nil {
				return
			}

			if res.Unix() != test.Expect.Unix() {
				t.Errorf("error strToTime, got %s, want %s", res.String(), test.Expect.String())
			}
		})
	}
}
