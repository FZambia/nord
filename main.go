package main

import (
	"net/http"

	"github.com/FZambia/nord/libnord"
)

func main() {
	config := parseCommandLineOptions()
	m := libnord.GetHandler(config.Config)
	config.Config.Logger.Println("running Nord on", config.Addr)
	if err := http.ListenAndServe(config.Addr, m); err != nil {
		config.Config.Logger.Fatal(err)
	}
}
