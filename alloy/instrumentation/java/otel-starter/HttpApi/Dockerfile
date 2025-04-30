FROM openjdk:24-slim AS run
COPY /target/http-api-*.jar app.jar
CMD ["java","-jar","app.jar"]