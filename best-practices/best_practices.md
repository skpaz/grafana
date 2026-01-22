# Best Practices

An index of Grafana's best practice documentation, as well as additional best practice documentation for adjacent projects like OpenTelemetry. This index is focused on Grafana Cloud customers. It will contain references to Grafana's open source products, but the same best practices can be applied to Cloud.

You can also [search for "best practices" in Grafana's technical documentation](https://grafana.com/search/?query=best+practices&type=doc) and filter by specific product or feature.

## Grafana

> _[Grafana Alloy](https://grafana.com/docs/alloy/latest/) combines the strengths of the leading collectors into one place. Whether observing applications, infrastructure, or both, Grafana Alloy can collect, process, and export telemetry signals to scale and future-proof your observability approach._

- [Clustering](https://grafana.com/docs/alloy/latest/get-started/clustering/#best-practices) - _scrape target size and quantity._
- [Data pipelines](https://grafana.com/docs/alloy/latest/get-started/components/build-pipelines/#best-practices) - _pipeline scope, labels, secrets, and tests._
- [Module security](https://grafana.com/docs/alloy/latest/get-started/modules/#security) - _secure module use and maintenance._

### Cloud

> _[Grafana Cloud](https://grafana.com/docs/grafana-cloud/) is a highly available, performant, and scalable observability platform for your applications and infrastructure. It provides a centralized view over all of your observability data, whether the data lives in Grafana Cloud Metrics services or in your own bare-metal and cloud environments. With native support for many popular data sources such as Prometheus, Elasticsearch, and Amazon CloudWatch, all you have to do to begin creating dashboards and querying metrics data is to configure data sources in Grafana Cloud._

#### Adaptive Telemetry

> _[Adaptive Telemetry](https://grafana.com/docs/grafana-cloud/adaptive-telemetry/) optimizes which signals get stored and ensures that your teams only keep the most valuable telemetry._

##### Adaptive Traces

> _[Adaptive Traces](https://grafana.com/docs/grafana-cloud/adaptive-telemetry/adaptive-traces/) helps you automatically identify and retain your most valuable traces, so that you can get the insights you need into application performance and availability, while optimizing your overall observability costs._

- [Policies](https://grafana.com/docs/grafana-cloud/adaptive-telemetry/adaptive-traces/guides/best-practices-policies/) - _policy strategy, revision, and refinement._

#### Alerts & IRM

> _[Alerts & IRM](https://grafana.com/docs/grafana-cloud/alerting-and-irm/) enables teams to efficiently detect, respond, and learn from incidents within one centralized platform._

#### Alerting

> _[Grafana Alerting](https://grafana.com/docs/grafana-cloud/alerting-and-irm/alerting/) allows you to learn about problems in your systems moments after they occur._

- [Alerting](https://grafana.com/docs/grafana-cloud/alerting-and-irm/alerting/guides/best-practices/) - _prioritization, escalation, scope, clarity, and ownership._ 

##### SLO

> _With [Grafana SLO](https://grafana.com/docs/grafana-cloud/alerting-and-irm/slo/), you can create metrics to measure the quality of the service you provide users._

- [SLOs](https://grafana.com/docs/grafana-cloud/alerting-and-irm/slo/best-practices/) - _building good SLOs, team alignment, simple queries, alerts, and labels._

#### Dashboards

> _[Dashboards](https://grafana.com/docs/grafana-cloud/visualizations/dashboards/) allow you to query, transform, visualize, and understand your data no matter where it’s stored._

- [Dashboards](https://grafana.com/docs/grafana-cloud/visualizations/dashboards/build-dashboards/best-practices/) - _observability strategies, maturity model, creation, and management._

#### Private data source connect (PDC)

> _[Private data source connect](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/), or PDC, is a way for you to establish a private, secured connection between a Grafana Cloud instance, or stack, and data sources secured within a private network._

- [Data source configuration](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/data-source-best-practice/) - _data source configuration._

### k6

> _[Grafana k6](https://grafana.com/docs/k6/latest/) is an open-source, developer-friendly, and extensible load testing tool. k6 allows you to prevent performance issues and proactively improve reliability._

- [JavaScript API](https://grafana.com/docs/k6/latest/javascript-api/) - _best practices by module/function._

### Loki (Logs)

> _[Grafana Loki](https://grafana.com/docs/loki/latest/) is a set of open source components that can be composed into a fully featured logging stack. A small index and highly compressed chunks simplifies the operation and significantly lowers the cost of Loki._

- [Queries](https://grafana.com/docs/loki/latest/query/bp-query/) - _label selectors, time range, filters, text parsers, and recording rules._
- [Labels](https://grafana.com/docs/loki/latest/get-started/labels/bp-labels/) - _static labels, dynamic labels, bounded values, and labels applied by clients._

### Mimir (Metrics)

> _[Grafana Mimir](https://grafana.com/docs/mimir/latest/) is an open source software project that provides horizontally scalable, highly available, multi-tenant, long-term storage for Prometheus and OpenTelemetry metrics._

- [Queries](https://grafana.com/docs/mimir/latest/query/query-best-practices/) - _label matchers, time range, evaluation frequency, filters, and recording rules._
- Label best practices are similiar to Loki. The main difference is application. With Loki, the focus is on indexed labels and search performance. Mimir's main concern is the total number of active series.

### Tempo (Traces)

> _[Grafana Tempo](https://grafana.com/docs/tempo/latest/) is an open-source, easy-to-use, and high-scale distributed tracing backend. Tempo lets you search for traces, generate metrics from spans, and link your tracing data with logs and metrics._

- [Trace instrumentation](https://grafana.com/docs/tempo/latest/set-up-for-tracing/instrument-send/best-practices/) - _span and resource attributes, where to add spans, and span length._

## OpenTelemetry

Grafana Labs is an active user of and contributor to the [OpenTelemetry](https://opentelemetry.io/) project.

- [OpenTelemetry best practices: A user’s guide to getting started with OpenTelemetry](https://grafana.com/blog/opentelemetry-best-practices-a-users-guide-to-getting-started-with-opentelemetry/)

### APIs & SDKs

#### .NET

- [Grafana .NET SDK](https://github.com/grafana/grafana-opentelemetry-dotnet) - _we [recommend the use](https://grafana.com/docs/opentelemetry/instrument/grafana-dotnet/) of our .NET SDK with Grafana._
- [Metrics](https://opentelemetry.io/docs/languages/dotnet/metrics/best-practices/) - _instrument types, memory management, correlation, and enrichment._

#### Java

- [Grafana Java Agent](https://github.com/grafana/grafana-opentelemetry-java) - _we [recommend the use](https://grafana.com/docs/opentelemetry/instrument/grafana-java/) of our Java agent with Grafana._
