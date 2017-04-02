package ledger

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
)

// Client is the discovery client
var (
	Client       consul.Client
	registration *consul.AgentServiceRegistration
	serviceName  string
	serviceTags  []string
	serviceURL   string
)

// ServiceConfig is what Setup uses to scaffold a service
type ServiceConfig struct {
	Name string
	Tags []string
	URL  string
	Port int
}

// Setup does some basic setup for a service
func Setup(config ServiceConfig) {
	serviceName = config.Name
	if len(config.Tags) > 0 {
		for i := range config.Tags {
			serviceTags = append(serviceTags, config.Tags[i])
		}
	}
	serviceURL = config.URL
	registration = &consul.AgentServiceRegistration{
		ID:   serviceName,
		Name: serviceName,
		Port: config.Port,
		Tags: serviceTags,
	}
}

// Register sets up discovery
func Register() {
	fmt.Println("[service register]")
	client, err := consul.NewClient(&consul.Config{Address: "127.0.0.1:8500"})
	if err != nil {
		fmt.Println(err.Error())
	}
	Client = *client
	Client.Agent().ServiceRegister(registration)
	Client.KV().Put(&consul.KVPair{Key: serviceName + ":healthy", Value: []byte("true")}, nil)
	Client.KV().Put(&consul.KVPair{Key: serviceName + ":url", Value: []byte(serviceURL)}, nil)
}

// Deregister notifies consul we're offline
func Deregister() {
	fmt.Printf("[service deregister] %v\n", serviceName)
	Client.Agent().ServiceRegister(registration)
}

// LookupKV lets a client look up a KV pair
func LookupKV(key string) (*[]byte, error) {
	defaultQueryOpts := consul.QueryOptions{}
	v, _, e := Client.KV().Get(key, &defaultQueryOpts)
	if e != nil {
		return nil, e
	}
	return &v.Value, nil
}

// CheckHealth checks a given service to see if it is healthy
func CheckHealth(service, tag string, passingOnly bool) (healthy bool) {
	checks, _, err := Client.Health().Service(service, tag, passingOnly, nil)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	for i := range checks {
		status := checks[i].Checks.AggregatedStatus()
		if status != consul.HealthPassing {
			return false
		}
	}
	return true
}
