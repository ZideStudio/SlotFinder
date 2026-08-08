package server

import (
	"app/config"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_PanicsOnListenError(t *testing.T) {
	// Occupy a port so Init's r.Run fails synchronously with "address
	// already in use", instead of blocking forever on the success path.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	defer func() { _ = listener.Close() }()

	addr := listener.Addr().(*net.TCPAddr)

	origHost := os.Getenv("APP_HOST")
	origPort := os.Getenv("APP_PORT")
	defer func() {
		_ = os.Setenv("APP_HOST", origHost)
		_ = os.Setenv("APP_PORT", origPort)
		config.Init()
	}()

	_ = os.Setenv("APP_HOST", addr.IP.String())
	_ = os.Setenv("APP_PORT", strconv.Itoa(addr.Port))
	config.Init()

	assert.Panics(t, func() {
		Init()
	})
}
