#!/bin/bash

# 启动微服务监控UI界面，web版，访问地址：http://localhost:8082
# --registry 指定注册中心种类
# --registry_address 指定注册地址
micro --registry=etcd --registry_address=127.0.0.1:2379 web
