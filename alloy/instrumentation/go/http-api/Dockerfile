# dockerfile-app
FROM golang:latest AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /http-api

FROM ubuntu:latest AS run
WORKDIR /
COPY --from=build /http-api /http-api
EXPOSE 8080
CMD ["/http-api"]
