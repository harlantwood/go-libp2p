package mux

import (
	"bytes"
	"testing"

	inet "github.com/jbenet/go-ipfs/p2p/net"
	protocol "github.com/jbenet/go-ipfs/p2p/protocol"
)

var testCases = map[string]string{
	"/bitswap": "\u0009/bitswap\n",
	"/dht":     "\u0005/dht\n",
	"/ipfs":    "\u0006/ipfs\n",
	"/ipfs/dksnafkasnfkdajfkdajfdsjadosiaaodj": ")/ipfs/dksnafkasnfkdajfkdajfdsjadosiaaodj\n",
}

func TestWrite(t *testing.T) {
	for k, v := range testCases {
		var buf bytes.Buffer
		if err := protocol.WriteHeader(&buf, protocol.ID(k)); err != nil {
			t.Fatal(err)
		}

		v2 := buf.Bytes()
		if !bytes.Equal(v2, []byte(v)) {
			t.Errorf("failed: %s - %v != %v", k, []byte(v), v2)
		}
	}
}

func TestHandler(t *testing.T) {

	outs := make(chan string, 10)

	h := func(n string) func(s inet.Stream) {
		return func(s inet.Stream) {
			outs <- n
		}
	}

	m := Mux{Handlers: StreamHandlerMap{}}
	m.Default = h("default")
	m.Handlers["/dht"] = h("bitswap")
	// m.Handlers["/ipfs"] = h("bitswap") // default!
	m.Handlers["/bitswap"] = h("bitswap")
	m.Handlers["/ipfs/dksnafkasnfkdajfkdajfdsjadosiaaodj"] = h("bitswap")

	for k, v := range testCases {
		var buf bytes.Buffer
		if _, err := buf.Write([]byte(v)); err != nil {
			t.Error(err)
			continue
		}

		name, err := protocol.ReadHeader(&buf)
		if err != nil {
			t.Error(err)
			continue
		}

		if name != protocol.ID(k) {
			t.Errorf("name mismatch: %s != %s", k, name)
			continue
		}
	}

}