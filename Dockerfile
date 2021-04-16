FROM golang:1.16.3-alpine AS builder
WORKDIR /go/src/github.com/charleyzhu/FastDNS/
COPY . .
RUN apk add --no-cache gcc build-base
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build --ldflags "-extldflags -static"


FROM golang:1.16.3
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
WORKDIR /FastDNS

COPY --from=builder /go/src/github.com/charleyzhu/FastDNS/FastDNS .
COPY --from=builder /go/src/github.com/charleyzhu/FastDNS/config/config.yaml ./config/
RUN chmod 755 ./FastDNS
VOLUME /FastDNS/config

CMD ["./FastDNS"]