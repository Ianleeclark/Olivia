package incomingNetwork

import (
	"bufio"
	"fmt"
	"github.com/GrappigPanda/Olivia/bloomfilter"
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/config"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

var BASENODE = "127.0.0.1:5454"
var CONFIG = config.ReadConfig()

func sendCommand(command string, t *testing.T) string {
	conn, err := net.DialTimeout("tcp", BASENODE, 1*time.Second)
	if err != nil {
		t.Errorf("%v", err)
	}
	conn.Write([]byte(fmt.Sprintf("%s\n", command)))

	reader := bufio.NewReader(conn)
	str, err := reader.ReadString('\n')
	if err != nil {
		t.Errorf("%v", err)
	}

	return str
}

func TestGetBloomfilter(t *testing.T) {
	// All this test needs to do is assure us that we get a properly formatted
	// bloom filter in the return. We don't need to worry about if it matches
	// anything, as we have tests for that already.
	str := sendCommand("REQUEST bloomfilter\n", t)

	bf_str := strings.Split(str, " ")
	inputStr := strings.TrimSpace(bf_str[1])
	_, err := bloomfilter.Deserialize(inputStr, uint(CONFIG.BloomfilterSize))
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestPingPong(t *testing.T) {
	pong := sendCommand("PING 0", t)
	if pong != "0:PONG 1\n" {
		t.Errorf("Expected 0:PONG 1, got %v", pong)
	}
}

func TestSetKey(t *testing.T) {
	expectedReturn := ":SAT key1:value1\n"
	retVal := sendCommand("SET key1:value1\n", t)

	if expectedReturn != retVal {
		t.Errorf("Expected %v, got %v", expectedReturn, retVal)
	}

	// I intentionally don't test sending in multiple keys here, as I test it
	// in the bloom filter update below.
}

func TestSetKeyUpdatesBloomFilter(t *testing.T) {
	sendCommand("SET key1:value1,key2:value2,key3:value\n", t)

	str := sendCommand("REQUEST bloomfilter\n", t)

	bf_str := strings.Split(str, " ")
	inputStr := strings.TrimSpace(bf_str[1])
	bf, err := bloomfilter.Deserialize(inputStr, uint(CONFIG.BloomfilterSize))
	if err != nil {
		t.Errorf("%v", err)
	}

	if ok, _ := bf.HasKey([]byte("key1")); !ok {
		t.Errorf("Bloom filter doesn't contain correct keys %v", bf.Serialize())
	}

	if ok, _ := bf.HasKey([]byte("key2")); !ok {
		t.Errorf("Bloom filter doesn't contain correct keys %v", bf.Serialize())
	}

	if ok, _ := bf.HasKey([]byte("key3")); !ok {
		t.Errorf("Bloom filter doesn't contain correct keys %v", bf.Serialize())
	}
}

func TestGetKeyFromRemoteNode(t *testing.T) {
	// It's not yet possible to do this, for this to be done, we need to first
	// add listening ports and base nodes to be a part of the config file.
}

func TestMain(m *testing.M) {
	mh := message_handler.NewMessageHandler()
	CONFIG.IsTesting = true
	cache := cache.NewCache(mh, CONFIG)
	config := config.ReadConfig()
	stopchan := StartNetworkRouter(mh, cache, config)

	os.Exit(m.Run())

	stopchan <- struct{}{}
}
