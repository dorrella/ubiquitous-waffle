FROM golang:1.22.5-alpine3.20 as builder
WORKDIR /build
RUN <<EOF
export CGO_ENABLED=0
export GOOS=linux
go install go.opentelemetry.io/collector/cmd/builder@latest
cat > otelcol-builder.yaml << EOCAT
dist:
  name: otel-custom
  description: "loki and prometheus exporters"
  output_path: /build
receivers:
  - gomod:
      go.opentelemetry.io/collector/receiver/otlpreceiver v0.104.0
  - gomod:
      github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.104.0
  - gomod:
      github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.104.0
  - gomod:
      github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.104.0
exporters:
  - gomod:
      github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter v0.104.0
  - gomod:
      github.com/open-telemetry/opentelemetry-collector-contrib/exporter/lokiexporter v0.104.0
  - gomod:
      go.opentelemetry.io/collector/exporter/otlpexporter v0.104.0
  - gomod:
      go.opentelemetry.io/collector/exporter/debugexporter v0.104.0
extensions:
  - gomod:
      github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.104.0
  - gomod:
      go.opentelemetry.io/collector/extension/memorylimiterextension v0.104.0
processors:
  - gomod:
      go.opentelemetry.io/collector/processor/batchprocessor v0.104.0
EOCAT
builder --config=otelcol-builder.yaml
EOF

FROM otel/opentelemetry-collector:0.104.0
COPY --from=builder /build/otel-custom /otelcol
