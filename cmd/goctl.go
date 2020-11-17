package main

import (
    "flag"
    "gocache/client"
)


func main() {
    remote := flag.String("remote", ":2020", "remote address to connect to")
    flag.Parse()
    cli := client.NewClient(*remote)
    cli.Serve()
}
