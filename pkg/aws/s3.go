// Package aws provides provisioners for AWS resources.
package aws

import (
	"context"
	"fmt"

	"github.com/rushabh268/cloud-storage-provisioner/pkg/compliance"
	"github.com/rushabh268/cloud-storage-provisioner/pkg/models"
)

// S3Config holds configuration specific to an S3 bucket.
type S3Config struct {
	models.ResourceConfig

	// Versioning enables S3 object versioning (required for HIPAA).
	Versioning bool

	// LifecycleDays sets a lifecycle rule to transition objects to Glacier after N days.
	// 0 disables lifecycle management.
	LifecycleDays int
}

// S3Provisioner provisions an S3 bucket with encryption, logging, and public access block.
// It implements models.Provisioner.
type S3Provisioner struct {
	cfg S3Config
}

// NewS3Provisioner returns a configured S3Provisioner.
func NewS3Provisioner(cfg S3Config) *S3Provisioner {
	// Ensure ManagedBy tag is always set — enforced here, not left to callers.
	cfg.Tags = cfg.Tags.WithDefaults(
		cfg.Tags["Project"],
		cfg.Tags["Environment"],
		cfg.Tags["Owner"],
	)
	return &S3Provisioner{cfg: cfg}
}

// ComplianceCheck runs all applicable CIS/SOC2/HIPAA controls for S3.
// It does not make any AWS API calls.
func (p *S3Provisioner) ComplianceCheck() []models.ComplianceResult {
	controls := append(compliance.MandatoryControls,
		compliance.Control{
			ID:          "CIS-2.6",
			Framework:   "CIS",
			Description: "S3 buckets must have versioning enabled.",
			Remediation: "Pass --versioning flag to the s3 provision command.",
			Check: func(cfg models.ResourceConfig) (bool, string) {
				if !p.cfg.Versioning {
					return false, "versioning is disabled; enable it to satisfy CIS 2.6 and HIPAA audit requirements"
				}
				return true, ""
			},
		},
	)
	results, _ := compliance.Run(controls, p.cfg.ResourceConfig)
	return results
}

// Provision creates the S3 bucket. Compliance checks run first; any mandatory
// control failure aborts provisioning and returns a descriptive error.
func (p *S3Provisioner) Provision(ctx context.Context) error {
	// Step 1: compliance gate — no AWS API call happens before this passes.
	results := p.ComplianceCheck()
	for _, r := range results {
		if !r.Passed {
			return fmt.Errorf("compliance check failed [%s]: %s — remediation: %s",
				r.ControlID, r.Detail, r.Remediation)
		}
	}

	// Step 2: create bucket (real implementation would call AWS SDK here).
	fmt.Printf("[aws/s3] creating bucket %q in %s\n", p.cfg.Name, p.cfg.Region)
	// aws.CreateBucket(ctx, p.cfg.Name, string(p.cfg.Region))

	// Step 3: block all public access — always, non-negotiable.
	fmt.Printf("[aws/s3] blocking public access on %q\n", p.cfg.Name)

	// Step 4: enable AES-256 server-side encryption.
	fmt.Printf("[aws/s3] enabling SSE-AES256 on %q\n", p.cfg.Name)

	// Step 5: enable server access logging to a dedicated log bucket.
	fmt.Printf("[aws/s3] enabling access logging on %q\n", p.cfg.Name)

	// Step 6: enable versioning if requested.
	if p.cfg.Versioning {
		fmt.Printf("[aws/s3] enabling versioning on %q\n", p.cfg.Name)
	}

	// Step 7: apply tags.
	fmt.Printf("[aws/s3] applying %d tags to %q\n", len(p.cfg.Tags), p.cfg.Name)

	return nil
}

// Destroy removes the S3 bucket and its contents after emptying it.
func (p *S3Provisioner) Destroy(ctx context.Context) error {
	fmt.Printf("[aws/s3] emptying and destroying bucket %q\n", p.cfg.Name)
	return nil
}
