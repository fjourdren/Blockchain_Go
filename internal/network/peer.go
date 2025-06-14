package network

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"time"
)

// Peer represents a network peer in the blockchain network
type Peer struct {
	ID         string `json:"id"`
	Popularity int    `json:"popularity"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
}

// NewPeer creates a new peer instance
func NewPeer(id string, popularity int, host string, port int) *Peer {
	return &Peer{
		ID:         id,
		Popularity: popularity,
		Host:       host,
		Port:       port,
	}
}

// GeneratePeerID generates a unique peer ID
func GeneratePeerID() string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("%d", timestamp)

	hasher := sha512.New()
	hasher.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:16]
}

// GetAddress returns the full address of the peer
func (p *Peer) GetAddress() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

// SendTCP sends data to the peer via TCP and returns the response
func (p *Peer) SendTCP(data []byte) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", p.GetAddress(), 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to peer %s: %w", p.GetAddress(), err)
	}
	defer conn.Close()

	// Set write deadline
	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return nil, fmt.Errorf("failed to set write deadline: %w", err)
	}

	// Send data
	if _, err := conn.Write(data); err != nil {
		return nil, fmt.Errorf("failed to send data to peer: %w", err)
	}

	// Set read deadline
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	// Read response
	buffer := make([]byte, 8192)
	n, err := conn.Read(buffer)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		
		return nil, fmt.Errorf("failed to read response from peer: %w", err)
	}

	return buffer[:n], nil
}

// IsEqual checks if two peers are the same
func (p *Peer) IsEqual(other *Peer) bool {
	if p == nil || other == nil {
		return false
	}
	
	return p.ID == other.ID && p.Host == other.Host && p.Port == other.Port
}

// String returns a string representation of the peer
func (p *Peer) String() string {
	if p == nil {
		return "Peer <nil>"
	}
	
	return fmt.Sprintf("Peer %s (%s:%d, popularity: %d)", p.ID, p.Host, p.Port, p.Popularity)
}
