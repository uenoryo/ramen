package ramen

import (
	"testing"

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
