package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/yosefsatrioaji/go-ovirt-client/v3"
)

var vnicProfileSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Name of the VNIC profile.",
		ValidateDiagFunc: validateNonEmpty,
	},
	"network_id": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "ID of the network to which the VNIC profile is attached.",
		ValidateDiagFunc: validateUUID,
	},
	"pass_through": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Indicates whether the VNIC profile is pass-through.",
	},
	"port_mirroring": {
		Type:        schema.TypeBool,
		Required:    true,
		ForceNew:    true,
		Description: "Indicates whether port mirroring is enabled for the VNIC profile.",
	},
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func (p *provider) vnicProfileResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.vnicProfileCreate,
		ReadContext:   p.vnicProfileRead,
		DeleteContext: p.vnicProfileDelete,
		// Importer: &schema.ResourceImporter{
		// 	StateContext: p.vnicProfileImport,
		// },
		Schema:      vnicProfileSchema,
		Description: "The ovirt_vnic_profile resource creates VNIC profiles in oVirt.",
	}
}

func (p *provider) vnicProfileCreate(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	name := data.Get("name").(string)
	networkID := data.Get("network_id").(string)
	passThrough := data.Get("pass_through").(string)
	portMirroring := data.Get("port_mirroring").(bool)
	params := ovirtclient.CreateVNICProfileParams()
	if passThrough != "" {
		params = params.WithPassThrough(passThrough)
	}
	if portMirroring {
		params = params.WithPortMirroring(portMirroring)
	}
	vnicProfile, err := client.CreateVNICProfile(name, ovirtclient.NetworkID(networkID), params)
	if err != nil {
		return errorToDiags("create VNIC profile", err)
	}
	return vnicProfileResourceUpdate(vnicProfile, data)
}

func (p *provider) vnicProfileRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	vnicProfileID := data.Id()
	vnicProfile, err := client.GetVNICProfile(ovirtclient.VNICProfileID(vnicProfileID))
	if err != nil {
		return errorToDiags("read VNIC profile", err)
	}
	return vnicProfileResourceUpdate(vnicProfile, data)
}

func (p *provider) vnicProfileDelete(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	vnicProfileID := data.Id()
	if err := client.RemoveVNICProfile(ovirtclient.VNICProfileID(vnicProfileID)); err != nil {
		if !isNotFound(err) {
			return errorToDiags("delete VNIC profile", err)
		}
	}
	data.SetId("")
	return nil
}

func vnicProfileResourceUpdate(vnicProfile ovirtclient.VNICProfile, data *schema.ResourceData) diag.Diagnostics {
	if err := data.Set("name", vnicProfile.Name()); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("network_id", string(vnicProfile.NetworkID())); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("pass_through", vnicProfile.PassThrough()); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("port_mirroring", vnicProfile.PortMirroring()); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("id", string(vnicProfile.ID())); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(string(vnicProfile.ID()))
	return nil
}
