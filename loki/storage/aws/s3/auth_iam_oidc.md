# AWS S3 authentication via IAM OIDC for Grafana Loki

## Assumptions

- You're familiar with AWS EKS, IAM, and S3.
- You're installing Loki/Grafana Enterprise Logs (GEL) on an EKS cluster in Auto Mode.
- You're deploying Loki/GEL in Single Scalable mode.
- You're familiar with `awscli` and `eksctl`; this doc doesn't contain any Console examples.
- You're familiar with `kubectl`.
- Your kubeconfig is [properly set up](https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html).

## Resources

- [IAM roles for service accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html)

## IAM OIDC

There are several methods we can use to authenticate to S3.

1. Hard-code an IAM user's access key and secret in the `storage` section of the Helm chart.
2. [IAM roles for service accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) via OIDC.
3. [Pod-level access to AWS services](https://docs.aws.amazon.com/eks/latest/userguide/pod-identities.html) via Pod Identities.

The first option is is not recommended since it introduces the risk that those credentials could be leaked. This document
 covers **IAM roles for service accounts via OIDC**.

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

### 2. Create IAM access policy

#### Create access policy content

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

### 3. Verify or create an IAM OIDC provider

#### Check for an existing OIDC provider

See if your EKS cluster has an OIDC provider set up:

```txt
CLUSTER_NAME=<CLUSTER_NAME> && \
OIDC_ID=$(aws eks describe-cluster \
  --name ${CLUSTER_NAME} \
  --region <AWS_REGION> \
  --query "cluster.identity.oidc.issuer" \
  --output text | cut -d '/' -f 5) && \
aws iam list-open-id-connect-providers | grep ${OIDC_ID} | sed -nr 's/.+"Arn":.+"(.+)"/\n\1/p'
```

Replace `<CLUSTER_NAME>` with your cluster name and `<AWS_REGION>` with the region the cluster is located in.

Example:

```txt
CLUSTER_NAME=loki-lab-cluster && \
OIDC_ID=$(aws eks describe-cluster \
  --name ${CLUSTER_NAME} \
  --region us-west-1 \
  --query "cluster.identity.oidc.issuer" \
  --output text | cut -d '/' -f 5) && \
aws iam list-open-id-connect-providers | grep ${OIDC_ID} | sed -nr 's/.+"Arn":.+"(.+)"/\n\1/p'
```

The command should return a provider ARN that matches this pattern:

```txt
arn:aws:iam::{AWS_ACCOUNTID}:oidc-provider/oidc.eks.{AWS_REGION}.amazonaws.com/id/{OIDC_PROVIDER_ID}
```

If an ARN is not returned, we'll need to create a provider. Otherwise, you can proceed to **4. Create an IAM role and trust policy**.

#### Create an IAM OIDC provider if one does not exist

`eksctl` makes this simple:

```txt
eksctl utils associate-iam-oidc-provider \
  --cluster <CLUSTER_NAME> \
  --region <AWS_REGOIN> \
  --approve
```

Replace `<CLUSTER_NAME>` with your cluster name and `<AWS_REGION>` with the region the cluster is located in.

Example:

```txt
eksctl utils associate-iam-oidc-provider \
  --cluster loki-lab-cluster \
  --region us-west-1 \
  --approve
```

You'll see output that looks like this:

```txt
2024-12-04 11:04:55 [ℹ]  will create IAM Open ID Connect provider for cluster "loki-lab-cluster" in "us-west-1"
2024-12-04 11:04:55 [✔]  created IAM Open ID Connect provider for cluster "loki-lab-cluster" in "us-west-1"
```

Re-run the OIDC provider check commands from the previous step. You should now see a provider ARN.

### 4. Create an IAM role and trust policy

Loki/GEL uses two Kubernetes service accounts: a general service account that applies to Loki's various services, and a specific service account that is
 used by the `tokengen` job. We can accomodate both KSAs with one role. They can also be split into seperate roles if desired; that is not covered here.

#### Create trust policy content

Create and save `iam-oidc-trust-policy.json` with the contents below:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "<OIDC_PROVIDER_ARN>"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "oidc.eks.<AWS_REGION>.amazonaws.com/id/<OIDC_PROVIDER_ID>:aud": "sts.amazonaws.com",
          "oidc.eks.<AWS_REGION>.amazonaws.com/id/<OIDC_PROVIDER_ID>:sub": [
            "system:serviceaccount:<NAMESPACE>:<KSA_NAME>",
            "system:serviceaccount:<NAMESPACE>:enterprise-logs-tokengen"
          ]
        }
      }
    }
  ]
}
```

Replace `<OIDC_PROVIDER_ARN>` with the OIDC provider ARN, `<AWS_REGION>` with the region the cluster is located in, `<OIDC_PROVIDER_ID>` with the
 OIDC provider ID (which is the bit after `id/` in the ARN), and `<NAMESPACE>` with the Kubernetes namespace. Decide on a name for the general KSA
 you want to create later, and replace `<KSA_NAME>` with it. The name for the `tokengen` KSA is pre-determined.

> [!NOTE]
> The `tokengen` KSA names are unique to Loki and GEL. Loki is `loki-tokengen`, while GEL is `enterprise-logs-tokengen`.

Example:

```txt
cat << EOF > iam-oidc-trust-policy.json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::123456789000:oidc-provider/oidc.eks.us-west-1.amazonaws.com/id/XL456F6J1NX0VD04QMMJU47PUVY082K7"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "oidc.eks.us-west-1.amazonaws.com/id/XL456F6J1NX0VD04QMMJU47PUVY082K7:aud": "sts.amazonaws.com",
          "oidc.eks.us-west-1.amazonaws.com/id/XL456F6J1NX0VD04QMMJU47PUVY082K7:sub": [
            "system:serviceaccount:gel:loki-lab-sa",
            "system:serviceaccount:gel:enterprise-logs-tokengen"
          ]
        }
      }
    }
  ]
}
EOF
```

#### Create an IAM role for the KSAs

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
  --assume-role-policy-document file://iam-oidc-trust-policy.json
```

Capture the ARN of the role created. We'll use it later in the Helm chart.

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

## Helm

### values.yaml

A Simple Scalable deployment breaks up Loki/GEL operations into distinct `read`, `write`, and `backend` services that can be scaled independently.
 If `enterprise.enabled` is set to `true`, there'll also be a `tokengen` job and an `adminApi` pod that need to connect to the storage backend.

The example below does not represent a complete `values.yaml` file, only the parameters that need to be updated.

<!-- ### TODO: FINALIZE YAML ### -->

```yaml
# GEL ONLY: enterprise.*
enterprise:
  tokengen:
    annotations:
      "eks.amazonaws.com/role-arn": "FIX_ME_ROLE_ARN"

serviceAccount:
  create: true
  name: FIXME_KSA_NAME
  annotations:
    "eks.amazonaws.com/role-arn": "FIX_ME_ROLE_ARN"

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
 your cluster is located, and each `FIXME_*CHUNKS*_BUCKET_NAME` with the name of the buckets created in the Create GCS Bucket(s) step above.
