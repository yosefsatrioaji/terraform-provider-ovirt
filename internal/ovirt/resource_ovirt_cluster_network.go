package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/yosefsatrioaji/go-ovirt-client/v3"
)

var clusterNetworkSchema = map[string]*schema.Schema{
	"cluster_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "ID of the cluster to which the network is attached.",
		ValidateDiagFunc: validateUUID,
		ForceNew:         true,
	},
	"network_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "ID of the network to attach to the cluster.",
		ValidateDiagFunc: validateUUID,
		ForceNew:         true,
	},
	"required": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Indicates whether the network is required for the cluster.",
		Default:     false,
		ForceNew:    true,
	},
}

func (p *provider) clusterNetworkResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.clusterNetworkCreate,
		ReadContext:   p.clusterNetworkRead,
		DeleteContext: p.clusterNetworkDelete,
		Schema:        clusterNetworkSchema,
		Description:   "The ovirt_cluster_network resource attaches a network to a cluster in oVirt.",
	}
}

func (p *provider) clusterNetworkCreate(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	clusterID := data.Get("cluster_id").(string)
	networkID := data.Get("network_id").(string)
	required := data.Get("required").(bool)
	clusterNetwork, err := client.CreateClusterNetwork(
		ovirtclient.ClusterID(clusterID),
		ovirtclient.NetworkID(networkID),
		required,
	)
	if err != nil {
		return errorToDiags("create cluster network", err)
	}
	data.SetId(string(clusterNetwork.ClusterID()) + "_" + string(clusterNetwork.NetworkID()))
	return p.clusterNetworkRead(ctx, data, nil)
}

func (p *provider) clusterNetworkRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	clusterID := data.Get("cluster_id").(string)
	networkID := data.Get("network_id").(string)
	clusterNetwork, err := client.ClusterNetworkGet(
		ovirtclient.ClusterID(clusterID),
		ovirtclient.NetworkID(networkID),
	)
	if err != nil {
		return errorToDiags("get cluster network", err)
	}
	return clusterNetworkResourceUpdate(clusterNetwork, data)
}

func (p *provider) clusterNetworkDelete(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	clusterID := data.Get("cluster_id").(string)
	networkID := data.Get("network_id").(string)
	err := client.RemoveClusterNetwork(
		ovirtclient.ClusterID(clusterID),
		ovirtclient.NetworkID(networkID),
	)
	if err != nil {
		return errorToDiags("remove cluster network", err)
	}
	data.SetId("")
	return nil
}

func clusterNetworkResourceUpdate(clusterNetwork ovirtclient.ClusterNetwork, data *schema.ResourceData) diag.Diagnostics {
	if err := data.Set("cluster_id", string(clusterNetwork.ClusterID())); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("network_id", string(clusterNetwork.NetworkID())); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("required", clusterNetwork.Required()); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(string(clusterNetwork.ClusterID()) + "_" + string(clusterNetwork.NetworkID()))
	return nil
}
