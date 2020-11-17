package server

import (
    "bufio"
    "fmt"
    "gocache/driver"
    "gocache/global"
    "log"
    "net"
    "strings"
)

const (
    LOGIN = 1
    COMMAND = 2
)

type Server struct {
    ListenAddr string
    AuthToken string
    Cache driver.Cache
}



func NewServer(listenaddr string, authtoken string) *Server {
    var cache driver.Cache
    switch global.Config.Evict {
    case "simple":
        cache = driver.New(global.Config.Capacity).Simple().Build()
    case "lru":
        cache = driver.New(global.Config.Capacity).LRU().Build()
    case "lfu":
        cache = driver.New(global.Config.Capacity).LFU().Build()
    case "arc":
        cache = driver.New(global.Config.Capacity).ARC().Build()
    default:
        log.Fatal("error evict type")
    }
    return &Server{
        ListenAddr: listenaddr,
        AuthToken: authtoken,
        Cache: cache,
    }
}


func (s *Server) Serve() {
    ln, err := net.Listen("tcp", s.ListenAddr)
    if err != nil {
        fmt.Println("listen error", err)
        return
    }
    for {
        client, err := ln.Accept()
        if err != nil {
            fmt.Println("accept error", err)
        }
        go s.HandleClient(client)
    }
}


func (s *Server) readCmd(client net.Conn) ([]string, error) {
    reader := bufio.NewReader(client)
    cmd, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("read error", err)
        client.Close()
        return nil, err
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
    return purecmd, nil
}


func (s *Server) HandleClient(client net.Conn) {
    defer client.Close()
    var state = LOGIN
    for state == LOGIN {
        cmd, err := s.readCmd(client)
        if err != nil {
            fmt.Println("parse cmd error", err)
            return
        }
        if len(cmd) != 2 || cmd[0] != "auth" {
            s.outputMsg(client, "need authenticate!")
        }
        if cmd[1] != s.AuthToken {
            s.outputMsg(client, "authtoken error!")
        } else {
            s.outputMsg(client, "login succeed!")
            state = COMMAND
        }
    }

    for {
        cmd, err := s.readCmd(client)
        if err != nil {
            fmt.Println("read error", err)
            return
        }
        switch cmd[0] {
        case "logout":
            return
        case "put":
            if len(cmd) < 3 {
                s.outputMsg(client, "parse cmd error")
            } else {
                value := strings.Join(cmd[2:], " ")
                err := s.putCmd(cmd[1], value)
                if err != nil {
                    s.outputMsg(client, "put failed, try again")
                } else {
                    s.outputMsg(client, value)
                }
            }
        case "get":
            if len(cmd) != 2 {
                s.outputMsg(client, "parse cmd error")
            } else {
                value, err := s.getCmd(cmd[1])
                if err != nil {
                    fmt.Println("get error", err)
                    s.outputMsg(client, "get failed, try again")
                } else {
                    s.outputMsg(client, value.(string))
                }
            }
        }
    }
}


func (s *Server) outputMsg(client net.Conn, str string) {
    fmt.Println(str)
    client.Write([]byte(str + "\n"))
}


func (s *Server) putCmd(key string, value interface{}) error {
    return s.Cache.Set(key, value)
}


func (s *Server) getCmd(key string) (interface{}, error) {
    return s.Cache.Get(key)
}
