# Terraform Provider for VirtualBox

This is a awesome Terraform provider for VirtualBox. Is written in `Go` and you need to read this Readme until the end to understand some key concepts. You also need VirtualBox up and running. 

## Configure the development environment

### Go and Hashicorp libraries

1- You need `go` up and running. Follow this [tutorial](https://golang.org/doc/install). **Important: documentation talks about go 1.8, but go 1.9 is better - TRUST ME**. 

2- You need to define `GOPATH` and `GOBIN` env vars. I use `$HOME/go` for `GOPATH` and `$HOME/work/bin` for `GOBIN`, but we are living in a free country and you can choose any folder you want to use. Usually a simple `export GOPATH=$HOME/go` and `export GOBIN=$HOME/work/bin` at your `.bashrc` file solves this, but if you uses something more exoteric than `bash` ask Google to help you. DuckDuckGo and Bing are options also. Don't forget to add `GOBIN` path to you path (`export PATH=$PATH:$GOBIN`). 

3- Did I talked about clone this repo? You know how to do it, right? One point, this should be inside the `GOPATH/src` folder. 

4- Install Hashicorp required packages:

```
$ go get github.com/hashicorp/terraform
$ go get github.com/terraform-providers/terraform-provider-template
$ go install github.com/hashicorp/terraform
$ go install github.com/terraform-providers/terraform-provider-template
```


### Managing dependencies

The [Terraform guide for plugins](https://www.terraform.io/docs/plugins/provider.html) says "Note that some packages from the Terraform repository are used as library dependencies by providers, such as github.com/hashicorp/terraform/helper/schema; it is recommended to use [govendor](https://github.com/kardianos/govendor) to create a local vendor copy of the relevant packages in the provider repository, as can be seen in the repositories within the terraform-providers GitHub organization". 

To install `govendor` just run: 

```
$ go get -u github.com/kardianos/govendor
```

This will be installed at the `GOBIN` path as defined at step 2 in the previous section of this tutorial. You can check if everything is ok using:

```
$ govendor list
```

To vendorize all dependencies in a **vendor** folder just run:

```
$ govendor add +external
```

I removed the *vendor* folder from this repository, so you need to run 

```
$ govendor init
$ govendor add +external
```

Now run the old-but-amazing `make` to check if everything is on the right place. 

```
$ make
rm -f terraform-provider-vbox
cd /home/marcelo/go/src/vbox-provider; \
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.VERSION=0.0.1" -o terraform-provider-vbox .
```

### Check if VirtualBox is up and running

A VirtualBox or VB is a software virtualization package that installs on an operating system as an application. VirtualBox allows additional operating systems to be installed on it, as a Guest OS, and run in a virtual environment. This will be used to create resources with our Terraform Provider. To do the magic we will use the awesome tool called `vboxmanage`, but I'll not explain how to install VirtualBox here. Just Google it (or DuckDuckGo, or Bing, you choose). Check if `vboxmanage` is available at your environment using:

```
$ vboxmanage --version
5.0.40_Ubuntur115130
```

## Writing the Provider

Basically a provider is composed of two parts: the **provider** itself and some **resources**. Since we are using Terraform to manage our AWS infrastructure I believe you already saw something like this:

```
provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region     = "${var.aws_region}"
}


resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}
```

Here, the `aws` provider part is taking care of setting up your AWS Client and authenticate you to the AWS API. Then the `aws_vpc` resource will create an AWS VPC with the correct CIDR block for you, using this AWS Client.

If you take a look at the existing [providers](https://github.com/terraform-providers), you will notice that the structure is almost always something like :

* **provider.go**: Implement the *core* of the Provider.
* **config.go**: Configure the API client with the credentials from the Provider.
* **resource_<resource_name>.go**: Implement a specific resource handler with the CRUD functions.
* **import_<resource_name>.go**: Make possible to import existing resources. 
* **data_source_<resource_name>.go**: Used to fetch data from outside of Terraform to be used in other resources. For example, you will be able to fetch the latest AWS AMI ID and use it for an AWS instance. Same as *import*.

I tried to keep this POC as closer as possible to this structure. This is not a pattern, like the old-and-good MVC, but should be since every provider implements something closer to this. 

### The Schema

At our `provider.go` I defined two properties for fake authentication (VirtualBox don't use it, is just for fun). The idea is simulate a provider with a basic authentication with `user` and `token`. Both are required and are defined as [schema.Schema types](https://godoc.org/github.com/hashicorp/terraform/helper/schema#Schema). 

```
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
```

### The Resources 

I created two resources for this POC, one called `vbox_disk` and another called `vbox_instance`. The idea is create a disk to allocate space, and create a machine and attach the disk. I'll not start the machine or install any OS, just provide infrastructure.

The `vbox_disk` have two attributes: **size** and **name**. Both are required by `vboxmanage` command to create a basic disk and register at VirtualBox, **size** is specified in MB (1024 means 1GB) and **name** will be the file name for .vdi file and is required by the instance.

The `vbox_instance` have three attributes: **name**, **diskname** and **osname**. You can define any name to the instance, for example *banana*. To **diskname** you should use the same disk name used at disk creation and **osname** must be a valid VirtualBox os name, like *Ubuntu* or *ArchLinux* or any listed with the command `vboxmanage list ostypes`.

Resources are described using the [schema.Resource](https://godoc.org/github.com/hashicorp/terraform/helper/schema#Resource) structure. This structure has the following fields:

* Schema - The configuration schema for this resource. Schemas are covered in more detail below.
* Create, Read, Update, and Delete - These are the callback functions that implement CRUD operations for the resource. The only optional field is Update. If your resource doesn't support update, then you may keep that field nil.
* Importer - If this is non-nil, then this resource is importable. It is recommended to implement this.  

The CRUD operations in more detail, along with their contracts:

**Create** - This is called to create a new instance of the resource. Terraform guarantees that an existing ID is not set on the resource data. That is, you're working with a new resource. Therefore, you are responsible for calling SetId on your schema.ResourceData using a value suitable for your resource. This ensures whatever resource state you set on schema.ResourceData will be persisted in local state. If you neglect to SetId, no resource state will be persisted.

**Read** - This is called to resync the local state with the remote state. Terraform guarantees that an existing ID will be set. This ID should be used to look up the resource. Any remote data should be updated into the local data. No changes to the remote resource are to be made.

**Update** - This is called to update properties of an existing resource. Terraform guarantees that an existing ID will be set. Additionally, the only changed attributes are guaranteed to be those that support update, as specified by the schema. Be careful to read about partial states below.

**Delete** - This is called to delete the resource. Terraform guarantees an existing ID will be set.

**Exists** - This is called to verify a resource still exists. It is called prior to Read, and lowers the burden of Read to be able to assume the resource exists. If the resource is no longer present in remote state, calling SetId with an empty string will signal its removal.

