package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"hexblox/internal/blockchain"
	"hexblox/internal/wallet"
	"net/http"
	"time"
)

type Node struct {
	Blockchain      *blockchain.Blockchain
	Wallet          *wallet.Wallet
	TransactionPool *wallet.TransactionPool

	HttpServer *gin.Engine

	ctx     context.Context
	host    host.Host
	notifee *Notifee

	gossipPubSub  *pubsub.PubSub
	topics        map[string]*pubsub.Topic
	subscriptions map[string]*pubsub.Subscription
}

func Run(httpPort string, hostPort string) *Node {
	node := &Node{
		Blockchain:      blockchain.NewBlockchain(),
		Wallet:          wallet.NewWallet(),
		TransactionPool: wallet.NewTransactionPool(),
		ctx:             context.Background(),
		topics:          make(map[string]*pubsub.Topic),
		subscriptions:   make(map[string]*pubsub.Subscription),
	}
	node.initHost(hostPort)
	node.initSub()
	node.PropagateChain()
	node.initHttpServer(httpPort)

	return node
}

func (node *Node) initHttpServer(httpPort string) {
	httpServer := gin.Default()
	node.HttpServer = httpServer
	SetBlockchainRoutes(node)
	SetTransactionRoutes(node)
	httpServer.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from HTTP server")
	})

	go func() {
		if err := httpServer.Run(":" + httpPort); err != nil {
			panic(err)
		}
	}()
}

func (node *Node) initHost(hostPort string) {
	nodeHost, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", hostPort)),
	)
	if err != nil {
		fmt.Println("Failed to create host:", err)
		panic(err)
	}
	fmt.Printf("Host %s initialized on port: %s\n", nodeHost.ID(), hostPort)

	node.host = nodeHost
	node.notifee = &Notifee{h: nodeHost}
}

func (node *Node) initSub() {
	gossipPubSub, err := pubsub.NewGossipSub(node.ctx, node.host)
	if err != nil {
		panic(err)
	}
	node.gossipPubSub = gossipPubSub

	if err := node.setupDiscovery(); err != nil {
		panic(err)
	}

	// TODO: see if there is a way to implicitly wait for completion of setupDiscovery
	time.Sleep(5 * time.Second)

	node.joinRooms()
}

func (node *Node) joinRooms() {
	node.joinRoom("hexblox")
	node.joinRoom("hexblox-transaction-pool")
}

func (node *Node) joinRoom(room string) {
	topic, err := node.gossipPubSub.Join(room)
	if err != nil {
		panic(err)
	}
	node.topics[room] = topic

	subscriber, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}
	fmt.Println("Subscribed to topic:", topic.String())
	node.subscriptions[room] = subscriber

	go node.subscribe(room)
}

func (node *Node) subscribe(topic string) {
	subscription := node.subscriptions[topic]
	for {
		msg, err := subscription.Next(node.ctx)
		if err != nil {
			panic(err)
		}

		// only consider messages delivered by other peers
		if msg.ReceivedFrom == node.host.ID() {
			continue
		}

		fmt.Printf("Message reviced from: %s\n", msg.ReceivedFrom.String())
		switch topic {
		case "hexblox":
			node.syncChain(msg)
		case "hexblox-transaction-pool":
			node.syncTransactionPool(msg)
		}
	}
}

func (node *Node) setupDiscovery() error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(node.host, "hexblox-pubsub", node.notifee)
	return s.Start()
}

func (node *Node) sendMessage(topic string, message string) error {
	err := node.topics[topic].Publish(node.ctx, []byte(message))
	if err != nil {
		fmt.Println("Failed to sent message:", err)
		return err
	}
	fmt.Println("Message sent to topic:", topic)
	return nil
}

func (node *Node) PropagateChain() {
	jsonBlockchain, err := json.Marshal(node.Blockchain.Chain())
	if err != nil {
		panic(err)
	}

	err = node.sendMessage("hexblox", string(jsonBlockchain))
	if err != nil {
		panic(err)
	}
}

func (node *Node) syncChain(message *pubsub.Message) {
	messageData := string(message.Data)
	var newChain []*blockchain.Block

	if err := json.Unmarshal([]byte(messageData), &newChain); err != nil {
		fmt.Println("Error:", err)
		return
	}
	node.Blockchain.ReplaceChain(newChain)
}

func (node *Node) PropagateTransaction(transaction *wallet.Transaction) {
	jsonTransaction, err := json.Marshal(transaction)
	if err != nil {
		panic(err)
	}

	err = node.sendMessage("hexblox-transaction-pool", string(jsonTransaction))
}

func (node *Node) syncTransactionPool(message *pubsub.Message) {
	messageData := string(message.Data)
	var newTransaction *wallet.Transaction

	if err := json.Unmarshal([]byte(messageData), &newTransaction); err != nil {
		fmt.Println("Error:", err)
		return
	}

	node.TransactionPool.Transactions = append(node.TransactionPool.Transactions, newTransaction)
	fmt.Println("Transaction pool successfully updated.")
}
