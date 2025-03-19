FROM golang:alpine AS builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLE=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o rforum

FROM scratch

WORKDIR /rforum

COPY --from=builder ./build/rforum ./
COPY ./settings ./settings

EXPOSE 8888

ENTRYPOINT [ "./rforum", "-f", "./settings" ]
