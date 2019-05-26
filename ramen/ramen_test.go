package ramen

import "testing"

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
