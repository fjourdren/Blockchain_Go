package network

import (
	"blockchain-go/internal/blockchain"
	"blockchain-go/internal/config"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

// Manager handles network communication and peer management
type Manager struct {
	mu               sync.RWMutex
	me               *Peer
	blockchain       *blockchain.Blockchain
	peers            []*Peer
	lastBlockIndex   int
	config           config.NetworkConfig
	broadcastManager *BroadcastManager
}

// NewManager creates a new network manager
func NewManager(cfg config.NetworkConfig, bc *blockchain.Blockchain) *Manager {
	peerID := GeneratePeerID()
	me := NewPeer(peerID, 0, cfg.Host, cfg.Port)

	return &Manager{
		me:               me,
		blockchain:       bc,
		peers:            make([]*Peer, 0),
		config:           cfg,
		broadcastManager: NewBroadcastManager(),
	}
}

// NewJoiningManager creates a network manager that joins an existing network
func NewJoiningManager(cfg config.NetworkConfig, bc *blockchain.Blockchain, initHost string, initPort int) *Manager {
	manager := NewManager(cfg, bc)

	// Add initial peer
	initPeer := NewPeer("0", 0, initHost, initPort)
	manager.AddPeer(initPeer)

	// Join the network
	if err := manager.joinNetwork(initPeer); err != nil {
		log.Printf("Failed to join network: %v", err)
		return manager
	}

	// Sync the full chain from the peer
	if err := manager.SyncFullChainFromPeer(initPeer); err != nil {
		log.Fatalf("Failed to sync blockchain from peer: %v", err)
	}

	return manager
}

// StartServer starts the TCP server
func (m *Manager) StartServer() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(m.me.Port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server listening on port %d", m.me.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go m.handleConnection(conn)
	}
}

// handleConnection handles an incoming TCP connection
func (m *Manager) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Set read deadline
	if err := conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		return
	}

	// Read data
	buffer := make([]byte, 8192)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	// Process packet
	response, err := m.processPacket(buffer[:n])
	if err != nil {
		response = []byte{}
	}

	// Always send a response
	if response == nil {
		response = []byte{}
	}

	// Send response
	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return
	}

	conn.Write(response)
}

// processPacket processes an incoming packet and returns a response
func (m *Manager) processPacket(data []byte) ([]byte, error) {
	packet, err := FromJSON(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse packet: %w", err)
	}

	switch packet.Type {
	case PacketTypeSingle:
		return m.handleSinglePacket(packet)
	case PacketTypeBroadcast:
		return m.handleBroadcastPacket(packet)
	default:
		return nil, fmt.Errorf("unknown packet type: %s", packet.Type)
	}
}

// handleSinglePacket handles single-target packets
func (m *Manager) handleSinglePacket(packet *Packet) ([]byte, error) {
	switch packet.Name {
	case PacketNameJoin:
		return m.handleJoin(packet)
	case PacketNameJoinAnswer:
		return m.handleJoinAnswer(packet)
	case PacketNameGetLatestBlock:
		return m.handleGetLatestBlock(packet)
	case PacketNameGetLatestBlockAnswer:
		return m.handleGetLatestBlockAnswer(packet)
	case PacketNameDownloadBlock:
		return m.handleDownloadBlock(packet)
	case PacketNameDownloadBlockAnswer:
		return m.handleDownloadBlockAnswer(packet)
	default:
		return nil, fmt.Errorf("unknown packet name: %s", packet.Name)
	}
}

// handleBroadcastPacket handles broadcast packets
func (m *Manager) handleBroadcastPacket(packet *Packet) ([]byte, error) {
	if m.broadcastManager.HasPacket(packet.Index) {
		return []byte{}, nil
	}

	m.broadcastManager.AddPacket(packet.Index)

	switch packet.Name {
	case PacketNameFoundBlock:
		return m.handleFoundBlock(packet)
	default:
		return []byte{}, nil
	}
}

// handleJoin handles a join request
func (m *Manager) handleJoin(packet *Packet) ([]byte, error) {
	if !m.HasPeer(packet.Sender) {
		m.AddPeer(packet.Sender)
	}

	// Update last block index
	latestBlock, err := m.blockchain.GetLatestBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}
	m.lastBlockIndex = latestBlock.Index

	// Send join answer
	managerData, err := m.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize manager: %w", err)
	}

	response := NewPacket(m.me, PacketTypeSingle, PacketNameJoinAnswer, managerData)
	return response.ToJSON()
}

