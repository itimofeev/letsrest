package main

import (
	"github.com/itimofeev/letsrest"
)

func main() {
	framework := letsrest.IrisHandler(letsrest.NewHTTPRequester(), letsrest.NewRequestStore())
	framework.Listen(":8080")
}
