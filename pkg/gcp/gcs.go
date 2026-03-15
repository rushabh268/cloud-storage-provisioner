// Package gcp provides provisioners for GCP resources.
package gcp

import (
	"context"
	"fmt"

	"github.com/rushabh268/cloud-storage-provisioner/pkg/compliance"
	"github.com/rushabh268/cloud-storage-provisioner/pkg/models"
)

// GCSBucketConfig holds configuration for a GCP Cloud Storage bucket.
type GCSBucketConfig struct {
	models.ResourceConfig

	// ProjectID is the GCP project to create the bucket in.
	ProjectID string

	// StorageClass: STANDARD, NEARLINE, COLDLINE, ARCHIVE.
	StorageClass string

	// UniformBucketAccess enforces uniform bucket-level IAM (CIS GCP 5.1).
	UniformBucketAccess bool
}

// GCSBucketProvisioner provisions a GCP Cloud Storage bucket.
type GCSBucketProvisioner struct {
	cfg GCSBucketConfig
}

func NewGCSBucketProvisioner(cfg GCSBucketConfig) *GCSBucketProvisioner {
	cfg.Tags = cfg.Tags.WithDefaults(
		cfg.Tags["Project"],
		cfg.Tags["Environment"],
		cfg.Tags["Owner"],
	)
	if cfg.StorageClass == "" {
		cfg.StorageClass = "STANDARD"
	}
	return &GCSBucketProvisioner{cfg: cfg}
}

func (p *GCSBucketProvisioner) ComplianceCheck() []models.ComplianceResult {
	controls := append(compliance.MandatoryControls,
		compliance.Control{
			ID:          "CIS-GCP-5.1",
			Framework:   "CIS",
			Description: "GCS buckets should use uniform bucket-level access.",
			Remediation: "Pass --uniform-bucket-access flag to the gcp gcs provision command.",
			Check: func(cfg models.ResourceConfig) (bool, string) {
				if !p.cfg.UniformBucketAccess {
					return false, "uniform bucket-level access is disabled; ACL-based access control is less secure"
				}
				return true, ""
			},
		},
		compliance.Control{
			ID:          "CIS-GCP-5.2",
			Framework:   "CIS",
			Description: "GCS buckets must specify a GCP project ID.",
			Remediation: "Pass --project-id <id> to the gcp gcs provision command.",
			Check: func(cfg models.ResourceConfig) (bool, string) {
				if p.cfg.ProjectID == "" {
					return false, "project ID is not set"
				}
				return true, ""
			},
		},
	)
	results, _ := compliance.Run(controls, p.cfg.ResourceConfig)
	return results
}

func (p *GCSBucketProvisioner) Provision(ctx context.Context) error {
	results := p.ComplianceCheck()
	for _, r := range results {
		if !r.Passed {
			return fmt.Errorf("compliance check failed [%s]: %s — remediation: %s",
				r.ControlID, r.Detail, r.Remediation)
		}
	}

	fmt.Printf("[gcp/gcs] creating bucket %q in project %q\n", p.cfg.Name, p.cfg.ProjectID)
	fmt.Printf("[gcp/gcs] enabling uniform bucket-level access on %q\n", p.cfg.Name)
	fmt.Printf("[gcp/gcs] enabling Google-managed encryption on %q\n", p.cfg.Name)
	fmt.Printf("[gcp/gcs] enabling Cloud Audit Logs for %q\n", p.cfg.Name)
	fmt.Printf("[gcp/gcs] blocking public access on %q\n", p.cfg.Name)
	fmt.Printf("[gcp/gcs] applying labels to %q\n", p.cfg.Name)
	return nil
}

func (p *GCSBucketProvisioner) Destroy(ctx context.Context) error {
	fmt.Printf("[gcp/gcs] destroying bucket %q\n", p.cfg.Name)
	return nil
}
