package beanstalk

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: CreateUser,
		Read:   ReadUser,
		Update: UpdateUser,
		Delete: DeleteUser,

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"account_admin": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"timezone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "London",
			},

			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func CreateUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	email := d.Get("email").(string)

	req := &InvitationCreateRequestWrap{
		InvitationCreateRequest{
			User{
				Name:  d.Get("name").(string),
				Email: email,
			},
		},
	}
	res := &InvitationWrap{}

	err := client.Post([]string{"invitations"}, req, res)
	if err != nil {
		return err
	}

	// By creating an invitation we also created a user, but the
	// Beanstalk API doesn't give us the id of the user in the
	// response so we have to go hunt for it in the user list,
	// using the email address (which is guaranteed unique).

	id := 0
	var ures []UserWrap
	pathParts := []string{"users"}
	queryArgs := map[string]string{
		"per_page": "50",
		"page":     "0",
	}
	pageIdx := 1
Pages:
	for id == 0 {
		queryArgs["page"] = strconv.Itoa(pageIdx)
		err := client.Get(pathParts, queryArgs, &ures)
		if err != nil {
			return err
		}
		if len(ures) == 0 {
			return fmt.Errorf("invited user %v is not in user list", email)
		}
		for _, userWrap := range ures {
			if userWrap.User.Email == email {
				id = userWrap.User.ID
				break Pages
			}
		}

		pageIdx++
	}

	d.SetId(strconv.Itoa(id))
	d.Set("id", id)

	return UpdateUser(d, meta)
}

func ReadUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	var res UserWrap
	err := client.Get([]string{"users", d.Id()}, nil, &res)
	if err != nil {
		if _, ok := err.(*NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("username", res.User.Username)
	d.Set("email", res.User.Email)
	d.Set("name", res.User.Name)
	d.Set("account_admin", res.User.IsAccountAdmin)
	d.Set("timezone", res.User.Timezone)
	d.Set("first_name", res.User.FirstName)
	d.Set("last_name", res.User.LastName)

	return nil
}

func UpdateUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	userWrap := &UserWrap{
		User: User{
			Username: d.Get("username").(string),
			Email: d.Get("email").(string),
			Name: d.Get("name").(string),
			IsAccountAdmin: d.Get("account_admin").(bool),
			Timezone: d.Get("timezone").(string),
		},
	}

	err := client.Put([]string{"users", d.Id()}, userWrap, nil)
	if err != nil {
		return err
	}

	return ReadUser(d, meta)
}

func DeleteUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	err := client.Delete([]string{"users", d.Id()})
	if err == nil {
		d.SetId("")
	}
	return err
}

type Invitation struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type InvitationWrap struct {
	Invitation Invitation `json:"invitation"`
}

type InvitationCreateRequest struct {
	User User `json:"user"`
}

type InvitationCreateRequestWrap struct {
	Invitation InvitationCreateRequest `json:"invitation"`
}

type User struct {
	ID             int    `json:"id,omitempty"`
	Username       string `json:"login"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	Timezone       string `json:"timezone"`
	IsAccountAdmin bool   `json:"admin"`
	IsAccountOwner bool   `json:"owner"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
}

type UserWrap struct {
	User User `json:"user"`
}
