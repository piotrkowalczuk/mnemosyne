// Package discovery is under development.
package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type service struct {
	Address string `json:"ServiceAddress"`
	Port    int    `json:"ServicePort"`
}

// DiscoverHTTP ...
func DiscoverHTTP(endpoint string) ([]string, error) {
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("discovery: request failure: %s", err.Error())
	}
	defer res.Body.Close()
	var (
		tmp      []service
		services []string
	)
	if err := json.NewDecoder(res.Body).Decode(&tmp); err != nil {
		return nil, fmt.Errorf("discovery: request payload decoding failure: %s", err.Error())
	}
	for _, s := range tmp {
		services = append(services, fmt.Sprintf("%s:%d", s.Address, s.Port))
	}
	return services, nil
}

// DiscoverDNS ...
func DiscoverDNS(address string) ([]string, error) {
	_, addresses, err := net.LookupSRV("mnemosyned", "grpc", address)
	if err != nil {
		return nil, err
	}

	if len(addresses) == 0 {
		return nil, errors.New("discovery: srv lookup retured nothing")
	}
	var (
		tmp      []service
		services []string
	)
	for _, addr := range addresses {
		tmp = append(tmp, service{Address: addr.Target, Port: int(addr.Port)})
	}
	for _, s := range tmp {
		services = append(services, fmt.Sprintf("%s:%d", s.Address, s.Port))
	}
	return services, nil
}
