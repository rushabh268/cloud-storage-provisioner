package azure

import (
	"context"
	"testing"

	"github.com/rushabh268/cloud-storage-provisioner/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validAzureConfig() StorageAccountConfig {
	return StorageAccountConfig{
		ResourceConfig: models.ResourceConfig{
			Name:     "myappstorage",
			Region:   "eastus",
			Provider: models.ProviderAzure,
			Tags: models.Tags{
				"Project":     "my-app",
				"Environment": "prod",
				"Owner":       "rushabh268@gmail.com",
			},
		},
		ResourceGroup:   "my-rg",
		SKU:             "GRS",
		EnableHTTPSOnly: true,
	}
}

func TestStorageAccount_ComplianceCheck_PassesForValidConfig(t *testing.T) {
	p := NewStorageAccountProvisioner(validAzureConfig())
	results := p.ComplianceCheck()
	for _, r := range results {
		assert.True(t, r.Passed, "expected control %s to pass, got: %s", r.ControlID, r.Detail)
	}
}

func TestStorageAccount_ComplianceCheck_FailsMissingTags(t *testing.T) {
	cfg := validAzureConfig()
	delete(cfg.Tags, "Project")

	p := NewStorageAccountProvisioner(cfg)
	var failed []models.ComplianceResult
	for _, r := range p.ComplianceCheck() {
		if !r.Passed {
			failed = append(failed, r)
		}
	}
	require.Len(t, failed, 1)
	assert.Equal(t, "CIS-1.1", failed[0].ControlID)
}

func TestStorageAccount_ComplianceCheck_FailsHTTPSNotEnabled(t *testing.T) {
	cfg := validAzureConfig()
	cfg.EnableHTTPSOnly = false

	p := NewStorageAccountProvisioner(cfg)
	var failed []models.ComplianceResult
	for _, r := range p.ComplianceCheck() {
		if !r.Passed {
			failed = append(failed, r)
		}
	}
	require.Len(t, failed, 1)
	assert.Equal(t, "CIS-AZ-3.1", failed[0].ControlID)
}

func TestStorageAccount_ComplianceCheck_FailsMissingResourceGroup(t *testing.T) {
	cfg := validAzureConfig()
	cfg.ResourceGroup = ""

	p := NewStorageAccountProvisioner(cfg)
	var failed []models.ComplianceResult
	for _, r := range p.ComplianceCheck() {
		if !r.Passed {
			failed = append(failed, r)
		}
	}
	require.Len(t, failed, 1)
	assert.Equal(t, "CIS-AZ-3.2", failed[0].ControlID)
}

func TestStorageAccount_ComplianceCheck_FailsPublicAccess(t *testing.T) {
	cfg := validAzureConfig()
	cfg.AllowPublic = true

	p := NewStorageAccountProvisioner(cfg)
	var failed []models.ComplianceResult
	for _, r := range p.ComplianceCheck() {
		if !r.Passed {
			failed = append(failed, r)
		}
	}
	require.Len(t, failed, 1)
	assert.Equal(t, "CIS-2.1", failed[0].ControlID)
}

func TestStorageAccount_DefaultSKU_IsGRS(t *testing.T) {
	cfg := validAzureConfig()
	cfg.SKU = ""

	p := NewStorageAccountProvisioner(cfg)
	assert.Equal(t, "GRS", p.cfg.SKU)
}

func TestStorageAccount_Provision_AbortsOnComplianceFailure(t *testing.T) {
	cfg := validAzureConfig()
	cfg.EnableHTTPSOnly = false

	p := NewStorageAccountProvisioner(cfg)
	err := p.Provision(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CIS-AZ-3.1")
}

func TestStorageAccount_Provision_SucceedsForValidConfig(t *testing.T) {
	p := NewStorageAccountProvisioner(validAzureConfig())
	err := p.Provision(context.Background())
	assert.NoError(t, err)
}
