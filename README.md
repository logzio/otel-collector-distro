# Logzio opentelemetry collector distro
[![Go Report Card](https://goreportcard.com/badge/github.com/logzio/otel-collector-distro/logzio/exporter/logzioexporter)](https://goreportcard.com/report/github.com/logzio/otel-collector-distro/logzio/exporter/logzioexporter)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Logz.io distribution of the OpenTelemetry collector. It provides a simple and unified solution to collect, process, and ship telemetry data to logz.io

## Quick Start
#### macOS
Start collecting system logs from your macOS machine:
* Download the `otelcol-logzio-darwin_amd64` binary from the latest release
* Run `chmod +x otelcol-logzio-darwin_amd64`
* Create configuration file named `macos-logs.yml` you can find the appropriate configuration in [otel-config](/otel-config/macos-logs.yml)
* Run logzio otel collector distro:
```shell
export TRACING_TOKEN=<<logzio_logs_token>> 
export LOGZIO_REGION=<<logzio_region>>
export LOGZIO_TYPE=<<logzio_log_type>>
./otelcol-logzio-darwin_amd64 --config macos-logs.yml
```


## Configuration

The logz.io openTelemetry collector distro uses Standard openTelemetry configuration.
For the default and some example configs, see the [otel-config](/otel-config/) directory.
For general configuration help, see the [openTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).
#### Logz.io opentelemetry collector distro components

| Receiver                        | Processor                     | Exporter                           | Extensions             |
|---------------------------------|-------------------------------|------------------------------------|------------------------|
| otlpreceiver                    | attributesprocessor           | `logzioexporter`                   | ballastextension       |
| awscontainerinsightreceiver     | resourceprocessor             | `jsonlogexporter`                  | zpagesextension        |
| awsecscontainermetricsreceiver  | batchprocessor                | loggingexporter                    | bearertokenauthextension|
| awsxrayreceiver                 | memorylimiterprocessor        | otlpexporter                       | healthcheckextension   |
| carbonreceiver                  | probabilisticsamplerprocessor | fileexporter                       | oidcauthextension      |
| collectdreceiver                | metricstransformprocessor     | otlphttpexporter                   | pprofextension         |
| dockerstatsreceiver             | spanprocessor                 | prometheusexporter                 |                        |
| dotnetdiagnosticsreceiver       | filterprocessor               | prometheusremotewriteexporter      |                        |
| filelogreceiver                 | resourcedetectionprocessor    |                                    |                        |
| fluentforwardreceiver           | groupbyattrsprocessor         |                                    |                        |
| googlecloudspannerreceiver      | groupbytraceprocessor         |                                    |                        |
| hostmetricsreceiver             | routingprocessor              |                                    |                        |
| jaegerreceiver                  | spanmetricsprocessor          |                                    |                        |
| jmxreceiver                     | tailsamplingprocessor         |                                    |                        |
| journaldreceiver                |                               |                                    |                        |
| kafkametricsreceiver            |                               |                                    |                        |
| kafkareceiver                   |                               |                                    |                        |
| opencensusreceiver              |                               |                                    |                        |
| podmanreceiver                  |                               |                                    |                        |
| prometheusreceiver              |                               |                                    |                        |
| receivercreator                 |                               |                                    |                        |
| redisreceiver                   |                               |                                    |                        |
| sapmreceiver                    |                               |                                    |                        |
| signalfxreceiver                |                               |                                    |                        |
| simpleprometheusreceiver        |                               |                                    |                        |
| splunkhecreceiver               |                               |                                    |                        |
| statsdreceiver                  |                               |                                    |                        |
| syslogreceiver                  |                               |                                    |                        |
| tcplogreceiver                  |                               |                                    |                        |
| udplogreceiver                  |                               |                                    |                        |
| wavefrontreceiver               |                               |                                    |                        |
| windowsperfcountersreceiver     |                               |                                    |                        |
| zipkinreceiver                  |                               |                                    |                        |
| zookeeperreceiver               |                               |                                    |                        |
