package dht

import (
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"testing"
)

func TestaddCommandToMessageHandler(t *testing.T) {
	hash := hashRequest("testmd5")
	ch := make(chan string)
	mh := message_handler.NewMessageHandler()
	addCommandToMessageHandler(hash, ch, mh)
}

// Oh, the things we'll do for those sweet, sweet coverage points
func TesthashRequest(t *testing.T) {
	expectedReturn := "32269AE63A25306BB46A03D6F38BD2B7"
	hash := hashRequest("testmd5")

	if hash != expectedReturn {
		t.Fatalf("Google has somehow managed to break md5. Abandon Go.")
	}

}
