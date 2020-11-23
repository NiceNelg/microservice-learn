#!/bin/bash

# 启动micro api网关 
# --namespace 指定服务的命名空间
# --type 指定服务访问方式，api：暴露api服务，用于面向公众提供服务；service：暴露后端内部服务，使其可以通过http://{host}:{post}/{serverName}/{apiName}/{methodName}的方式访问，非常规模式
micro api --namespace=go.micro --type=service --handler=http
