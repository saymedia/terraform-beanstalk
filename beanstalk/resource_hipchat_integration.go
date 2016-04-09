package beanstalk

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHipchatIntegration() *schema.Resource {
	integrationType := &integrationType{
		Name: "HipchatIntegration",
		Attributes: map[string]*schema.Schema{
			"service_access_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						hash := sha1.Sum([]byte(v.(string)))
						return hex.EncodeToString(hash[:])
					default:
						return ""
					}
				},
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
