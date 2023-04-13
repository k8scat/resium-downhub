FROM golang:1.20.3-alpine3.17 AS builder
WORKDIR /data/resium-downhub
COPY . .
RUN apk add --no-cache git && \
    go mod download && \
    go build -trimpath -o bin/downhub main.go

FROM alpine:3.17
LABEL maintainer="K8sCat <k8scat@gmail.com>"
WORKDIR /data/resium-downhub
EXPOSE 8080
COPY --from=builder /data/resium-downhub/bin/downhub .
CMD ["./downhub"]
