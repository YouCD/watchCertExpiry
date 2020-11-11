FROM alpine:latest

ENV TZ Asia/Shanghai
#
##apk 切换到国内源 && filebeat install
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update && apk add --no-cache ca-certificates bash tzdata libc6-compat \
    && apk add --no-cache --virtual=build-dependencies wget  \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && apk del build-dependencies \
    && rm -rf /var/cache/apk/* \
    && mkdir -p /app \
    && mkdir -p /app/conf

MAINTAINER YCD "hnyoucd@gmail.com"

ADD ./bin/watchCertExpiry /app


#暴露端口
EXPOSE 8080

#最终运行docker的命令
WORKDIR /app
CMD ["./watchCertExpiry"]