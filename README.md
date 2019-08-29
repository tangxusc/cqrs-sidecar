# cqrs-sidecar

cqrs的目标是在各种语言中实现:

1. 收到grpc推送的消息;
2. 开启事务;
3. 执行cqrs中查询侧视图更新;
4. 使用存储过程ack 消息(这样 在一个事务中就可以实现ack和更新视图,不再需要事务管理器,来分别管理mq和db事务)
5. 提交事务;

> 以上动作,顺序执行,出现错误直接回滚,grpc会重复推送错误消息

## 特性

- [x] pulsar接入(key_share订阅)
- [x] rpc推送event
- [x] proxy代理查询数据库
- [x] 存储过程实现ack event
- [ ] 重启后恢复推送未发送的event
- [ ] metrics实现
- [ ] grpc推送时如何不重复?
- [ ] opentracing接入

## 架构

![](./document/img_1.png)

## 参照

```shell
github.com/apache/pulsar/pulsar-client-go
github.com/go-sql-driver/mysql
github.com/golang/protobuf
github.com/sirupsen/logrus
github.com/spf13/cobra
google.golang.org/grpc
```

