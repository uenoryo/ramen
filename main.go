package main

import (
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
	"github.com/uenoryo/ramen/ramen"
	yaml "gopkg.in/yaml.v2"
)

const (
	configPath = "./config.yml"
)

func main() {
	if err := _main(); err != nil {
		log.Fatal("Whoops!! ", err.Error())
	}
}

func _main() error {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errors.Wrapf(err, "read config:%s failed", configPath)
	}

	cnf := ramen.Config{}
	if err := yaml.Unmarshal(buf, &cnf); err != nil {
		return errors.Wrap(err, "yaml unmershal failed")
	}

	ramen := ramen.New(cnf)

	if err := ramen.Run(); err != nil {
		return errors.Wrap(err, "ramen launch failed")
	}
	return nil
}
