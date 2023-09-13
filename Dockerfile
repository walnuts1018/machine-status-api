FROM golang:1.21 as builder
ENV ROOT=/build
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

COPY ./ ./
RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o main $ROOT/main.go && chmod +x ./main

FROM alpine:3
WORKDIR /app

COPY --from=builder /build/main ./
CMD ["./main"]
LABEL org.opencontainers.image.source=https://github.com/walnuts1018/machine-status-api
