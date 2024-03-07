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
	"hexblox/internal/api/routes"
	"hexblox/internal/blockchain"
	"net/http"
	"time"
)

type Node struct {
	blockchain *blockchain.Blockchain[string]

	httpServer *gin.Engine

	ctx     context.Context
	host    host.Host
	notifee *Notifee

	gossipPubSub  *pubsub.PubSub
	topics        map[string]*pubsub.Topic
	subscriptions map[string]*pubsub.Subscription
}

func Run(httpPort string, hostPort string) *Node {
	node := &Node{
		blockchain:    blockchain.NewBlockchain[string](),
		ctx:           context.Background(),
		topics:        make(map[string]*pubsub.Topic),
		subscriptions: make(map[string]*pubsub.Subscription),
	}
	node.initHttpServer(httpPort)
	node.initHost(hostPort)
	node.initSub()
	node.propagateChain()

	return node
}

func (node *Node) initHttpServer(httpPort string) {
	httpServer := gin.Default()
	routes.SetBlockchainRoutes(httpServer, node.blockchain)
	httpServer.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from HTTP server")
	})

	go func() {
		if err := httpServer.Run(":" + httpPort); err != nil {
			panic(err)
		}
	}()

	node.httpServer = httpServer
}

func (node *Node) initHost(hostPort string) {
	nodeHost, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", hostPort)),
	)
	if err != nil {
		fmt.Println("Failed to create host:", err)
		panic(err)
	}
	fmt.Println(fmt.Sprintf("Host %s initialized on port: %s", nodeHost.ID(), hostPort))

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

	room := "hexblox"
	topic, err := gossipPubSub.Join(room)
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

		fmt.Printf("got message: %s, from: %s\n", string(msg.Data), msg.ReceivedFrom.String())
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

func (node *Node) propagateChain() {
	jsonBlockchain, err := json.Marshal(node.blockchain.Chain())
	if err != nil {
		panic(err)
	}

	err = node.sendMessage("hexblox", string(jsonBlockchain))
	if err != nil {
		panic(err)
	}
}
