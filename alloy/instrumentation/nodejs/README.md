# http-api

Build with Node.js v18.20.5 and Express v5.1.0. It uses
 [auto-instrumentations-node](https://www.npmjs.com/package/@opentelemetry/auto-instrumentations-node)
 and the [OTel SDK for Node](https://opentelemetry.io/docs/languages/js/getting-started/nodejs/)
 to instrument traces.

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

- [Node.js | OpenTelemetry](https://opentelemetry.io/docs/languages/js/getting-started/nodejs/)
- [Node.js RESTful API | tutorialspoint](https://www.tutorialspoint.com/nodejs/nodejs_restful_api.htm)
