FROM golang:1.18-alpine as builder

COPY . /image-pipeline-server

RUN apk update \
  && apk add git make curl jq \
  && cd /image-pipeline-server \
  && make build

FROM alpine 

EXPOSE 7001

# tzdata 安装所有时区配置或可根据需要只添加所需时区

RUN addgroup -g 1000 go \
  && adduser -u 1000 -G go -s /bin/sh -D go \
  && apk add --no-cache ca-certificates tzdata

COPY --from=builder /image-pipeline-server/image-pipeline-server /usr/local/bin/image-pipeline-server
COPY --from=builder /image-pipeline-server/entrypoint.sh /entrypoint.sh

USER go

WORKDIR /home/go

HEALTHCHECK --timeout=10s --interval=10s CMD [ "wget", "http://127.0.0.1:7001/ping", "-q", "-O", "-"]

CMD ["image-pipeline-server"]

ENTRYPOINT ["/entrypoint.sh"]
