package network

import (
	"sync"
	"time"
)

// BroadcastManager manages broadcast packet deduplication
type BroadcastManager struct {
	mu      sync.RWMutex
	packets map[int]time.Time
}

// NewBroadcastManager creates a new broadcast manager
func NewBroadcastManager() *BroadcastManager {
	return &BroadcastManager{
		packets: make(map[int]time.Time),
	}
}

// HasPacket checks if a packet with the given index has been processed
func (bm *BroadcastManager) HasPacket(index int) bool {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	_, exists := bm.packets[index]
	return exists
}

// AddPacket adds a packet to the processed list
func (bm *BroadcastManager) AddPacket(index int) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	bm.packets[index] = time.Now()
}
