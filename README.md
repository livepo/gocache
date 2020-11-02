gocache是一个用go写的kv系统


TODO:

- [ ] 缓存淘汰策略
- [ ] 分片
- [ ] 副本
- [ ] 分布式
- [ ] 租约
- [ ] 事务
- [ ] 主从同步


服务端:

go build cmd/gocache.go


客户端:

go build cmd/goctl.go


![](https://github.com/livepo/gocache/blob/master/20201102145120.png)
