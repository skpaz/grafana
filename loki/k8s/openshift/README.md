# Red Hat OpenShift

Misc. notes when installing GEL or Loki OSS on OpenShift.

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
