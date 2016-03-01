package beanstalk

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type integrationType struct {
	Name       string
	Attributes map[string]*schema.Schema
}

func (it *integrationType) resource() *schema.Resource {
	resourceSchema := map[string]*schema.Schema{
		"repository_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	}

	for k, v := range it.Attributes {
		resourceSchema[k] = v
	}

	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			client := meta.(*Client)
			return it.Read(d, client)
		},
		Create: func(d *schema.ResourceData, meta interface{}) error {
			client := meta.(*Client)
			return it.Create(d, client)
		},
		Update: func(d *schema.ResourceData, meta interface{}) error {
			client := meta.(*Client)
			return it.Update(d, client)
		},
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			client := meta.(*Client)
			return it.Delete(d, client)
		},

		Schema: resourceSchema,
	}
}

func (it *integrationType) Read(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func (it *integrationType) Create(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func (it *integrationType) Update(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func (it *integrationType) Delete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
