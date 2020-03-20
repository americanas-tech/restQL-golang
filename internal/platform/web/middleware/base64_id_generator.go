package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"sync/atomic"
)

type Base64IdGenerator struct {
	prefix  string
	counter *uint64
}

func NewBase64IdGenerator() *Base64IdGenerator {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.URLEncoding.EncodeToString(buf[:])
	}

	prefix := fmt.Sprintf("%s/%s", hostname, b64[0:10])
	var initialCounter uint64

	return &Base64IdGenerator{prefix: prefix, counter: &initialCounter}
}

func (i *Base64IdGenerator) Run() string {
	c := atomic.AddUint64(i.counter, 1)
	return fmt.Sprintf("%s-%06d", i.prefix, c)
}
