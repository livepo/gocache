package server

import (
    "fmt"
    "gocache/global"
    "net"
    "bufio"
    "strings"
    "errors"
    "sync"
)


type Server struct {
    ListenAddr string
    AuthToken string
    LRU *LRUCache
    Lock *sync.RWMutex
}


func NewServer(listenaddr string, authtoken string) *Server {
    return &Server{
        ListenAddr: listenaddr,
        AuthToken: authtoken,
        LRU: NewLRUCache(global.CAPACITY),
        Lock: &sync.RWMutex{},
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
    var state = "login"
    for state == "login" {
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
            state = "cmd"
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
                    s.outputMsg(client, value)
                }
            }
        }
    }
}


func (s *Server) outputMsg(client net.Conn, str string) {
    fmt.Println(str)
    client.Write([]byte(str + "\n"))
}


func (s *Server) putCmd(key string, value string) error {
    s.Lock.Lock()
    s.Data[key] = value
    s.Lock.Unlock()
    return nil
}


func (s *Server) getCmd(key string) (string, error) {
    s.Lock.RLock()
    defer s.Lock.RUnlock()
    if v, ok := s.Data[key]; ok {
        return v, nil
    } else {
        return "", errors.New(fmt.Sprintf("key %s not found", key))
    }
}
