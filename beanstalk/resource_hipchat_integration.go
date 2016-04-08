package beanstalk

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHipchatIntegration() *schema.Resource {
	integrationType := &integrationType{
		Name: "HipchatIntegration",
		Attributes: map[string]*schema.Schema{
			"service_access_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HIPCHAT_API_TOKEN", nil),
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
	}

	return integrationType.resource()
}
