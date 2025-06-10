package ipfs

import (
	"fmt"
	"net/http"
)

type IpfsService struct {
	jwt     string
	gateway string
	client  *http.Client
}

func NewIpfsService(jwt string) (*IpfsService, error) {
	ipfsService := &IpfsService{
		jwt:    jwt,
		client: &http.Client{},
	}
	gateway, err := ipfsService.getPinataGateway()
	if err != nil {
		return nil, fmt.Errorf("error getting pinata gateway: %w", err)
	}
	ipfsService.gateway = gateway

	return ipfsService, nil
}
