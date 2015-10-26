package main

import (
    "github.com/ipedrazas/botd/api"
    "fmt"
    "net/http"
)

func main() {
    fmt.Println("Server starting")
    http.ListenAndServe(":8008", api.Handlers())
}
