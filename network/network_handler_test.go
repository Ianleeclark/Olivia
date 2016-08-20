package networkHandler

import (
	"bufio"
	"fmt"
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/config"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"net"
	"os"
	"testing"
	"time"
)

func sendCommand(command string, nodeIP string, t *testing.T) string {
	conn, err := net.DialTimeout("tcp", nodeIP, 1*time.Second)
	if err != nil {
		os.Exit(-1)
	}
	conn.Write([]byte(fmt.Sprintf("%s\n", command)))

	reader := bufio.NewReader(conn)
	str, err := reader.ReadString('\n')
	if err != nil {
		os.Exit(-2)
	}

	return str
}

func TestGetRemoteKey(t *testing.T) {
	sendCommand("SET key1:value1,key2:value2,key3:value3\n", "127.0.0.1:5555", t)
	response := sendCommand("GET key1\n", "127.0.0.1:5454", t)

	expectedReturn := "GOT key1:value1"

	if response != expectedReturn {
		t.Fatalf("Expected %v, got %v", expectedReturn, response)
	}
}

func TestMain(m *testing.M) {
	stopChan := make(chan struct{})
	stopChan2 := make(chan struct{})
	mh := message_handler.NewMessageHandler()
	cache := cache.NewCache()
	cfg := config.ReadConfig()
	secondCfg := config.ReadConfig()
	secondCfg.BaseNode = false
	secondCfg.RemotePeers = append(secondCfg.RemotePeers, "127.0.0.1:5454")
	secondCfg.ListenPort = 5555

	go StartIncomingNetwork(mh, cache, cfg, stopChan)
	go StartIncomingNetwork(mh, cache, secondCfg, stopChan2)
	os.Exit(m.Run())
	stopChan <- struct{}{}
	stopChan2 <- struct{}{}
}
