FROM openjdk:24-slim AS run
WORKDIR /app
RUN apt-get update \
    && apt-get install -y wget jq \
    && rm -rf /var/lib/apt/lists/*
RUN wget -O agent.jar $( \
      wget -q -O - https://api.github.com/repos/grafana/grafana-opentelemetry-java/releases/latest \
      | jq -r '.assets[] | select(.name | contains ("jar")) | .browser_download_url' \
    )
COPY /target/http-api-*.jar app.jar
CMD ["java","-jar","app.jar","-javaagent:agent.jar"]