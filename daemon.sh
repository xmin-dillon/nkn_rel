#!/bin/bash

# 监控 进程状态
while :
do
    # 1s 1m 1h
    sleep 5s 
    
    # 查找进程id
    PID=$(ps -e|grep nknd|awk '{printf $1}')
    if [ -z $PID ]; then
        echo "NKN 客户端异常退出，重新启动 ..."
        #导入密码
        export NKN_WALLET_PASSWORD=yang208
        nohup ./nknd --no-nat >>/dev/null 2>>/dev/null  &
    else
        echo "NKN 客户端正常 $PID ..."
        continue
    fi
done
