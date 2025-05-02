# HttpApi

Built with Maven 3.9.9 and OpenJDK 24. This uses the
 [OpenTelemetry Spring Boot starter](https://opentelemetry.io/docs/zero-code/java/spring-boot-starter/)
 to automatically instrument metrics, logs, and traces.

## Build

Create .jar file.

```plaintext
mvn clean install
```

Build containers:

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

## Alloy

You will need to update `/alloy/config.alloy` with your OTLP endpoint, tenant ID, and token.

## References

- [OpenTelemetry Java Spring Boot starter](https://opentelemetry.io/docs/zero-code/java/spring-boot-starter/)
- [dev.to: Building REST APIs in Java](https://dev.to/respect17/building-rest-apis-in-java-a-beginners-guidehey-devto-community-121)
- [Spring Boot OpenTelemetry dependency error](https://stackoverflow.com/questions/79528816/springboot-gradle-open-telemetry-dependency-error)
