# Hexblox

**Hexblox** is a from-scratch implementation of a proof-of-work (PoW) blockchain system written in Go. It provides a minimal and educational yet functional blockchain node with core features such as block mining, peer-to-peer networking, and an HTTP API for interaction.

The project is designed to demonstrate fundamental blockchain concepts, network communication using libp2p, and building decentralized systems without relying on external blockchain frameworks.

---

## Features

- Proof-of-Work consensus algorithm
- Block structure with transactions and timestamps
- Mining logic and nonce validation
- Full peer-to-peer (P2P) communication using libp2p
  - PubSub via Gossipsub
  - Peer discovery using mDNS
- HTTP API for external interaction
  - Submit new transactions
  - View blockchain state
  - Trigger mining
- Stateless architecture, with logs printed to terminal

---

## Technologies & Libraries

- **Go** – core implementation language
- **Gin Gonic** – lightweight HTTP server framework used for exposing the REST API
- **libp2p** – peer-to-peer networking stack for node discovery and messaging
  - Gossipsub – message propagation across the P2P network
  - mDNS – local peer discovery on the same network

---

## Getting Started

### Prerequisites

- Go 1.20 or later installed
- Open ports (or different values) for HTTP and P2P services

### Installation & Running

Clone the repository:

```bash
git clone https://github.com/cuturic01/hexblox.git
cd hexblox
```

Run a node by specifying both the HTTP and P2P ports:

```bash
go run main.go <http-port> <p2p-port>
```

Example:

```bash
go run main.go 8080 9000
```
You can now interact with the node via http://localhost:8080.

To run multiple nodes locally, start each instance with a different combination of ports.
