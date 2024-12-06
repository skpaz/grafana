# AWS S3 authentication via Pod Identities for Grafana Loki

## Assumptions

- You're familiar with AWS EKS, IAM, and S3.
- You're installing Loki/Grafana Enterprise Logs (GEL) on an EKS cluster in Auto Mode.
- You're deploying Loki/GEL in Single Scalable mode.
- You're familiar with `awscli`; this doc doesn't contain any Console examples.

## Resources

- [Learn how EKS Pod Identity grants pods access to AWS services](https://docs.aws.amazon.com/eks/latest/userguide/pod-identities.html)

## Pod Identities

There are several methods we can use to authenticate to S3.

1. Hard-code an IAM user's access key and secret in the `storage` section of the Helm chart.
2. [IAM roles for service accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) via OIDC.
3. [Pod-level access to AWS services](https://docs.aws.amazon.com/eks/latest/userguide/pod-identities.html) via Pod Identities.

The first option is is not recommended since it introduces the risk that those credentials could be leaked. This document
 covers **Pod-level access to AWS services via Pod Identities**.

## Instructions

> [!IMPORTANT]
> The roles and policies below are the _minimum_ required for this setup. They do not take into account your organizations security requirements.
 It is your responsibility to secure your environment.

### 1. Create S3 bucket(s)

Loki/GEL needs an S3 bucket for admin, chunks, and ruler. The command below will create a private S3 bucket in `<AWS_REGION>` with
 _Bucket Owner Enforced_ object ownership.

```txt
aws s3api create-bucket \
  --bucket <BUCKET_NAME> \
  --create-bucket-configuration LocationConstraint=<AWS_REGION>
```

Replace `<BUCKET_NAME>` with the desired bucket name and `<AWS_REGION>` with the region where you'd like to create the bucket. If you don't specify
 a region, it will default to us-east-1.

Example:

```txt
aws s3api create-bucket \
  --bucket loki-lab-admin \
  --create-bucket-configuration LocationConstraint=us-west-1

aws s3api create-bucket \
  --bucket loki-lab-chunks \
  --create-bucket-configuration LocationConstraint=us-west-1

aws s3api create-bucket \
  --bucket loki-lab-ruler \
  --create-bucket-configuration LocationConstraint=us-west-1
```

Repeat for each bucket.

### 2. Create an IAM access policy

#### IAM access policy content

First, we need to create the content of the policy. Create and save `s3-access-policy.json` with the contents below:

```json
{ 
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:ListBucket",
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeletObject"
      ],
      "Effect": "Allow",
      "Resource": [
        "arn:aws:s3:::<BUCKET_NAME>",
        "arn:aws:s3:::<BUCKET_NAME>/*",
        "arn:aws:s3:::<BUCKET_NAME>",
        "arn:aws:s3:::<BUCKET_NAME>/*",
        "arn:aws:s3:::<BUCKET_NAME>",
        "arn:aws:s3:::<BUCKET_NAME>/*"
      ]
    }
  ]
}
```

Replace each pair of `arn:aws:s3::` resources with the name of your admin, chunk, and rule buckets.

Example:

```txt
cat << EOF > s3-access-policy.json
{ 
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:ListBucket",
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeletObject"
      ],
      "Effect": "Allow",
      "Resource": [
        "arn:aws:s3:::loki-lab-admin",
        "arn:aws:s3:::loki-lab-admin/*",
        "arn:aws:s3:::loki-lab-chunks",
        "arn:aws:s3:::loki-lab-chunks/*",
        "arn:aws:s3:::loki-lab-ruler",
        "arn:aws:s3:::loki-lab-ruler/*"
      ]
    }
  ]
}
EOF
```

#### Create the IAM policy

```txt
aws iam create-policy \
  --policy-name <POLICY_NAME> \
  --policy-document file://<POLICY_JSON_FILE>
```

Replace `<POLICY_NAME>` with your desired policy name and `<POLICY_JSON_FILE>` with the path to the JSON file created in the previous step.

Example:

```txt
aws iam create-policy \
  --policy-name loki-lab-s3-access-policy \
  --policy-document file://s3-access-policy.json
```

Capture the policy ARN that the CLI returns. We'll need it later.

### 3. Create an IAM role and trust policy

#### Create trust policy content

Create and save `pods-ident-trust-policy.json` with the contents below:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "sts:AssumeRole",
        "sts:TagSession"
      ],
      "Effect": "Allow",
      "Principal": {
        "Service": "pods.eks.amazonaws.com"
      }
    }
  ]
}
```

Example:

```txt
cat << EOF > pods-ident-trust-policy.json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "sts:AssumeRole",
        "sts:TagSession"
      ],
      "Effect": "Allow",
      "Principal": {
        "Service": "pods.eks.amazonaws.com"
      }
    }
  ]
}
EOF
```

#### Create IAM role for the KSAs

Create a role for the KSAs to use:

```txt
aws iam create-role \
  --role-name <ROLE_NAME> \
  --assume-role-policy-document file://<POLICY_JSON_FILE>
