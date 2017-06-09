package main

import (
	"github.com/itimofeev/letsrest"
)

func main() {
	framework := letsrest.IrisHandler(letsrest.NewRequestStore(letsrest.NewHTTPRequester()))
	framework.Listen(":8080")
}
