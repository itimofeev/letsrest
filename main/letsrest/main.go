package main

import (
	"github.com/itimofeev/letsrest"
)

func main() {
	config := letsrest.ReadConfigFromEnv()
	pool := letsrest.NewWorkerPool(letsrest.NewHTTPRequester())
	framework := letsrest.IrisHandler(letsrest.NewDataStore(config, pool))
	framework.Listen(":8080")
}
