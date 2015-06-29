package beanstalk

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Create: CreateRepository,
		Read:   ReadRepository,
		Update: UpdateRepository,
		Delete: DeleteRepository,

		Schema: map[string]*schema.Schema{
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"color_label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "white",
			},

			"default_git_branch": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "master",
			},

			"vcs": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "git",
				ForceNew: true,
			},

			"create_svn_structure": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func CreateRepository(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	req := &Repository{
		Title:              d.Get("title").(string),
		Name:               d.Get("name").(string),
		TypeID:             d.Get("vcs").(string),
		CreateSVNStructure: d.Get("create_svn_structure").(bool),
	}

	res := &RepositoryWrap{}

	err := client.Post([]string{"repositories"}, req, res)
	if err != nil {
		return err
	}

	id := strconv.Itoa(res.Repository.ID)
	d.SetId(id)
	d.Set("id", id)

	return UpdateRepository(d, meta)
}

func ReadRepository(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	res := &RepositoryWrap{}

	err := client.Get([]string{"repositories", d.Id()}, res)
	if err != nil {
		if _, ok := err.(*NotFoundError); ok {
			d.SetId("")
			d.Set("id", "")
			return nil
		} else {
			return err
		}
	}

	d.Set("title", res.Repository.Title)
	d.Set("name", res.Repository.Name)
	d.Set("color_label", res.Repository.ColorLabel)
	d.Set("default_git_branch", res.Repository.DefaultGitBranch)
	d.Set("vcs", res.Repository.VCS)
	d.Set("id", res.Repository.ID)
	d.Set("url", res.Repository.URL)

	return nil
}

func RenameRepository(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	req := &RepositoryWrap{
		Repository: Repository{
			Name: d.Get("name").(string),
		},
	}

	err := client.Put([]string{"repositories", d.Id(), "rename"}, req, nil)
	if err != nil {
		return err
	}

	return nil
}

func UpdateRepository(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	d.Partial(true)

	if d.HasChange("name") {
		// Renaming has its own special operation.
		err := RenameRepository(d, meta)
		if err != nil {
			return err
		}

		d.SetPartial("name")
	}

	req := &RepositoryWrap{
		Repository: Repository{
			Title:            d.Get("title").(string),
			ColorLabel:       d.Get("color_label").(string),
			DefaultGitBranch: d.Get("default_git_branch").(string),
		},
	}

	err := client.Put([]string{"repositories", d.Id()}, req, nil)
	if err != nil {
		return err
	}

	d.Partial(false)

	return ReadRepository(d, meta)
}

func DeleteRepository(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Beanstalk does not allow repositories to be deleted via its API. Delete this repository via the UI and run 'terraform refresh' to make Terraform notice it's gone.")
}

type Repository struct {
	ID                 int    `json:"id,omitempty"`
	Title              string `json:"title"`
	Name               string `json:"name"`
	ColorLabel         string `json:"color_label,omitempty"`
	DefaultGitBranch   string `json:"default_branch,omitempty"`
	TypeID             string `json:"type_id,omitempty"`
	VCS                string `json:"vcs,omitempty"`
	CreateSVNStructure bool   `json:"create_structure"`
	URL                string `json:"repository_url,omitempty"`
}

type RepositoryWrap struct {
	Repository Repository `json:"repository"`
}
