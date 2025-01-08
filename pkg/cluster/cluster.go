package cluster

import (
	"hash/crc32"
	"sort"
	"sync"
)

type Node struct {
	ID      string
	Address string
}

type ConsistentHash struct {
	nodes     map[uint32]Node
	sortedIDs []uint32
	mu        sync.RWMutex
}

func NewConsistentHash() *ConsistentHash {
	return &ConsistentHash{
		nodes: make(map[uint32]Node),
	}
}

func (ch *ConsistentHash) HashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (ch *ConsistentHash) AddNode(node Node) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	hash := ch.HashKey(node.ID)
	ch.nodes[hash] = node
	ch.sortedIDs = append(ch.sortedIDs, hash)
	sort.Slice(ch.sortedIDs, func(i, j int) bool { return ch.sortedIDs[i] < ch.sortedIDs[j] })
}

func (ch *ConsistentHash) GetNode(key string) (Node, bool) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if len(ch.sortedIDs) == 0 {
		return Node{}, false
	}

	keyHash := ch.HashKey(key)
	idx := sort.Search(len(ch.sortedIDs), func(i int) bool {
		return ch.sortedIDs[i] >= keyHash
	})

	if idx == len(ch.sortedIDs) {
		idx = 0
	}

	return ch.nodes[ch.sortedIDs[idx]], true
}

func (ch *ConsistentHash) RemoveNode(nodeID string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	hash := ch.HashKey(nodeID)
	delete(ch.nodes, hash)

	for i, id := range ch.sortedIDs {
		if id == hash {
			ch.sortedIDs = append(ch.sortedIDs[:i], ch.sortedIDs[i+1:]...)
			break
		}
	}
}
