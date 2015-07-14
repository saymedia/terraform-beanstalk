package beanstalk

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Create: CreateTeam,
		Read:   ReadTeam,
		Update: UpdateTeam,
		Delete: DeleteTeam,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"color_label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "white",
			},

			"user_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Set: func(v interface{}) int {
					return v.(int)
				},
			},

			"repository_permissions": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repository_id": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  "white",
						},
						"repository_title": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"can_write": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"can_deploy": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"can_configure_deployments": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
				Set: hashRepositoryPermissions,
			},
		},
	}
}

func CreateTeam(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	req := teamFromResourceData(d)

	res := &TeamReadWrap{}

	err := client.Post([]string{"teams"}, req, res)
	if err != nil {
		return err
	}

	updateResourceDataFromTeam(&res.Team, d)

	return nil
}

func UpdateTeam(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	req := teamFromResourceData(d)

	res := &TeamReadWrap{}

	err := client.Put([]string{"teams", d.Id()}, req, res)
	if err != nil {
		return err
	}

	updateResourceDataFromTeam(&res.Team, d)

	return nil
}

func DeleteTeam(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	err := client.Delete([]string{"teams", d.Id()})
	if err == nil {
		d.SetId("")
	}
	return err
}

func ReadTeam(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	var res TeamReadWrap
	err := client.Get([]string{"teams", d.Id()}, nil, &res)
	if err != nil {
		if _, ok := err.(*NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	updateResourceDataFromTeam(&res.Team, d)

	return nil
}

func teamFromResourceData(d *schema.ResourceData) *TeamWrite {
	userIdsSet := d.Get("user_ids").(*schema.Set)
	userIdsI := userIdsSet.List()

	permissionsesSet := d.Get("repository_permissions").(*schema.Set)
	permissionsesI := permissionsesSet.List()

	team := &TeamWrite{
		Name:        d.Get("name").(string),
		ColorLabel:  d.Get("color_label").(string),
		UserIDs:     make([]int, len(userIdsI)),
		Permissions: map[string]TeamRepositoryPermissionsWrite{},
	}

	for i, userIdI := range userIdsI {
		team.UserIDs[i] = userIdI.(int)
	}

	for _, permissionsI := range permissionsesI {
		permissionsMap := permissionsI.(map[string]interface{})
		repositoryId := permissionsMap["repository_id"].(int)

		team.Permissions[strconv.Itoa(repositoryId)] = TeamRepositoryPermissionsWrite{
			CanWrite: permissionsMap["can_write"].(bool),
			CanDeploy: permissionsMap["can_deploy"].(bool),
			CanConfigureDeployments: permissionsMap["can_configure_deployments"].(bool),
		}
	}

	return team
}

func updateResourceDataFromTeam(team *TeamRead, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(team.ID))
	d.Set("id", team.ID)
	d.Set("name", team.Name)
	d.Set("color_label", team.ColorLabel)

	userIds := make([]int, len(team.Users))
	for i, user := range team.Users {
		userIds[i] = user.ID
	}
	d.Set("user_ids", userIds)

	permissionses := make([]map[string]interface{}, len(team.Permissions))
	for i, permissions := range team.Permissions {
		permissionses[i] = map[string]interface{}{
			"repository_id":             permissions.RepositoryID,
			"repository_title":          permissions.RepositoryTitle,
			"can_write":                 permissions.CanWrite,
			"can_deploy":                permissions.CanDeploy,
			"can_configure_deployments": permissions.CanConfigureDeployments,
		}
	}
	d.Set("repository_permissions", permissionses)
}

func hashRepositoryPermissions(v interface{}) int {
	m := v.(map[string]interface{})
	hashInput := fmt.Sprintf(
		"%v %v %v %v",
		m["repository_id"].(int),
		m["can_write"].(bool),
		m["can_deploy"].(bool),
		m["can_configure_deployments"].(bool),
	)
	return hashcode.String(hashInput)
}

type TeamWrite struct {
	ID          int                                       `json:"id,omitempty"`
	Name        string                                    `json:"name"`
	ColorLabel  string                                    `json:"color_label,omitempty"`
	UserIDs     []int                                     `json:"users"`
	Permissions map[string]TeamRepositoryPermissionsWrite `json:"permissions"`
}

type TeamRead struct {
	ID          int                             `json:"id,omitempty"`
	Name        string                          `json:"name"`
	ColorLabel  string                          `json:"color_label,omitempty"`
	Users       []User                          `json:"users"`
	Permissions []TeamRepositoryPermissionsRead `json:"permissions"`
}

type TeamReadWrap struct {
	Team TeamRead `json:"team"`
}

type TeamRepositoryPermissionsWrite struct {
	CanWrite                bool `json:"write"`
	CanDeploy               bool `json:"deploy"`
	CanConfigureDeployments bool `json:"configure_deployments"`
}

type TeamRepositoryPermissionsRead struct {
	RepositoryID            int    `json:"repository_id"`
	RepositoryTitle         string `json:"repository_title"`
	CanWrite                bool   `json:"write"`
	CanDeploy               bool   `json:"deploy"`
	CanConfigureDeployments bool   `json:"configure_deployments"`
}
