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
				ForceNew: true,
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
				Type:     schema.TypeString,
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

			"account_owner": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func CreateUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	req := &InvitationCreateRequestWrap{
		InvitationCreateRequest{
			User{
				Name:  d.Get("name").(string),
				Email: d.Get("email").(string),
			},
		},
	}
	res := &InvitationWrap{}

	err := client.Post([]string{"invitations"}, req, res)
	if err != nil {
		return err
	}

	id := strconv.Itoa(res.Invitation.ID)
	d.SetId(id)
	d.Set("id", id)

	return nil
}

func ReadUser(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*Client)

	return fmt.Errorf("ReadUser not yet implemented")
}

func UpdateUser(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*Client)

	return fmt.Errorf("UpdateUser not yet implemented")
}

func DeleteUser(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*Client)

	return fmt.Errorf("DeleteUser not yet implemented")
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
	Username       string `json:"username"`
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
