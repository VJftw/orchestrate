package connection

import "github.com/vjftw/orchestrate/cadet/node"

type Message struct {
	Key  string     `json:"key"`
	Data *node.Node `json:"data"`
}

func NewMessage(key string, nodeInfo *node.Node) *Message {
	return &Message{
		Key:  key,
		Data: nodeInfo,
	}
}
