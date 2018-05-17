package cloudprovider

type Interface interface {
	CreateInstance() (string, error)
	DeleteInstance(id string) error
	WaitForInstance(id string)
	TagInstance(id, key, value string) error
	GetInstanceIP(id string) (string, error)
	GetCredentials() interface{}
}
