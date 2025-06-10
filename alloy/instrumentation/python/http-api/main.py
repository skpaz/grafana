import fastapi
import json
from opentelemetry import trace
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
import os
from pydantic import BaseModel

# add service name to resource attributes
resource = Resource(
  attributes = {
    SERVICE_NAME: os.environ['SERVICE_NAME']
  }
)

# define trace provider
traceProvider = TracerProvider(resource=resource)
# add span processor to trace provider
traceProvider.add_span_processor(BatchSpanProcessor(OTLPSpanExporter()))
# set trace provider
trace.set_tracer_provider(traceProvider)

# define app
app = fastapi.FastAPI()

# define class to store city info, used for PUT
class City(BaseModel):
  name: str
  state: str
  county: str
  founded: int
  population: int

# read cities.json into memory
with open('cities.json', encoding="utf-8") as file:
  cities = json.load(file)

# GET cities
@app.get("/cities")
async def get_cities():
  return cities

# GET specific city
@app.get("/cities/{city_id}")
async def get_city(city_id: str):
  return cities[city_id]

# PUT new city
@app.put("/cities")
async def put_city(city: City):
  # fetch list of numeric keys
  keys = []
  for k, v in cities.items():
    keys.append(int(k))
  # add 1 to max() key, combine with city
  data = { str(max(keys)+1):city}
  cities.update(data)
  return cities

# start the fastapi instrumentor 
FastAPIInstrumentor.instrument_app(app)
