FROM golang:1.23.9-alpine as builder
ADD . /src
WORKDIR /src/otelbuilder/
ENV CGO_ENABLED=0
RUN apk --update add make gcc g++ curl git
RUN make install-builder
RUN make otelcol-logzio-linux_amd64

FROM alpine:3.15.4 as certs
RUN apk --update add ca-certificates

FROM alpine:3.15.4 as directories
RUN mkdir /etc/otel/
COPY --from=builder /src/otel-config/default.yml /etc/otel/config.yaml

FROM scratch
ARG BUILD_TAG=latest
ENV TAG=$BUILD_TAG
ARG USER_UID=10001
USER ${USER_UID}
ENV HOME=/etc/otel/

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /src/otelbuilder/cmd/otelcol-logzio-linux_amd64 /otelcol-logzio
COPY --from=directories --chown=${USER_UID}:${USER_UID} /etc/otel/ /etc/otel/
EXPOSE 4317 55680 55679 8888 6060 7276 9411 9943 1234 6831 6832 14250 14268 4318 8888
ENTRYPOINT ["/otelcol-logzio"]
CMD ["--config", "/etc/otel/config.yaml"]