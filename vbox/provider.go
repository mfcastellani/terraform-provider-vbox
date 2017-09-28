package vbox

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		// Schema is where you list the parameters of your provider. For instance, with AWS
		// we have the access_key and the secret_key. For this example we are using some fake
		// data to simulate authentication (VirtualBox don't need auth).
		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("API_USER", nil),
				Description: "Fake API User",
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("API_TOKEN", nil),
				Description: "Fake API Token",
			},
		},

		// ResourceMap is the list of the resources
		// managed by your provider.
		ResourcesMap: map[string]*schema.Resource{
			"vbox_disk":     resourceVboxDisk(),
			"vbox_instance": resourceVboxInstance(),
		},

		// ConfigureFunc is the function which, among
		// other things, instantiates and configures
		// the client you use to interact with the targeted
		// API (AWS SDK for example).
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	// Just read the values and dispose
	user := d.Get("user").(string)
	token := d.Get("token").(string)
	_ = user
	_ = token

	return nil, nil
}
