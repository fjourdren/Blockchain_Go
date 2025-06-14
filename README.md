# Blockchain Go

A clean, production-ready blockchain implementation in Go with P2P networking, mining, and simplified architecture.

## Features

- **Blockchain Core**: Complete blockchain implementation with proof-of-work mining
- **P2P Network**: Distributed peer-to-peer network for block synchronization
- **Dynamic Difficulty**: Automatic difficulty adjustment based on mining rate
- **Configuration Management**: YAML-based configuration with sensible defaults
- **Clean Architecture**: Well-organized, modular code structure
- **Thread Safety**: Proper synchronization for concurrent operations

## Project Structure

```
blockchain-go/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── blockchain/
│   │   ├── block.go           # Block implementation
│   │   └── blockchain.go      # Blockchain core logic
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── miner/
│   │   └── miner.go           # Mining implementation
│   └── network/
│       ├── broadcast.go       # Broadcast packet management
│       ├── manager.go         # Network manager
│       ├── packet.go          # Network packet definitions
│       └── peer.go            # Peer implementation
├── config.yaml                # Default configuration
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd blockchain-go
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
# Start a new blockchain network
go run cmd/main.go

# Join an existing network
go run cmd/main.go -init-host 127.0.0.1 -init-port 8080 -port 8081

# Use custom configuration
go run cmd/main.go -config custom-config.yaml
```

## Configuration

The application uses YAML configuration files. See `config.yaml` for the default configuration:

### Blockchain Configuration
- `difficulty_calculation_blocks`: Number of blocks to consider for difficulty calculation
- `target_block_time`: Target time between blocks in seconds

### Network Configuration
- `host`: Network host address
- `port`: Network port

### Miner Configuration
- `network_sync_interval`: Interval for network synchronization during mining
- `max_nonce`: Maximum nonce value for mining

## Usage

### Starting a New Network

```bash
go run cmd/main.go
```

This will:
1. Create a genesis block
2. Start the P2P network server
3. Begin mining new blocks

### Joining an Existing Network

```bash
go run cmd/main.go -init-host 127.0.0.1 -init-port 8080 -port 8081
```

This will:
1. Connect to the specified peer
2. Download the existing blockchain
3. Start participating in the network

### Command Line Options

- `-config <path>`: Path to configuration file (default: config.yaml)
- `-init-host <host>`: Initial peer host for joining network
- `-init-port <port>`: Initial peer port for joining network
- `-port <port>`: Port to listen on (default: 8080)

## Architecture

### Blockchain Package
- **Block**: Represents a single block with validation and mining capabilities
- **Blockchain**: Manages the chain of blocks with difficulty adjustment and validation

### Network Package
- **Peer**: Represents a network peer with TCP communication
- **Packet**: Network packet definitions for P2P communication
- **Manager**: Handles network operations, peer management, and synchronization
- **BroadcastManager**: Manages broadcast packet deduplication

### Miner Package
- **Miner**: Implements the proof-of-work mining algorithm with network synchronization

### Config Package
- **Config**: Manages application configuration with file loading and defaults

## Key Features

### Code Organization
- **Modular Design**: Clear separation of concerns with dedicated packages
- **Clean Interfaces**: Well-defined interfaces and abstractions
- **Proper Naming**: Go naming conventions throughout the codebase

### Thread Safety
- **Mutex Protection**: Proper synchronization for concurrent operations
- **Channel Communication**: Safe communication between goroutines

### Configuration Management
- **YAML Configuration**: Human-readable configuration format
- **Default Values**: Sensible defaults for all settings
- **Runtime Override**: Command-line options for configuration override

### Network Protocol
- **Structured Packets**: Well-defined packet types and formats
- **Reliable Communication**: TCP-based reliable communication
- **Broadcast Deduplication**: Prevents duplicate broadcast processing

## Development

### Adding New Features

1. **New Block Types**: Extend the `Block` struct in `internal/blockchain/block.go`
2. **New Network Messages**: Add packet types in `internal/network/packet.go`
3. **New Configuration**: Extend the config structs in `internal/config/config.go`

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/blockchain
```

### Building

```bash
# Build for current platform
go build -o blockchain-go cmd/main.go

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o blockchain-go cmd/main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request