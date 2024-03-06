package p2p

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"hexblox/internal/api/routes"
	"hexblox/internal/blockchain"
	"net/http"
	"time"
)

type Node struct {
	blockchain *blockchain.Blockchain[string]
	host       host.Host
	gossipPub  *pubsub.PubSub
	httpServer *gin.Engine
}

func Run(httpPort string, hostPort string) *Node {
	bc := blockchain.NewBlockchain[string]()
	httpServer := createHttpServer(httpPort, bc)
	nodeHost, gossipSub, _ := createHost(hostPort)

	return &Node{
		blockchain: bc,
		host:       nodeHost,
		gossipPub:  gossipSub,
		httpServer: httpServer,
	}
}

func createHttpServer(httpPort string, bc *blockchain.Blockchain[string]) *gin.Engine {
	httpServer := gin.Default()
	routes.SetBlockchainRoutes(httpServer, bc)
	httpServer.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from HTTP server")
	})

	go func() {
		if err := httpServer.Run(":" + httpPort); err != nil {
			panic(err)
		}
	}()

	return httpServer
}

func createHost(hostPort string) (host.Host, *pubsub.PubSub, *pubsub.Subscription) {
	ctx := context.Background()
	nodeHost, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", hostPort)),
	)
	if err != nil {
		fmt.Println("Failed to create host:", err)
		panic(err)
	}

	gossipSub, subscription := createSub(ctx, nodeHost)

	nodeHost.SetStreamHandler("/test", handleStream)

	return nodeHost, gossipSub, subscription
}

func handleStream(stream network.Stream) {
	fmt.Println("New stream received")

	// Read the incoming message from the stream
	reader := bufio.NewReader(stream)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read message from stream:", err)
		return
	}

	// Process the incoming message
	fmt.Println("Received message:", message)

	// Optionally, you can send a response back to the peer
	// response := "Response from host"
	// _, err = stream.Write([]byte(response))
	// if err != nil {
	//     fmt.Println("Failed to send response:", err)
	//     return
	// }

	// Close the stream
	err = stream.Close()
	if err != nil {
		fmt.Println("Failed to close stream:", err)
	}
}

func createSub(ctx context.Context, host host.Host) (*pubsub.PubSub, *pubsub.Subscription) {
	gossipSub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}
	if err := setupDiscovery(host); err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	room := "hexblox"
	topic, err := gossipSub.Join(room)
	if err != nil {
		panic(err)
	}
	subscriber, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}
	fmt.Println("Subscribed to topic:", topic.String())
	go subscribe(subscriber, ctx, host.ID())

	err = topic.Publish(ctx, []byte("Hello, peer!"))
	if err != nil {
		panic(err)
	}

	return gossipSub, subscriber
}

func subscribe(subscriber *pubsub.Subscription, ctx context.Context, hostID peer.ID) {
	for {
		msg, err := subscriber.Next(ctx)
		if err != nil {
			panic(err)
		}

		// only consider messages delivered by other peers
		if msg.ReceivedFrom == hostID {
			continue
		}

		fmt.Printf("got message: %s, from: %s\n", string(msg.Data), msg.ReceivedFrom.String())
	}
}

func setupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, "hexblox-pubsub", &Notifee{h: h})
	return s.Start()
}
