package main


import "gocache/client"


func main() {
    cli := client.NewClient(":2021")
    cli.Serve()
}
