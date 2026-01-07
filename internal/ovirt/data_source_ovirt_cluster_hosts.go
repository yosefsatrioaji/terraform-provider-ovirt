package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) clusterHostsDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.clusterHostsDataSourceRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "oVirt cluster ID in the Data Center.",
				ValidateDiagFunc: validateUUID,
			},
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the host.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status of the host.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the host.",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Comment of the host.",
						},
						"nics": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of host nics",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of nic host",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of nic host",
									},
								},
							},
						},
					},
				},
			},
		},
		Description: `A set of all hosts of a Cluster.`,
	}
}

func (p *provider) clusterHostsDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	clusterID := data.Get("cluster_id").(string)
	allHosts, err := client.ListHosts()

	if err != nil {
		return errorToDiags("list all hosts", err)
	}

	hosts := make([]map[string]interface{}, 0)

	for _, host := range allHosts {
		if string(host.ClusterID()) == clusterID {
			hostMap := make(map[string]interface{}, 0)
			hostMap["id"] = host.ID()
			hostMap["status"] = host.Status()
			hostMap["name"] = host.Name()
			hostMap["comment"] = host.Comment()
			nics := make([]map[string]interface{}, 0)
			hostNics, err := host.HostNICs()
			if err != nil {
				return errorToDiags("list host nics", err)
			}
			for _, nic := range hostNics {
				nicMap := make(map[string]interface{}, 0)
				nicMap["id"] = nic.ID()
				nicMap["name"] = nic.Name()
				nics = append(nics, nicMap)
			}
			hostMap["nics"] = nics
			hosts = append(hosts, hostMap)
		}
	}

	if err := data.Set("hosts", hosts); err != nil {
		return errorToDiags("set hosts", err)
	}

	data.SetId(clusterID)

	return nil
}
