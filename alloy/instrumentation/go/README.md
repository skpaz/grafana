# http-api

Built with Go v1.24.0 and Gin v1.10.0. It uses
 [otelgin](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin)
 to easily instrument traces.

## Build

```plaintext
docker compose up --build
```

## Test

```plaintext
curl http://localhost:8080/cities

curl http://localhost:8080/cities \
   --include \
   --header "Content-Type: application/json" \
   --request "POST" \
   --data '{"id":4,"name":"Boston","state":"MA","county":"Suffolk","founded":1625,"population":675647}'

curl http://localhost:8080/cities/{:id}
```

## Alloy

You will need to update `/alloy/config.alloy` with your OTLP endpoint, tenant ID, and token.

## References

- [Implementing OpenTelemetry in a Gin application | SigNoz](https://signoz.io/blog/opentelemetry-gin/)
- [Tutorial: Developing a RESTful API with Go and Gin | Go](https://go.dev/doc/tutorial/web-service-gin)
- [Go | OpenTelemetry](https://opentelemetry.io/docs/languages/go/)
