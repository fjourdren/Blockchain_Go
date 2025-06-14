package blockchain

import (
	"fmt"
	"log"
	"sync"
)

// Blockchain represents the main blockchain structure
type Blockchain struct {
	mu                          sync.RWMutex
	difficultyCalculationBlocks int
	targetBlockTime             int
	chain                       []*Block
}

// New creates a new blockchain instance
func New(difficultyCalculationBlocks, targetBlockTime int) *Blockchain {
	return &Blockchain{
		difficultyCalculationBlocks: difficultyCalculationBlocks,
		targetBlockTime:             targetBlockTime,
		chain:                       make([]*Block, 0),
	}
}

// CreateGenesisBlock creates and mines the genesis block
func (bc *Blockchain) CreateGenesisBlock() *Block {
	log.Println("Mining genesis block...")

	genesis := NewBlock(0, 2, 2, "Genesis Block", "")
	genesis.Mine()
	
	bc.AddBlockWithoutVerification(genesis)

	log.Printf("Genesis block created: %s", genesis.Hash)
	return genesis
}

// CanAddBlock checks if a block can be added to the chain
func (bc *Blockchain) CanAddBlock(block *Block) error {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if len(bc.chain) == 0 {
		return fmt.Errorf("blockchain is empty")
	}

	latestBlock := bc.chain[len(bc.chain)-1]

	// Check difficulty
	if block.Difficulty != latestBlock.NextBlockDifficulty {
		return fmt.Errorf("block difficulty %d does not match expected %d",
			block.Difficulty, latestBlock.NextBlockDifficulty)
	}

	// Check index
	if block.Index != latestBlock.Index+1 {
		return fmt.Errorf("block index %d is not sequential", block.Index)
	}

	// Check previous hash
	if block.PreviousHash != latestBlock.Hash {
		return fmt.Errorf("block previous hash does not match latest block hash")
	}

	// Validate block
	if err := block.IsValid(); err != nil {
		return fmt.Errorf("block validation failed: %w", err)
	}

	return nil
}

// AddBlock adds a block to the chain with verification
func (bc *Blockchain) AddBlock(block *Block) error {
	if err := bc.CanAddBlock(block); err != nil {
		return fmt.Errorf("cannot add block: %w", err)
	}

	bc.AddBlockWithoutVerification(block)

	log.Printf("Added block #%d to chain - Hash: %s, Difficulty: %d, Nonce: %d", 
		block.Index, block.Hash, block.Difficulty, block.Nonce)
	return nil
}

// AddBlockWithoutVerification adds a block without validation (for syncing)
func (bc *Blockchain) AddBlockWithoutVerification(block *Block) {
	bc.mu.Lock()
	bc.chain = append(bc.chain, block)
	bc.mu.Unlock()
}

// HasBlock checks if a block with the given index exists
func (bc *Blockchain) HasBlock(index int) bool {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return index >= 0 && index < len(bc.chain)
}

// GetBlock returns the block at the given index
func (bc *Blockchain) GetBlock(index int) (*Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if index < 0 || index >= len(bc.chain) {
		return nil, fmt.Errorf("block index %d out of range", index)
	}

	return bc.chain[index], nil
}

// GetLatestBlock returns the most recent block
func (bc *Blockchain) GetLatestBlock() (*Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if len(bc.chain) == 0 {
		return nil, fmt.Errorf("blockchain is empty")
	}

	return bc.chain[len(bc.chain)-1], nil
}

// GetChainLength returns the number of blocks in the chain
func (bc *Blockchain) GetChainLength() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return len(bc.chain)
}

// IsValid validates the entire blockchain
func (bc *Blockchain) IsValid() error {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if len(bc.chain) == 0 {
		return fmt.Errorf("blockchain is empty")
	}

	for i := 1; i < len(bc.chain); i++ {
		currentBlock := bc.chain[i]
		previousBlock := bc.chain[i-1]

		// Validate individual block
		if err := currentBlock.IsValid(); err != nil {
			return fmt.Errorf("block #%d is invalid: %w", currentBlock.Index, err)
		}

		// Check chain integrity
		if currentBlock.PreviousHash != previousBlock.Hash {
			return fmt.Errorf("chain integrity broken at block #%d", currentBlock.Index)
		}
	}

	return nil
}

// CalculateAverageMiningTime calculates the average mining time for a range of blocks
func (bc *Blockchain) CalculateAverageMiningTime(startBlock, stopBlock int) (int, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if startBlock < 0 || stopBlock >= len(bc.chain) || startBlock >= stopBlock {
		return 0, fmt.Errorf("invalid block range: %d to %d", startBlock, stopBlock)
	}

	totalTime := int64(0)
	blockCount := stopBlock - startBlock

	for i := startBlock + 1; i <= stopBlock; i++ {
		currentBlock := bc.chain[i]
		previousBlock := bc.chain[i-1]
		totalTime += currentBlock.Timestamp - previousBlock.Timestamp
	}

	return int(totalTime / int64(blockCount)), nil
}

// ShouldRecalculateDifficulty checks if difficulty should be recalculated
func (bc *Blockchain) ShouldRecalculateDifficulty() bool {
	latestBlock, err := bc.GetLatestBlock()
	if err != nil {
		return false
	}

	return latestBlock.Index > 0 &&
		latestBlock.Index%bc.difficultyCalculationBlocks == 0
}

// CalculateNewDifficulty calculates the new difficulty based on recent mining times
func (bc *Blockchain) CalculateNewDifficulty() (int, error) {
	latestBlock, err := bc.GetLatestBlock()
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block: %w", err)
	}

	startIndex := latestBlock.Index - bc.difficultyCalculationBlocks
	if startIndex < 0 {
		startIndex = 0
	}

	averageTime, err := bc.CalculateAverageMiningTime(startIndex, latestBlock.Index)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate average mining time: %w", err)
	}

	currentDifficulty := latestBlock.Difficulty
	targetTime := bc.targetBlockTime

	// Adjust difficulty based on average mining time
	if averageTime > int(float64(targetTime)*1.25) {
		// Mining is too slow, decrease difficulty
		if currentDifficulty > 0 {
			currentDifficulty--
		}
	} else if averageTime < int(float64(targetTime)*0.75) {
		// Mining is too fast, increase difficulty
		currentDifficulty++
	}

	return currentDifficulty, nil
}

// GetBlocks returns a slice of blocks in the specified range
func (bc *Blockchain) GetBlocks(startIndex, endIndex int) ([]*Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if startIndex < 0 || endIndex >= len(bc.chain) || startIndex > endIndex {
		return nil, fmt.Errorf("invalid block range: %d to %d", startIndex, endIndex)
	}

	blocks := make([]*Block, endIndex-startIndex+1)
	copy(blocks, bc.chain[startIndex:endIndex+1])
	return blocks, nil
}

// String returns a string representation of the blockchain
func (bc *Blockchain) String() string {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return fmt.Sprintf("Blockchain with %d blocks", len(bc.chain))
}

// ReplaceChain replaces the local blockchain with a new chain (used for syncing from a peer)
func (bc *Blockchain) ReplaceChain(newChain []*Block) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	
	bc.chain = make([]*Block, len(newChain))
	copy(bc.chain, newChain)
}
