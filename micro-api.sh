#!/bin/bash

# 启动micro api网关 
micro --registry=etcd --registry_address=127.0.0.1:2379 api --namespace=go.micro --type=api --handler=http
