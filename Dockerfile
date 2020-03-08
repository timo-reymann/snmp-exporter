FROM golang:1.14-alpine as auth-builder
WORKDIR /workspace
COPY snmp_auth.go main.go
RUN go build -o entrypoint main.go

FROM golang:1.14-alpine as exporter-builder
WORKDIR /workspace
RUN apk add --no-cache git
RUN git clone --depth=1 https://github.com/prometheus/snmp_exporter .
RUN go build -o exporter .

FROM alpine
COPY --from=auth-builder /workspace/entrypoint /bin/entrypoint
COPY --from=exporter-builder /workspace/exporter /bin/snmp_exporter
RUN mkdir -p /etc/snmp_exporter
ENTRYPOINT ["/bin/entrypoint"]

