# otel-collector-distro
Logz.io distribution of opentelemery collector

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
