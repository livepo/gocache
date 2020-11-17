package main


import (
    "flag"
    "gocache/server"
)


func main() {
    listenaddr := flag.String("listen", ":2020", "listen address")
    authtoken := flag.String("auth", "loginuser", "auth login user")
    flag.Parse()
    serv := server.NewServer(*listenaddr, *authtoken)
    serv.Serve()
}
