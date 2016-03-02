package beanstalk

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceModularWebhookIntegration() *schema.Resource {
	integrationType := &integrationType{
		Name: "ModularWebHooksIntegration",
		Attributes: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"triggers": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"commit": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"push": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"deploy": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"comment": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"create_branch": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"delete_branch": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"create_tag": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"delete_tag": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}

	return integrationType.resource()
}
