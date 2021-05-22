### demo已上传https://github.com/day112/gqlgen_example

`1   mkdir graphql_demo`

`2   cd graphql_demo`

`3   go mod init`

`4   go get github.com/99designs/gqlgen`

`5   gqlgen init 生成 graph 目录结构, 删除 graph 目录下文件， 并创建 schema.graphqls`

`6   再次执行  gqlgen init 生成新的文件`

`7   go run server.go`

# 消息中间件

# kafka demo

docker pull wurstmeister/zookeeper

docker run -d --name zookeeper -p 2181:2181 -v /etc/localtime:/etc/localtime wurstmeister/zookeeper

docker pull wurstmeister/kafka

docker run -d --name kafka -p 9092:9092 -e KAFKA_BROKER_ID=0 -e KAFKA_ZOOKEEPER_CONNECT=192.168.8.104:2181/kafka -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://192.168.8.104:9092 -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 -v /etc/localtime:/etc/localtime wurstmeister/kafka


# nats



