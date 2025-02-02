package main

import (
	"fmt"
	"net/http"
)

var config = Config{}
var router = ConfigureRouter()

func main() {
	config.init()

	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), router)
}
