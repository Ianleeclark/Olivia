package dht

import (
	"log"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/GrappigPanda/Olivia/bloomfilter"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"github.com/GrappigPanda/Olivia/network/receiver"
	"github.com/GrappigPanda/Olivia/parser"
	"net"
	"time"
)

// State represents the state that the remote peer is in.
type State int

const (
	// Disconnected signifies that the remote node is not yet connected.
	Disconnected State = iota
	// Connected signifies that there is a working connection between the
	// remote peer and our current node.
	Connected
	// Timeout signifies that the remote node has timed out and a
	// connection couldn't be established.
	Timeout
)

// Peer Houses the state for remote Peers
type Peer struct {
	Status      State
	Conn        *net.Conn
	IPPort      string
	BloomFilter *olilib.BloomFilter
	MessageBus  *message_handler.MessageHandler
}

// NewPeer handles creating a new peer to be used in communicating between nodes
func NewPeer(conn *net.Conn, mh *message_handler.MessageHandler) *Peer {
	ipPort := (*conn).RemoteAddr().String()
	log.Println("New peer connected: %v", ipPort)

	return &Peer{
		Disconnected,
		conn,
		ipPort,
		nil,
		mh,
	}
}

// NewPeerByIP handles creating a peer by its ip, opening a connection, &c.
func NewPeerByIP(ipPort string, mh *message_handler.MessageHandler) *Peer {
	newPeer := &Peer{
		Disconnected,
		nil,
		ipPort,
		nil,
		mh,
	}

	return newPeer
}

// Connect opens a connection to a remote peer
func (p *Peer) Connect() error {
	conn, err := net.DialTimeout("tcp", p.IPPort, 5*time.Second)
	if err != nil {
		if err, _ := err.(net.Error); err.Timeout() {
			p.Status = Timeout
		}
		return err
	}

	p.Conn = &conn
	p.Status = Connected
	p.GetBloomFilter()

	return nil
}

// Disconnect closes a connection to a remote peer.
func (p *Peer) Disconnect() {
	(*p.Conn).Close()
}

// SendCommand Handles sending a command to a remote node. Command is like this
// "hash:Command"
func (p *Peer) SendCommand(Command string) {
	(*p.Conn).Write([]byte(Command))
}

// SendRequest handles taking in a peer object and a command and sending a
// command which will be responded to the calling channel once the request has
// been fulfilled
func (p *Peer) SendRequest(Command string, responseChannel chan string, mh *message_handler.MessageHandler) {
	receiver := network_receiver.NewReceiver(mh, p.Conn)

	hash := hashRequest(Command)
	addCommandToMessageHandler(hash, responseChannel, mh)

	go func() {
		receiver.Run()
	}()

	p.SendCommand(fmt.Sprintf("%s:%s\n", hash, Command))
}

// GetBloomFilter handles retrieving a remote node's bloom filter.
func (p *Peer) GetBloomFilter() {
	responseChannel := make(chan string)

	go p.SendRequest(
		parser.GET_REMOTE_BLOOMFILTER,
		responseChannel,
		p.MessageBus,
	)

	go func() {
		parser := parser.NewParser(p.MessageBus)
		response := <-responseChannel

		responseData, err := parser.Parse(response, p.Conn)
		if err != nil {
			log.Println(err)
			return
		}

		for k, _ := range responseData.Args {
			bf, err := olilib.ConvertStringtoBF(responseData.Args[k])
			if err != nil {
				p.BloomFilter = nil
			}
			p.BloomFilter = bf
			break
		}

	}()
}

// GetPeerListAsync handles retrieving all known peers from a remote node.
func (p *Peer) GetPeerList(responseChannel chan string) {
	p.SendRequest(parser.GET_REMOTE_PEERLIST, responseChannel, p.MessageBus)
}


// addCommandToMessageHandler send a command to the message container to store
// the callback channel.
func addCommandToMessageHandler(hash string, responseChannel chan string, mh *message_handler.MessageHandler) {
	keyVal := message_handler.NewKeyValPair(hash, responseChannel, nil)

	mh.AddKeyChannel <- keyVal
}

// hashRequest hashes the command so that later the channel can be responded to
// from the message container
func hashRequest(Command string) string {
	hasher := md5.New()
	hasher.Write([]byte(time.Now().UTC().String()))
	hasher.Write([]byte(Command))

	return hex.EncodeToString(hasher.Sum(nil))
}
