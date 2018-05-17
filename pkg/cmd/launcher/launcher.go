package launcher

import (
	"flag"

	"github.com/spf13/cobra"
)

func NewCommand(name string) *cobra.Command {
	c := &cobra.Command{
		Use:   name,
		Short: "Start a set of instances for use as OpenShift cluster nodes",
	}

	c.AddCommand(
		NewStartCommand(),
	)

	// add the glog flags
	c.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	flag.CommandLine.Parse([]string{})
	return c
}
