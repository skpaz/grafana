# http-api

Built with Python 3.13.5 and [FastAPI](https://fastapi.tiangolo.com/) 0.115.12.
 It uses [opentelemetry-instrumentation-fastapi](https://opentelemetry-python-contrib.readthedocs.io/en/latest/instrumentation/fastapi/fastapi.html),
 the [OTel SDK for Python](https://opentelemetry.io/docs/languages/python/),
 and some additional OTel components to instrumented traces.

## Build

```plaintext
docker compose up --build
```

## Test

```plaintext
curl http://localhost:8080/cities

curl -X 'PUT' 'http://localhost:8080/cities' \
  -H 'Content-Type: application/json' \
  -d '{"name":"Boston","state":"MA","county":"Suffolk","founded":1625,"population":675647}'

curl http://localhost:8080/cities/{:id}
```

## Alloy

You will need to update `/alloy/config.alloy` with your OTLP endpoint, tenant ID, and token.

## References

- [Python | OpenTelemetry](https://opentelemetry.io/docs/languages/python/)
- [OpenTelemetry FastAPI Instrumentation](https://opentelemetry-python-contrib.readthedocs.io/en/latest/instrumentation/fastapi/fastapi.html)
- [FastAPI Tutorial](https://fastapi.tiangolo.com/tutorial/)
