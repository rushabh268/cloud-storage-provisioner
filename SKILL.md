---
name: cloud-provisioner
description: >
  Provision cloud resources (AWS S3, Azure Storage, GCP GCS, and more) with automatic CIS Benchmark,
  SOC 2, and HIPAA compliance checks. Aborts if any mandatory control fails and prints remediation guidance.
  Use for: creating new cloud resources, running compliance checks on planned resources, or auditing
  existing infrastructure against CIS controls.
---

# Claude Skill: Cloud Provisioner

Provisions cloud resources on AWS, Azure, and GCP with built-in compliance enforcement.

## Usage

```
/provision <provider> <resource-type> [flags]
```

## Supported Resources

| Provider | Resource         | Compliance Controls        |
|----------|-----------------|---------------------------|
| aws      | s3              | CIS 1.1, 2.1, 2.6         |
| azure    | storage         | CIS-AZ 3.1, 3.2            |
| gcp      | gcs             | CIS-GCP 5.1, 5.2           |

## Examples

```bash
# Provision a compliant S3 bucket
/provision aws s3 \
  --name my-app-data \
  --region us-east-1 \
  --project my-app \
  --environment prod \
  --owner rushabh268@gmail.com \
  --versioning

# Provision an Azure Storage Account
/provision azure storage \
  --name myappstorage \
  --region eastus \
  --resource-group my-rg \
  --project my-app \
  --environment prod \
  --owner rushabh268@gmail.com \
  --https-only

# Run compliance check without provisioning
/provision aws s3 --dry-run --name my-bucket ...
```

## Compliance Output

```
Running compliance checks...
  [PASS] CIS-1.1: All resources must have mandatory tags
  [PASS] CIS-2.1: Public access must not be enabled
  [PASS] CIS-2.6: S3 buckets must have versioning enabled

All compliance checks passed. Provisioning...
[aws/s3] creating bucket "my-app-data" in us-east-1
[aws/s3] blocking public access on "my-app-data"
[aws/s3] enabling SSE-AES256 on "my-app-data"
[aws/s3] enabling access logging on "my-app-data"
[aws/s3] enabling versioning on "my-app-data"
[aws/s3] applying 4 tags to "my-app-data"
```
