package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/yosefsatrioaji/go-ovirt-client/v3"
)

var resourceNetworkAttachmentSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"host_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "ID of the host to which the network is attached.",
		ValidateDiagFunc: validateUUID,
		ForceNew:         true,
	},
	"network_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "ID of the network to attach to the host.",
		ValidateDiagFunc: validateUUID,
		ForceNew:         true,
	},
	"nic_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the NIC associated with this network attachment.",
		ForceNew:    true,
	},
}

func (p *provider) resourceNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.resourceNetworkAttachmentCreate,
		ReadContext:   p.resourceNetworkAttachmentRead,
		DeleteContext: p.resourceNetworkAttachmentDelete,
		Schema:        resourceNetworkAttachmentSchema,
		Description:   "The ovirt_resource_network_attachment resource creates network attachments in oVirt.",
	}
}

func (p *provider) resourceNetworkAttachmentCreate(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	hostID := data.Get("host_id").(string)
	networkID := data.Get("network_id").(string)
	nicName := data.Get("nic_name").(string)
	networkAttachment, err := client.AttachNetworkToHost(
		ovirtclient.HostID(hostID),
		ovirtclient.NetworkID(networkID),
		nicName,
	)
	if err != nil {
		return errorToDiags("create network attachment", err)
	}
	return resourceNetworkAttachmentUpdate(networkAttachment, data)
}

func (p *provider) resourceNetworkAttachmentRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	networkAttachmentID := data.Id()
	hostID := data.Get("host_id").(string)
	nicName := data.Get("nic_name").(string)
	networkAttachment, err := client.GetNetworkAttachment(ovirtclient.NetworkAttachmentID(networkAttachmentID), ovirtclient.HostID(hostID), nicName)
	if err != nil {
		return errorToDiags("get network attachment", err)
	}
	return resourceNetworkAttachmentUpdate(networkAttachment, data)
}

func (p *provider) resourceNetworkAttachmentDelete(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	networkAttachmentID := data.Id()
	hostID := data.Get("host_id").(string)
	nicName := data.Get("nic_name").(string)
	err := client.DetachNetworkFromHost(ovirtclient.NetworkAttachmentID(networkAttachmentID), ovirtclient.HostID(hostID), nicName)
	if err != nil {
		return errorToDiags("delete network attachment", err)
	}
	data.SetId("")
	return nil
}

func resourceNetworkAttachmentUpdate(networkAttachment ovirtclient.NetworkAttachment, data *schema.ResourceData) diag.Diagnostics {
	if err := data.Set("host_id", string(networkAttachment.HostID())); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("network_id", string(networkAttachment.NetworkID())); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("nic_name", networkAttachment.NicName()); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("id", string(networkAttachment.ID())); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(string(networkAttachment.ID()))
	return nil
}
