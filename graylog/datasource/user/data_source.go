package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: read,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the user. Either user_id or username must be set.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The username. Either user_id or username must be set.",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of the user.",
			},
			"full_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full name of the user.",
			},
			"first_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The first name of the user.",
			},
			"last_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last name of the user.",
			},
			"timezone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timezone of the user.",
			},
			"session_timeout_ms": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The session timeout in milliseconds.",
			},
			"roles": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The roles assigned to the user.",
			},
			"permissions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The permissions assigned to the user.",
			},
			"external": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the user is external (e.g., LDAP).",
			},
			"read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the user is read-only.",
			},
			"session_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the user has an active session.",
			},
			"last_activity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the user's last activity.",
			},
			"client_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The client address of the user's last session.",
			},
			"account_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account status (enabled, disabled, deleted).",
			},
			"service_account": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is a service account.",
			},
		},
	}
}
