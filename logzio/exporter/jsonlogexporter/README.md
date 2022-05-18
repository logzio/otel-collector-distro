# Json log Exporter

This exporter supports sending log data over http as json byte array

**Note:** This exporter developed by logz.io and works best with logz.io as logs backend, but you can use it to send log data as json byte array to any backend that supports the format

### The following configuration options are supported:

* `custom_endpoint` (Optional): Custom endpoint (example: https://example-api.com)
* `token` (Optional): Your logz.io account token for your tracing account.
* `region` (Optional): Your logz.io account [region code](https://docs.logz.io/user-guide/accounts/account-region.html#available-regions). Defaults to `us`. Required only if your logz.io region is different than US.
* `drain_interval` (Optional): Queue drain interval in seconds. Defaults to `3`.
* `queue_capacity` (Optional): Queue capacity in bytes. Defaults to `20 * 1024 * 1024` ~ 20mb.
* `queue_max_length` (Optional): Max number of items allowed in the queue. Defaults to `500000`.

### Examples
Sending log data to logz.io:
```yaml
exporters:
  jsonlog:
    token: "LOGZIOtraceTOKEN"
    region: "eu"
    drain_interval: 3
    queue_capacity: 5000000
    queue_max_length: 500000
```
Sending log data to http endpoint:
```yaml
exporters:
  jsonlog:
    custom_endpoint: "https://example-api.com"
    drain_interval: 3
    queue_capacity: 5000000
    queue_max_length: 500000
```
### Data conversion
This exporter converts individual `plog.LogRecord` in to `map[string]interface{}` that represents json object.

### Full configuration examples
Sending log data to logz.io:
```yaml
receivers:
  filelog:
    include: [ "/private/var/log/*.log" ] #macos system log directory
    attributes:
      type: myType #logzio type 

exporters:
  jsonlog:
    token: "logzio_token"
    region: "us"

extensions:
  pprof:
    endpoint: :1777
  zpages:
    endpoint: :55679
service:
  extensions: [pprof, zpages]
  pipelines:
    logs:
      receivers: [filelog]
      exporters: [jsonlog]
  telemetry:
    logs:
      level: "debug"
```
Sending log data to http endpoint:
```yaml
receivers:
  filelog:
    include: [ "/private/var/log/*.log" ] #macos system log directory

exporters:
  jsonlog:
    custom_endpoint: "https://example-api.com"

extensions:
  pprof:
    endpoint: :1777
  zpages:
    endpoint: :55679
service:
  extensions: [pprof, zpages]
  pipelines:
    logs:
      receivers: [filelog]
      exporters: [jsonlog]
  telemetry:
    logs:
      level: "debug"
```

