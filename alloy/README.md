# Grafana Alloy

## Ports

| Port | Protocol | Component | Note |
| - | - | - | - |
| 12345 | TCP | Alloy web UI       | HTTP           |
| 4317  | TCP | OTLP               | gRPC           |
| 4318  | TCP | OTLP               | HTTP           |
| 9999  | TCP | Prometheus         | HTTP           |
| 14250 | TCP | OTLP Jaeger        | gRPC           |
| 6832  | TCP | OTLP Jaeger        | Thrift         |
| 6831  | TCP | OTLP Jaeger        | Thrift Compact |
| 14268 | TCP | OTLP Jaeger        | Thrift HTTP    |
| 9411  | TCP | OTel Load Balancer |                |

## K8s Monitor Helm Chart

### Feature Map

### v2

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
