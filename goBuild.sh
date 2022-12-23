#!/bin/bash

# 构建docker image
docker build -f Dockerfile -t golang1.17 .

## 删除旧容器1
docker stop GINCHAT
docker rm GINCHAT

# 构建docker container
docker run -itd -p 8081:8081 -v /Users/liaoshaoyu/GoProject/GINCHAT:/build --name GINCHAT golang1.17