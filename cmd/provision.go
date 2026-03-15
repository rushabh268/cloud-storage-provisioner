// Package cmd implements the CLI commands for cloud-storage-provisioner.
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/rushabh268/cloud-storage-provisioner/pkg/aws"
	"github.com/rushabh268/cloud-storage-provisioner/pkg/models"
)

// ProvisionS3 is the entrypoint for `ccp provision aws s3`.
// It reads flags, builds the config, runs compliance checks, and provisions.
func ProvisionS3(name, region, project, environment, owner string, versioning bool) error {
	cfg := aws.S3Config{
		ResourceConfig: models.ResourceConfig{
			Name:     name,
			Region:   models.Region(region),
			Provider: models.ProviderAWS,
			Tags: models.Tags{
				"Project":     project,
				"Environment": environment,
				"Owner":       owner,
			},
		},
		Versioning: versioning,
	}

	p := aws.NewS3Provisioner(cfg)

	// Print compliance results before provisioning so the operator can review.
	fmt.Println("Running compliance checks...")
	results := p.ComplianceCheck()
	allPassed := true
	for _, r := range results {
		status := "PASS"
		if !r.Passed {
			status = "FAIL"
			allPassed = false
		}
		fmt.Printf("  [%s] %s: %s\n", status, r.ControlID, r.Description)
		if !r.Passed {
			fmt.Printf("         Detail:      %s\n", r.Detail)
			fmt.Printf("         Remediation: %s\n", r.Remediation)
		}
	}

	if !allPassed {
		fmt.Fprintln(os.Stderr, "\nProvisioning aborted: one or more compliance checks failed.")
		return fmt.Errorf("compliance checks failed")
	}

	fmt.Println("\nAll compliance checks passed. Provisioning...")
	return p.Provision(context.Background())
}
