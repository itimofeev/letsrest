package main

import (
	"github.com/itimofeev/letsrest"
)

func main() {
	framework, srv := letsrest.IrisHandler(letsrest.NewHTTPRequester(), letsrest.NewRequestStore())
	go srv.ListenForTasks()
	framework.Listen(":8080")
}
