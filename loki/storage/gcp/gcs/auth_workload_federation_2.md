# GCP GCS Authentication via Workload Federation and Service Account for Grafana Loki

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

## Workload Federation

There are two methods to enable WIF for GKE.

1. IAM principal authorization via Kubernetes `ServiceAccount`.
2. Use a Kubernetes `ServiceAccount` to impersonate an Google service account.

This document covers **service account impersonation**, where you link a Kubernetes service account to IAM in order to impersonate a Google Service Account in order to authenticate to the GCS API and access GCS resources.

> [!WARNING]
> Per GCP, service account impersonation should only be used if the limitations imposed by principal authorization causes problems.
 See **Resources > Authenticate to Google Cloud APIs from GKE workloads** for more information.

## Instructions

> [!IMPORTANT]
> The roles and policies below are the _minimum_ required for this setup. They do not take into account your organizations security requirements.
 It is your responsibility to secure your environment.

### 1. Enable Workload Identity on your GKE Cluster

Enable WIF on your GKE cluster if needed. GKE Autopilot clusters enable WIF automatically, while Standard clusters must enable it explicitly.

```txt
gcloud container clusters update <CLUSTER_NAME> \
  --location=<REGION> \
  --workload-pool=<PROJECT_NAME>.svc.id.goog
```

Replace `<CLUSTER_NAME>` with the name of your GKE cluster, `<REGION>` with the region your cluster is located in, and `<PROJECT_NAME>` with your GCP
project name.

Example:

```txt
gcloud container clusters update loki-lab-cluster \
  --location=us-west1 \
  --workload-pool=loki-lab.svc.id.goog
```

### 2. Create GCS Bucket(s)

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

### 3. Create a Google Service Account

Create a Google Service Account:

```txt
gcloud iam service-accounts create <GSA_NAME>
```

Replace `<GSA_NAME>` with the name of the service account you want to create.

Example:

```txt
gcloud iam service-accounts create loki-lab-gsa
```

### 4. Bind GSA to Bucket(s)

> [!INFO]
> The [pre-defined `role/storage.objectUser` role](https://cloud.google.com/storage/docs/access-control/iam-roles) is sufficient for Loki / GEL to
 operate. See [IAM permissions for Cloud Storage](https://cloud.google.com/storage/docs/access-control/iam-permissions) for details about each individual
 permission. You can use this predefined role or create your own with matching permissions.

Bind your GSA to your GCS bucket(s):

```txt
gcloud storage buckets add-iam-policy-binding gs://<BUCKET_NAME> \
  --member=serviceAccount:<GSA_NAME>@<PROJECT_NAME>.iam.gserviceaccount.com \
  --role=roles/storage.objectUser
```

Replace `<BUCKET_NAME>` with the name of the bucket(s) created above, `<GSA_NAME>` with the name of the GSA created above, and `<PROJECT_NAME>` with the
 name of your GCP project.

Unlike the `buckets create` command, this command only supports one bucket at a time.

Example:

```txt
gcloud storage buckets add-iam-policy-binding gs://loki-lab-chunks \
  --member=serviceAccount:loki-lab-gsa@loki-lab-project.iam.gserviceaccount.com \
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

### 6. Create Kubernetes Service Account

Create a KSA on your K8s cluster:

```txt
kubectl create serviceaccount <KSA_NAME> \
  --namespace <NAMESPACE>
```

Replace `<KSA_NAME>` with the name of the KSA created above, and `<NAMESPACE>` with the namespace where your Loki/GEL workloads are located.

Example:

```txt
kubectl create serviceaccount loki-lab-ksa \
  --namespace loki
```

### 7. Annotate the KSA to Impersonate the GSA

Add an annotation to the KSA so it can impersonate the GSA:

```txt
kubectl annotate serviceaccount <KSA_NAME> \
  --namespace <NAMESPACE> \
  iam.gke.io/gcp-service-account=<GSA_NAME>@<PROJECT_NAME>.iam.gserviceaccount.com
```

Replace `<KSA_NAME>` with the name of the KSA created above, `<NAMESPACE>` with the namespace where your Loki/GEL workloads are located,
 `<GSA_NAME>` with the name of the GSA created above, and `<PROJECT_NAME>` with the name of your GCP project.

Example:

```txt
kubectl annotate serviceaccount loki-lab-ksa \
  --namespace loki \
  iam.gke.io/gcp-service-account=loki-lab-gsa@loki-lab.iam.gserviceaccount.com
```

For GEL installs, there will be another service account that needs access to the GCS buckets that were created. This KSA is created via the Loki helm chart, so we can't annotate it now. Instead, we'll annotate it via the Helm chart `values.yaml` file, outlined below. You can see where the `tokengen` KSA gets its annotations via the Helm chart [here](https://github.com/grafana/loki/blob/52a8ef8a50397574457ef586722bfec222e914de/production/helm/loki/templates/tokengen/serviceaccount-tokengen.yaml).

While it may be possible to modify the auto-generated `-tokengen` KSA, it's easier to simply bind it in the same way as the KSA you created.

### 8. Bind the GSA to the KSA

<!-- ### TODO: RENAME THIS ### -->

Bind the GSA to the KSA:

```txt
gcloud projects add-iam-policy-binding <PROJECT_NAME> \
  --member="serviceAccount:<GSA_NAME>@<PROJECT_NAME>.iam.gserviceaccount.com" \
  --role="roles/iam.workloadIdentityUser"
```

Replace `<GSA_NAME>` with your GCP project name, and `<PROJECT_NAME>` with the name of the GSA created above.

Example:

```txt
gcloud projects add-iam-policy-binding loki-lab \
  --member="serviceAccount:loki-lab-gsa@loki-lab.iam.gserviceaccount.com" \
  --role="roles/iam.workloadIdentityUser"
```

## Helm

### values.yaml

A Simple Scalable deployment breaks up Loki/GEL operations into distinct `read`, `write`, and `backend` services that can be scaled independently.
 If `enterprise.enabled` is set to `true`, there'll also be a `tokengen` job and an `adminApi` pod that need to connect to the storage backend.

The example below does not represent a complete `values.yaml` file, only the parameters that need to be updated to run Loki/GEL with workload federation.

The `serviceAccount` is automatically applied to `read`, `write`, `backend`, and `adminApi`. The `tokenGen` job creates its own KSA,
which necessitates the additional annotations described the **Annotate the KSA to Impersonate the GSA** step above and reflected in the
 `values.yaml` example below.

<!-- ### TODO: FINALIZE YAML ### -->

```yaml
# GEL ONLY: enterprise.*
enterprise:
  tokengen:
    # Add this annotation to bind the tokengen KSA to the GSA you created.
    # Example: iam.gke.io/gcp-service-account=loki-lab-gsa@loki-lab.iam.gserviceaccount.com
    annotations: 
      - iam.gke.io/gcp-service-account=<GSA_NAME>@<PROJECT_NAME>.iam.gserviceaccount.com

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
