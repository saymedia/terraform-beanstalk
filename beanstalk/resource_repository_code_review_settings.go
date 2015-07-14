package beanstalk

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRepositoryCodeReviewSettings() *schema.Resource {
	return &schema.Resource{
		Create: CreateRepositoryCodeReviewSettings,
		Read:   ReadRepositoryCodeReviewSettings,
		Update: UpdateRepositoryCodeReviewSettings,
		Delete: DeleteRepositoryCodeReviewSettings,

		Schema: map[string]*schema.Schema{
			"repository_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"unanimous_approval": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"auto_reopen": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"default_assignee_user_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"default_watching_user_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"default_watching_team_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func CreateRepositoryCodeReviewSettings(d *schema.ResourceData, meta interface{}) error {
	id := strconv.Itoa(d.Get("repository_id").(int))
	d.SetId(id)
	return UpdateRepositoryCodeReviewSettings(d, meta)
}

func ReadRepositoryCodeReviewSettings(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	res := &RepositoryCodeReview{}

	err := client.Get([]string{d.Id(), "code_reviews", "settings"}, nil, res)
	if err != nil {
		if _, ok := err.(*NotFoundError); ok {
			d.SetId("")
			d.Set("id", "")
			return nil
		} else {
			return err
		}
	}

	d.Set("unanimous_approval", res.UnanimousApproval)
	d.Set("auto_reopen", res.AutoReopen)

	assigneeUserIds := make([]int, len(res.DefaultAssignees))
	for i, item := range res.DefaultAssignees {
		assigneeUserIds[i] = item.ID
	}
	d.Set("default_assignee_user_ids", assigneeUserIds)

	watcherUserIds := make([]int, 0, len(res.DefaultWatchers))
	watcherTeamIds := make([]int, 0, len(res.DefaultWatchers))
	for _, item := range res.DefaultWatchers {
		switch item.Type {
		case "User":
			watcherUserIds = append(watcherUserIds, item.ID)
		case "Team":
			watcherTeamIds = append(watcherTeamIds, item.ID)
		default:
			log.Printf("Ignored watcher of unknown type %v", item.Type)
		}
	}
	d.Set("default_watching_user_ids", watcherUserIds)
	d.Set("default_watching_team_ids", watcherTeamIds)

	return nil
}

func UpdateRepositoryCodeReviewSettings(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	sliceToInts := func(in []interface{}) []int {
		ret := make([]int, len(in))
		for i, si := range in {
			ret[i] = si.(int)
		}
		return ret
	}

	assigneeUserIds := sliceToInts(d.Get("default_assignee_user_ids").([]interface{}))
	watcherUserIds := sliceToInts(d.Get("default_watching_user_ids").([]interface{}))
	watcherTeamIds := sliceToInts(d.Get("default_watching_team_ids").([]interface{}))

	req := &RepositoryCodeReview{
		UnanimousApproval:      d.Get("unanimous_approval").(bool),
		AutoReopen:             d.Get("auto_reopen").(bool),
		DefaultAssigneeUserIDs: assigneeUserIds,
		DefaultWatcherUserIDs:  watcherUserIds,
		DefaultWatcherTeamIDs:  watcherTeamIds,
	}

	return client.Put([]string{d.Id(), "code_reviews", "settings"}, req, nil)
}

func DeleteRepositoryCodeReviewSettings(d *schema.ResourceData, meta interface{}) error {
	// Deleting this really just means no longer managing it with Terraform,
	// so we don't need to take any action here.
	return nil
}

type RepositoryCodeReview struct {
	UnanimousApproval      bool                          `json:"unanimous_approval"`
	AutoReopen             bool                          `json:"auto_reopen"`
	DefaultWatchers        []RepositoryCodeReviewWatcher `json:"default_watchers,omitempty"`
	DefaultAssignees       []RepositoryCodeReviewWatcher `json:"default_assignees,omitempty"`
	DefaultAssigneeUserIDs []int                         `json:"default_assignees"`
	DefaultWatcherUserIDs  []int                         `json:"default_watchers_user_ids"`
	DefaultWatcherTeamIDs  []int                         `json:"default_watchers_team_ids"`
}

type RepositoryCodeReviewWatcher struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Username string `json:"login"`
	Email    string `json:"email"`
}
