package vbox

import (
	"github.com/hashicorp/terraform/helper/schema"
	"os/exec"
	"strings"
)

func resourceVboxInstance() *schema.Resource {
	return &schema.Resource{
		Create:   resourceVboxInstanceCreate,
		Read:     resourceVboxInstanceRead,
		Update:   resourceVboxInstanceUpdate,
		Delete:   resourceVboxInstanceDelete,
		Exists:   resourceVboxInstanceExists,
		Importer: nil,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Virtual Machine Name",
			},
			"diskname": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Disk Name",
			},
			"osname": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Operational System Name",
			},
		},
	}
}

func resourceVboxInstanceExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	return true, nil
}

func resourceVboxInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	var diskName string
	if v, ok := d.GetOk("diskname"); ok {
		diskName = v.(string)
	}

	var osName string
	if v, ok := d.GetOk("osname"); ok {
		diskName = v.(string)
	}

	// This should, obviously, be improved...
	// create VM
	var err error
	cmd := exec.Command("vboxmanage", "createvm", "--name", name, "--ostype", osName, "--register")
	if err = cmd.Start(); err != nil {
		return err
	}
	_ = cmd.Wait()

	// Add SATA Controller
	cmd = exec.Command("vboxmanage", "storagectl", name, "--name", "SATAController", "--add", "sata", "--controller", "IntelAHCI")
	if err = cmd.Start(); err != nil {
		return err
	}
	_ = cmd.Wait()

	// Add Disk
	localDiskName := strings.Join([]string{diskName, ".vdi"}, "")
	cmd = exec.Command("vboxmanage", "storageattach", name, "--storagectl", "SATAController", "--port", "0",
		"--device", "0", "--type", "hdd", "--medium", localDiskName)
	if err = cmd.Start(); err != nil {
		return err
	}
	_ = cmd.Wait()

	d.SetId(diskName)

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
