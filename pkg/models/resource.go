// Package models defines shared interfaces and types used across all cloud providers.
package models

import "context"

// Region is a typed string for cloud regions to prevent passing arbitrary strings.
type Region string

// Provider identifies a cloud provider.
type Provider string

const (
	ProviderAWS   Provider = "aws"
	ProviderAzure Provider = "azure"
	ProviderGCP   Provider = "gcp"
)

// Tags are required on every provisioned resource (CIS control: asset management).
// ManagedBy is always set to "cloud-storage-provisioner" automatically.
type Tags map[string]string

func (t Tags) WithDefaults(project, environment, owner string) Tags {
	out := make(Tags, len(t)+4)
	for k, v := range t {
		out[k] = v
	}
	out["Project"] = project
	out["Environment"] = environment
	out["Owner"] = owner
	out["ManagedBy"] = "cloud-storage-provisioner"
	return out
}

// ResourceConfig holds the common configuration shared by all resource types.
type ResourceConfig struct {
	Name        string
	Region      Region
	Provider    Provider
	Tags        Tags
	AllowPublic bool // Must be explicitly set; defaults to false.
}

// ComplianceResult reports the outcome of a compliance check.
type ComplianceResult struct {
	ControlID   string
	Description string
	Passed      bool
	Detail      string
	Remediation string
}

// Provisioner is the interface every cloud resource provisioner must implement.
// Compliance checks run inside Provision() before any API call is made.
type Provisioner interface {
	// Provision creates the resource. It runs compliance checks internally
	// and returns an error if any mandatory control fails.
	Provision(ctx context.Context) error

	// Destroy removes the resource.
	Destroy(ctx context.Context) error

	// ComplianceCheck runs all applicable controls and returns results.
	// It does not provision anything.
	ComplianceCheck() []ComplianceResult
}
