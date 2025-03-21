# Grafana Alloy

## Ports

| Port | Protocol | Component | Note |
| - | - | - | - |
| 12345 | TCP | Alloy web UI                  | HTTP           |
| 4317  | TCP | `otelcol.receiver.otlp`       | gRPC           |
| 4318  | TCP | `otelcol.receiver.otlp`       | HTTP           |
| 9999  | TCP | `prometheus.receive_http`     | HTTP           |
| 14250 | TCP | `otel.receiver.jaeger`        | gRPC           |
| 6832  | TCP | `otel.receiver.jaeger`        | Thrift         |
| 6831  | TCP | `otel.receiver.jaeger`        | Thrift Compact |
| 14268 | TCP | `otel.receiver.jaeger`        | Thrift HTTP    |
| 9411  | TCP | `otel.exporter.loadbalancing` |                |

## K8s Monitor Helm Chart

| Feature                                    | values.yaml                           | K8s Pod            | K8s Type    |
| ------------------------------------------ | ------------------------------------- | ------------------ | ----------- |
| Cluster metrics                            | `clusterMetrics`                      | alloy-metrics      | StatefulSet |
| -                                          | `clusterMetrics.node-exporter`        | node-exporter      | DaemonSet   |
| -                                          | `clusterMetrics.windows-exporter`     | windows-exporter   | DaemonSet   |
| -                                          | `clustermetrics.kube-state-metrics`   | kube-state-metrics | Deployment  |
| Cost metrics                               | `clusterMetrics.opencost`             | alloy-metrics      | -           |
| Energy metrics                             | `clusterMetrics.kepler`               | alloy-metrics      | -           |
| Autodiscovery with annotations             | `annotationAutodiscovery`             | alloy-metrics      | -           |
| Prometheus Operator objects                | `prometheusOperatorObjects`           | alloy-metrics      | -           |
| Cluster events                             | `clusterEvents`                       | alloy-singleton    | Deployment  |
| Node logs                                  | `nodeLogs`                            | alloy-logs         | DaemonSet   |
| Pod logs                                   | `podLogs`                             | alloy-logs         | -           |
| Application receivers                      | `applicationObservability`            | alloy-receiver     | DaemonSet   |
| Grafana application observability          | `applicationObservability.connectors` | alloy-receiver     | -           |
| Metrics & traces of inbound/outbound calls | `autoInstrumentation`                 | alloy-metrics      | -           |
| -                                          | `autoInstrumentation.beyla`           | beyla              | DaemonSet   |
| CPU profiling                              | `profiling`                           | alloy-profiles     | DaemonSet   |
