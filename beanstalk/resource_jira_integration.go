package beanstalk

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceJiraIntegration() *schema.Resource {
	integrationType := &integrationType{
		Name: "JiraIntegration",
		Attributes: map[string]*schema.Schema{
			"service_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_login": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: hashForState,
			},
			"service_password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: hashForState,
			},
			"service_project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		WriteOnlyAttributes: []string{"service_login", "service_password"},
	}

	return integrationType.resource()
}
