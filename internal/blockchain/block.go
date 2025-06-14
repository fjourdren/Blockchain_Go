package blockchain

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Block represents a single block in the blockchain
type Block struct {
	Index               int    `json:"index"`
	Timestamp           int64  `json:"timestamp"`
	Difficulty          int    `json:"difficulty"`
	NextBlockDifficulty int    `json:"next_block_difficulty"`
	Data                string `json:"data"`
	Hash                string `json:"hash"`
	PreviousHash        string `json:"previous_hash"`
	Nonce               int    `json:"nonce"`
}

// NewBlock creates a new block with the given parameters
func NewBlock(index int, difficulty, nextDifficulty int, data, previousHash string) *Block {
	return &Block{
		Index:               index,
		Timestamp:           time.Now().Unix(),
		Difficulty:          difficulty,
		NextBlockDifficulty: nextDifficulty,
		Data:                data,
		PreviousHash:        previousHash,
		Nonce:               0,
	}
}

// ComputeHash computes the hash of the block
func (b *Block) ComputeHash() string {
	data := fmt.Sprintf("%d%d%s%d%d%d%s",
		b.Index, b.Nonce, b.PreviousHash, b.Difficulty,
		b.NextBlockDifficulty, b.Timestamp, b.Data)

	hasher := sha512.New()
	hasher.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// IsHashValid checks if the block's hash meets the difficulty requirement
func (b *Block) IsHashValid(hash string) bool {
	if b.Difficulty == 0 {
		return true
	}

	prefix := strings.Repeat("0", b.Difficulty)
	return strings.HasPrefix(hash, prefix)
}

// Mine performs proof-of-work mining on the block
func (b *Block) Mine() {
	for !b.IsHashValid(b.Hash) {
		b.Timestamp = time.Now().Unix()
		b.Nonce = int(time.Now().UnixNano() % 4294967296)
		b.Hash = b.ComputeHash()
	}
}

// IsValid validates the block integrity
func (b *Block) IsValid() error {
	// Check if calculated hash matches stored hash
	if calculatedHash := b.ComputeHash(); calculatedHash != b.Hash {
		return fmt.Errorf("block hash mismatch: calculated %s, stored %s", calculatedHash, b.Hash)
	}

	// Check if hash meets difficulty requirement
	if !b.IsHashValid(b.Hash) {
		return fmt.Errorf("block hash does not meet difficulty requirement")
	}

	// Check block size limit (2MB)
	if blockSize := len(b.ToJSON()); blockSize > 2097152 {
		return fmt.Errorf("block size %d exceeds limit of 2MB", blockSize)
	}

	return nil
}

// ToJSON serializes the block to JSON
func (b *Block) ToJSON() []byte {
	data, _ := json.Marshal(b)
	return data
}

// FromJSON deserializes a block from JSON
func FromJSON(data []byte) (*Block, error) {
	var block Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %w", err)
	}
	return &block, nil
}

// String returns a string representation of the block
func (b *Block) String() string {
	return fmt.Sprintf("Block #%d (Hash: %s, Nonce: %d)", b.Index, b.Hash, b.Nonce)
}
