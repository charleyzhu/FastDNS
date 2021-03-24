FROM golang:1.15.6 AS builder
WORKDIR /go/src/github.com/charleyzhu/FastDNS/
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build


FROM golang:1.15.6
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
WORKDIR /FastDNS

COPY --from=builder /go/src/github.com/charleyzhu/FastDNS/FastDNS .
COPY --from=builder /go/src/github.com/charleyzhu/FastDNS/config/config.yaml ./config/
RUN chmod 755 ./FastDNS
VOLUME /FastDNS/config

CMD ["./FastDNS"]