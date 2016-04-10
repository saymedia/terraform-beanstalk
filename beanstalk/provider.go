package beanstalk

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"beanstalk_hipchat_integration":             resourceHipchatIntegration(),
			"beanstalk_jira_integration":                resourceJiraIntegration(),
			"beanstalk_modular_webhook_integration":     resourceModularWebhookIntegration(),
			"beanstalk_repository":                      resourceRepository(),
			"beanstalk_repository_code_review_settings": resourceRepositoryCodeReviewSettings(),
			"beanstalk_user":                            resourceUser(),
			"beanstalk_team":                            resourceTeam(),
		},

		Schema: map[string]*schema.Schema{
			"account_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BEANSTALK_USERNAME", nil),
			},
			"access_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BEANSTALK_ACCESS_TOKEN", nil),
			},
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &ClientConfig{
		AccountName: d.Get("account_name").(string),
		Username:    d.Get("username").(string),
		AccessToken: d.Get("access_token").(string),
	}
	return NewClient(config)
}
