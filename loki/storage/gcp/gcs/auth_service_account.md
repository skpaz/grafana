# GCP GCS with Authentication via Service Account

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

## GCP

### IAM

#### Create a Service Account

```txt
gcloud iam service-accounts create <SA_NAME> \
  --description="<DESCRIPTION>" \
  --display-name="<DISPLAY_NAME>"
```

#### Create a Service Account Key

```txt
gcloud iam service-accounts keys create <PATH> \
  --iam-account=<SA_NAME>@<PROJECT_ID>.iam.gserviceaccount.com
```

The token will be saved to the `<PATH>` specified, i.e. `/path/to/private_key.json`. The token can either be saved to a K8s secret and accessed
 via `GOOGLE_APPLICATION_CREDENTIALS` or used as the value for `service_account`. See below for additional information.

#### Permissions

The [pre-defined `role/storage.objectUser` role](https://cloud.google.com/storage/docs/access-control/iam-roles) is sufficient for Loki / GEL to
 operate. See [IAM permissions for Cloud Storage](https://cloud.google.com/storage/docs/access-control/iam-permissions) for details about each individual
 permission. You can use this predefined role or create your own with matching permissions.

### GCS

#### Create Bucket(s)

Create your GCS bucket(s). Multiple buckets can be created at the same time.

```txt
gcloud storage buckets create gs://<BUCKET_NAME> [gs://<BUCKET_NAME>] \
  --location=<REGION> \
  --default-storage-class=STANDARD \
  --public-access-prevention \
  --uniform-bucket-level-access \
  --soft-delete-duration=7d
```

You will need 2-3 buckets: chunks, ruler, and admin (GEL only). The bucket names must be unique, and should not be named "chunks," "ruler," or "admin."

#### Add IAM Policy to Bucket(s)

Create an IAM policy binding on the bucket(s) using the service account created previously and the role(s) of your choice. One command per bucket.

```txt
gcloud storage buckets add-iam-policy-binding gs://<BUCKET_NAME> \
  --member=serviceAccount:<SA_NAME>@<PROJECT_ID>.iam.gserviceaccount.com \
  --role=roles/storage.objectUser
```

## Kubernetes

### Create a Kubernetes Secret

```txt
kubectl create secret generic <SECRET_NAME> \
  --namespace <NAMESPACE> \
  --from-file=key.json=<PATH_TO_JSON_KEY>
```

Replace `<SECRET_NAME>` with a name for your secret, `<NAMESPACE>` with the K8s namespace Loki/GEL was installed in (if any), and `<PATH_TO_JSON_KEY>`
 with the path to the service account key JSON file. This is the same file created in the **Create a Service Account Key** step above.

## Helm

### values.yaml

Loki/GEL services support the `service_account` method (not covered below / not recommended) as well as ADC. For more information, see
 [How Application Default Credentials works](https://cloud.google.com/docs/authentication/application-default-credentials). In the example below,
  we use the `GOOGLE_APPLICATION_CREDENTIALS` method.

> [!WARNING]
> This increases the risk that the JSON key could be leaked. For this reason, it's not recommended to use service account-based
> authentication in production environments until ADC is supported across all services.

The example below does not represent a complete `values.yaml` file, only the parameters that need to be updated to run Loki/GEL with service
 acccount-based authentication.

A Simple Scalable deployment breaks up Loki/GEL operations into distinct `read`, `write`, and `backend` services that can be scaled independently.
 If `enterprise.enabled` is set to `true`, there'll also be a `tokengen` job and and `adminApi` pod that need to connect to the storage backend.

```yaml
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

# GEL ONLY: tokenGen
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
 and each `<*_BUCKET_NAME>` with the name of the buckets created in the **GCS > Create Bucket(s)** step above.
