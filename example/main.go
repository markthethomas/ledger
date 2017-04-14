package main

import (
	"fmt"

	"github.com/markthethomas/ledger"
	"github.com/spf13/viper"
)

func main() {
	serviceConfig := ledger.ServiceConfig{
		Name: "fidelius",
		Tags: []string{"service", "api", "fidelius"},
		URL:  "http://localhost:3000",
		Port: viper.GetInt("SERVICE_PORT"),
	}
	ledger.Setup(serviceConfig)
	ledger.Register()

	r, err := ledger.LookupKV("fidelius:url")
	fmt.Println(string(*r))
	fmt.Println(err)
	defer ledger.Deregister()
}
