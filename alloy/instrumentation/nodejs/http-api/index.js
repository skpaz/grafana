// -- otel
// credit: https://opentelemetry.io/docs/languages/js/getting-started/nodejs/

const { NodeSDK } = require('@opentelemetry/sdk-node');
const { ConsoleSpanExporter } = require ('@opentelemetry/sdk-trace-node');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-proto');
const { resourceFromAttributes } = require('@opentelemetry/resources');
const { 
  ATTR_SERVICE_NAME,
  ATTR_SERVICE_VERSION
} = require('@opentelemetry/semantic-conventions');

// debug
//const { diag, DiagConsoleLogger, DiagLogLevel } = require('@opentelemetry/api');
//diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.DEBUG);
//

const serviceName = process.env.SERVICE_NAME
const collectorURL = process.env.OTEL_EXPORTER_OTLP_ENDPOINT

const sdk = new NodeSDK({
  traceExporter: new OTLPTraceExporter({
    url: collectorURL + '/v1/traces',
  }),
  resource: resourceFromAttributes({
    [ ATTR_SERVICE_NAME ]: serviceName,
    [ ATTR_SERVICE_VERSION ]: "1.0.0",
  }),
  instrumentations: [
    getNodeAutoInstrumentations()
  ]
});

sdk.start();

// -- http-api
// credit: https://www.tutorialspoint.com/nodejs/nodejs_restful_api.htm

var express = require('express');
var app = express();
var fs = require('fs');
var onFinished = require('on-finished')

app.get('/cities', function(req,res) {
  fs.readFile(__dirname + "/" + "cities.json", 'utf-8', function(err,data) {
    res.end(data);
    onFinished(res,function(err,res) {
      console.log("[Express] GET\t/cities\t\t\t\t\t\t["+res.statusCode+"]")
    })
  });
});

app.get('/cities/:id', function(req,res) {
  fs.readFile(__dirname + "/" + "cities.json", 'utf-8', function(err,data) {
    var cities = JSON.parse(data);
    var city = cities[req.params.id]
    res.end(JSON.stringify(city));
    onFinished(res,function(err,res) {
      console.log("[Express] GET\t/cities/:id ("+req.params.id+")\t\t\t\t\t["+res.statusCode+"]")
    })
  });
})

var bodyParser = require('body-parser');
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extend:true}));

app.post('/cities', function(req,res) {
  fs.readFile(__dirname + "/" + "cities.json", 'utf-8', function(err,data) {
    var cities = JSON.parse(data);
    var maxCityId = Math.max.apply(Math, cities.map(mId => mId.id))
    var city = req.body
    cities[maxCityId+1] = city
    res.end(JSON.stringify(cities));
    onFinished(res,function(err,res) {
      console.log("[Express] POST\t/cities\t\t\t\t\t\t["+res.statusCode+"]")
    })
  });
})

var server = app.listen(8080, function() {
  console.log("http-api listening at http://localhost:8080")
})
