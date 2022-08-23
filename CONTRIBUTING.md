## Adding new components
**Before making a contribution please open an issue describing the new addition and the need for it**
* Fork the repo
* Create a new branch for the new components
* Add your new component in `./logzio/<component_type>` (exporter, receiver...)
* Add `Makefile` for your component, you can check out `logzioexporter` makefile for example
* Go to `./otelbuilder/otelcol-builder.yaml` and add the path to your component. Example:
```yaml
exporters:
    # New exporter component
  - gomod: "github.com/logzio/otel-collector-distro/logzio/exporter/logzioexporter v0.0.1"
    path: ./../logzio/exporter/logzioexporter
```
* For testing the build you can go to `./otelbuilder` and run `make otelcol-logzio-<<os>>_<<arch>>` ( macOS example: `otelcol-logzio-darwin_amd64` ). This command will build the binary for your machine. Now you can run and test the usage.
* Push changes to your branch and open a PR
