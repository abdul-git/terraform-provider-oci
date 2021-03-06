// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/oracle/bmcs-go-sdk"

	"github.com/oracle/terraform-provider-oci/options"

	"github.com/oracle/terraform-provider-oci/crud"
)

func VirtualNetworkDatasource() *schema.Resource {
	return &schema.Resource{
		Read: readVirtualNetworks,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"compartment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"virtual_networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     VirtualNetworkResource(),
			},
		},
	}
}

func readVirtualNetworks(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(*OracleClients)
	sync := &VirtualNetworkDatasourceCrud{}
	sync.D = d
	sync.Client = client.client
	return crud.ReadResource(sync)
}

type VirtualNetworkDatasourceCrud struct {
	crud.BaseCrud
	Res *baremetal.ListVirtualNetworks
}

func (s *VirtualNetworkDatasourceCrud) Get() (e error) {
	compartmentID := s.D.Get("compartment_id").(string)

	opts := &baremetal.ListOptions{}
	options.SetListOptions(s.D, opts)

	s.Res = &baremetal.ListVirtualNetworks{
		VirtualNetworks: []baremetal.VirtualNetwork{},
	}

	for {
		var list *baremetal.ListVirtualNetworks
		if list, e = s.Client.ListVirtualNetworks(compartmentID, opts); e != nil {
			break
		}

		s.Res.VirtualNetworks = append(s.Res.VirtualNetworks, list.VirtualNetworks...)

		if hasNextPage := options.SetNextPageOption(list.NextPage, &opts.PageListOptions); !hasNextPage {
			break
		}
	}

	return
}

func (s *VirtualNetworkDatasourceCrud) SetData() {
	if s.Res == nil {
		return
	}
	// Important, if you don't have an ID, make one up for your datasource
	// or things will end in tears
	s.D.SetId(time.Now().UTC().String())
	resources := []map[string]interface{}{}
	for _, v := range s.Res.VirtualNetworks {
		res := map[string]interface{}{
			"cidr_block":               v.CidrBlock,
			"compartment_id":           v.CompartmentID,
			"dns_label":                v.DnsLabel,
			"default_dhcp_options_id":  v.DefaultDHCPOptionsID,
			"default_route_table_id":   v.DefaultRouteTableID,
			"default_security_list_id": v.DefaultSecurityListID,
			"display_name":             v.DisplayName,
			"id":                       v.ID,
			"state":                    v.State,
			"time_created":             v.TimeCreated.String(),
		}
		resources = append(resources, res)
	}

	if f, fOk := s.D.GetOk("filter"); fOk {
		resources = ApplyFilters(f.(*schema.Set), resources)
	}

	if err := s.D.Set("virtual_networks", resources); err != nil {
		panic(err)
	}

	return
}
