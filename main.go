package main

import (
	"io/ioutil"

	"github.com/uenoryo/ramen/ramen"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	buf, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}

	cnf := ramen.Config{}
	err = yaml.Unmarshal(buf, &cnf)
	if err != nil {
		panic(err)
	}

	ramen := ramen.New(cnf)
	ramen.Run()
}
