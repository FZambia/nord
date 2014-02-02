package main

import (
	"net/http"

	"github.com/FZambia/nord/libnord"
)

func main() {
	config := parseCommandLineOptions()
	m := libnord.GetHandler(config.Config)
	http.ListenAndServe(config.Addr, m)
}
