package main

import (
	"log"

	"github.com/uenoryo/ramen/ramen"
)

func main() {
	if err := ramen.Init(); err != nil {
		log.Printf("error ramen init, %s", err.Error())
	}
}
