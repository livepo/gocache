package client

import (
    "fmt"
    "net"
    "os"
    "bufio"
    "strings"
)


type Client struct {
    remoteAddr string
    inputCh chan []string
    remoteCh chan []string
    remote net.Conn
}


func NewClient(remoteAddr string) *Client {
    remote, err := net.Dial("tcp", remoteAddr)
    if err != nil {
        fmt.Sprintf("connecting %s failed", remoteAddr)
        os.Exit(1)
    }
    inputCh := make(chan []string)
    remoteCh := make(chan []string)
    return &Client{
        remoteAddr: remoteAddr,
        inputCh: inputCh,
        remoteCh: remoteCh,
        remote: remote,
    }
}


func (c *Client) Close() {
    c.remote.Close()
    return
}

func (c *Client) parseCmd(reader *bufio.Reader, input chan []string) {
    for {
        cmd, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("read error", err)
            return
        }
        cmd = strings.TrimSpace(cmd)
        strarray := strings.Split(cmd, " ")
        purecmd := []string{}
        for _, s := range strarray {
            strimd := strings.Trim(s, " ")
            if len(strimd) > 0 {
                purecmd = append(purecmd, strimd)
            }
        }
        input <- purecmd
    }
}


func checkMsgLen(msg []string, length int) {
    if len(msg) != length {
        fmt.Println("msg length error", msg)
    }
}


func checkMsgLenGte(msg []string, length int) {
    if len(msg) < length {
        fmt.Println("msg length error", msg)
    }
}


func (c *Client) Serve() {
    inputCh, remoteCh := make(chan []string), make(chan []string)
    stdinReader := bufio.NewReader(os.Stdin)
    remoteReader := bufio.NewReader(c.remote)
    go c.parseCmd(stdinReader, inputCh)
    go c.parseCmd(remoteReader, remoteCh)

    fmt.Print("goctl> ")
    for {
        select {
        case msg := <-inputCh:
            if len(msg) == 0 {
                fmt.Print("goctl> ")
                continue
            }
            switch msg[0] {
            case "auth":
                checkMsgLen(msg, 2)
                c.Auth(msg[1])
            case "put":
                checkMsgLenGte(msg, 3)
                value := strings.Join(msg[2:], " ")
                c.Put(msg[1], value)
            case "get":
                checkMsgLen(msg, 2)
                c.Get(msg[1])
            case "logout":
                c.Close()
                c.Logout()
                return
            default:
                fmt.Println("cannot parse command")
            }
        case msg := <-remoteCh:
            fmt.Println(strings.Join(msg, " "))
            fmt.Print("goctl> ")
        }
    }
}


func (c *Client) Auth(logintoken string) {
    c.remote.Write([]byte(fmt.Sprintf("auth %s\n", logintoken)))
}


func (c *Client) Put(key string, value string) {
    c.remote.Write([]byte(fmt.Sprintf("put %s %s\n", key, value)))
}


func (c *Client) Get(key string) {
    c.remote.Write([]byte(fmt.Sprintf("get %s\n", key)))
}


func (c *Client) Logout() {
    c.remote.Write([]byte("logout\n"))
}
