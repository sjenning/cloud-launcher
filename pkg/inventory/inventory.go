package inventory

import (
	"fmt"
	"os"

	"github.com/golang/go/src/text/template"
)

type Role string

const (
	RoleMaster  = "master"
	RoleInfra   = "infra"
	RoleCompute = "compute"
)

type Interface interface {
	AddNode(ip string, role Role)
	Render() error
}

type Config struct {
	Version      string
	Token        string
	ClusterName  string
	AWSAccessKey string
	AWSSecretKey string
	TemplateFile string
	OutputFile   string
}

type Label struct {
	Key   string
	Value string
}

type Host struct {
	IP        string
	Labels    []Label
	NodeGroup string
}

type inventory struct {
	Config
	Etcd    []Host
	Masters []Host
	Nodes   []Host
}

var _ Interface = &inventory{}

func New(config Config) Interface {
	return &inventory{
		Config: config,
	}
}

func (i *inventory) AddNode(ip string, role Role) {
	switch role {
	case RoleMaster:
		i.Etcd = append(i.Etcd, Host{IP: ip})
		i.Masters = append(i.Masters, Host{IP: ip})
		i.Nodes = append(i.Nodes, Host{
			IP:        ip,
			NodeGroup: "node-config-master",
		})
	case RoleInfra:
		i.Nodes = append(i.Nodes, Host{
			IP:        ip,
			NodeGroup: "node-config-infra",
			Labels: []Label{
				{
					Key:   "region",
					Value: "infra",
				},
			},
		})
	case RoleCompute:
		i.Nodes = append(i.Nodes, Host{
			IP:        ip,
			NodeGroup: "node-config-compute",
		})
	}

}

func (i *inventory) Render() error {
	file, err := os.Create(i.OutputFile)
	if err != nil {
		return fmt.Errorf("Unable to create output file %s: %v", i.OutputFile, err)
	}
	defer file.Close()

	template, err := template.ParseFiles(i.TemplateFile)
	if err != nil {
		return fmt.Errorf("Unable to parse template file %s: %v", i.TemplateFile, err)
	}
	err = template.Execute(file, i)
	if err != nil {
		return fmt.Errorf("Unable to render template: %v", err)
	}
	return nil
}
