package main

import (
	"github.com/itimofeev/letsrest"
)

func main() {
	framework := letsrest.IrisHandler(letsrest.NewDataStore(letsrest.NewHTTPRequester()))
	framework.Listen(":8080")
}
