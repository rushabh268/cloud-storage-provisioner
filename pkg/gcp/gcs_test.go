package gcp

import (
	"context"
	"testing"

	"github.com/rushabh268/cloud-storage-provisioner/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validGCSConfig() GCSBucketConfig {
	return GCSBucketConfig{
		ResourceConfig: models.ResourceConfig{
			Name:     "my-app-bucket",
			Region:   "us-central1",
			Provider: models.ProviderGCP,
			Tags: models.Tags{
				"Project":     "my-app",
				"Environment": "prod",
				"Owner":       "rushabh268@gmail.com",
			},
		},
		ProjectID:           "my-gcp-project",
		StorageClass:        "STANDARD",
		UniformBucketAccess: true,
	}
}

func TestGCSBucket_ComplianceCheck_PassesForValidConfig(t *testing.T) {
	p := NewGCSBucketProvisioner(validGCSConfig())
	results := p.ComplianceCheck()
	for _, r := range results {
		assert.True(t, r.Passed, "expected control %s to pass, got: %s", r.ControlID, r.Detail)
	}
}

func TestGCSBucket_ComplianceCheck_FailsMissingTags(t *testing.T) {
	cfg := validGCSConfig()
	delete(cfg.Tags, "Environment")

	p := NewGCSBucketProvisioner(cfg)
	var failed []models.ComplianceResult
	for _, r := range p.ComplianceCheck() {
		if !r.Passed {
			failed = append(failed, r)
		}
	}
	require.Len(t, failed, 1)
	assert.Equal(t, "CIS-1.1", failed[0].ControlID)
	assert.Contains(t, failed[0].Detail, "Environment")
}

func TestGCSBucket_ComplianceCheck_FailsUniformAccessDisabled(t *testing.T) {
	cfg := validGCSConfig()
	cfg.UniformBucketAccess = false

	p := NewGCSBucketProvisioner(cfg)
	var failed []models.ComplianceResult
	for _, r := range p.ComplianceCheck() {
		if !r.Passed {
			failed = append(failed, r)
		}
	}
	require.Len(t, failed, 1)
	assert.Equal(t, "CIS-GCP-5.1", failed[0].ControlID)
}

func TestGCSBucket_ComplianceCheck_FailsMissingProjectID(t *testing.T) {
	cfg := validGCSConfig()
	cfg.ProjectID = ""

	p := NewGCSBucketProvisioner(cfg)
	var failed []models.ComplianceResult
	for _, r := range p.ComplianceCheck() {
		if !r.Passed {
			failed = append(failed, r)
		}
	}
	require.Len(t, failed, 1)
	assert.Equal(t, "CIS-GCP-5.2", failed[0].ControlID)
}

func TestGCSBucket_ComplianceCheck_FailsPublicAccess(t *testing.T) {
	cfg := validGCSConfig()
	cfg.AllowPublic = true

	p := NewGCSBucketProvisioner(cfg)
	var failed []models.ComplianceResult
	for _, r := range p.ComplianceCheck() {
		if !r.Passed {
			failed = append(failed, r)
		}
	}
	require.Len(t, failed, 1)
	assert.Equal(t, "CIS-2.1", failed[0].ControlID)
}

func TestGCSBucket_DefaultStorageClass_IsStandard(t *testing.T) {
	cfg := validGCSConfig()
	cfg.StorageClass = ""

	p := NewGCSBucketProvisioner(cfg)
	assert.Equal(t, "STANDARD", p.cfg.StorageClass)
}

func TestGCSBucket_Provision_AbortsOnComplianceFailure(t *testing.T) {
	cfg := validGCSConfig()
	cfg.UniformBucketAccess = false

	p := NewGCSBucketProvisioner(cfg)
	err := p.Provision(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CIS-GCP-5.1")
}

func TestGCSBucket_Provision_SucceedsForValidConfig(t *testing.T) {
	p := NewGCSBucketProvisioner(validGCSConfig())
	err := p.Provision(context.Background())
	assert.NoError(t, err)
}
