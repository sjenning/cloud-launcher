package launcher

import (
	"fmt"
	"os"

	"github.com/sjenning/cloud-launcher/pkg/cloudprovider/aws"
	"github.com/sjenning/cloud-launcher/pkg/cmd"
	"github.com/sjenning/cloud-launcher/pkg/inventory"
	"github.com/spf13/cobra"
)

type StartOptions struct {
	clusterName         string
	version             string
	token               string
	region              string
	imageID             string
	instanceType        string
	subnetID            string
	keyName             string
	inventoryTemplate   string
	inventoryOutputFile string
}

const defaultVersion = "3.10"

func NewStartCommand() *cobra.Command {
	o := &StartOptions{
		version:           defaultVersion,
		token:             "",
		region:            "us-east-1",
		imageID:           "ami-0af8ebba55e71b379",
		instanceType:      "m4.large",
		subnetID:          "subnet-cf57c596",
		keyName:           "libra",
		inventoryTemplate: os.Getenv("HOME") + "/.aws-launcher-template",
	}

	c := &cobra.Command{
		Use:   "start --cluster-name=CLUSTER-NAME --token=TOKEN",
		Short: "Start instances for a cluster",
		Run: func(c *cobra.Command, args []string) {
			cmd.CheckError(o.Run(c))
		},
	}

	flags := c.Flags()
	flags.StringVar(&o.clusterName, "cluster-name", o.clusterName, "name of the cluster, required")
	flags.StringVar(&o.token, "token", o.token, "token for accessing reg-aws registry, required")
	flags.StringVar(&o.version, "version", o.version, "version of OpenShift to be installed")
	flags.StringVar(&o.region, "region", o.region, "cloud provider region in which to create instances")
	flags.StringVar(&o.imageID, "image-id", o.imageID, "image ID to use for instances")
	flags.StringVar(&o.instanceType, "instance-type", o.instanceType, "instance type to use for instances")
	flags.StringVar(&o.subnetID, "subnet-id", o.subnetID, "subnetID to use for instances")
	flags.StringVar(&o.keyName, "key-name", o.keyName, "ssh key to install on instnaces")
	flags.StringVar(&o.inventoryTemplate, "inventory-template", o.inventoryTemplate, "inventory template file")
	flags.StringVar(&o.inventoryOutputFile, "inventory-output-file", o.inventoryOutputFile, "output inventory file (default <cluster-name>.inventory")

	return c
}

func (o *StartOptions) Run(c *cobra.Command) error {
	if o.clusterName == "" {
		return fmt.Errorf("command requires --cluster-name to be provided")
	}
	if o.token == "" {
		return fmt.Errorf("command requires --token to be provided")
	}
	if o.inventoryOutputFile == "" {
		o.inventoryOutputFile = o.clusterName + ".inventory"
	}

	provider, err := aws.NewCloudProvider(&aws.Config{
		Region:       o.region,
		ImageID:      o.imageID,
		InstanceType: o.instanceType,
		SubnetID:     o.subnetID,
		KeyName:      o.keyName,
	})
	if err != nil {
		return err
	}
	fmt.Println("starting cluster named", o.clusterName)
	master, err := provider.CreateInstance()
	if err != nil {
		return err
	}
	provider.TagInstance(master, "Name", o.clusterName+"-master")
	provider.TagInstance(master, "kubernetes.io/cluster/"+o.clusterName, "true")
	fmt.Println("created master", master)

	infra, err := provider.CreateInstance()
	if err != nil {
		return err
	}

	provider.TagInstance(infra, "Name", o.clusterName+"-infra")
	provider.TagInstance(infra, "kubernetes.io/cluster/"+o.clusterName, "true")
	fmt.Println("created infra", infra)

	node, err := provider.CreateInstance()
	if err != nil {
		return err
	}
	provider.TagInstance(node, "Name", o.clusterName+"-node")
	provider.TagInstance(node, "kubernetes.io/cluster/"+o.clusterName, "true")
	fmt.Println("created node", node)

	fmt.Println("waiting for all instances to be ready")
	provider.WaitForInstance(master)
	provider.WaitForInstance(infra)
	provider.WaitForInstance(node)

	accesskey := "accesskey"
	secretkey := "secretkey"
	credentials, ok := provider.GetCredentials().(aws.Credentials)
	if ok {
		accesskey = credentials.AccessKeyID
		secretkey = credentials.SecretAccessKey
	}

	inv := inventory.New(inventory.Config{
		Version:      o.version,
		Token:        o.token,
		ClusterName:  o.clusterName,
		AWSAccessKey: accesskey,
		AWSSecretKey: secretkey,
		TemplateFile: o.inventoryTemplate,
		OutputFile:   o.inventoryOutputFile,
	})

	fmt.Println("getting instance IP information")
	masterIP, err := provider.GetInstanceIP(master)
	if err != nil {
		return err
	}
	inv.AddNode(masterIP, inventory.RoleMaster)
	infraIP, err := provider.GetInstanceIP(infra)
	if err != nil {
		return err
	}
	inv.AddNode(infraIP, inventory.RoleInfra)
	nodeIP, err := provider.GetInstanceIP(node)
	if err != nil {
		return err
	}
	inv.AddNode(nodeIP, inventory.RoleCompute)

	fmt.Println("generating inventory file at", o.inventoryOutputFile)
	err = inv.Render()
	if err != nil {
		return err
	}

	fmt.Println("done!")
	return nil
}
