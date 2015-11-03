package main

import (
	"fmt"
	"github.com/ipedrazas/botd/api"
	"net/http"
)

func main() {
	fmt.Println("Server starting")
	http.ListenAndServe(":9090", api.Handlers())
}
