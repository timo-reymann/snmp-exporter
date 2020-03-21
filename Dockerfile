FROM golang:1.14-alpine as auth-builder
ENV CGO_ENABLED=0
WORKDIR /workspace
COPY snmp_auth.go main.go
RUN go build -a -ldflags '-s' -o entrypoint main.go

FROM golang:1.14-alpine as exporter-builder
ENV CGO_ENABLED=0
WORKDIR /workspace
RUN apk add --no-cache git
RUN git clone --depth=1 https://github.com/prometheus/snmp_exporter .
RUN go build -a -ldflags '-s' -o exporter .

FROM scratch
COPY --from=auth-builder /workspace/entrypoint /bin/entrypoint
COPY --from=exporter-builder /workspace/exporter /bin/snmp_exporter
VOLUME ["/etc/snmp_exporter"]
ENTRYPOINT ["/bin/entrypoint"]
