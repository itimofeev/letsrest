package main

import (
	"github.com/itimofeev/letsrest"
)

func main() {
	config := letsrest.ReadConfigFromEnv()
	pool := letsrest.NewWorkerPool(letsrest.NewHTTPRequester())//пользователь
	framework := letsrest.IrisHandler(letsrest.NewDataStore(config, pool))//создаётся МонгоДС
	framework.Listen(":8080")
}
