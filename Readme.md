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

## Check if VirtualBox is up and running

A VirtualBox or VB is a software virtualization package that installs on an operating system as an application. VirtualBox allows additional operating systems to be installed on it, as a Guest OS, and run in a virtual environment. This will be used to create resources with our Terraform Provider. To do the magic we will use the awesome tool called `vboxmanage`. Check if is available at your environment using:

```
$ vboxmanage --version
5.0.40_Ubuntur115130
```

I'll not explain how to install VirtualBox here. Just Google it (or DuckDuckGo, or Bing, you choose).


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

### The Resource 

