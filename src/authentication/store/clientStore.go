package store

import "sync"

type ClientStore struct {
	clients map[string]struct{}
	mutex   sync.RWMutex
}

func (c *ClientStore) Has(entry string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.clients[entry]
	return ok
}

func (c *ClientStore) Add(entry string) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.clients[entry] = struct{}{}
}

func (c *ClientStore) Delete(entry string) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	delete(c.clients, entry)
}

func NewClientStore() ClientStore {
	return ClientStore{
		clients: make(map[string]struct{}),
	}
}
