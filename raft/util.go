package raft

import (
	"crypto/sha1"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"os"
	"time"
)

// OpenPort creates a listener on the specified port.
func OpenPort(port int) (net.Listener, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	conn, err := net.Listen("tcp4", fmt.Sprintf("%v:%v", hostname, port))
	return conn, err
}

// AddrToId converts a network address to a Raft node ID of specified length.
func AddrToId(addr string, length int) string {
	h := sha1.New()
	h.Write([]byte(addr))
	v := h.Sum(nil)
	keyInt := big.Int{}
	keyInt.SetBytes(v[:length])
	return keyInt.String()
}

// randomTimeout uses time.After to create a timeout between minTimeout and 2x that.
func randomTimeout(minTimeout time.Duration) <-chan time.Time {
	maxTimeout := 2 * int(minTimeout)
	rand.Seed(time.Now().UTC().UnixNano()) // rand seed
	dur := time.Duration(rand.Intn(maxTimeout-int(minTimeout)) + int(minTimeout))
	return time.After(dur)
}

// createCacheId creates a unique ID to store a client request and corresponding
// reply in cache.
func createCacheId(clientId, sequenceNum uint64) string {
	return fmt.Sprintf("%v-%v", clientId, sequenceNum)
}
