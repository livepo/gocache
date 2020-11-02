package main


import (
    "gocache/server"
)


func main() {
    serv := server.NewServer(":2021", "loginuser")
    serv.Serve()
}
