package launcher

import (
	"fmt"

	"github.com/sjenning/cloud-launcher/pkg/cloudprovider/aws"
	"github.com/sjenning/cloud-launcher/pkg/cmd"
	"github.com/spf13/cobra"
)

type StopOptions struct {
	clusterName string
	region      string
}

func NewStopCommand() *cobra.Command {
	t := &StopOptions{
		region: defaultRegion,
	}
	c := &cobra.Command{
		Use:   "stop",
		Short: "Stop instances for a cluster",
		Run: func(c *cobra.Command, args []string) {
			cmd.CheckError(t.Run(c))
		},
	}
	flags := c.Flags()
	flags.StringVar(&t.region, "region", t.region, "cluster region")
	flags.StringVar(&t.clusterName, "cluster-name", t.clusterName, "name of the cluster, required")
	return c
}

func (t *StopOptions) Run(c *cobra.Command) error {
	if len(t.clusterName) == 0 {
		return fmt.Errorf("--cluster-name must be specified")
	}
	provider, err := aws.NewCloudProvider(&aws.Config{
		Region: t.region,
	})
	if err != nil {
		return err
	}
	ids, err := provider.GetInstanceIDsByClusterName(t.clusterName)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		fmt.Printf("No instances found\n")
	} else {
		for _, id := range ids {
			fmt.Printf("Stopping %s\n", id)
			provider.TagInstance(id, "Name", "-terminate")
		}
	}
	return nil
}