```

Replace `<ROLE_NAME>` with your desired role name, and `<POLICY_JSON_FILE>` with the path to the JSON files created in the previous step.

Example:

```txt
aws iam create-role \
  --role-name loki-lab-sa-role \
  --assume-role-policy-document file://pods-ident-trust-policy.json
```

Capture the role's ARN once it's created. We'll use it in a later step.

#### Attach IAM access policy to role

Attach the access policy we created earlier to the KSA IAM role:

```txt
aws iam attach-role-policy \
  --role-name <ROLE_NAME> \
  --policy-arn <POLICY_ARN>
```

Replace `<ROLE_NAME>` with the name of the role created in the last step, and `<POLICY_ARN>` with the S3 access policy ARN you saved from the
 **Create an IAM access policy** step above.

Example:

```txt
aws iam attach-role-policy \
  --role-name loki-lab-sa-role \
  --policy-arn arn:aws:iam::123456789000:policy/loki-lab-s3-access-policy
```

### 4. Create Pod Identity Association

Create a Pod Identity Association for each KSA:

```txt
aws eks create-pod-identity-association \
--cluster-name <CLUSTER_NAME> \
--region <AWS_REGION> \
--namespace <NAMESPACE> \
--role-arn <ROLE_ARN> \
--service-account <KSA_NAME>
```

Replace `<CLUSTER_NAME>` with the name of your cluster, `<AWS_REGION>` with the name of the region where the cluster is located, `<NAMESPACE>` with
 the desired namespace, `<ROLE_ARN>` with the role ARN that was captured in the **Create IAM role for the KSAs** step above. Last, decide on a name for
 the general KSA you want to create later and replace `<KSA_NAME>` with it.

Example:

```txt
aws eks create-pod-identity-association \
--cluster-name loki-lab-cluster \
--region us-west-1 \
--namespace gel \
--role-arn arn:aws:iam::123456789000:policy/loki-lab-s3-access-policy \
--service-account loki-lab-sa

aws eks create-pod-identity-association \
--cluster-name loki-lab-cluster \
--region us-west-1 \
--namespace gel \
--role-arn arn:aws:iam::123456789000:policy/loki-lab-s3-access-policy \
--service-account enterprise-logs-tokengen
```

> [!NOTE]
> The `tokengen` KSA names are unique to Loki and GEL. Loki is `loki-tokengen`, while GEL is `enterprise-logs-tokengen`.

## Helm

### values.yaml

A Simple Scalable deployment breaks up Loki/GEL operations into distinct `read`, `write`, and `backend` services that can be scaled independently.
 If `enterprise.enabled` is set to `true`, there'll also be a `tokengen` job and an `adminApi` pod that need to connect to the storage backend.

The example below does not represent a complete `values.yaml` file, only the parameters that need to be updated.

<!-- ### TODO: FINALIZE YAML ### -->

```yaml
serviceAccount:
  create: true
  name: FIXME_KSA_NAME

loki:
  schemaConfig:
    configs:
      - from: 2020-07-01
        store: tsdb
        object_store: s3
        schema: v13
        index:
          prefix: index_
          period: 24h

  storage:
    type: s3
    s3:
      region: FIXME_AWS_REGION
    bucketNames:
      chunks: FIXME_CHUNKS_BUCKET_NAME
      ruler: FIXME_RULER_BUCKET_NAME
      # GEL ONLY: bucketNames.admin
      admin: FIXME_ADMIN_BUCKET_NAME
    type: s3
```

Replace `FIXME_KSA_NAME` with the name of the KSA created in the **Create trust policy content** step above, `FIXME_AWS_REGION` with the AWS region where
 your cluster is located, and each `FIXME_*CHUNKS*_BUCKET_NAME` with the name of the buckets created in the **Create S3 Bucket(s)** step above.
