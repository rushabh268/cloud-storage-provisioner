# Progress: AWS S3 Provisioner

**Context saved before compaction — read this first after resuming.**

## What was completed

- `pkg/models/resource.go` — `Provisioner` interface, `ResourceConfig`, `Tags`, `ComplianceResult`
- `pkg/compliance/controls.go` — `Control` struct, `MandatoryControls` (CIS-1.1 tags, CIS-2.1 public access), `Run()`
- `pkg/aws/s3.go` — `S3Provisioner` implementing `Provisioner`, compliance gate in `Provision()`
- `pkg/aws/s3_test.go` — 6 table-driven tests, 92% coverage

## What is next

- [ ] Add `pkg/aws/vpc.go` — VPC with private subnets, flow logs enabled (CIS-3.9)
- [ ] Add `pkg/aws/rds.go` — RDS with encryption, Multi-AZ, no public endpoint (CIS-2.3.1)
- [ ] Wire up cobra commands in `cmd/` so provisioners are callable from CLI
- [ ] Add `pkg/azure/vm.go` and `pkg/gcp/compute.go`

## Decisions made

- Compliance check runs inside `Provision()`, not as a separate pre-step, so callers cannot skip it.
- `Tags.WithDefaults()` always adds `ManagedBy=cloud-storage-provisioner` — callers cannot forget it.
- `AllowPublic` defaults to false and must be explicitly set; this is enforced by the CIS-2.1 control.

## Open questions

- Should the CLI print compliance results in JSON format for CI pipeline consumption? Consider `--output json` flag.
- Should we support importing existing resources into state, or provision-only for now?
