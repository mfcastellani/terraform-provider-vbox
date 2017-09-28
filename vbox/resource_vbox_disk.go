package vbox

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"os/exec"
	"strings"
)

func resourceVboxDisk() *schema.Resource {
	return &schema.Resource{
		Create:   resourceVboxDiskCreate,
		Read:     resourceVboxDiskRead,
		Update:   resourceVboxDiskUpdate,
		Delete:   resourceVboxDiskDelete,
		Exists:   resourceVboxDiskExists,
		Importer: nil,

		Schema: map[string]*schema.Schema{
			"size": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Disk Size",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Disk Name",
			},
		},
	}
}

func resourceVboxDiskExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	return false, nil
}

func resourceVboxDiskCreate(d *schema.ResourceData, meta interface{}) error {
	var diskSize int
	if v, ok := d.GetOk("size"); ok {
		diskSize = v.(int)
	}

	var diskName string
	if v, ok := d.GetOk("name"); ok {
		diskName = v.(string)
	}

	var err error
	localDiskName := strings.Join([]string{diskName, ".vdi"}, "")
	cmd := exec.Command("vboxmanage", "createhd", "--filename", localDiskName, "--size", fmt.Sprint(diskSize))
	if err = cmd.Run(); err != nil {
		return err
	}

	d.SetId(diskName)

	return nil
}

func resourceVboxDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVboxDiskRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVboxDiskDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
