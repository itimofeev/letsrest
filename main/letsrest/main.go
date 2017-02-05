package main

import (
	"github.com/itimofeev/letsrest"
)

func main() {
	framework := letsrest.IrisHandler(&letsrest.HTTPRequester{}, letsrest.NewRequestStore())
	framework.Listen(":6111")
}
