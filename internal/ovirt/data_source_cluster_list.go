package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) clusterListDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.clusterListDataSourceRead,
		Schema: map[string]*schema.Schema{
			"clusters": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ID of the cluster.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the cluster.",
							Computed:    true,
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "Comment of the cluster.",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the cluster.",
							Computed:    true,
						},
					},
				},
			},
		},
		Description: `This data source retrieves a list of clusters`,
	}
}

func (p *provider) clusterListDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	allClusters, err := client.ListClusters()
	if err != nil {
		return errorToDiags("list all clusters", err)
	}
	clusterList := make([]map[string]interface{}, 0)
	for _, cluster := range allClusters {
		clusterMap := make(map[string]interface{}, 0)
		clusterMap["id"] = string(cluster.ID())
		clusterMap["name"] = cluster.Name()
		clusterMap["comment"] = cluster.Comment()
		clusterMap["description"] = cluster.Description()
		clusterList = append(clusterList, clusterMap)
	}
	if err := data.Set("clusters", clusterList); err != nil {
		return diag.FromErr(err)
	}
	data.SetId("clusters")
	return nil
}
