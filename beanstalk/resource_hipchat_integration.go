package beanstalk

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHipchatIntegration() *schema.Resource {
	integrationType := &integrationType{
		Name: "HipchatIntegration",
		Attributes: map[string]*schema.Schema{
			"service_access_token": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: hashForState,
			},
			"service_room_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"listen_commits": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"listen_deployments": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		WriteOnlyAttributes: []string{"service_access_token"},
	}

	return integrationType.resource()
}
