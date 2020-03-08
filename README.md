snmp-exporter
===
[![GitHub Release](https://img.shields.io/github/v/release/timo-reymann/snmp-exporter.svg?label=version)](https://github.com/timo-reymann/snmp-exporter/releases)
![buildx](https://github.com/timo-reymann/snmp-exporter/workflows/buildx/badge.svg)

Docker-Image for the [snmp_exporter](https://github.com/prometheus/snmp_exporter) with added security token.

Intended for usage with nas devices that have a docker integration

## Usage

### Generate snmp.yml

Generate your `snmp.yml` using the generator, to do that place a `generator.yml` in the current directory:

```bash
docker run --rm -it -v "${PWD}:/opt" prom/snmp-generator generate
```

More info in the [official repo](https://github.com/prometheus/snmp_exporter/tree/master/generator).

### Spin up exporter

Place your `snmp.yml` in the current directory.

```bash
docker run --rm  -v $PWD:/etc/snmp_exporter  -it -p 3000:3000 timoreymann/snmp_exporter:latest-arm64v8 \
	-token yourToken
```

Prometheus config:

```yaml
  - job_name: nas-exporter
    scrape_interval: 35s
    scrape_timeout: 30s
    static_configs:
      - targets: [exporterHost:3000]
    params:
      token: [yourToken]
      module: [module_name]
      target: [targetAddress]

```

