package vbox

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVboxInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceVboxInstanceCreate,
		Read:   resourceVboxInstanceRead,
		Update: resourceVboxInstanceUpdate,
		Delete: resourceVboxInstanceDelete,
		Exists: resourceVboxInstanceExists,

		Schema: map[string]*schema.Schema{},
	}
}

func resourceVboxInstanceExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	return true, nil
}

func resourceVboxInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVboxInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVboxInstanceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVboxInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
