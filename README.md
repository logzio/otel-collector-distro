# Logzio opentelemetry collector distro
[![Go Report Card](https://goreportcard.com/badge/github.com/logzio/otel-collector-distro/logzio/exporter/logzioexporter)](https://goreportcard.com/report/github.com/logzio/otel-collector-distro/logzio/exporter/logzioexporter)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![unit-tests](https://github.com/logzio/otel-collector-distro/actions/workflows/test-go.yml/badge.svg)](https://github.com/logzio/otel-collector-distro/actions/workflows/test-go.yml)
[![build](https://github.com/logzio/otel-collector-distro/actions/workflows/release-artifacts.yml/badge.svg)](https://github.com/logzio/otel-collector-distro/actions/workflows/release-artifacts.yml)


Logz.io distribution of the OpenTelemetry collector. It provides a simple and unified solution to collect, process, and ship telemetry data to logz.io

## Quick Start
### MacOs
Start collecting system logs from your macOS machine.

The script below will perform the following steps:
- Download the `otelcol-logzio-darwin_amd64` binary from the latest release
- Download Macos quickstart configuration file
- Configure the related environment variables and run the binary

```shell
curl -L  https://github.com/logzio/otel-collector-distro/releases/download/v0.66.0/otelcol-logzio-darwin_amd64 > otelcol-logzio-darwin_amd64
chmod +x otelcol-logzio-darwin_amd64
curl -L https://raw.githubusercontent.com/logzio/otel-collector-distro/master/otel-config/macos.yml > macos.yml
export LOGS_TOKEN=<<logzio_logs_token>> 
export LOGZIO_REGION=<<logzio_region>>
export LOGZIO_TYPE=<<logzio_log_type>>
./otelcol-logzio-darwin_amd64 --config macos.yml
```

### Ubuntu linux
Start collecting system logs from your ubuntu linux machine.

The script below will perform the following steps:
- Download the `otelcol-logzio-linux_amd64` binary from the latest release
- Download linux quickstart configuration file
- Configure the related environment variables and run the binary


```shell
curl -L  https://github.com/logzio/otel-collector-distro/releases/download/v0.66.0/otelcol-logzio-linux_amd64 > otelcol-logzio-linux_amd64
chmod +x otelcol-logzio-linux_amd64
curl -L https://raw.githubusercontent.com/logzio/otel-collector-distro/development/otel-config/linux.yml > linux.yml
export LOGS_TOKEN=<<logzio_logs_token>> 
export LOGZIO_REGION=<<logzio_region>>
export LOGZIO_TYPE=<<logzio_log_type>>
./otelcol-logzio-linux_amd64 --config linux.yml
```

### Docker
Start collecting traces with jaeger hotrod demo application.

The script below will perform the following steps:
- Download a `docker-compose.yml` file that will deploy:
  - `logzio-otel-collector` container
  - `hotrod` application container to generate traces
- Configure the related environment variables and run the compose file

```shell
curl -L  https://raw.githubusercontent.com/logzio/otel-collector-distro/master/otel-config/docker-compose.yml > docker-compose.yml
export TRACING_TOKEN=<<logzio_tracing_token>> 
export LOGZIO_REGION=<<logzio_region>>
docker compose up 
```
- Go to `http://localhost:8080` and start generating traces
- Check out your traces in logz.io


## Configuration

The logz.io openTelemetry collector distro uses Standard openTelemetry configuration.
For the default and some example configs, see the [otel-config](/otel-config/) directory.
For general configuration help, see the [openTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).
#### Logz.io opentelemetry collector distro components

| Receiver                       | Processor                     | Exporter                      | Extensions               |
|--------------------------------|-------------------------------|-------------------------------|--------------------------|
| otlpreceiver                   | attributesprocessor           | `logzioexporter`              | ballastextension         |
| awscontainerinsightreceiver    | resourceprocessor             | `jsonlogexporter`             | zpagesextension          |
| awsecscontainermetricsreceiver | batchprocessor                | loggingexporter               | bearertokenauthextension |
| awsxrayreceiver                | memorylimiterprocessor        | otlpexporter                  | healthcheckextension     |
| carbonreceiver                 | probabilisticsamplerprocessor | fileexporter                  | oidcauthextension        |
| collectdreceiver               | metricstransformprocessor     | otlphttpexporter              | pprofextension           |
| dockerstatsreceiver            | spanprocessor                 | prometheusexporter            |                          |
| dotnetdiagnosticsreceiver      | filterprocessor               | prometheusremotewriteexporter |                          |
| filelogreceiver                | resourcedetectionprocessor    |                               |                          |
| fluentforwardreceiver          | groupbyattrsprocessor         |                               |                          |
| googlecloudspannerreceiver     | groupbytraceprocessor         |                               |                          |
| hostmetricsreceiver            | routingprocessor              |                               |                          |
| jaegerreceiver                 | spanmetricsprocessor          |                               |                          |
| jmxreceiver                    | tailsamplingprocessor         |                               |                          |
| journaldreceiver               |                               |                               |                          |
| k8seventsreceiver              |                               |                               |                          |
| kafkametricsreceiver           |                               |                               |                          |
| kafkareceiver                  |                               |                               |                          |
| opencensusreceiver             |                               |                               |                          |
| podmanreceiver                 |                               |                               |                          |
| prometheusreceiver             |                               |                               |                          |
| receivercreator                |                               |                               |                          |
| redisreceiver                  |                               |                               |                          |
| sapmreceiver                   |                               |                               |                          |
| signalfxreceiver               |                               |                               |                          |
| simpleprometheusreceiver       |                               |                               |                          |
| splunkhecreceiver              |                               |                               |                          |
| statsdreceiver                 |                               |                               |                          |
| syslogreceiver                 |                               |                               |                          |
| tcplogreceiver                 |                               |                               |                          |
| udplogreceiver                 |                               |                               |                          |
| wavefrontreceiver              |                               |                               |                          |
| windowsperfcountersreceiver    |                               |                               |                          |
| windowseventlogreceiver        |                               |                               |                          |
| zipkinreceiver                 |                               |                               |                          |
| zookeeperreceiver              |                               |                               |                          |


