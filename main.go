package main

import (
	"github.com/FZambia/nord/libnord"
	"net/http"
)

func main() {
	config := parseCommandLineOptions()
	m := libnord.GetHandler(config.Config)
	http.ListenAndServe(config.Addr, m)
}
