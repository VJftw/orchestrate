package cadet

import (
	"time"

	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
)

// Cadet - a node to run containers on
type Cadet struct {
	ID           uint      `json:"-" gorm:"primary_key"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	CadetGroupID uint      `json:"-"`
	UUID         string    `json:"uuid" gorm:"unique"`
	Key          string    `json:"key" gorm:"not null"`
}

// GetUUID - Returns the UUID
func (c Cadet) GetUUID() string {
	return c.UUID
}

type Message struct {
	Key  string `json:"key"`
	Data *Node  `json:"data"`
}

type Node struct {
	Metrics    []Metric     `json:"metrics"`
	Containers []*Container `json:"containers"`
}

func NewMessage() *Message {
	n := &Node{
		Metrics:    []Metric{},
		Containers: []*Container{},
	}
	return &Message{
		Data: n,
	}
}

type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type Container struct {
	Name             string                    `json:"name"`
	Config           *container.Config         `json:"config"`
	HostConfig       *container.HostConfig     `json:"hostConfig"`
	NetworkingConfig *network.NetworkingConfig `json:"networkingConfig"`
}
