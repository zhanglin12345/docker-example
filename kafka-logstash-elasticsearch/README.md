# kafka->logstash->elasticsearch

## 使用方法（mac系统）

1. 解压文件： tar zxvf kafka-logstash-elasticsearch.tar.gz

2. 进入目录：cd kafka-logstash-elasticsearch

3. 获取本机IP，放在环境变量里：export DOCKER_HOST_IP=$(./bin/get_local_ip)

4. 启动：docker-compose up -d

   注意：因./bin/get_local_ip 获取本机ip脚本工具在mac下build因此只支持mac操作系统，其它系统请将docker-compose.yml中的${DOCKER_HOST_IP}替换为本机ip


##  服务连接方式

以下为各服务本地开发环境的连接方式，均不需要进行账号密码认证。

```
zookeeper: 
IP端口：127.0.0.1:2181

kafka:
IP端口：127.0.0.1:9092   
topcic: app-visits-r1p1    
注意：topic是在docker-compose.yml中通过环境变量的形式创建的，KAFKA_CREATE_TOPICS: "app-visits-r1p1:1:1"

elasticsearch: 
IP端口：127.0.0.1:9200
elasticsearch indices: app-visits-r1p1-%{+YYYY.MM.dd}
                       每天会生成一个app-visits-r1p1-2018.12.11 这样日期变量格式的indice，可以在api中使用 app-visits-r1p1-* 通配符形式进行查询

elasticsearch中的indices是通过config/logstash-kafka-es.conf中的index => "app-visits-r1p1-%{+YYYY.MM.dd}"创建的
```



## python kafka producer  demo

```
#-*- coding:utf-8 -*-
from kafka import KafkaProducer
import json
import time
import random

'''
    生产者demo
    向app-visits-r1p1主题中循环写入测试json数据
    注意事项：要写入json数据需加上value_serializer参数，如下代码
'''
producer = KafkaProducer(
                            value_serializer=lambda v: json.dumps(v).encode('utf-8'),
                            bootstrap_servers=['127.0.0.1:9092']
                         )
for i in range(1,200):
    time.sleep(0.2)
    data={
        "url": "http://test.com/sss/asa111.html",
        "appId": "app_"+str(i),
        "channel": random.randint(0,9),
        "vin": i+100000,
        "date": int(time.time()),
        "status": random.randint(0,9)%2
    }
    producer.send('application-r1p1', data)
    producer.close()
```



### Elasticsearch curl demo

```
shaoqian.wang@CNshaoqian~ $ curl http://127.0.0.1:9200/app-visits-r1p1-*/_search?pretty
{
  "took" : 6,
  "timed_out" : false,
  "_shards" : {
    "total" : 5,
    "successful" : 5,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : 30,
    "max_score" : 1.0,
    "hits" : [
      {
        "_index" : "app-visits-r1p1-2018.12.11",
        "_type" : "doc",
        "_id" : "cjVLm2cBFquu6bj-5Got",
        "_score" : 1.0,
        "_source" : {
          "@version" : "1",
          "status" : 1,
          "vin" : 103004,
          "url" : "http://test.com/sss/asa111.html",
          "appId" : "app_3004",
          "date" : 1544498700,
          "channel" : 4,
          "@timestamp" : "2018-12-11T03:25:00.434Z"
        }
      },
      ...
```

### Reference
https://stackoverflow.com/questions/35418939/kafka-with-docker-dynamic-advertised-host-name