// handleJoinAnswer handles a join answer
func (m *Manager) handleJoinAnswer(packet *Packet) ([]byte, error) {
	var responseData struct {
		Me             *Peer `json:"me"`
		LastBlockIndex int   `json:"last_block_index"`
	}
	
	if err := json.Unmarshal(packet.Content, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal manager data: %w", err)
	}

	// Check if we have peers and the response contains valid peer data
	if len(m.peers) > 0 && responseData.Me != nil {
		m.UpdatePeer(m.peers[0], responseData.Me)
	}

	return []byte{}, nil
}

// handleGetLatestBlock handles a get latest block request
func (m *Manager) handleGetLatestBlock(packet *Packet) ([]byte, error) {
	latestBlock, err := m.blockchain.GetLatestBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	blockData, err := json.Marshal(latestBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal block: %w", err)
	}

	response := NewPacket(m.me, PacketTypeSingle, PacketNameGetLatestBlockAnswer, blockData)
	return response.ToJSON()
}

// handleGetLatestBlockAnswer handles a get latest block answer
func (m *Manager) handleGetLatestBlockAnswer(packet *Packet) ([]byte, error) {
	return []byte{}, nil
}

// handleDownloadBlock handles a download block request
func (m *Manager) handleDownloadBlock(packet *Packet) ([]byte, error) {
	var request struct {
		StartIndex int `json:"start_index"`
		EndIndex   int `json:"end_index"`
	}

	if err := json.Unmarshal(packet.Content, &request); err != nil {
		return nil, fmt.Errorf("failed to unmarshal download request: %w", err)
	}

	blocks, err := m.blockchain.GetBlocks(request.StartIndex, request.EndIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get blocks: %w", err)
	}

	blocksData, err := json.Marshal(blocks)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal blocks: %w", err)
	}

	response := NewPacket(m.me, PacketTypeSingle, PacketNameDownloadBlockAnswer, blocksData)
	return response.ToJSON()
}

// handleDownloadBlockAnswer handles a download block answer
func (m *Manager) handleDownloadBlockAnswer(packet *Packet) ([]byte, error) {
	return []byte{}, nil
}

// handleFoundBlock handles a found block broadcast
func (m *Manager) handleFoundBlock(packet *Packet) ([]byte, error) {
	blockIndexStr := string(packet.Content)
	blockIndex, err := strconv.Atoi(blockIndexStr)
	if err != nil {
		return []byte{}, nil
	}

	// Check if we already have this block or a newer one
	latestBlock, err := m.blockchain.GetLatestBlock()
	if err == nil && latestBlock.Index >= blockIndex {
		return []byte{}, nil
	}

	// Try to sync chain if we're behind
	m.SyncChain(packet.Sender, blockIndex)

	return []byte{}, nil
}

// joinNetwork joins an existing network
func (m *Manager) joinNetwork(initPeer *Peer) error {
	// Send join request
	joinData, err := NewPacket(m.me, PacketTypeSingle, PacketNameJoin, []byte("{}")).ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize join packet: %w", err)
	}

	response, err := initPeer.SendTCP(joinData)
	if err != nil {
		return fmt.Errorf("failed to send join request: %w", err)
	}

	// Process join answer
	_, err = m.processPacket(response)
	if err != nil {
		return fmt.Errorf("failed to process join answer: %w", err)
	}

	return nil
}

// AddPeer adds a peer to the network
func (m *Manager) AddPeer(peer *Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.hasPeerInternal(peer) {
		m.peers = append(m.peers, peer)
		m.me.Popularity = len(m.peers)
	}
}

// RemovePeer removes a peer from the network
func (m *Manager) RemovePeer(peer *Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, p := range m.peers {
		if p.IsEqual(peer) {
			m.peers = append(m.peers[:i], m.peers[i+1:]...)
			m.me.Popularity = len(m.peers)
			return
		}
	}
}

// UpdatePeer updates peer information
func (m *Manager) UpdatePeer(oldPeer, newPeer *Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, p := range m.peers {
		if p.IsEqual(oldPeer) {
			m.peers[i] = newPeer
			return
		}
	}
}

// HasPeer checks if a peer exists in the network
func (m *Manager) HasPeer(peer *Peer) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.hasPeerInternal(peer)
}

