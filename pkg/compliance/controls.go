// Package compliance implements CIS Benchmark, SOC 2, and HIPAA controls
// for each supported cloud resource type.
package compliance

import "github.com/rushabh268/cloud-storage-provisioner/pkg/models"

// Control is a single compliance check.
type Control struct {
	ID          string
	Framework   string // "CIS", "SOC2", "HIPAA"
	Description string
	Remediation string
	// Check returns (passed, detail). detail explains why a check failed.
	Check func(cfg models.ResourceConfig) (bool, string)
}

// MandatoryControls are controls that abort provisioning if they fail.
// These correspond to CIS Level 1 controls.
var MandatoryControls = []Control{
	{
		ID:          "CIS-1.1",
		Framework:   "CIS",
		Description: "All resources must have mandatory tags: Project, Environment, Owner, ManagedBy.",
		Remediation: "Pass --tag Project=<name> --tag Environment=<env> --tag Owner=<email> to the provision command.",
		Check: func(cfg models.ResourceConfig) (bool, string) {
			required := []string{"Project", "Environment", "Owner", "ManagedBy"}
			for _, key := range required {
				if cfg.Tags[key] == "" {
					return false, "missing required tag: " + key
				}
			}
			return true, ""
		},
	},
	{
		ID:          "CIS-2.1",
		Framework:   "CIS",
		Description: "Public access must not be enabled on storage resources without explicit override.",
		Remediation: "Remove --allow-public flag or justify its use in the PR description.",
		Check: func(cfg models.ResourceConfig) (bool, string) {
			if cfg.AllowPublic {
				return false, "AllowPublic is set; this resource will be publicly accessible"
			}
			return true, ""
		},
	},
}

// Run executes all provided controls against cfg and returns the full result set.
// If any mandatory control fails, the first failure is returned as an error.
func Run(controls []Control, cfg models.ResourceConfig) ([]models.ComplianceResult, error) {
	results := make([]models.ComplianceResult, 0, len(controls))
	for _, ctrl := range controls {
		passed, detail := ctrl.Check(cfg)
		results = append(results, models.ComplianceResult{
			ControlID:   ctrl.ID,
			Description: ctrl.Description,
			Passed:      passed,
			Detail:      detail,
			Remediation: ctrl.Remediation,
		})
	}
	return results, nil
}
