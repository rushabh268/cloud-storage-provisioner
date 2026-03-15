# CLAUDE.md — cloud-storage-provisioner

This file is the persistent context for Claude Code working in this repository.
It is checked into git and updated by the team whenever Claude makes a mistake or a new pattern is established.
Keep this file under 200 lines. Every rule here exists because something went wrong without it.

---

## Project Overview

`cloud-storage-provisioner` is a Go CLI tool and Claude Skill that provisions cloud resources on AWS,
Azure, and GCP following CIS Benchmark, SOC 2, and HIPAA compliance controls.

**Language:** Go 1.22
**Style:** Standard Go idioms. No framework magic. Prefer `errors.Is` / `errors.As` over string matching.
**Build:** `make build` — must pass before any commit.
**Tests:** `make test` — 80% coverage minimum enforced by CI.
**Lint:** `golangci-lint run` — zero warnings policy.

---

## Non-Negotiable Rules

### Code Quality
- **Never add abstraction layers** unless they are used in at least two call sites. The default is flat, readable code.
- **Never remove comments** you do not understand. If a comment seems wrong, flag it — do not silently delete it.
- **Dead code must be deleted**, not commented out. If it should be kept, add a `// TODO(username):` explaining why.
- **Error handling is explicit.** Never swallow errors with `_`. If an error is intentionally ignored, add `//nolint:errcheck // reason`.
- **Keep functions under 50 lines.** If a function grows past this, split it before continuing.

### Testing
- Every exported function must have a test. Use `testify/assert` — not raw `t.Error`.
- Table-driven tests are preferred. Avoid copy-pasted test cases.
- Mocks go in `*_mock_test.go` files, not in production code.
- When adding a feature, add the test first, then the implementation. Claude should never skip this.

### Cloud & Compliance
- **All resources must be tagged** with at minimum: `Project`, `Environment`, `Owner`, `ManagedBy=cloud-storage-provisioner`.
- **Encryption at rest is always enabled.** Never create an S3 bucket, Azure Storage Account, or GCP bucket without encryption.
- **Public access is always blocked** on storage resources unless the user explicitly passes `--allow-public` and confirms.
- **Logging must be enabled** on all resources. CloudTrail for AWS, Azure Monitor for Azure, Cloud Audit Logs for GCP.
- **CIS Benchmark checks run automatically** before resource creation. If a check fails, abort and print the failing control.

### Git
- Commit messages follow Conventional Commits: `feat(aws): add VPC with private subnets`.
- **Never commit directly to main.** All changes go through a PR with at least one review.
- CLAUDE.md is updated as part of the PR if a new pattern was established.

---

## Project Structure

```
cloud-storage-provisioner/
├── cmd/                    # CLI entrypoints (cobra commands)
├── pkg/
│   ├── aws/                # AWS provisioners (EC2, S3, VPC, RDS)
│   ├── azure/              # Azure provisioners (VM, Storage, VNet)
│   ├── gcp/                # GCP provisioners (Compute, GCS, VPC)
│   ├── compliance/         # CIS/SOC2/HIPAA check engine
│   └── models/             # Shared resource models and interfaces
├── .claude/commands/       # Reusable slash commands for this repo
├── .github/workflows/      # CI: build, test, lint, security scan
└── CLAUDE.md               # This file
```

---

## Common Tasks

### Add a new AWS resource provisioner
1. Create `pkg/aws/<resource>.go` implementing `models.Provisioner`.
2. Add compliance checks in `pkg/compliance/<resource>_checks.go`.
3. Add a cobra command in `cmd/<resource>.go`.
4. Add tests with at least 80% coverage.
5. Update `CLAUDE.md` if any new patterns were established.

### Run the full compliance check
```bash
make compliance-check PROVIDER=aws RESOURCE=s3
```

### Add a new CIS control
Controls live in `pkg/compliance/controls.go`. Each control is a struct with ID, description,
remediation text, and a check function. The check function receives the resource config and returns
`(passed bool, detail string)`.

---

## Context Compaction Checkpoint

When context is getting long (Claude will note this), save progress to:
`docs/progress/<feature-name>.md` before compacting. This file should contain:
- What was completed
- What is next
- Any decisions made and why
- Open questions

After compaction, Claude should read this file first.

---

## What Not to Do (Learned from Mistakes)

- **Do not use `interface{}` anywhere.** Use typed interfaces or generics.
- **Do not call `os.Exit` outside of `main()`.** Return errors up the call stack.
- **Do not hardcode region strings.** Use the `models.Region` type.
- **Do not create resources without running compliance checks first.** The `compliance.Check()` call in `Provision()` is not optional.
- **Do not use `sync.Mutex` directly.** Use the `sync.RWMutex` wrapper in `pkg/models/safe_map.go`.
- **Do not write multi-cloud logic in a single function.** AWS, Azure, and GCP provisioners are always separate files.

---

## Resources
- CIS Benchmarks: https://www.cisecurity.org/cis-benchmarks
- AWS Well-Architected: https://aws.amazon.com/architecture/well-architected/
- Azure Security Baseline: https://learn.microsoft.com/en-us/security/benchmark/azure/
- GCP Security Best Practices: https://cloud.google.com/security/best-practices
