package main

import (
	"flag"
	"log"

	"blockchain-go/internal/blockchain"
	"blockchain-go/internal/config"
	"blockchain-go/internal/miner"
	"blockchain-go/internal/network"
)

func main() {
	// Parse command line flags
	var (
		configPath = flag.String("config", "config.yaml", "Path to configuration file")
		initHost   = flag.String("init-host", "", "Initial peer host for joining network")
		initPort   = flag.Int("init-port", 0, "Initial peer port for joining network")
		port       = flag.Int("port", 8080, "Port to listen on")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Printf("Failed to load config: %v, using defaults", err)
		cfg = config.Default()
	}

	// Override port if specified
	if *port != 8080 {
		cfg.Network.Port = *port
	}

	// Create blockchain instance
	bc := blockchain.New(cfg.Blockchain.DifficultyCalculationBlocks, cfg.Blockchain.TargetBlockTime)

	// Create network manager
	var nm *network.Manager
	if *initHost != "" && *initPort != 0 {
		// Join existing network
		nm = network.NewJoiningManager(cfg.Network, bc, *initHost, *initPort)
	} else {
		// Create genesis block and start new network
		genesisBlock := bc.CreateGenesisBlock()
		log.Printf("Created genesis block: %s", genesisBlock.Hash)
		nm = network.NewManager(cfg.Network, bc)
	}

	// Start network server
	go nm.StartServer()
	log.Printf("P2P Network started on port %d", cfg.Network.Port)

	// Start mining
	miner.Start(nm, cfg.Miner)
}
