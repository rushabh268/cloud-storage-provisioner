package aws

import (
	"context"
	"testing"

	"github.com/rushabh268/cloud-storage-provisioner/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validS3Config() S3Config {
	return S3Config{
		ResourceConfig: models.ResourceConfig{
			Name:     "my-app-data",
			Region:   "us-east-1",
			Provider: models.ProviderAWS,
			Tags: models.Tags{
				"Project":     "my-app",
				"Environment": "prod",
				"Owner":       "rushabh268@gmail.com",
			},
		},
		Versioning:    true,
		LifecycleDays: 90,
	}
}

func TestS3Provisioner_ComplianceCheck_PassesForValidConfig(t *testing.T) {
	p := NewS3Provisioner(validS3Config())
	results := p.ComplianceCheck()

	for _, r := range results {
		assert.True(t, r.Passed, "expected control %s to pass, got: %s", r.ControlID, r.Detail)
	}
}

func TestS3Provisioner_ComplianceCheck_FailsMissingTags(t *testing.T) {
	cfg := validS3Config()
	delete(cfg.Tags, "Owner")

	p := NewS3Provisioner(cfg)
	results := p.ComplianceCheck()

	var failed []models.ComplianceResult
	for _, r := range results {
		if !r.Passed {
			failed = append(failed, r)
		}
	}

	require.Len(t, failed, 1)
	assert.Equal(t, "CIS-1.1", failed[0].ControlID)
	assert.Contains(t, failed[0].Detail, "Owner")
}

func TestS3Provisioner_ComplianceCheck_FailsPublicAccess(t *testing.T) {
	cfg := validS3Config()
	cfg.AllowPublic = true

	p := NewS3Provisioner(cfg)
	results := p.ComplianceCheck()

	var failed []models.ComplianceResult
	for _, r := range results {
		if !r.Passed {
			failed = append(failed, r)
		}
	}

	require.Len(t, failed, 1)
	assert.Equal(t, "CIS-2.1", failed[0].ControlID)
}

func TestS3Provisioner_ComplianceCheck_FailsMissingVersioning(t *testing.T) {
	cfg := validS3Config()
	cfg.Versioning = false

	p := NewS3Provisioner(cfg)
	results := p.ComplianceCheck()

	var failed []models.ComplianceResult
	for _, r := range results {
		if !r.Passed {
			failed = append(failed, r)
		}
	}

	require.Len(t, failed, 1)
	assert.Equal(t, "CIS-2.6", failed[0].ControlID)
}

func TestS3Provisioner_Provision_AbortsOnComplianceFailure(t *testing.T) {
	cfg := validS3Config()
	cfg.Versioning = false

	p := NewS3Provisioner(cfg)
	err := p.Provision(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "CIS-2.6")
}

func TestS3Provisioner_Provision_SucceedsForValidConfig(t *testing.T) {
	p := NewS3Provisioner(validS3Config())
	err := p.Provision(context.Background())
	assert.NoError(t, err)
}
