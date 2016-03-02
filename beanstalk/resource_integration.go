package beanstalk

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

type integrationType struct {
	Name       string
	Attributes map[string]*schema.Schema
}

// The Beanstalk API has "integration" as a concept, but it is an abstract
// type that has a separate subtype for each integration type. Thus the
// implementation here is abstract and is instantiated for each of the
// physical resources.

func (it *integrationType) resource() *schema.Resource {
	resourceSchema := map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"repository_name": &schema.Schema{
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

func (it *integrationType) Read(d *schema.ResourceData, client *Client) error {
	repositoryName := d.Get("repository_name").(string)
	integrationId := d.Id()

	data := map[string]interface{}{}

	err := client.Get([]string{"repositories", repositoryName, "integrations", integrationId}, nil, &data)
	if err != nil {
		return err
	}

	it.refreshFromJSON(d, data["integration"].(map[string]interface{}))

	return nil
}

func (it *integrationType) Create(d *schema.ResourceData, client *Client) error {
	req := it.prepareForJSON(d)

	type responseIntegration struct {
		Id int `json:"id"`
	}
	type response struct {
		Integration responseIntegration `json:"integration"`
	}

	res := &response{}

	repositoryName := d.Get("repository_name").(string)

	err := client.Post([]string{"repositories", repositoryName, "integrations"}, req, res)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(res.Integration.Id))

	return nil
}

func (it *integrationType) Update(d *schema.ResourceData, client *Client) error {
	req := it.prepareForJSON(d)

	repositoryName := d.Get("repository_name").(string)
	integrationId := d.Id()

	return client.Put([]string{"repositories", repositoryName, "integrations", integrationId}, req, nil)
}

func (it *integrationType) Delete(d *schema.ResourceData, client *Client) error {
	repositoryName := d.Get("repository_name").(string)
	integrationId := d.Id()

	return client.Delete([]string{"repositories", repositoryName, "integrations", integrationId})
}

func (it *integrationType) prepareForJSON(d *schema.ResourceData) map[string]interface{} {
	ret := map[string]interface{}{}

	ret["type"] = it.Name

	for k, s := range it.Attributes {
		ret[k] = prepareForJSON(s, d.Get(k))
	}

	return map[string]interface{}{
		"integration": ret,
	}
}

func (it *integrationType) refreshFromJSON(d *schema.ResourceData, data map[string]interface{}) {
	for k, s := range it.Attributes {
		d.Set(k, decodeFromJSON(s, data[k]))
	}
}

func prepareForJSON(s *schema.Schema, value interface{}) interface{} {
	// This supports only what's required for the structures used by
	// Beanstalk's integrations. Almost all fields are primitive types,
	// but the modular webhook integration uses a nested object.
	switch s.Type {
	case schema.TypeList:
		if s.MaxItems == 1 {
			elem := s.Elem
			if resource, ok := elem.(*schema.Resource); ok {
				ret := map[string]interface{}{}
				values := value.([]interface{})
				valueMap := values[0].(map[string]interface{})
				for k, s := range resource.Schema {
					ret[k] = prepareForJSON(s, valueMap[k])
				}
				return ret
			} else {
				return prepareForJSON(elem.(*schema.Schema), value)
			}
		}
		break
	default:
		// No transformation for other types
		return value
	}

	// Unreachable
	return nil
}

func decodeFromJSON(s *schema.Schema, value interface{}) interface{} {
	// This supports only what's required for the structures used by
	// Beanstalk's integrations. Almost all fields are primitive types,
	// but the modular webhook integration uses a nested object.
	switch s.Type {
	case schema.TypeList:
		if s.MaxItems == 1 {
			elem := s.Elem
			if resource, ok := elem.(*schema.Resource); ok {
				ret := map[string]interface{}{}
				valueMap := value.(map[string]interface{})
				for k, s := range resource.Schema {
					ret[k] = decodeFromJSON(s, valueMap[k])
				}
				return []interface{}{ret}
			} else {
				return value
			}
		}
		break
	default:
		// No transformation for other types
		return value
	}

	// Unreachable
	return nil
}
