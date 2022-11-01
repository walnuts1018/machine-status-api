FROM golang:1.17.6 as builder
ENV ROOT=/build
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

COPY ./src ./
RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o main $ROOT/main.go && chmod +x ./main

FROM alpine:3.15
WORKDIR /app

COPY --from=builder /build/main ./
RUN id
CMD ["./main"]
LABEL org.opencontainers.image.source=https://github.com/walnuts1018/machine-status-api
