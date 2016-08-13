package node

import (
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
)

type Node struct {
	Metrics    []Metric     `json:"metrics"`
	Containers []*Container `json:"containers"`
}

func NewNode() *Node {
	return &Node{
		Metrics:    []Metric{},
		Containers: []*Container{},
	}
}

type Metric interface {
	Update()
}

type Container struct {
	Name             string                    `json:"name"`
	Config           *container.Config         `json:"config"`
	HostConfig       *container.HostConfig     `json:"hostConfig"`
	NetworkingConfig *network.NetworkingConfig `json:"networkingConfig"`
}
