# Send Logs to Loki

## Test Logs

### Kubernetes Service Port Forward

If Loki/GEL is installed on K8s, a port-forward may be necessary to access the gateway:

```txt
kubectl port-forward -n <NAMESPACE> service/<NAME_OF_GATEWAY_SERVICE> <LOCAL_PORT>:80
```

Replace `<NAMESPACE>` with the name of the namespace where you installed Loki/GE. By default, the `<NAME_OF_GATEWAY_SERVICE>` for Loki and GEL should be `loki-gateway` and `enterprise-logs-gateway` respectively. Pick a random, free/open `<LOCAL_PORT>` to forward traffic to the gateway on TCP/80.

Example:

```txt
kubectl port-forward -n gel service/enterprise-logs-gateway 8080:80
```

### Send Test Logs

Send test logs:

```txt
curl -v -H 'Content-Type: application/json' \
  -H "Authorization: Basic $(echo '<TENANT_NAME>:<TOKEN>' | base64)" \
  -s -X POST 'http://<URL>>:<PORT>/loki/api/v1/push' \
  --data-raw '{"streams":[{"stream":{"app":"test","env":"prod"},"values":[["<UNIX_EPOCH_IN_NANOSECONDS>", "test log message " ],["<UNIX_EPOCH_IN_NANOSECONDS>","test log message 2" ]]}]}'
```

Example:

```txt
curl -v -H 'Content-Type: application/json' \
  -H "Authorization: Basic $(echo 'primary:tokenbarfhere' | base64)" \
  -s -X POST 'http://localhost:8080/loki/api/v1/push' \
  --data-raw '{"streams":[{"stream":{"app":"test","env":"prod"},"values":[["1732736384000000000", "fizzbuzz1" ],["173273638500000000","fizzbuzz2" ]]}]}'
```

If you don't want to include the token in each `curl` command, you could also assign it to variable.

```txt
TOKEN=$(echo '<TENANT_NAME>:<TOKEN>' | base64)
curl -v -H 'Content-Type: application/json' \
  -H "Authorization: Basic ${TOKEN}" \
  -s -X POST 'http://<URL>>:<PORT>/loki/api/v1/push' \
  --data-raw '{"streams":[{"stream":{"app":"test","env":"prod"},"values":[["<UNIX_EPOCH_IN_NANOSECONDS>", "fizzbuzz1" ],["<UNIX_EPOCH_IN_NANOSECONDS>","fizzbuzz2" ]]}]}'
```

#### UNIX Timestamp

You can fetch a Unix timestamp in nanoseconds from macOS / Terminal:

```txt
echo $(($(date +%s) * 1000000000))
```

> NOTE: `%N` is not available in macOS.

It should be similar in other *nix-based OSes, but `+%N` may be able to produce nanoseconds natively without math involved.

For Windows / PowerShell, try:

```powershell
[int](Get-Date -UFormat %s) * 1000000000 -as [decimal]
```

You can also use a tool like [EpochConverter](https://www.epochconverter.com/) to generate a timestamp.
