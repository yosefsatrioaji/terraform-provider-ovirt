package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/marifwicaksana/go-ovirt-client/v3"
)

func (p *provider) storageDomainDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.storageDomainDataSourceRead,
		Schema: map[string]*schema.Schema{
			"storage_domain_id": {
				Type:             schema.TypeString,
				Description:      "ID of the oVirt VM.",
				Required:         true,
				ValidateDiagFunc: validateUUID,
			},
			"storage_domain_name": {
				Type:        schema.TypeString,
				Description: "Name of the oVirt Storage Domain.",
				Computed:    true,
			},
			"available": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Available space in the Storage Domain.",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the Storage Domain.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Storage Domain.",
			},
			"external_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External status of the Storage Domain.",
			},
		},
		Description: `This data source get storage domain details.`,
	}
}

func (p *provider) storageDomainDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	storageDomainId := data.Get("storage_domain_id").(string)

	storageDomain, err := client.GetStorageDomain(ovirtclient.StorageDomainID(storageDomainId))
	if err != nil {
		return errorToDiags("getting storage domain", err)
	}
	data.Set("storage_domain_name", storageDomain.Name())
	data.Set("available", storageDomain.Available())
	data.Set("storage_type", storageDomain.StorageType())
	data.Set("status", storageDomain.Status())
	data.Set("external_status", storageDomain.ExternalStatus())
	data.SetId(storageDomainId)
	return nil
}
