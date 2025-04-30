# HttpApi

Built with Maven 3.9.9 and OpenJDK 24. This uses Grafana's distribution of the
 [OpenTelemetry's Java instrumentation agent](https://github.com/grafana/grafana-opentelemetry-java),
 to automatically instrument metrics, logs, and traces.

THe Grafana distribution has been optimized for
 [Application Observability](https://grafana.com/products/cloud/application-observability/).

OTel recommends the agent over
 [Spring Boot starter](https://opentelemetry.io/docs/zero-code/java/spring-boot-starter/)
 since it has more out of the box instrumentation.

## Build

```plaintext
docker compose up --build
```

## Test

```plaintext
curl localhost:8080/cities

curl -X POST http://localhost:8080/cities \
  -H "Content-Type: application/json" \
  -d '{"name":"Boston","state":"MA","county":"Suffolk","founded":1625,"population":675647}'
```

## References

- [The Grafana OpenTelemetry Distribution for Java: Optimized for Application Observability](https://grafana.com/blog/2023/11/16/the-grafana-opentelemetry-distribution-for-java-optimized-for-application-observability/)
- [OpenTelemetry Java instrumentation agent](https://opentelemetry.io/docs/zero-code/java/agent/)
- [open-telemetry/opentelemetry-java-instrumentation](https://github.com/open-telemetry/opentelemetry-java-instrumentation)
- [dev.to: Building REST APIs in Java](https://dev.to/respect17/building-rest-apis-in-java-a-beginners-guidehey-devto-community-121)
- [Spring Boot OpenTelemetry dependency error](https://stackoverflow.com/questions/79528816/springboot-gradle-open-telemetry-dependency-error)
