package main

import (
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/sjenning/cloud-launcher/pkg/cmd"
	"github.com/sjenning/cloud-launcher/pkg/cmd/launcher"
)

func main() {
	defer glog.Flush()
	baseName := filepath.Base(os.Args[0])
	err := launcher.NewCommand(baseName).Execute()
	cmd.CheckError(err)
}
