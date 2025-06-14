package miner

import (
	"blockchain-go/internal/blockchain"
	"blockchain-go/internal/config"
	"blockchain-go/internal/network"
	"log"
	"strconv"
	"time"
)

// Miner represents the mining process
type Miner struct {
	networkManager *network.Manager
	config         config.MinerConfig
	stopChan       chan struct{}
}

// NewMiner creates a new miner instance
func NewMiner(nm *network.Manager, cfg config.MinerConfig) *Miner {
	return &Miner{
		networkManager: nm,
		config:         cfg,
		stopChan:       make(chan struct{}),
	}
}

// Start begins the mining process
func Start(nm *network.Manager, cfg config.MinerConfig) {
	miner := NewMiner(nm, cfg)
	miner.Start()
}

// Start begins the mining process
func (m *Miner) Start() {
	log.Println("Starting miner...")

	mineCount := 0
	startTime := time.Now()
	restart := false

	for {
		select {
		case <-m.stopChan:
			log.Println("Miner stopped")
			return
		default:
			restart = false

			// Get latest block
			latestBlock, err := m.networkManager.GetBlockchain().GetLatestBlock()
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			// Calculate next block difficulty
			nextBlockDifficulty := latestBlock.NextBlockDifficulty
			if m.networkManager.GetBlockchain().ShouldRecalculateDifficulty() {
				newDifficulty, err := m.networkManager.GetBlockchain().CalculateNewDifficulty()
				if err != nil {
					log.Printf("Failed to calculate new difficulty: %v", err)
				} else {
					nextBlockDifficulty = newDifficulty
					log.Printf("New difficulty: %d", newDifficulty)
				}
			}

			// Ensure minimum difficulty
			difficultyForThisBlock := latestBlock.NextBlockDifficulty
			if difficultyForThisBlock < 1 {
				difficultyForThisBlock = 1
			}

			// Create new block
			block := blockchain.NewBlock(
				latestBlock.Index+1,
				difficultyForThisBlock,
				nextBlockDifficulty,
				"data",
				latestBlock.Hash,
			)

			// Mine the block
			for !block.IsHashValid(block.Hash) && !restart {
				block.Timestamp = time.Now().Unix()
				block.Nonce = int(time.Now().UnixNano() % int64(m.config.MaxNonce))
				block.Hash = block.ComputeHash()

				mineCount++

				// Check if we should sync with network
				if time.Since(startTime) > time.Duration(m.config.NetworkSyncInterval)*time.Second {
					restart = true

					timeSinceLastMiningRateCalc := time.Since(startTime)
					hashRate := float64(mineCount) / timeSinceLastMiningRateCalc.Seconds()
					log.Printf("Mining rate: %.2f H/s", hashRate)

					mineCount = 0
					startTime = time.Now()
				}
			}

			// Add mined block to blockchain
			if !restart {
				if err := m.networkManager.GetBlockchain().AddBlock(block); err != nil {
					log.Printf("Failed to add block: %v", err)
				} else {
					log.Printf("Mined block #%d (Hash: %s, Nonce: %d)",
						block.Index, block.Hash, block.Nonce)

					// Broadcast found block
					m.broadcastFoundBlock(block.Index)
				}
			}
		}
	}
}

// Stop stops the mining process
func (m *Miner) Stop() {
	close(m.stopChan)
}

// broadcastFoundBlock broadcasts a found block to the network
func (m *Miner) broadcastFoundBlock(blockIndex int) {
	blockIndexStr := strconv.Itoa(blockIndex)
	packet := network.NewBroadcastPacket(
		m.networkManager.GetMe(),
		network.PacketNameFoundBlock,
		[]byte(blockIndexStr),
		blockIndex,
	)

	packetData, err := packet.ToJSON()
	if err != nil {
		log.Printf("Failed to serialize found block packet: %v", err)
		return
	}

	m.networkManager.Broadcast(packetData)
	log.Printf("Broadcasted found block #%d", blockIndex)
}
