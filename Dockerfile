FROM golang:1.15.2 as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo -o bot main.go

FROM alpine:latest
RUN apk update && apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /build/bot .
CMD ["./bot"]
