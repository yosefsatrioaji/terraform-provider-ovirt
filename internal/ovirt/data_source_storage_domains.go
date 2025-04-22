//nolint:dupl,revive
package ovirt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) storageDomainsDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.storageDomainsDataSourceRead,
		Schema: map[string]*schema.Schema{
			"storage_domains": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "A set of all Storage Domains.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "oVirt ID of the Storage Domain.",
						},
						"storage_domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Storage Domain.",
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
				},
			},
		},
	}
}

func (p *provider) storageDomainsDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	storageDomains, err := client.ListStorageDomains()
	if err != nil {
		return errorToDiags(fmt.Sprintf("list storage domain"), err)
	}

	domains := make([]map[string]interface{}, 0)

	for _, storageDomain := range storageDomains {
		domain := make(map[string]interface{}, 0)

		domain["storage_domain_id"] = storageDomain.ID()
		domain["storage_domain_name"] = storageDomain.Name()
		domain["available"] = storageDomain.Available()
		domain["storage_type"] = storageDomain.StorageType()
		domain["status"] = storageDomain.Status()
		domain["external_status"] = storageDomain.ExternalStatus()

		domains = append(domains, domain)
	}

	if err := data.Set("storage_domains", domains); err != nil {
		return errorToDiags("set storage domains", err)
	}
	return nil
}
