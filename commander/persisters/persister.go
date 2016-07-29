package persisters

import (
	"log"
	"net"
	"time"
)

// Persister - interface defining what we can do with persistable structs
type Persister interface {
	Save(Persistable) error
	GetInto(Persistable, interface{}, ...interface{}) error
	Delete(Persistable) error
}

// Persistable - interface defining what is persistable
type Persistable interface {
	GetUUID() string
}

func waitForService(address string) bool {

	for i := 0; i < 12; i++ {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Println("Connection error:", err)
		} else {
			conn.Close()
			log.Println("Got connection!")
			return true
		}
		time.Sleep(5 * time.Second)
	}

	return false
}
