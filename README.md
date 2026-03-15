# cloud-storage-provisioner

A Go CLI tool that provisions cloud storage resources on AWS, Azure, and GCP with automatic
[CIS Benchmark](https://www.cisecurity.org/cis-benchmarks), SOC 2, and HIPAA compliance enforcement.
Compliance checks run before every provisioning operation — if a check fails, the operation is aborted
and remediation guidance is printed.

## Features

- **Multi-cloud:** AWS S3, Azure Storage Accounts, GCP Cloud Storage
- **Compliance-first:** CIS Benchmark, SOC 2, and HIPAA controls checked automatically
- **Secure by default:** encryption at rest, public access blocked, logging enabled on every resource
- **Mandatory tagging:** all resources tagged with `Project`, `Environment`, `Owner`, `ManagedBy`
- **Dry-run mode:** validate compliance without provisioning

## Supported Resources

| Provider | Resource        | Compliance Controls     |
|----------|-----------------|-------------------------|
| aws      | s3              | CIS 1.1, 2.1, 2.6       |
| azure    | storage         | CIS-AZ 3.1, 3.2         |
| gcp      | gcs             | CIS-GCP 5.1, 5.2        |

## Installation

```bash
git clone https://github.com/rushabh268/cloud-storage-provisioner.git
cd cloud-storage-provisioner
make build
```

Requires Go 1.22+.

## Usage

```bash
# Provision a compliant AWS S3 bucket
ccp provision aws s3 \
  --name my-app-data \
  --region us-east-1 \
  --project my-app \
  --environment prod \
  --owner you@example.com \
  --versioning

# Provision an Azure Storage Account
ccp provision azure storage \
  --name myappstorage \
  --region eastus \
  --resource-group my-rg \
  --project my-app \
  --environment prod \
  --owner you@example.com \
  --https-only

# Provision a GCP Cloud Storage bucket
ccp provision gcp gcs \
  --name my-gcs-bucket \
  --region us-central1 \
  --project my-app \
  --environment prod \
  --owner you@example.com

# Dry-run: run compliance checks without provisioning
ccp provision aws s3 --dry-run --name my-bucket ...
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

If a check fails, provisioning is aborted:

```
Running compliance checks...
  [FAIL] CIS-2.1: Public access must not be enabled
         Detail:      AllowPublicAccess is set to true
         Remediation: Remove --allow-public or disable public access in your config

Provisioning aborted: one or more compliance checks failed.
```

## Development

```bash
make build          # compile
make test           # run tests (80% coverage minimum enforced)
golangci-lint run   # lint (zero warnings policy)
make compliance-check PROVIDER=aws RESOURCE=s3  # run CIS checks standalone
```

## Project Structure

```
cloud-storage-provisioner/
├── cmd/                    # CLI entrypoints
├── pkg/
│   ├── aws/                # AWS provisioners
│   ├── azure/              # Azure provisioners
│   ├── gcp/                # GCP provisioners
│   ├── compliance/         # CIS/SOC2/HIPAA check engine
│   └── models/             # Shared resource models and interfaces
└── .github/workflows/      # CI: build, test, lint, security scan
```

## License

[MIT](LICENSE)
