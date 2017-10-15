package main

import (
	"github.com/aletheia7/ul"
	"log"
)

func setup_log() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(ul.New_object(`godoc-proxy`, ``))
}
