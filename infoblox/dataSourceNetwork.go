package infoblox

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	ibclient "github.com/infobloxopen/infoblox-go-client"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkRead,
		Schema: map[string]*schema.Schema{
			"network_view_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("network_view_name", "default"),
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("tenantID", nil),
				Description: "Unique identifier of your tenant in cloud.",
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	networkViewName := d.Get("network_view_name").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := meta.(*ibclient.Connector)

	cidr := d.Get("cidr").(string)
	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	ea := make(ibclient.EA)
	network, err := objMgr.GetNetwork(networkViewName, cidr, ea)
	if err != nil {
		return fmt.Errorf("reading infoblox_network: %w", err)
	}

	d.SetId(network.Ref)

	d.Set("name", network.Ea["Network Name"])
	d.Set("tenant_id", tenantID)

	return nil
}
