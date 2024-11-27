# GCP GCS with Authentication via Workload Federation

## Assumptions

- You're installing Loki/Grafana Enterprise Logs (GEL) on a GKE cluster in Autopilot mode.
  - Standard GKE clusters require additional steps. See **Resources > Authenticate to Google Cloud APIs from GKE workloads** for more information.
- You're deplpoying Loki/GEL in Simple Scalable mode.
- You're familiar with `gcloud`; this doc doesn't contain any Console examples.
- You're using [the `grafana/loki` Helm chart](https://github.com/grafana/loki/tree/main/production/helm/loki).

## Resources

- [Authenticate to Google Cloud APIs from GKE workloads](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity)
- [IAM roles for Cloud Storage](https://cloud.google.com/storage/docs/access-control/iam-roles)
- [IAM permissions for Cloud Storage](https://cloud.google.com/storage/docs/access-control/iam-permissions)
- [Creating buckets](https://cloud.google.com/storage/docs/creating-buckets)
- [How Application Default Credentials (ADC) work](https://cloud.google.com/docs/authentication/application-default-credentials)

## GCP

There are two methods to enable WIF for GKE.

1. IAM principal authorization via Kubernetes `ServiceAccount`.
2. Use Kubernetes `ServiceAccount` to impersonate an IAM service account.

This document covers **IAM principal authorization**. Per GCP, linked service accounts should only be used if the limitations imposed by principal authorization causes problems. See **Resources > Authenticate to Google Cloud APIs from GKE workloads** for more information.

### GKE

#### Create Service Account

Create a K8s service account:

```txt
kubectl create serviceaccount <KSA_NAME> \
  --namespace <NAMESPACE>
```

Replace `<KSA_NAME>` with a service acccount name and `<NAMESPACE>` with the namespace where Loki/GEL is installed.

### GCS

Create your GCS bucket(s). Multiple buckets can be created at the same time.

```txt
gcloud storage buckets create gs://<BUCKET_NAME> [gs://<BUCKET_NAME>] \
  --location=<REGION> \
  --default-storage-class=STANDARD \
  --public-access-prevention \
  --uniform-bucket-level-access \
  --soft-delete-duration=7d
```

Replace `<BUCKET_NAME>` with your bucket name(s) and `<REGION>` with a region (ex. us-west1). You can also adjust the other options to suit your needs.

You will need 2-3 buckets: chunks, ruler, and admin (GEL only). The bucket names must be unique, and should not be named "chunks," "ruler," or "admin."

### IAM

#### Permissions

The [pre-defined `role/storage.objectUser` role](https://cloud.google.com/storage/docs/access-control/iam-roles) is sufficient for Loki / GEL to
 operate. See [IAM permissions for Cloud Storage](https://cloud.google.com/storage/docs/access-control/iam-permissions) for details about each individual
 permission. You can use this predefined role or create your own with matching permissions.

#### Add IAM Policy to Bucket(s)

Create an IAM policy binding on the bucket(s) using the KSA created previously and the role(s) of your choice. One command per bucket.

```txt
gcloud storage buckets add-iam-policy-binding gs://<BUCKET_NAME> \
  --role=roles/storage.objectViewer \
  --member=principal://iam.googleapis.com/projects/<PROJECT_NUMBER>/locations/global/workloadIdentityPools/<PROJECT_ID>.svc.id.goog/subject/ns/<NAMESPACE>/sa/<KSA_NAME> \
  --condition=None
```

Replace `<PROJECT_ID>` with the GCP project ID (ex. project-name), `<PROJECT_NUMBER>` with the project number (ex. 1234567890), `<NAMESPACE>` with the namespace where Loki/GEL is installed, and `<KSA_NAME>` with the name of the KSA you created above.

For GEL, you'll also need bind the `tokenGen` KSA as well:

```txt
gcloud storage buckets add-iam-policy-binding gs://<BUCKET_NAME> \
  # --role=roles/storage.objectUser \
  --member=principal://iam.googleapis.com/projects/<PROJECT_NUMBER>/locations/global/workloadIdentityPools/<PROJECT_ID>.svc.id.goog/subject/ns/<NAMESPACE>/sa/enterprise-logs-tokengen \
  --condition=None
```

The `tokenGen` KSA uses a fixed name: `enterprise-logs-tokengen`. It's defined [here](https://github.com/grafana/loki/blob/4b5925a28e61f29a20aaabda3a159386a8ba7638/production/helm/loki/templates/tokengen/_helpers.yaml), which is based on `loki.name` defined [here](https://github.com/grafana/loki/blob/716d54e2a9617a80c2496a46e9c4cbf8ed51a5d9/production/helm/loki/templates/_helpers.tpl). While it may be possible to modify the auto-generated `-tokengen` KSA, it's currently easier to simply grant it the same permissions as the KSA you created.

## Helm

### values.yaml

A Simple Scalable deployment breaks up Loki/GEL operations into distinct `read`, `write`, and `backend` services that can be scaled independently.
 If `enterprise.enabled` is set to `true`, there'll also be a `tokengen` job and an `adminApi` pod that need to connect to the storage backend.

The `serviceAccount` is automatically applied to `read`, `write`, `backend`, and `adminApi`. The `tokenGen` job creates its own KSA, which necessitates the additional IAM policy binding in the **IAM > Add IAM Policy to Bucket(s)** step above.

```yaml
serviceAccount:
  create: false
  name: <KSA_NAME>

loki:
  schemaConfig:
    configs:
      - from: 2020-07-01
        store: tsdb
        object_store: gcs
        schema: v13
        index:
          prefix: index_
          period: 24h

  storage:
    bucketNames:
      chunks: <CHUNKS_BUCKET_NAME>
      ruler: <RULER_BUCKET_NAME>
      # GEL ONLY: admin bucket
      admin: <ADMIN_BUCKET_NAME>
    type: gcs
```

Replace `<KSA_NAME>` with the name of the KSA created in the **GKE > Create Service Account** and each `<*_BUCKET_NAME>` with the name of the buckets created in the **GCS > Create Bucket(s)** step above.