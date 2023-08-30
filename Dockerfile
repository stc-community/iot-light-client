FROM golang:1.19-alpine AS builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY . .

RUN go build -o client .


FROM alpine AS final
WORKDIR /app
COPY --from=builder /build/client /app/

ENTRYPOINT ["/app/client"]
