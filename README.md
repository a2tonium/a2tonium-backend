# A2Tonium Backend

A2Tonium Backend is a robust backend service designed to manage and process courses created on the TON Blockchain. This backend is built with resilience and simplicity in mind: it can be started at any time, automatically detects its last state, and continues processing seamlessly without relying on a traditional database. All data is stored permanently on the blockchain, ensuring immutability and transparency.

---

## Table of Contents

- [Project Overview](#project-overview)  
- [Features](#features)  
- [Architecture](#architecture)  
- [Configuration](#configuration)  
- [Getting Started](#getting-started)  
- [Usage](#usage)  
- [Makefile](#makefile)  
- [Project Structure](#project-structure)  
- [Logging](#logging)  
- [Contributing](#contributing)  
- [License](#license)  

---

## Project Overview

The A2Tonium Backend acts as the core processing engine for courses deployed on the TON Blockchain. It integrates with IPFS for decentralized file storage and uses TON wallets for blockchain interactions. The backend is designed to be stateless with no external database dependency; all state and data are retrieved and stored on the blockchain, making it highly fault-tolerant and scalable.

---

## Features

- **Blockchain-native:** All course data is stored on the TON Blockchain, ensuring permanence and security.
- **No external database:** Eliminates the need for a database by leveraging blockchain as the single source of truth.
- **Resumable processing:** The backend can be restarted anytime and will resume processing from the last known state.
- **IPFS integration:** Uses Pinata IPFS service for decentralized file pinning and retrieval.
- **TON wallet support:** Supports keypair generation from mnemonic phrases and interacts with the TON blockchain.
- **Configurable logging:** Supports different log levels and optional file logging with rotation.

---

## Architecture

The backend is composed of several internal services:

- **A2Tonium Service:** Core business logic for processing courses.
- **IPFS Service:** Handles file uploads and retrievals via Pinata IPFS.
- **TON Service:** Manages wallet initialization and blockchain interactions.
- **JSON Generator:** Generates JSON metadata for courses or related data.

The main application initializes these services and orchestrates their interactions.

---

## Configuration

The backend uses a YAML configuration file to manage its settings. Below is an example configuration structure:

```yaml
configs:
  # Logger configuration
  log_level: "info"           # Options: debug, info, warn, error
  log_into_file: false        # Enable/disable logging to a file
  log_file_path: ""           # Path to the log file if enabled
  log_file_max_size: 100      # Max log file size in MB
  log_file_backups: 15        # Number of backup log files to keep
  log_file_max_age: 1         # Max age of log files in days

  # IPFS configuration
  pinata_jwt_token: "your_jwt_token"  # JWT token for Pinata IPFS API

  # TON Wallet configuration
  mnemonic_phrase: "your_mnemonic_phrase"  # Seed phrase for wallet keypair generation
```

---

## Getting Started

### Prerequisites

- Go 1.18+ installed
- Access to a Pinata IPFS account for JWT token
- TON wallet mnemonic phrase

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-repo/a2tonium-backend.git
   cd a2tonium-backend
   ```

2. Configure your environment by copying and editing the config file:

   ```bash
   cp config/config_example.yml config/config_local.yml
   # Edit config/config_local.yml and fill in your Pinata JWT token and mnemonic phrase
   ```

3. Build the project:

   ```bash
   go build -o a2tonium-backend ./cmd/a2tonium
   ```

---

## Usage

This project includes a convenient `Makefile` to simplify running the backend and generating the public key. Using the Makefile ensures consistent command usage and reduces the chance of errors.

### Run the backend

To start the backend service with the local configuration, simply run:

```bash
make run
```

This command runs the backend with the `-config=local` flag, loading the configuration from `config/config_local.yml`.

### Generate Public Key

To generate the base64-encoded public key from your mnemonic phrase, use:

```bash
make generatePublicKey
```

This will output your public key derived from the mnemonic phrase specified in the configuration.

---

## Makefile

The `Makefile` contains the following targets:

```makefile
.PHONY: run
run:
	$(info Running...)
	go run ./cmd/a2tonium -config=local

generatePublicKey:
	$(info Public Key Generation...)
	go run ./cmd/a2tonium -config=local --generatePublicKey
```

- `make run` — runs the backend with the local config.
- `make generatePublicKey` — generates and prints the public key from the mnemonic.

---

## Project Structure

```
.
├── cmd
│   └── a2tonium           # Main application entrypoint
├── config                 # Configuration files
├── internal
│   ├── app                # Core application modules
│   │   ├── a2tonium       # Business logic for course processing
│   │   ├── ipfs           # IPFS integration
│   │   ├── json_generator # JSON metadata generation
│   │   └── ton            # TON blockchain interaction
│   └── infrastructure     # Supporting infrastructure code
├── pkg                    # Shared packages (config, logger, ton crypto)
├── Makefile               # Build and utility commands
├── go.mod                 # Go module file
└── README.md              # This file
```

---

## Logging

The backend supports configurable logging levels and optional logging to a rotating file. Configure logging in the YAML config file under the `configs` section.

---

## Contributing

Contributions are welcome! Please open issues or pull requests for bug fixes, features, or improvements.

---

## License

[MIT License](LICENSE)

---

If you have any questions or need further assistance, feel free to reach out!

---

**Enjoy building on the TON Blockchain with A2Tonium!** 🚀
