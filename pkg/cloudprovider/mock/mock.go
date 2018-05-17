package mock

import (
	"fmt"
	"strconv"

	"github.com/sjenning/cloud-launcher/pkg/cloudprovider"
)

type mockCloudProvider struct {
	instanceIndex int
}

var _ cloudprovider.Interface = &mockCloudProvider{}

func NewCloudProvider() (*mockCloudProvider, error) {
	return &mockCloudProvider{}, nil
}

func (p *mockCloudProvider) CreateInstance() (string, error) {
	p.instanceIndex++
	fmt.Println("Created instance", p.instanceIndex)
	return strconv.Itoa(p.instanceIndex), nil
}

func (p *mockCloudProvider) DeleteInstance(id string) error {
	fmt.Println("Deleted instance", id)
	return nil
}

func (p *mockCloudProvider) WaitForInstance(id string) {
}

func (p *mockCloudProvider) TagInstance(id, key, value string) error {
	fmt.Printf("Tagging instance %s with %s=%s\n", id, key, value)
	return nil
}

func (p *mockCloudProvider) GetInstanceIP(id string) (string, error) {
	return "10.0.0." + id, nil
}

func (p *mockCloudProvider) GetCredentials() interface{} {
	return ""
}