// hasPeerInternal checks if a peer exists (internal use, no locking)
func (m *Manager) hasPeerInternal(peer *Peer) bool {
	if peer == nil {
		return false
	}

	for _, p := range m.peers {
		if p.IsEqual(peer) {
			return true
		}
	}
	return false
}

// GetPeerIndex returns the index of a peer
func (m *Manager) GetPeerIndex(peer *Peer) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if peer == nil {
		return -1
	}
	
	for i, p := range m.peers {
		if p.IsEqual(peer) {
			return i
		}
	}
	return -1
}

// Broadcast sends a message to all peers
func (m *Manager) Broadcast(data []byte) {
	m.mu.RLock()
	peers := make([]*Peer, len(m.peers))
	copy(peers, m.peers)
	m.mu.RUnlock()

	for _, peer := range peers {
		if peer == nil {
			continue
		}
		
		go func(p *Peer) {
			if _, err := p.SendTCP(data); err != nil {
				log.Printf("Failed to broadcast to peer %s: %v", p.String(), err)
			}
		}(peer)
	}
}

// SyncChain synchronizes the blockchain with a peer
func (m *Manager) SyncChain(peer *Peer, targetIndex int) bool {
	if peer == nil {
		return false
	}
	
	latestBlock, err := m.blockchain.GetLatestBlock()
	if err != nil {
		return false
	}

	if latestBlock.Index >= targetIndex {
		return false
	}

	// Download missing blocks
	blocks, err := m.DownloadBlocks(peer, latestBlock.Index+1, targetIndex)
	if err != nil {
		return false
	}

	// Add blocks to chain
	for _, block := range blocks {
		m.blockchain.AddBlockWithoutVerification(block)
	}

	return true
}

// DownloadBlocks downloads blocks from a peer
func (m *Manager) DownloadBlocks(peer *Peer, startIndex, endIndex int) ([]*blockchain.Block, error) {
	if peer == nil {
		return nil, fmt.Errorf("cannot download blocks from nil peer")
	}
	
	request := struct {
		StartIndex int `json:"start_index"`
		EndIndex   int `json:"end_index"`
	}{
		StartIndex: startIndex,
		EndIndex:   endIndex,
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	packet := NewPacket(m.me, PacketTypeSingle, PacketNameDownloadBlock, requestData)
	packetData, err := packet.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize packet: %w", err)
	}

	response, err := peer.SendTCP(packetData)
	if err != nil {
		return nil, fmt.Errorf("failed to send download request: %w", err)
	}

	responsePacket, err := FromJSON(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var blocks []*blockchain.Block
	if err := json.Unmarshal(responsePacket.Content, &blocks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal blocks: %w", err)
	}

	return blocks, nil
}

// ToJSON serializes the manager to JSON
func (m *Manager) ToJSON() ([]byte, error) {
	managerData := struct {
		Me             *Peer `json:"me"`
		LastBlockIndex int   `json:"last_block_index"`
	}{
		Me:             m.me,
		LastBlockIndex: m.lastBlockIndex,
	}

	return json.Marshal(managerData)
}

// GetMe returns the local peer
func (m *Manager) GetMe() *Peer {
	return m.me
}

// GetBlockchain returns the blockchain instance
func (m *Manager) GetBlockchain() *blockchain.Blockchain {
	return m.blockchain
}

// SyncFullChainFromPeer synchronizes the entire blockchain from a peer
func (m *Manager) SyncFullChainFromPeer(peer *Peer) error {
	// Get latest block from peer
	latestBlockData, err := NewPacket(m.me, PacketTypeSingle, PacketNameGetLatestBlock, nil).ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize latest block request: %w", err)
	}

	latestBlockResponse, err := peer.SendTCP(latestBlockData)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}

	responsePacket, err := FromJSON(latestBlockResponse)
	if err != nil {
		return fmt.Errorf("failed to parse latest block response: %w", err)
	}

	var block blockchain.Block
	if err := json.Unmarshal(responsePacket.Content, &block); err != nil {
		return fmt.Errorf("failed to unmarshal latest block: %w", err)
	}

	// Download all blocks
	blocks, err := m.DownloadBlocks(peer, 0, block.Index)
	if err != nil {
		return fmt.Errorf("failed to download blocks: %w", err)
	}

	// Replace the entire chain
	m.blockchain.ReplaceChain(blocks)

	return nil
}
