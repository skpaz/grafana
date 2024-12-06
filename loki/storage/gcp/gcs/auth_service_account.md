# GCP GCS Authentication via Service Account for Grafana Loki

> [!WARNING]
> GCP recommends [Workload Identity Federation](https://cloud.google.com/iam/docs/workload-identity-federation) over service account-based authentication.

## Assumptions

- You're installing Loki/Grafana Enterprise Logs (GEL) on a GKE cluster.
- You're deplpoying Loki/GEL in Simple Scalable mode.
- You're familiar with `gcloud`; this doc doesn't contain any Console examples.
- You're using [the `grafana/loki` Helm chart](https://github.com/grafana/loki/tree/main/production/helm/loki).

## Resources

- [Create and grant roles to service agents](https://cloud.google.com/iam/docs/create-service-agents)
- [Create service accounts](https://cloud.google.com/iam/docs/service-accounts-create)
- [Attach service accounts to resources](https://cloud.google.com/iam/docs/attach-service-accounts)
- [IAM roles for Cloud Storage](https://cloud.google.com/storage/docs/access-control/iam-roles)
- [IAM permissions for Cloud Storage](https://cloud.google.com/storage/docs/access-control/iam-permissions)
- [Creating buckets](https://cloud.google.com/storage/docs/creating-buckets)
- [How Application Default Credentials (ADC) work](https://cloud.google.com/docs/authentication/application-default-credentials)

## Service Account Authentication

This document covers **Google service account (GSA) authentication**, where a service account and secret key are used to authenticate to the GCS API
 in order to access GCS resources.

## Instructions

> [!IMPORTANT]
> The roles and policies below are the _minimum_ required for this setup. They do not take into account your organizations security requirements.
 It is your responsibility to secure your environment.

### 1. Create a Google Service Account

Create a Google Service Account:

```txt
gcloud iam service-accounts create <GSA_NAME>
```

Replace `<GSA_NAME>` with the name of the service account you want to create.

Example:

```txt
gcloud iam service-accounts create loki-lab-gsa
```

### 2. Create a Service Account Key

```txt
gcloud iam service-accounts keys create <PATH> \
  --iam-account=<GSA_NAME>@<PROJECT_ID>.iam.gserviceaccount.com
```

Replace `<PATH>` with the path to your JSON key file, `<GSA_NAME>` with your GSA name, and `<PROJECT_ID>` with your GCP project ID.

Example:

```txt
gcloud iam service-accounts keys create ~/Downloads/key.json \
  --iam-account=loki-lab-gsa@loki-lab.iam.gserviceaccount.com
```

Once downloaded, the token can either be saved to a K8s secret and accessed via `GOOGLE_APPLICATION_CREDENTIALS` or used as the value for `service_account`. See below for additional information.

### 3. Create GCS Bucket(s)

Create your GCS bucket(s). Multiple buckets can be created with a single command.

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

Example:

```txt
gcloud storage buckets create gs://loki-lab-chunks gs://loki-lab-ruler gs://loki-lab-admin \
  --location=us-west1 \
  --default-storage-class=STANDARD \
  --public-access-prevention \
  --uniform-bucket-level-access \
  --soft-delete-duration=7d
```

### 4. Add IAM Policy to Bucket(s)

> [!INFO]
> The [pre-defined `role/storage.objectUser` role](https://cloud.google.com/storage/docs/access-control/iam-roles) is sufficient for Loki / GEL to
 operate. See [IAM permissions for Cloud Storage](https://cloud.google.com/storage/docs/access-control/iam-permissions) for details about each individual
 permission. You can use this predefined role or create your own with matching permissions.

Create an IAM policy binding on the bucket(s) using the service account created previously and the role(s) of your choice. One command per bucket.

```txt
gcloud storage buckets add-iam-policy-binding gs://<BUCKET_NAME> \
  --member=serviceAccount:<GSA_NAME>@<PROJECT_ID>.iam.gserviceaccount.com \
  --role=roles/storage.objectUser
```

Replace `<BUCKET_NAME>` with the name of the GCS bucket you'd like to create, `<GSA_NAME>` with the GSA name, and `<PROJECT_ID>` with your GCP project ID.

Example:

```txt
gcloud storage buckets add-iam-policy-binding gs://loki-lab-chunks \
  --member=serviceAccount:loki-lab-gsa@loki-lab.iam.gserviceaccount.com \
  --role=roles/storage.objectUser
```

### 5. Create a Kubernetes Namespace

Create a K8s namespace where you'll install your Loki/GEL workloads:

```txt
kubectl create namespace <NAMESPACE>
```

Replace `<NAMESPACE>` with the namespace where your Loki/GEL workloads will be located.

Example:

```txt
kubectl create namespace loki
```

### 6. Create a Kubernetes Secret

Create a K8s secret in the same namespace as your Loki/GEL workloads:

```txt
kubectl create secret generic <SECRET_NAME> \
  --namespace <NAMESPACE> \
  --from-file=key.json=<PATH_TO_JSON_KEY>
```

Replace `<SECRET_NAME>` with a name for your secret, `<NAMESPACE>` with the K8s namespace Loki/GEL was installed in (if any), and `<PATH_TO_JSON_KEY>`
 with the path to the service account key JSON file. This is the same file created in the **Create a Service Account Key** step above.

Example:

```txt
kubectl create secret generic loki-lab-secret \
  --namespace loki \
  --from-file=key.json=~/Downloads/key.json
```

## Helm

### values.yaml

In the example below we use the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.

> [!WARNING]
> There is an alternate method that uses the `service_account` parameter, which is documented elsewhere. Use of the `service_account`
 method increases the risk that the JSON key could be leaked. For this reason, it's not recommended to use service account-based authentication
 in production environments.

The example below does not represent a complete `values.yaml` file, only the parameters that need to be updated to run Loki/GEL with service
 acccount-based authentication.

A Simple Scalable deployment breaks up Loki/GEL operations into distinct `read`, `write`, and `backend` services that can be scaled independently.
 If `enterprise.enabled` is set to `true`, there'll also be a `tokengen` job and an `adminApi` pod that need to connect to the storage backend.

```yaml
# GEL ONLY: enterprise.*
enterprise:
  tokengen:
    # Note the use of `env` instead of `extraEnv` here.
    env:
      - name: GOOGLE_APPLICATION_CREDENTIALS
        value: /etc/secrets/key.json
    extraVolumeMounts:
      - mountPath: /etc/secrets
        name: <MOUNT_NAME>
    extraVolumes:
      - name: <MOUNT_NAME>
        secret:
          secretName: <SECRET_NAME>

# GEL ONLY: adminApi
adminApi:
  # Note the use of `env` instead of `extraEnv` here.
  env:
    - name: GOOGLE_APPLICATION_CREDENTIALS
      value: /etc/secrets/key.json
  extraVolumeMounts:
    - mountPath: /etc/secrets
      name: <MOUNT_NAME>
  extraVolumes:
    - name: <MOUNT_NAME>
      secret:
        secretName: <SECRET_NAME>

backend:
  extraEnv:
    - name: GOOGLE_APPLICATION_CREDENTIALS
      value: /etc/secrets/key.json
  extraVolumeMounts:
    - mountPath: /etc/secrets
      name: <MOUNT_NAME>
  extraVolumes:
    - name: <MOUNT_NAME>
      secret:
        secretName: <SECRET_NAME>

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
      # GEL ONLY: bucketNames.admin
      admin: <ADMIN_BUCKET_NAME>
    type: gcs

read:
  extraEnv:
    - name: GOOGLE_APPLICATION_CREDENTIALS
      value: /etc/secrets/key.json
  extraVolumeMounts:
    - mountPath: /etc/secrets
      name: <MOUNT_NAME>
  extraVolumes:
    - name: <MOUNT_NAME>
      secret:
        secretName: <SECRET_NAME>

write:
  extraEnv:
    - name: GOOGLE_APPLICATION_CREDENTIALS
      value: /etc/secrets/key.json
  extraVolumeMounts:
    - mountPath: /etc/secrets
      name: <MOUNT_NAME>
  extraVolumes:
    - name: <MOUNT_NAME>
      secret:
        secretName: <SECRET_NAME>
```

Replace `<MOUNT_NAME>` with a name for your volume mount, `<SECRET_NAME>` with the secret name you defined in the `kubectl create secret` command above,
 and each `<*_BUCKET_NAME>` with the name of the buckets created in the **Create GCS Bucket(s)** step above.
