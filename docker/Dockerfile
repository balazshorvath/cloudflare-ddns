FROM golang:1.14 AS builder

WORKDIR /opt/cfddns
COPY . .
RUN ls -la
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o cfddns .

FROM alpine:latest

WORKDIR /opt/cfddns
COPY --from=builder /opt/cfddns/cfddns .

CMD ["./cfddns"]
