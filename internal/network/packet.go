package network

import (
	"encoding/json"
	"fmt"
)

// PacketType represents the type of network packet
type PacketType string

const (
	PacketTypeSingle    PacketType = "SINGLE"
	PacketTypeBroadcast PacketType = "BROADCAST"
)

// PacketName represents the specific operation of a packet
type PacketName string

const (
	PacketNameJoin                 PacketName = "JOIN"
	PacketNameJoinAnswer           PacketName = "JOINANSWER"
	PacketNameGetLatestBlock       PacketName = "GETLATESTBLOCK"
	PacketNameGetLatestBlockAnswer PacketName = "GETLATESTBLOCKANSWER"
	PacketNameDownloadBlock        PacketName = "DOWNLOADBLOCK"
	PacketNameDownloadBlockAnswer  PacketName = "DOWNLOADBLOCKANSWER"
	PacketNameFoundBlock           PacketName = "FOUNDBLOCK"
)

// Packet represents a network packet for communication between peers
type Packet struct {
	Sender  *Peer      `json:"sender"`
	Type    PacketType `json:"type"`
	Name    PacketName `json:"name"`
	Content []byte     `json:"content"`
	Index   int        `json:"index"`
}

// NewPacket creates a new packet instance
func NewPacket(sender *Peer, packetType PacketType, name PacketName, content []byte) *Packet {
	return &Packet{
		Sender:  sender,
		Type:    packetType,
		Name:    name,
		Content: content,
		Index:   0,
	}
}

// NewBroadcastPacket creates a new broadcast packet
func NewBroadcastPacket(sender *Peer, name PacketName, content []byte, index int) *Packet {
	return &Packet{
		Sender:  sender,
		Type:    PacketTypeBroadcast,
		Name:    name,
		Content: content,
		Index:   index,
	}
}

// ToJSON serializes the packet to JSON
func (p *Packet) ToJSON() ([]byte, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal packet: %w", err)
	}
	return data, nil
}

// FromJSON deserializes a packet from JSON
func FromJSON(data []byte) (*Packet, error) {
	var packet Packet
	if err := json.Unmarshal(data, &packet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal packet: %w", err)
	}
	return &packet, nil
}

// String returns a string representation of the packet
func (p *Packet) String() string {
	return fmt.Sprintf("Packet{Type: %s, Name: %s, Sender: %s}", p.Type, p.Name, p.Sender.ID)
}
