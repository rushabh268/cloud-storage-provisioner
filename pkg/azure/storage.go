// Package azure provides provisioners for Azure resources.
package azure

import (
	"context"
	"fmt"

	"github.com/rushabh268/cloud-storage-provisioner/pkg/compliance"
	"github.com/rushabh268/cloud-storage-provisioner/pkg/models"
)

// StorageAccountConfig holds configuration for an Azure Storage Account.
type StorageAccountConfig struct {
	models.ResourceConfig

	// ResourceGroup is required for all Azure resources.
	ResourceGroup string

	// SKU sets the storage redundancy tier: LRS, ZRS, GRS, RAGRS.
	SKU string

	// EnableHTTPSOnly forces all traffic over HTTPS (required by CIS Azure 3.1).
	EnableHTTPSOnly bool
}

// StorageAccountProvisioner provisions an Azure Storage Account with encryption,
// HTTPS enforcement, and logging.
type StorageAccountProvisioner struct {
	cfg StorageAccountConfig
}

func NewStorageAccountProvisioner(cfg StorageAccountConfig) *StorageAccountProvisioner {
	cfg.Tags = cfg.Tags.WithDefaults(
		cfg.Tags["Project"],
		cfg.Tags["Environment"],
		cfg.Tags["Owner"],
	)
	if cfg.SKU == "" {
		cfg.SKU = "GRS" // default to geo-redundant storage
	}
	return &StorageAccountProvisioner{cfg: cfg}
}

func (p *StorageAccountProvisioner) ComplianceCheck() []models.ComplianceResult {
	controls := append(compliance.MandatoryControls,
		compliance.Control{
			ID:          "CIS-AZ-3.1",
			Framework:   "CIS",
			Description: "Azure Storage Accounts must enforce HTTPS-only traffic.",
			Remediation: "Pass --https-only flag to the azure storage provision command.",
			Check: func(cfg models.ResourceConfig) (bool, string) {
				if !p.cfg.EnableHTTPSOnly {
					return false, "HTTPS-only is not enabled; plain HTTP traffic would be permitted"
				}
				return true, ""
			},
		},
		compliance.Control{
			ID:          "CIS-AZ-3.2",
			Framework:   "CIS",
			Description: "Azure Storage Accounts must use a resource group.",
			Remediation: "Pass --resource-group <name> to the azure storage provision command.",
			Check: func(cfg models.ResourceConfig) (bool, string) {
				if p.cfg.ResourceGroup == "" {
					return false, "resource group is not set"
				}
				return true, ""
			},
		},
	)
	results, _ := compliance.Run(controls, p.cfg.ResourceConfig)
	return results
}

func (p *StorageAccountProvisioner) Provision(ctx context.Context) error {
	results := p.ComplianceCheck()
	for _, r := range results {
		if !r.Passed {
			return fmt.Errorf("compliance check failed [%s]: %s — remediation: %s",
				r.ControlID, r.Detail, r.Remediation)
		}
	}

	fmt.Printf("[azure/storage] creating storage account %q in resource group %q\n",
		p.cfg.Name, p.cfg.ResourceGroup)
	fmt.Printf("[azure/storage] enabling HTTPS-only on %q\n", p.cfg.Name)
	fmt.Printf("[azure/storage] enabling Microsoft-managed encryption on %q\n", p.cfg.Name)
	fmt.Printf("[azure/storage] enabling Azure Monitor logging on %q\n", p.cfg.Name)
	fmt.Printf("[azure/storage] blocking public blob access on %q\n", p.cfg.Name)
	fmt.Printf("[azure/storage] applying tags to %q\n", p.cfg.Name)
	return nil
}

func (p *StorageAccountProvisioner) Destroy(ctx context.Context) error {
	fmt.Printf("[azure/storage] destroying storage account %q\n", p.cfg.Name)
	return nil
}
