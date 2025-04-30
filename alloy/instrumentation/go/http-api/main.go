package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
  "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
	"net/http"
	"os"
	"strconv"
)

// -- otel
// credit: https://signoz.io/blog/opentelemetry-gin/

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
)

func initTracer() func(context.Context) error {

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("service.version", "1.0.0"),
		),
	)

	if err != nil {
		log.Printf("Count not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

// -- http-api
// credit: https://go.dev/doc/tutorial/web-service-gin

type city struct {
	ID         int	  `json:"id"`
	Name       string `json:"name"`
	State      string `json:"state"`
	County     string `json:"county"`
	Founded    int    `json:"founded"`
	Population int    `json:"population"`
}

var cities = []city{
	{ID: 0, Name: "Seattle", State: "WA", County: "King", Founded: 1851, Population: 737015},
	{ID: 1, Name: "Portland", State: "OR", County: "Multnomah", Founded: 1845, Population: 652503},
	{ID: 2, Name: "Los Angeles", State: "CA", County: "Los Angeles", Founded: 1781, Population: 3898747},
	{ID: 3, Name: "Phoenix", State: "AZ", County: "Maricopa", Founded: 1867, Population: 1608139},
}

func getCities(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, cities)
}

func postCities(c *gin.Context) {
	var newCity city
	if err := c.BindJSON(&newCity); err != nil {
		return
	}
	cities = append(cities, newCity)
	c.IndentedJSON(http.StatusCreated, newCity)
}

func getCityById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	for _, a := range cities {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message":"city not found"})
}

// --- main

func main() {
	cleanup := initTracer()
	defer cleanup(context.Background())
	router := gin.Default()
	router.Use(otelgin.Middleware(serviceName))
	router.GET("/cities", getCities)
	router.GET("/cities/:id", getCityById)
	router.POST("/cities", postCities)
	router.Run("0.0.0.0:8080")
}