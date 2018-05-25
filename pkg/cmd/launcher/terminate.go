package launcher

import (
	"fmt"
	"os"
	"path"

	"github.com/alyu/configparser"
	"github.com/sjenning/cloud-launcher/pkg/cloudprovider/aws"
	"github.com/sjenning/cloud-launcher/pkg/cmd"
	"github.com/spf13/cobra"
)

const defaultCloudLauncherSection = "cloud-launcher:ids"

type TerminateOptions struct {
	clusterName         string
	inventoryOutputFile string
	region              string
	terminateInstance   bool
}

func NewTerminateCommand() *cobra.Command {
	t := &TerminateOptions{
		region: "us-east-1",
	}
	c := &cobra.Command{
		Use:   "terminate",
		Short: "Terminate instances for a cluster",
		Run: func(c *cobra.Command, args []string) {
			cmd.CheckError(t.Run(c))
		},
	}
	flags := c.Flags()
	flags.StringVar(&t.clusterName, "cluster-name", t.clusterName, "name of the cluster, required")
	flags.StringVar(&t.region, "region", t.region, "cloud provider region in which to create instances")
	flags.StringVar(&t.inventoryOutputFile, "inventory-output-file", t.inventoryOutputFile, "output inventory file (default $HOME/<cluster-name>.inventory")
	flags.BoolVar(&t.terminateInstance, "terminate-instance", t.terminateInstance, "defaults to marking an instance for termination")
	return c
}

func (t *TerminateOptions) Run(c *cobra.Command) error {
	if t.inventoryOutputFile == "" {
		t.inventoryOutputFile = path.Join(os.Getenv("HOME"), t.clusterName+".inventory")
	}
	provider, err := aws.NewCloudProvider(&aws.Config{
		Region:            t.region,
		TerminateInstance: t.terminateInstance,
	})
	if err != nil {
		return err
	}
	cfg, err := configparser.Read(t.inventoryOutputFile)
	if err != nil {
		return err
	}
	ids, err := cfg.Sections(defaultCloudLauncherSection)
	if err != nil {
		return fmt.Errorf("%s section could not be found", defaultCloudLauncherSection)
	}
	for _, id := range ids[0].OptionNames() {
		if len(id) == 0 {
			continue
		}
		fmt.Printf("Deleting Instance %s\n", id)
		provider.DeleteInstance(id)
	}
	return nil
}
