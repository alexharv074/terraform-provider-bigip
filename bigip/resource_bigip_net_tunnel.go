/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipNetTunnel() *schema.Resource {

	return &schema.Resource{
		Create: resourceBigipNetTunnelCreate,
		Read:   resourceBigipNetTunnelRead,
		Update: resourceBigipNetTunnelUpdate,
		Delete: resourceBigipNetTunnelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the TUNNEL",
			},
			"app_service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The application service that the object belongs to",
			},
			"auto_last_hop": {
				Type:     schema.TypeString,
				Optional: true,
				//Default:     "default",
				Description: "Specifies whether auto lasthop is enabled or not",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User defined description",
			},
			"local_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies a local IP address. This option is required",
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				//Default:     "bidirectional",
				Description: "Specifies how the tunnel carries traffic",
			},
			"partition": {
				Type:     schema.TypeString,
				Optional: true,
				//Default:     "Common",
				Description: "Displays the admin-partition within which this component resides",
			},
			"profile": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the profile that you want to associate with the tunnel",
			},
			"remote_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a remote IP address",
			},
			"secondary_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a secondary non-floating IP address when the local-address is set to a floating address",
			},
			"tos": {
				Type:     schema.TypeString,
				Optional: true,
				//Default:     "65535",
				Description: "Specifies a value for insertion into the Type of Service (ToS) octet within the IP header of the encapsulating header of transmitted packets",
			},
			"traffic_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a traffic-group for use with the tunnel",
			},
			"transparent": {
				Type:     schema.TypeString,
				Optional: true,
				//Default:     "disabled",
				Description: "Enables or disables the tunnel to be transparent",
			},
			"use_pmtu": {
				Type:     schema.TypeString,
				Optional: true,
				//Default:     "enabled",
				Description: "Enables or disables the tunnel to use the PMTU (Path MTU) information provided by ICMP NeedFrag error messages",
			},
			"idle_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				//Default:     300,
				Description: "Specifies an idle timeout for wildcard tunnels in seconds",
			},
			/*"if_index": {
			        Type:        schema.TypeInt,
			        Optional:    true,
			        //Default:     0,
			        Description: "Displays the index assigned to this tunnel",
			},*/
			"key": {
				Type:     schema.TypeInt,
				Optional: true,
				//Default:     0,
				Description: "The key field may represent different values depending on the type of the tunnel",
			},
			"mtu": {
				Type:     schema.TypeInt,
				Optional: true,
				//  Default:     0,
				Description: "Specifies the maximum transmission unit (MTU) of the tunnel",
			},
		},
	}

}

func resourceBigipNetTunnelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	local_address := d.Get("local_address").(string)
	profile := d.Get("profile").(string)
	config := &bigip.Tunnel{
		Name:         name,
		LocalAddress: local_address,
		Profile:      profile,
	}
	err := client.CreateTunnel(config)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Tunnel %s %v :", name, err)
		return err
	}

	d.SetId(name)

	err = resourceBigipNetTunnelUpdate(d, meta)
	if err != nil {
		client.DeleteTunnel(name)
		return err
	}

	return resourceBigipNetTunnelRead(d, meta)
}

func resourceBigipNetTunnelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Reading TUNNEL %s", name)
	tunnel, err := client.GetTunnel(name)
	if err != nil {
		return err
	}
	if tunnel == nil {
		d.SetId("")
		return fmt.Errorf("[ERROR] Tunnel (%s) not found, removing from state", d.Id())
	}
	log.Printf("[DEBUG] Tunnel:%+v", tunnel)
	if err := d.Set("app_service", tunnel.AppService); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AppService to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("auto_last_hop", tunnel.AutoLasthop); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AutoLasthop to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("description", tunnel.Description); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Description to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("idle_timeout", tunnel.IdleTimeout); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IdleTimeout to state for Tunnel (%s): %s", d.Id(), err)
	}
	/*if err := d.Set("if_index", tunnel.IfIndex); err != nil {
	        return fmt.Errorf("[DEBUG] Error saving IfIndex to state for Tunnel (%s): %s", d.Id(), err)
	}*/
	if err := d.Set("key", tunnel.Key); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Key to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("local_address", tunnel.LocalAddress); err != nil {
		return fmt.Errorf("[DEBUG] Error saving LocalAddress to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("mode", tunnel.Mode); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Mode to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("mtu", tunnel.Mtu); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Mtu to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("partition", tunnel.Partition); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Partition to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("profile", tunnel.Profile); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Profile to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("remote_address", tunnel.RemoteAddress); err != nil {
		return fmt.Errorf("[DEBUG] Error saving RemoteAddress to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("secondary_address", tunnel.SecondaryAddress); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SecondaryAddress to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("tos", tunnel.Tos); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Tos to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("traffic_group", tunnel.TrafficGroup); err != nil {
		return fmt.Errorf("[DEBUG] Error saving TrafficGroup to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("transparent", tunnel.Transparent); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Transparent to state for Tunnel (%s): %s", d.Id(), err)
	}
	if err := d.Set("use_pmtu", tunnel.UsePmtu); err != nil {
		return fmt.Errorf("[DEBUG] Error saving UsePmtu to state for Tunnel (%s): %s", d.Id(), err)
	}
	_ = d.Set("name", name)
	return nil
}

func resourceBigipNetTunnelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Updating Tunnel %s", name)

	r := &bigip.Tunnel{
		Name:        name,
		AppService:  d.Get("app_service").(string),
		AutoLasthop: d.Get("auto_last_hop").(string),
		Description: d.Get("description").(string),
		IdleTimeout: d.Get("idle_timeout").(int),
		//IfIndex:d.Get("if_index").(int),
		Key:              d.Get("key").(int),
		LocalAddress:     d.Get("local_address").(string),
		Mode:             d.Get("mode").(string),
		Mtu:              d.Get("mtu").(int),
		Partition:        d.Get("partition").(string),
		Profile:          d.Get("profile").(string),
		RemoteAddress:    d.Get("remote_address").(string),
		SecondaryAddress: d.Get("secondary_address").(string),
		Tos:              d.Get("tos").(string),
		TrafficGroup:     d.Get("traffic_group").(string),
		Transparent:      d.Get("transparent").(string),
		UsePmtu:          d.Get("use_pmtu").(string),
	}

	err := client.ModifyTunnel(name, r)
	if err != nil {
		return fmt.Errorf("Error modifying TUNNEL %s: %v", name, err)
	}

	return resourceBigipNetTunnelRead(d, meta)
}

func resourceBigipNetTunnelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Deleting TUNNEL %s", name)

	err := client.DeleteTunnel(name)
	if err != nil {
		return fmt.Errorf("Error Deleting Tunnel : %s", err)
	}

	d.SetId("")
	return nil
}
