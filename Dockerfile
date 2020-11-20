FROM ccr.ccs.tencentyun.com/astatium.com/alpine:3.12-arm64
WORKDIR /app
COPY main /app/dynamic-ip-alpine
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

CMD ["/app/dynamic-ip-alpine"]