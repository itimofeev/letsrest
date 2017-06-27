package main

import (
	"github.com/itimofeev/letsrest"
)

func main() {
	pool := letsrest.NewWorkerPool(letsrest.NewHTTPRequester())
	framework := letsrest.IrisHandler(letsrest.NewDataStore(pool))
	framework.Listen(":8080")
}
