# Red Hat OpenShift

Misc. notes when installing GEL or Loki OSS on OpenShift.

## CustomResourceDefinitions (CRDs)

Use `--skip-crds` if there are native OpenStack CRDs built.

## SecurityContextConstraints (SCCs)

Enable SCC.

```yaml
rbac:
  pspEnabled: false
  sccEnabled: true
  namespaced: false
```

## External Access

Use the NGINX gateway and an ingress route for external access to GEL. OpenShift recommends the use of an Ingress Controller (i.e. routes) for HTTP/HTTPS, and "other" TLS-encrypted protocols. Presumably, a passthrough route would be ideal:

```plaintext
oc expose service SERVICE_NAME
```

## DNS

### Global

Update DNS service and namespace names.

```yaml
global:
  # dnsService: kube-dns
  dnsService: dns-default
  # dnsNamespace: kube-system
  dnsNamespace: openshift-dns
```

### Gateway

Update the `ConfigMap` with the OpenShift DNS resolver.

```yaml
...
apiVersion: v1
data:
  nginx.conf: |
    ...
    http {
      ...
      resolver dns-default.openshift-dns.svc.cluster.local.;
      ...
```

## Tokengen

The `tokengen` service account will need access to the base SCC created by the
 Helm chart.

```plaintext
oc adm policy add-scc-to-user FIXME_INSTALL_NAME -n FIXME_NAMESPACE -z FIXME_INSTALL_NAME-tokengen
```

## MinIO

If MinIO is used, the `minio-sa` service account will need access to the `minio`
 SCC created by the Helm chart.

```plaintext
oc adm policy add-scc-to-user FIXME_INSTALL_NAME-minio -n FIXME_NAMESPACE -z minio-sa
```
