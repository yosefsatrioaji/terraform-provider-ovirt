package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/yosefsatrioaji/go-ovirt-client/v3"
)

var networkSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Name of the network.",
		ValidateDiagFunc: validateNonEmpty,
	},
	"data_center_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "ID of the data center in which the network is created.",
		ValidateDiagFunc: validateUUID,
		ForceNew:         true,
	},
	"vlan_id": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "VLAN ID of the network.",
	},
	"description": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Description of the network.",
	},
	"comment": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Comment of the network.",
	},
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func (p *provider) networkResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.networkCreate,
		ReadContext:   p.networkRead,
		DeleteContext: p.networkDelete,
		UpdateContext: p.networkUpdate,
		Schema:        networkSchema,
		Description:   "The ovirt_network resource creates networks in oVirt.",
	}
}

func (p *provider) networkCreate(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	name := data.Get("name").(string)
	dataCenterID := data.Get("data_center_id").(string)
	vlanID := data.Get("vlan_id").(int)
	description := data.Get("description").(string)
	comment := data.Get("comment").(string)
	network, err := client.CreateNetwork(
		ovirtclient.DatacenterID(dataCenterID),
		name,
		description,
		comment,
		vlanID,
	)
	if err != nil {
		return errorToDiags("create network", err)
	}
	return networkResourceUpdate(network, data)
}

func (p *provider) networkRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	networkID := data.Id()
	network, err := client.GetNetwork(ovirtclient.NetworkID(networkID))
	if err != nil {
		return errorToDiags("get network", err)
	}
	return networkResourceUpdate(network, data)
}

func (p *provider) networkUpdate(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	networkID := data.Id()
	dataCenterID := data.Get("data_center_id").(string)
	name := data.Get("name").(string)
	vlanID := data.Get("vlan_id").(int)
	description := data.Get("description").(string)
	comment := data.Get("comment").(string)
	updatedNetwork, err := client.UpdateNetwork(
		ovirtclient.NetworkID(networkID),
		ovirtclient.DatacenterID(dataCenterID),
		name,
		description,
		comment,
		vlanID,
	)
	if err != nil {
		return errorToDiags("update network", err)
	}
	return networkResourceUpdate(updatedNetwork, data)
}

func (p *provider) networkDelete(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	networkID := data.Id()
	if err := client.RemoveNetwork(ovirtclient.NetworkID(networkID)); err != nil {
		if !isNotFound(err) {
			return errorToDiags("delete network", err)
		}
	}
	data.SetId("")
	return nil
}

func networkResourceUpdate(network ovirtclient.Network, data *schema.ResourceData) diag.Diagnostics {
	if err := data.Set("name", network.Name()); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("data_center_id", string(network.DatacenterID())); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("vlan_id", network.VlanID()); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("description", network.Description()); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("comment", network.Comment()); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("id", string(network.ID())); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(string(network.ID()))
	return nil
}
