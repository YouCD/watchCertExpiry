# watchCertExpiry
[![Build Status](https://travis-ci.com/YouCD/watchCertExpiry.svg?branch=main)](https://travis-ci.com/YouCD/watchCertExpiry)

`watchCertExpiry`是golang编写的一个https站点证书过期时间监控的prometheus exporter
# 使用
   
* 定义环境变量
多个站点用`,`号隔开
`CERT_MONITOR_URL_LIST=github.com,www.baidu.com`
   
```shell
./watchCertExpiry
```
查看
```shell
curl 127.0.0.1:8080/metrics
...
# HELP website_ssl_earliest_cert_expiry websie ssl cert expiry timestamp
# TYPE website_ssl_earliest_cert_expiry gauge
website_ssl_earliest_cert_expiry{website="https://github.com"} 1.652184e+09
....

```
# docker 

```shell
go build -o bin/watchCertExpiry .

docker build . -t watch_cert_expiry

docker run -d --name watch_cert_expiry -e CERT_MONITOR_URL_LIST=github.com,www.baidu.com -p 8080:8080 watch_cert_expiry
```

