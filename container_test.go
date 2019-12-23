package container_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/teamlint/container"
)

func TestContainer(t *testing.T) {
	var HTTPBundle = container.Bundle(
		container.Provide(ProvideAddr("0.0.0.0", "8080")),
		container.Provide(NewMux, container.As(new(http.Handler))),
		container.Provide(NewHTTPServer, container.Prototype(), container.WithName("server")),
	)

	c := container.New(HTTPBundle)

	var server1 *http.Server
	err := c.Extract(&server1, container.Name("server"))
	require.NoError(t, err)

	var server2 *http.Server
	err = c.Extract(&server2, container.Name("server"))
	require.NoError(t, err)

	err = c.Invoke(PrintAddr)
	require.NoError(t, err)
}

// Addr
type Addr string

// ProvideAddr
func ProvideAddr(host string, port string) func() Addr {
	return func() Addr {
		return Addr(net.JoinHostPort(host, port))
	}
}

// NewHTTPServer
func NewHTTPServer(addr Addr, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    string(addr),
		Handler: handler,
	}
}

// NewMux
func NewMux() *http.ServeMux {
	return &http.ServeMux{}
}

// PrintAddr
func PrintAddr(addr Addr) {
	fmt.Println(addr)
}
